package models

type Blackout struct {
	ID 				string
	StartDate 		string
	EndDate 		string
	Description 	string
	Type 			string
	InitiatorName 	string
	Source 			string
}

// BlackoutInfo represents detailed information about a service outage
// @Description Детальная информация об отключении услуги
type BlackoutInfo struct {
    // Тип отключения: hot_water, cold_water, electricity, heat
    Type string `json:"service" example:"hot_water"`
    // Дата начала отключения в формате YYYY-MM-DD
    StartDate string `json:"start_off" example:"2019-01-15"`
    // Дата окончания отключения в формате YYYY-MM-DD
    EndDate string `json:"end_off" example:"2019-01-16"`
    // Количество затронутых адресов/зданий
    BuildingCount int64 `json:"amount_addresses" example:"25"`
}