package organizations

import (
	"log/slog"
	"net/http"
	"vlru-prsch/internal/lib/api/response"
	"vlru-prsch/internal/lib/date"
	"vlru-prsch/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Request represents the request structure
// @Description Структура запроса для фильтрации организаций
type Request struct {
	Name string `json:"name"`
}

// Response represents the API response for organizations
// @Description Ответ с информацией об организациях и их отключениях
type Response struct {
	response.Response
	Organizations []OrganizationInfo `json:"organizations"`
}

// OrganizationInfo represents information about a specific organization
// @Description Информация об организации и её отключениях
type OrganizationInfo struct {
	// Название организации
	Name string `json:"name" example:"МУПВ ВПЭС (электрические сети)"`
	// Количество зданий организации
	CountBuildings int64 `json:"count_buildings" example:"106"`
	// Адрес последнего отключения
	LastAddress string `json:"last_address" example:"Карбышева ул. 54"`
	// Время последнего отключения в формате "2006-01-02 15:04:05"
	TimeLastBlackout string `json:"time_last_blackout" example:"2019-01-28 09:39:00"`
}

type OrganizationGiver interface {
	GetOrganizations(currentTime string) ([]string, error)
	GetBuildingsCountByOrgName(name string, currentTime string) (int64, error)
	GetLastAddressByOrgName(name string, currentTime string) (string, string, error)
}

// New godoc
// @Summary Получить информацию об организациях
// @Description Возвращает список организаций с информацией о количестве зданий и последних отключениях
// @Tags organizations
// @Accept json
// @Produce json
// @Param curr_time query string true "Текущее время в формате YYYY-MM-DDTHH:MM:SSZ или YYYY-MM-DD_HH:MM:SS" example(2024-01-15_14:30:00 или 2024-01-15T14:30:00Z)
// @Security ApiKeyAuth
// @Success 200 {object} Response "Успешный ответ с данными об организациях"
// @Failure 400 {object} response.Response "Неверный запрос - пример: {\"status\":\"ERROR\",\"error\":\"curr_time parameter is required\"}"
// @Failure 400 {object} response.Response "Неверный формат времени - пример: {\"status\":\"ERROR\",\"error\":\"invalid time format\"}"
// @Failure 500 {object} response.Response "Внутренняя ошибка сервера - пример: {\"status\":\"ERROR\",\"error\":\"failed to get organizations\"}"
// @Router /off/orgs [get]
func New(log *slog.Logger, giver OrganizationGiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.organization.get.New"

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

		orgNames, err := giver.GetOrganizations(currTimeParse)
		if err != nil {
			log.Error("failed to get organizations", sl.Err(err))
			render.JSON(w, r, response.Error("failed to get organizations"))
			return
		}

		var organizationsInfo []OrganizationInfo

		for _, orgName := range orgNames {
			countBuildings, err := giver.GetBuildingsCountByOrgName(orgName, currTimeParse)
			if err != nil {
				log.Error("failed to get buildings count",
					slog.String("org", orgName), sl.Err(err))
				continue
			}

			lastTime, lastAddress, err := giver.GetLastAddressByOrgName(orgName, currTimeParse)
			if err != nil {
				log.Error("failed to get last blackout",
					slog.String("org", orgName), sl.Err(err))
				lastTime = "unknown"
				lastAddress = "unknown"
			}

			info := OrganizationInfo{
				Name:             orgName,
				CountBuildings:   countBuildings,
				LastAddress:      lastAddress,
				TimeLastBlackout: lastTime,
			}

			organizationsInfo = append(organizationsInfo, info)
			
		}

		render.JSON(w, r, Response{
			Response:      response.Ok(),
			Organizations: organizationsInfo,
		})
	}
}