package blackouts

import (
	"log/slog"
	"math"
	"net/http"
	"vlru-prsch/internal/lib/api/response"
	"vlru-prsch/internal/lib/date"
	"vlru-prsch/internal/lib/logger/sl"
	"vlru-prsch/internal/models"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Response represents the API response for blackouts
// @Description Информация о текущих отключениях
type Response struct {
	response.Response
	Blackouts []BlackoutInfo `json:"blackouts"`
}

// BlackoutInfo represents information about a specific type of blackout
// @Description Информация об отключении конкретного типа
type BlackoutInfo struct {
	// Тип отключения: hot_water (горячая вода), cold_water (холодная вода), electricity (электричество), heat (отопление)
	Type string `json:"type" example:"hot_water"`
	// Количество затронутых зданий
	CountBuildings int64 `json:"count_buildings" example:"15"`
	// Доля затронутых зданий в процентах
	FractionBuildings float64 `json:"fraction_buildings" example:"25.5"`
	// Время последнего отключения в формате "2006-01-02 15:04:05"
	TimeLastBlackout string `json:"time_last_blackout" example:"2019-01-15 14:30:00"`
}

type BlackoutGiver interface {
	GetBlackouts(currentTime string) ([]models.Blackout, error)
	GetBuildingsCount() (int64, error)
	GetBuildingsCountByBlackoutType(blackoutType string, currentTime string) (int64, error)
	GetLastBlackoutTimeByType(blackoutType string, currentTime string) (string, error)
}

// New godoc
// @Summary Получить информацию об отключениях
// @Description Возвращает статистику по отключениям горячей/холодной воды, электричества и отопления
// @Tags blackouts
// @Accept json
// @Produce json
// @Param curr_time query string true "Текущее время в формате YYYY-MM-DDTHH:MM:SSZ или YYYY-MM-DD_HH:MM:SS" example(2019-01-15T14:30:00Z или 2019-01-15_14:30:00)// @Security ApiKeyAuth
// @Success 200 {object} Response "Успешный ответ"
// @Failure 400 {object} response.Response "Неверный формат времени или отсутствует параметр curr_time"
// @Failure 500 {object} response.Response "Ошибка при получении данных"
// @Router /off/blackouts [get]
// @Response 400 {object} response.Response "Пример: {\"status\":\"ERROR\",\"error\":\"curr_time parameter is required\"}"
// @Response 400 {object} response.Response "Пример: {\"status\":\"ERROR\",\"error\":\"invalid time format\"}"
// @Response 500 {object} response.Response "Пример: {\"status\":\"ERROR\",\"error\":\"failed to get buildings data\"}"
func New(log *slog.Logger, giver BlackoutGiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.blackout.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		currTime := r.URL.Query().Get("curr_time")
		if currTime == "" {
			log.Warn("curr_time parameter is empty")
			render.JSON(w, r, response.Error("curr_time parameter is required"))
			return
		}

		currTimeParse, err := date.ParseQueryDate(currTime)
		if err != nil {
			log.Warn("invalid time format", slog.String("curr_time", currTime), sl.Err(err))
			render.JSON(w, r, response.Error("invalid time format"))
			return
		}

		totalBuildings, err := giver.GetBuildingsCount()
		if err != nil {
			log.Error("failed to get total buildings count", sl.Err(err))
			render.JSON(w, r, response.Error("failed to get buildings data"))
			return
		}

		blackoutTypes := []string{"hot_water", "cold_water", "electricity", "heat"}
		
		var blackoutsInfo []BlackoutInfo

		for _, blackoutType := range blackoutTypes {
			affectedBuildings, err := giver.GetBuildingsCountByBlackoutType(blackoutType, currTimeParse)
			if err != nil {
				log.Error("failed to get buildings count for type", 
				slog.String("type", blackoutType), slog.Any("error", err))
				continue
			}

			lastBlackoutTime, err := giver.GetLastBlackoutTimeByType(blackoutType, currTimeParse)
			if err != nil {
				log.Error("failed to get last blackout time", 
				slog.String("type", blackoutType), slog.Any("error", err))
				lastBlackoutTime = "unknown"
			}

			var fraction float64
			if totalBuildings > 0 {
				fraction = float64(affectedBuildings) / float64(totalBuildings)
			}
			percentage := math.Round(fraction * 100 * 100) / 100

			info := BlackoutInfo{
				Type:              blackoutType,
				CountBuildings:    affectedBuildings,
				FractionBuildings: percentage,
				TimeLastBlackout:  lastBlackoutTime,
			}

			blackoutsInfo = append(blackoutsInfo, info)
		}

		render.JSON(w, r, Response{
			Response:  response.Ok(),
			Blackouts: blackoutsInfo,
		})
	}
}