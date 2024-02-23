package models

type WeatherData struct {
	City string `json:"city"`
	CurrentCondition []CurrentCondition `json:"current_conditon"`
}

type CurrentCondition struct {
	Humidity int `json:"humidity"`
	TempC int `json:"temp_C"`
	TempF int `json:"temp_F"`
	ObservationTime string `json:"observation_time"`
	WindSpeedKmph int `json:"windspeedKmph"`
	WindSpeedMiles int `json:"windspeedMiles"`
}