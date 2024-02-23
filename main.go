package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"go-kafka/cmd/producer"
	"go-kafka/pkg/models"
)

func getWeatherData(city string) (models.WeatherData, error) {
	route := fmt.Sprintf("http://wttr.in/%s?format=j1", city)

	resp, err := http.Get(route)
	if err != nil {
		return models.WeatherData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.WeatherData{}, err
	}

	var weatherData models.WeatherData
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return models.WeatherData{}, err
	}
	return weatherData, nil
}


func main() {
	prod, err := producer.SetupProducer()
	if err != nil {
		producer.ErrorLog.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	producer.InfoLog.Println("Kafka producer initialized successfully")
	defer prod.Close()

	cities := []string{"Abuja", "London", "Miami", "Tokyo", "Singapore"}

	for _, city := range cities {
		weatherData, err := getWeatherData(city)
		if err != nil {
			log.Println("Error fetching weather data for", city, ":", err)
			continue
		}

		err = producer.produceToKafka(prod, weatherData)
		if err != nil {
			log.Println("Error producing to Kafka for", city, ":", err)
			continue
		}

		fmt.Printf("Weather data for %s sent to Kafka.\n", city)
	}
}
