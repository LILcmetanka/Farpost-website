package calendar

import (
	"log/slog"
	"net/http"
	"vlru-prsch/internal/lib/api/response"
	"vlru-prsch/internal/lib/date"
	"vlru-prsch/internal/lib/logger/sl"
	"vlru-prsch/internal/models"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Response represents the API response for calendar day
// @Description Ответ с детальной информацией об отключениях за конкретный день
type Response struct {
    response.Response
    Blackouts []InfoOffs `json:"blackouts"`
}

// InfoOffs represents detailed information about service outages for a specific day
// @Description Детальная информация об отключениях услуг за конкретный день
type InfoOffs struct {
    // Тип отключенной услуги: hot_water, cold_water, electricity, heat
    Service string `json:"service" example:"hot_water"`
    // Дата и время начала отключения в формате YYYY-MM-DD HH:MM:SS
    StartOff string `json:"start_off" example:"2019-01-15 10:00:00"`
    // Дата и время окончания отключения в формате YYYY-MM-DD HH:MM:SS
    EndOff string `json:"end_off" example:"2019-01-15 18:00:00"`
    // Количество затронутых адресов/зданий
    AmountAddresses int64 `json:"amount_addresses" example:"25"`
}

type DayInfoGiver interface {
    GetBlackoutsWithBuildingsCount(targetDate string) ([]models.BlackoutInfo, error)
}

// New godoc
// @Summary Получить детальную информацию об отключениях за день
// @Description Возвращает детальную информацию об отключениях услуг за указанную дату, включая время отключений и количество затронутых адресов
// @Tags calendar
// @Accept json
// @Produce json
// @Param date query string true "Целевая дата в формате YYYY-MM-DD" example(2019-01-15)
// @Security ApiKeyAuth
// @Success 200 {object} Response "Успешный ответ с детальной информацией об отключениях"
// @Failure 400 {object} response.Response "Отсутствует параметр date - пример: {\"status\":\"ERROR\",\"error\":\"date parameter is required\"}"
// @Failure 500 {object} response.Response "Ошибка при получении данных - пример: {\"status\":\"ERROR\",\"error\":\"failed to get blackouts information\"}"
// @Router /off/calendar/day [get]
func New(log *slog.Logger, giver DayInfoGiver) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.calendar.day.get.New"

        log := log.With(
            slog.String("op", op),
            slog.String("request_id", middleware.GetReqID(r.Context())),
        )

        targetDate := r.URL.Query().Get("date")
        if targetDate == "" {
            log.Warn("date parameter is empty")
            render.JSON(w, r, response.Error("date parameter is required"))
            return
        }

        if !date.IsValidDate(targetDate) {
            log.Warn("invalid date format provided", 
                slog.String("date", targetDate),
                slog.String("expected_format", "YYYY-MM-DD"))
            render.JSON(w, r, response.Error("invalid date format, expected YYYY-MM-DD"))
            return
        }   

        blackoutsInfo, err := giver.GetBlackoutsWithBuildingsCount(targetDate)
        if err != nil {
            log.Error("failed to get blackouts info for date", 
                slog.String("date", targetDate), 
                sl.Err(err))
            render.JSON(w, r, response.Error("failed to get blackouts information"))
            return
        }

        var infoOffs []InfoOffs
        for _, blackout := range blackoutsInfo {
            info := InfoOffs{
                Service:         blackout.Type,
                StartOff:        blackout.StartDate,
                EndOff:          blackout.EndDate,
                AmountAddresses: blackout.BuildingCount,
            }
            infoOffs = append(infoOffs, info)
        }

        render.JSON(w, r, Response{
            Response:  response.Ok(),
            Blackouts: infoOffs,
        })
    }
}