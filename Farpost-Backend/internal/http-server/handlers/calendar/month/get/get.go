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

// Response represents the API response for calendar month
// @Description Ответ с данными для календаря отключений по дням месяца
type Response struct {
  	response.Response
 	Dates []DateInfo `json:"dates"`
}

// DateInfo represents information about outages for a specific date
// @Description Информация об отключениях для конкретной даты
type DateInfo struct {
  	// Дата в формате YYYY-MM-DD
  	Date string `json:"date" example:"2019-01-15"`
  	// Список типов отключений в эту дату: hot_water, cold_water, electricity, heat
  	Services []string `json:"services" example:"hot_water,cold_water"`
}

type DatesGiver interface {
    GetBlackouts(currentTime string) ([]models.Blackout, error)
}

// New godoc
// @Summary Получить данные для календаря отключений за месяц
// @Description Возвращает информацию об отключениях по дням указанного месяца для отображения в календаре
// @Tags calendar
// @Accept json
// @Produce json
// @Param month query string true "Первый месяц в формате YYYY-MM" example(2019-01)
// @Security ApiKeyAuth
// @Success 200 {object} Response "Успешный ответ с данными за месяц"
// @Failure 400 {object} response.Response "Отсутствует параметр month - пример: {\"status\":\"ERROR\",\"error\":\"month parameter is required\"}"
// @Failure 400 {object} response.Response "Неверный формат месяца - пример: {\"status\":\"ERROR\",\"error\":\"failed to process month dates\"}"
// @Failure 500 {object} response.Response "Ошибка при получении данных - пример: {\"status\":\"ERROR\",\"error\":\"failed to get blackouts data\"}"
// @Router /off/calendar [get]
func New(log *slog.Logger, giver DatesGiver) http.HandlerFunc {
  	return func(w http.ResponseWriter, r *http.Request) {
    	const op = "handlers.calendar.month.get.New"

    	log := log.With(
      		slog.String("op", op),
      		slog.String("request_id", middleware.GetReqID(r.Context())),
    	)

    	month := r.URL.Query().Get("month")
    	if month == "" {
      		log.Warn("fday parameter is empty")
      		render.JSON(w, r, response.Error("month parameter is required"))
      		return
    	}

   		monthDates, err := date.GetAllDatesInMonth(month)
    	if err != nil {
      		log.Error("failed to get dates in month", sl.Err(err))
      		render.JSON(w, r, response.Error("failed to process month dates"))
      		return
    	}

    	var dates []DateInfo

    	for _, dateStr := range monthDates {
     		blackouts, err := giver.GetBlackouts(dateStr)
      		if err != nil {
        		log.Error("failed to get blackouts for date", 
          		slog.String("date", dateStr), 
          		sl.Err(err))
        		continue
      		}

      		serviceTypes := make(map[string]bool)
      		for _, blackout := range blackouts {
        		serviceTypes[blackout.Type] = true
      		}

      		var services []string
      		for serviceType := range serviceTypes {
        		services = append(services, serviceType)
      		}

      		dateInfo := DateInfo{
        		Date:     dateStr,
        		Services: services,
      		}
      		dates = append(dates, dateInfo)
    	}

    	render.JSON(w, r, Response{
      		Response: response.Ok(),
      		Dates: dates,
    	})
  	}
}