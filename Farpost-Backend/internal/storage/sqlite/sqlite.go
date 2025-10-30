package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"unicode"
	"vlru-prsch/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s failed to open db: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) FindStreets(substr string) ([]string, error) {
	const op = "storage.sqlite.FindStreets"

	if len(substr) == 0 {
		return []string{}, nil
	}

	lowSubstr := strings.ToLower(substr)
	normSubstr := func(str string) string {
		str = strings.ToLower(str)
		runes := []rune(str)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	}(substr)

	rows, err := s.db.Query(`
        SELECT name FROM streets 
        WHERE name LIKE ? OR name LIKE ?`,
		"%"+lowSubstr+"%", normSubstr+"%")
	if err != nil {
		return []string{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var streets []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return []string{}, fmt.Errorf("%s: %w", op, err)
		}
		streets = append(streets, name)
	}

	if err := rows.Err(); err != nil {
		return []string{}, fmt.Errorf("%s: %w", op, err)
	}

	return streets, nil
}

func (s *Storage) GetBlackouts(currentTime string) ([]models.Blackout, error) {
	const op = "storage.sqlite.GetBlackouts"

	queryTime := currentTime
    if len(currentTime) == 10 {
        queryTime = currentTime + " 23:59:59"
		currentTime = currentTime + " 00:00:00"
    }

	rows, err := s.db.Query(`
        SELECT id, start_date, end_date, description, type, initiator_name, source 
        FROM blackouts 
        WHERE start_date <= ? AND (end_date >= ? OR end_date IS NULL)
        ORDER BY start_date DESC`,
		queryTime, currentTime)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var blackouts []models.Blackout
	for rows.Next() {
		var blackout models.Blackout
		var endDate sql.NullString
		var source sql.NullString

		err := rows.Scan(
			&blackout.ID,
			&blackout.StartDate,
			&endDate,
			&blackout.Description,
			&blackout.Type,
			&blackout.InitiatorName,
			&source,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if endDate.Valid {
			blackout.EndDate = endDate.String
		}

		if source.Valid {
			blackout.Source = source.String
		}

		blackouts = append(blackouts, blackout)
	}

	return blackouts, nil
}

func (s *Storage) GetBuildingsCount() (int64, error) {
	const op = "storage.sqlite.GetBuildingsCount"

	var count int64
	if err := s.db.QueryRow("SELECT COUNT(*) FROM buildings WHERE is_fake = 0").Scan(&count); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (s *Storage) GetBuildingsCountByBlackoutType(blackoutType string, currentTime string) (int64, error) {
	const op = "storage.sqlite.GetBuildingsCountByBlackoutType"

	queryTime := currentTime
    if len(currentTime) == 10 {
        queryTime = currentTime + " 23:59:59"
		currentTime = currentTime + " 00:00:00"
    }

	var count int64
	err := s.db.QueryRow(`
        SELECT COUNT(DISTINCT b.id) 
        FROM buildings b
        JOIN blackouts_buildings bb ON b.id = bb.building_id
        JOIN blackouts bl ON bb.blackout_id = bl.id
        WHERE bl.type = ? 
		AND b.is_fake = 0
        AND bl.start_date <= ? 
        AND (bl.end_date >= ? OR bl.end_date IS NULL)`,
		blackoutType, queryTime, currentTime).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return count, nil
}

func (s *Storage) GetLastBlackoutTimeByType(blackoutType string, currentTime string) (string, error) {
	const op = "storage.sqlite.GetLastBlackoutTimeByType"

	var lastTime string
	err := s.db.QueryRow(`
        SELECT MAX(start_date) 
        FROM blackouts 
        WHERE type = ? 
        AND start_date <= ?`,
		blackoutType, currentTime).Scan(&lastTime)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return lastTime, nil
}

func (s *Storage) GetComplaintsLastDay(endTime string) ([]models.ComplaintData, error) {
    const op = "storage.sqlite.GetComplaintsLastDay"

    endTimeParsed, err := time.Parse("2006-01-02 15:04:05", endTime)
    if err != nil {
        return nil, fmt.Errorf("%s: invalid time format: %w", op, err)
    }
    
    startTime := endTimeParsed.Add(-24 * time.Hour)
    startTimeStr := startTime.Format("2006-01-02 15:04:05")
    endTimeStr := endTimeParsed.Format("2006-01-02 15:04:05")

    rows, err := s.db.Query(`
        SELECT 
            strftime('%Y-%m-%d %H:00:00', start_date) as hour,
            type,
            COUNT(*) as count
        FROM blackouts 
        WHERE start_date >= ? AND start_date < ?
        GROUP BY strftime('%Y-%m-%d %H', start_date), type
        ORDER BY hour`,
        startTimeStr, endTimeStr)
    
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer rows.Close()

    dataMap := make(map[string]*models.ComplaintData)

    for i := 0; i < 24; i++ {
        hourTime := startTime.Add(time.Duration(i) * time.Hour)
        hourKey := hourTime.Format("2006-01-02 15:00:00")
        displayTime := hourTime.Format("15:04")
        dataMap[hourKey] = &models.ComplaintData{
            Time:        displayTime,
            HotWater:    0,
            ColdWater:   0,
            Electricity: 0,
            Heating:     0,
        }
    }

    for rows.Next() {
        var hour, blackoutType string
        var count int
        if err := rows.Scan(&hour, &blackoutType, &count); err != nil {
            return nil, fmt.Errorf("%s: %w", op, err)
        }

        if data, exists := dataMap[hour]; exists {
            switch blackoutType {
            case "hot_water":
                data.HotWater = count
            case "cold_water":
                data.ColdWater = count
            case "electricity":
                data.Electricity = count
            case "heat":
                data.Heating = count
            }
        }
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    var result []models.ComplaintData
    for i := 0; i < 24; i++ {
        hourTime := startTime.Add(time.Duration(i) * time.Hour)
        hourKey := hourTime.Format("2006-01-02 15:00:00")
        if data, exists := dataMap[hourKey]; exists {
            result = append(result, *data)
        }
    }

    return result, nil
}

func (s *Storage) GetComplaintsLastHour(endTime string) ([]models.ComplaintData, error) {
    const op = "storage.sqlite.GetComplaintsLastHour"

    endTimeParsed, err := time.Parse("2006-01-02 15:04:05", endTime)
    if err != nil {
        return nil, fmt.Errorf("%s: invalid time format: %w", op, err)
    }
    
    startTime := endTimeParsed.Add(-1 * time.Hour)
    startTimeStr := startTime.Format("2006-01-02 15:04:05")
    endTimeStr := endTimeParsed.Format("2006-01-02 15:04:05")

    rows, err := s.db.Query(`
        SELECT 
            strftime('%Y-%m-%d %H:%M:00', start_date) as minute,
            type,
            COUNT(*) as count
        FROM blackouts 
        WHERE start_date >= ? AND start_date < ?
        GROUP BY strftime('%Y-%m-%d %H:%M', start_date), type
        ORDER BY minute`,
        startTimeStr, endTimeStr)
    
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer rows.Close()

    dataMap := make(map[string]*models.ComplaintData)

    for i := 0; i < 60; i++ {
        minuteTime := startTime.Add(time.Duration(i) * time.Minute)
        minuteKey := minuteTime.Format("2006-01-02 15:04:00")
        displayTime := minuteTime.Format("15:04")
        dataMap[minuteKey] = &models.ComplaintData{
            Time:        displayTime,
            HotWater:    0,
            ColdWater:   0,
            Electricity: 0,
            Heating:     0,
        }
    }

    for rows.Next() {
        var minute, blackoutType string
        var count int
        if err := rows.Scan(&minute, &blackoutType, &count); err != nil {
            return nil, fmt.Errorf("%s: %w", op, err)
        }

        if data, exists := dataMap[minute]; exists {
            switch blackoutType {
            case "hot_water":
                data.HotWater = count
            case "cold_water":
                data.ColdWater = count
            case "electricity":
                data.Electricity = count
            case "heat":
                data.Heating = count
            }
        }
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    var result []models.ComplaintData
    for i := 0; i < 60; i++ {
        minuteTime := startTime.Add(time.Duration(i) * time.Minute)
        minuteKey := minuteTime.Format("2006-01-02 15:04:00")
        if data, exists := dataMap[minuteKey]; exists {
            result = append(result, *data)
        }
    }

    return result, nil
}

func (s *Storage) GetComplaintsLastWeek(endTime string) ([]models.ComplaintData, error) {
    const op = "storage.sqlite.GetComplaintsLastWeek"

    endTimeParsed, err := time.Parse("2006-01-02 15:04:05", endTime)
    if err != nil {
        return nil, fmt.Errorf("%s: invalid time format: %w", op, err)
    }
    
    startTime := endTimeParsed.Add(-7 * 24 * time.Hour)
    startTimeStr := startTime.Format("2006-01-02 15:04:05")
    endTimeStr := endTimeParsed.Format("2006-01-02 15:04:05")

    rows, err := s.db.Query(`
        SELECT 
            strftime('%Y-%m-%d', start_date) as day,
            type,
            COUNT(*) as count
        FROM blackouts 
        WHERE start_date >= ? AND start_date < ?
        GROUP BY strftime('%Y-%m-%d', start_date), type
        ORDER BY day`,
        startTimeStr, endTimeStr)
    
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer rows.Close()

    dataMap := make(map[string]*models.ComplaintData)

    for i := 0; i < 7; i++ {
        dayTime := startTime.Add(time.Duration(i) * 24 * time.Hour)
        dayKey := dayTime.Format("2006-01-02")
        displayTime := dayTime.Format("02.01")
        dataMap[dayKey] = &models.ComplaintData{
            Time:        displayTime,
            HotWater:    0,
            ColdWater:   0,
            Electricity: 0,
            Heating:     0,
        }
    }

    for rows.Next() {
        var day, blackoutType string
        var count int
        if err := rows.Scan(&day, &blackoutType, &count); err != nil {
            return nil, fmt.Errorf("%s: %w", op, err)
        }

        if data, exists := dataMap[day]; exists {
            switch blackoutType {
            case "hot_water":
                data.HotWater = count
            case "cold_water":
                data.ColdWater = count
            case "electricity":
                data.Electricity = count
            case "heat":
                data.Heating = count
            }
        }
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    var result []models.ComplaintData
    for i := 0; i < 7; i++ {
        dayTime := startTime.Add(time.Duration(i) * 24 * time.Hour)
        dayKey := dayTime.Format("2006-01-02")
        if data, exists := dataMap[dayKey]; exists {
            result = append(result, *data)
        }
    }

    return result, nil
}

func (s *Storage) GetComplaintsLastMonth(endTime string) ([]models.ComplaintData, error) {
    const op = "storage.sqlite.GetComplaintsLastMonth"

    endTimeParsed, err := time.Parse("2006-01-02 15:04:05", endTime)
    if err != nil {
        return nil, fmt.Errorf("%s: invalid time format: %w", op, err)
    }
    
    startTime := endTimeParsed.Add(-30 * 24 * time.Hour)
    startTimeStr := startTime.Format("2006-01-02 15:04:05")
    endTimeStr := endTimeParsed.Format("2006-01-02 15:04:05")

    rows, err := s.db.Query(`
        SELECT 
            strftime('%Y-%m-%d', start_date) as day,
            type,
            COUNT(*) as count
        FROM blackouts 
        WHERE start_date >= ? AND start_date < ?
        GROUP BY strftime('%Y-%m-%d', start_date), type
        ORDER BY day`,
        startTimeStr, endTimeStr)
    
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer rows.Close()

    dataMap := make(map[string]*models.ComplaintData)

    for i := 0; i < 30; i++ {
        dayTime := startTime.Add(time.Duration(i) * 24 * time.Hour)
        dayKey := dayTime.Format("2006-01-02")
        displayTime := dayTime.Format("02.01") 
        dataMap[dayKey] = &models.ComplaintData{
            Time:        displayTime,
            HotWater:    0,
            ColdWater:   0,
            Electricity: 0,
            Heating:     0,
        }
    }

    for rows.Next() {
        var day, blackoutType string
        var count int
        if err := rows.Scan(&day, &blackoutType, &count); err != nil {
            return nil, fmt.Errorf("%s: %w", op, err)
        }

        if data, exists := dataMap[day]; exists {
            switch blackoutType {
            case "hot_water":
                data.HotWater = count
            case "cold_water":
                data.ColdWater = count
            case "electricity":
                data.Electricity = count
            case "heat":
                data.Heating = count
            }
        }
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    var result []models.ComplaintData
    for i := 0; i < 30; i++ {
        dayTime := startTime.Add(time.Duration(i) * 24 * time.Hour)
        dayKey := dayTime.Format("2006-01-02")
        if data, exists := dataMap[dayKey]; exists {
            result = append(result, *data)
        }
    }

    return result, nil
}

func (s *Storage) GetOrganizations(currentTime string) ([]string, error) {
	const op = "storage.sqlite.GetOrganizations"

	rows, err := s.db.Query(`
        SELECT b.initiator_name
        FROM blackouts b
        JOIN blackouts_buildings bb ON b.id = bb.blackout_id
        WHERE b.start_date <= ? AND (b.end_date >= ? OR b.end_date IS NULL)
        GROUP BY b.initiator_name
        ORDER BY COUNT(DISTINCT bb.building_id) DESC
		LIMIT 4`,
		currentTime, currentTime)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var organizations []string
	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		organizations = append(organizations, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return organizations, nil
}

func (s *Storage) GetBuildingsCountByOrgName(name string, currentTime string) (int64, error) {
	const op = "storage.sqlite.GetBuildingsCountByOrgName"

	var count int64
	err := s.db.QueryRow(`
        SELECT COUNT(DISTINCT bb.building_id) 
        FROM blackouts b
        JOIN blackouts_buildings bb ON b.id = bb.blackout_id
        JOIN buildings bu ON bb.building_id = bu.id
        WHERE b.initiator_name = ? 
        AND b.start_date <= ? 
        AND (b.end_date >= ? OR b.end_date IS NULL)
        AND bu.is_fake = 0`,
		name, currentTime, currentTime).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return count, nil
}

func (s *Storage) GetLastAddressByOrgName(name string, currentTime string) (string, string, error) {
	const op = "storage.sqlite.GetLastAddressByOrgName"

	var lastTime, address string
	err := s.db.QueryRow(`
        SELECT b.start_date, s.name || ' ' || bg.number
        FROM blackouts b
        JOIN blackouts_buildings bb ON b.id = bb.blackout_id
        JOIN buildings bg ON bb.building_id = bg.id
        JOIN streets s ON bg.street_id = s.id
        WHERE b.initiator_name = ? 
        AND b.start_date <= ?
        AND bg.is_fake = 0
        ORDER BY b.start_date DESC 
        LIMIT 1`,
		name, currentTime).Scan(&lastTime, &address)

	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return lastTime, address, nil
}

func (s *Storage) GetBlackoutsWithBuildingsCount(targetDate string) ([]models.BlackoutInfo, error) {
    const op = "storage.sqlite.GetBlackoutsWithBuildingsCount"

    queryTime := targetDate
    if len(targetDate) == 10 {
        queryTime = targetDate + " 23:59:59"
        targetDate = targetDate + " 00:00:00"
    }

    rows, err := s.db.Query(`
        SELECT 
            bl.type,
            bl.start_date,
            bl.end_date,
            COUNT(DISTINCT b.id) as building_count
        FROM blackouts bl
        JOIN blackouts_buildings bb ON bl.id = bb.blackout_id
        JOIN buildings b ON bb.building_id = b.id
        WHERE b.is_fake = 0
        AND bl.start_date <= ? 
        AND (bl.end_date >= ? OR bl.end_date IS NULL)
        GROUP BY bl.id, bl.type, bl.start_date, bl.end_date
        ORDER BY bl.start_date DESC`,
        queryTime, targetDate)
    // "2006-01-02 15:04:05" -> "2006-01-02 15:04"
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
    defer rows.Close()

    var blackouts []models.BlackoutInfo
    for rows.Next() {
        var blackout models.BlackoutInfo
        var endDate sql.NullString
        var startDate string
        
        err := rows.Scan(
            &blackout.Type,
            &startDate,
            &endDate,
            &blackout.BuildingCount,
        )
        if err != nil {
            return nil, fmt.Errorf("%s: %w", op, err)
        }

        if len(startDate) > 16 { 
            blackout.StartDate = startDate[:16]
        } else {
            blackout.StartDate = startDate
        }

        if endDate.Valid {
            endDateStr := endDate.String
            if len(endDateStr) > 16 {
                blackout.EndDate = endDateStr[:16]
            } else {
                blackout.EndDate = endDateStr
            }
        } else {
            blackout.EndDate = ""
        }

        blackouts = append(blackouts, blackout)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return blackouts, nil
}