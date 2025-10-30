package models

// ComplaintData represents complaint statistics data by outage type
// @Description Данные жалоб по типам отключений для построения графиков
type ComplaintData struct {
	// Временная метка точки данных
	Time string `json:"time" example:"2019-01-15 14:00:00"`
	// Количество жалоб на отключение горячей воды
	HotWater int `json:"hot" example:"5"`
	// Количество жалоб на отключение холодной воды  
	ColdWater int `json:"cold" example:"3"`
	// Количество жалоб на отключение электричества
	Electricity int `json:"electricity" example:"8"`
	// Количество жалоб на отключение отопления
	Heating int `json:"heating" example:"2"`
}