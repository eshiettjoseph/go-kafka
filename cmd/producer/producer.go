package producer

import (
	"log"
	"encoding/json"
	"fmt"
	"os"
	"github.com/IBM/sarama"
)


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

var (
	InfoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)

const (
	KafkaTopic          = "Weather Data"
	ProducerPort        = ":4000"
)

func SetupProducer() (sarama.SyncProducer, error) {
	// kafkaAddress := os.Getenv("KAFKA_SERVER_ADDRESS") 
	kafkaAddress := "localhost:9092"
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

func produceToKafka(producer sarama.SyncProducer, weatherData WeatherData) error {
	dataJSON, err := json.Marshal(weatherData)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(weatherData.City), // City as the key
		Value: sarama.StringEncoder(dataJSON),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}