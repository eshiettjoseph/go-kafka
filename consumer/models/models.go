package models

import(
	"gorm.io/gorm"
)
type WeatherData struct {
	gorm.Model
	City string `json:"city"`
	CurrentCondition []CurrentCondition `json:"current_condition" gorm:"foreignkey:CurrentConditionID"`
}

type CurrentCondition struct {
	gorm.Model
	CurrentConditionID    uint
	Humidity         string           `json:"humidity"`
	LocalObsDateTime string           `json:"localObsDateTime"`
	ObservationTime  string           `json:"observation_time"`
	Pressure         string           `json:"pressure"`
	TempC            string           `json:"temp_C"`
	TempF            string           `json:"temp_F"`
	WindspeedKmph    string           `json:"windspeedKmph"`
	WindspeedMiles   string           `json:"windspeedMiles"`
}