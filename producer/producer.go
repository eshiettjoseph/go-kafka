package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/robfig/cron"
	"github.com/IBM/sarama"
	"time"
	"go-kafka/producer/models"
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
	cronJob := cron.New()
	fmt.Println("Starting Temperature Scheduler")
	cronJob.AddFunc("@hourly", func() {
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
	})
	cronJob.Start()
}

func initializeProducer() (sarama.SyncProducer, error) {
	kafkaAddress := os.Getenv("KAFKA_SERVER_ADDRESS")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	kafkaBrokers := []string{kafkaAddress} // Add your Kafka broker address here
	producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
	if err != nil {
		return nil, err
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
