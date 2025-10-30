package search

import (
	"log/slog"
	"net/http"
	"vlru-prsch/internal/lib/api/response"
	"vlru-prsch/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Request представляет запрос на поиск
// @Description Запрос для поиска улиц по подстроке
type Request struct {
	// Подстрока для поиска улиц (минимум 2 символа)
	Suggest string `json:"suggest" example:"ленин"`
}

// Response представляет ответ на поиск
// @Description Ответ с найденными улицами
type Response struct {
	response.Response
	// Список найденных улиц
	Streets []string `json:"streets" example:"ул. Ленина,пр. Ленинский,пл. Ленинская"`
}

type StreetsFinder interface {
	FindStreets(substr string) ([]string, error)
}

// New godoc
// @Summary Поиск улиц по подстроке
// @Description Поиск улиц по частичному совпадению названия. Возвращает список улиц, содержащих указанную подстроку
// @Tags search
// @Accept json
// @Produce json
// @Param request body Request true "Параметры поиска"
// @Security ApiKeyAuth
// @Success 200 {object} Response "Успешный ответ со списком найденных улиц"
// @Failure 400 {object} response.Response "Неверный запрос - пример: {\"status\":\"ERROR\",\"error\":\"failed to decode req\"}"
// @Failure 400 {object} response.Response "Слишком короткая строка поиска - пример: {\"status\":\"ERROR\",\"error\":\"suggest must be at least 2 characters\"}"
// @Failure 500 {object} response.Response "Ошибка поиска - пример: {\"status\":\"ERROR\",\"error\":\"search failed\"}"
// @Router /off/search [post]
func New(log *slog.Logger, finder StreetsFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.search.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode req body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode req"))
			return
		}

		if len(req.Suggest) < 2 {
			log.Warn("search string too short", slog.String("suggest", req.Suggest))
			render.JSON(w, r, response.Error("suggest must be at least 2 characters"))
			return
		}

		log.Info("req body decoded", slog.Any("request", req))

		var streets []string

		streets, err := finder.FindStreets(req.Suggest)
		if err != nil {
			log.Error("failed to find streets", slog.Any("error", err))
			render.JSON(w, r, response.Error("search failed"))
			return
		}

		render.JSON(w, r, Response{
			Response: response.Ok(),
			Streets:  streets,
		})
	}
}