package complaints

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

// Response represents the API response for complaints
// @Description Ответ с данными жалоб для построения графиков
type Response struct {
	response.Response
	Complaints []models.ComplaintData `json:"complaints"`
}

type ComplaintsGiver interface {
	GetComplaintsLastHour(currTimeParse string) ([]models.ComplaintData, error)
	GetComplaintsLastDay(currTimeParse string) ([]models.ComplaintData, error)
	GetComplaintsLastWeek(currTimeParse string) ([]models.ComplaintData, error)
	GetComplaintsLastMonth(currTimeParse string) ([]models.ComplaintData, error)
}

const (
	PeriodHour  = "hour"
	PeriodDay   = "day"
	PeriodWeek  = "week"
	PeriodMonth = "month"
)

// New godoc
// @Summary Получить данные жалоб для графиков
// @Description Возвращает статистику жалоб за указанный период для построения графиков и аналитики
// @Tags complaints
// @Accept json
// @Produce json
// @Param period query string true "Период для агрегации данных: hour (последний час), day (последние 24 часа), week (последние 7 дней), month (последние 30 дней)" example(day)
// @Param curr_time query string true "Текущее время в формате YYYY-MM-DDTHH:MM:SSZ или YYYY-MM-DD_HH:MM:SS" example(2019-01-15_14:30:00 или 2019-01-15T14:30:00Z)
// @Security ApiKeyAuth
// @Success 200 {object} Response "Успешный ответ с данными жалоб по типам отключений"
// @Failure 400 {object} response.Response "Отсутствует параметр period - пример: {\"status\":\"ERROR\",\"error\":\"period parameter is required\"}"
// @Failure 400 {object} response.Response "Отсутствует параметр curr_time - пример: {\"status\":\"ERROR\",\"error\":\"curr_time parameter is required\"}"
// @Failure 400 {object} response.Response "Неверный формат времени - пример: {\"status\":\"ERROR\",\"error\":\"invalid time format\"}"
// @Failure 400 {object} response.Response "Неверный период - пример: {\"status\":\"ERROR\",\"error\":\"invalid period, use: hour, day, week, month\"}"
// @Failure 500 {object} response.Response "Ошибка при получении данных - пример: {\"status\":\"ERROR\",\"error\":\"failed to get complaints data\"}"
// @Router /off/complaints [get]
// @Example period=hour "Получить данные за последний час"
// @Example period=day "Получить данные за последние 24 часа" 
// @Example period=week "Получить данные за последние 7 дней"
// @Example period=month "Получить данные за последние 30 дней"
func New(log *slog.Logger, giver ComplaintsGiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.complaints.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		period := r.URL.Query().Get("period")
		if period == "" {
			log.Warn("period parameter is empty")
			render.JSON(w, r, response.Error("period parameter is required"))
			return
		}

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

		validPeriods := map[string]bool{
			PeriodHour:  true,
			PeriodDay:   true,
			PeriodWeek:  true,
			PeriodMonth: true,
		}

		if !validPeriods[period] {
			log.Warn("invalid period", slog.String("period", period))
			render.JSON(w, r, response.Error("invalid period, use: hour, day, week, month"))
			return
		}

		var complaints []models.ComplaintData
		switch period {
		case PeriodHour:
			complaints, err = giver.GetComplaintsLastHour(currTimeParse)
		case PeriodDay:
			complaints, err = giver.GetComplaintsLastDay(currTimeParse)
		case PeriodWeek:
			complaints, err = giver.GetComplaintsLastWeek(currTimeParse)
		case PeriodMonth:
			complaints, err = giver.GetComplaintsLastMonth(currTimeParse)
		}

		if err != nil {
			log.Error("failed to get complaints data",
				slog.String("period", period),
				slog.String("curr_time", currTimeParse),
				sl.Err(err))
			render.JSON(w, r, response.Error("failed to get complaints data"))
			return
		}

		render.JSON(w, r, Response{
			Response:   response.Ok(),
			Complaints: complaints,
		})
	}
}
