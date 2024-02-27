package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/IBM/sarama"
	"time"
	"producer/models"
	"os"
)

const KafkaTopic = "Weather"

func getWeatherData(city string) ([]byte, string, error) {
	route := fmt.Sprintf("https://wttr.in/%s?format=j1", city)
	response, err := http.Get(route)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return nil, "", err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, "", err
	}
	return body, "", nil

}

func main(){
	// Run the script initially
	runScript()

	// Schedule the script to run every 1 hour
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		runScript()
	}
}


func runScript() {
	fmt.Println("Starting Temperature Scheduler")
	cities := []string{"Abuja", "London", "Miami", "Tokyo", "Singapore"}

	for _, city := range cities {
		body, _, err := getWeatherData(city)
		if err != nil {
			fmt.Printf("Could not get temperature for %s: %s\n", city, err)
			continue
		}
		
		var weatherData models.WeatherData
		err = json.Unmarshal(body, &weatherData)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		
		parsedWeatherData := models.WeatherData{
			City: city,
			CurrentCondition:[]models.CurrentCondition{
				{
					Humidity:         weatherData.CurrentCondition[0].Humidity,
					TempC:            weatherData.CurrentCondition[0].TempC,
					TempF:            weatherData.CurrentCondition[0].TempF,
					ObservationTime:  weatherData.CurrentCondition[0].ObservationTime,
					WindspeedKmph:    weatherData.CurrentCondition[0].WindspeedKmph,
					WindspeedMiles:   weatherData.CurrentCondition[0].WindspeedMiles,
				},
			},
		}

		// Initialize Kafka producer
		producer, err := initializeProducer()
		if err != nil {
			fmt.Println("Error initializing Kafka producer:", err)
			return
		}
		defer producer.Close()

		// Produce to Kafka
		err = produceToKafka(producer, parsedWeatherData)
		if err != nil {
			fmt.Println("Error producing to Kafka:", err)
			return
		}

		fmt.Printf("Weather data %s produced to Kafka.\n", city)
	}
}


func initializeProducer() (sarama.SyncProducer, error) {
	kafkaAddress := os.Getenv("KAFKA_SERVER_ADDRESS")
	if kafkaAddress == "" {
		return nil, fmt.Errorf("kafka address not set in environment variable KAFKA_SERVER_ADDRESS")
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafkaAddress}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return producer, nil
}

func produceToKafka(producer sarama.SyncProducer, parsedWeatherData models.WeatherData) error {
	dataJSON, err := json.Marshal(parsedWeatherData)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Value: sarama.StringEncoder(dataJSON),
	}
	// Produce the message to Kafka
	_, _, err = producer.SendMessage(msg)
	return err
}