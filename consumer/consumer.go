package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/IBM/sarama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
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

type DB struct {
	*gorm.DB
}

func main() {
	dsn := os.ExpandEnv("host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} port=${DB_PORT} sslmode=${DB_SSLMODE}")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&WeatherData{}, &CurrentCondition{})

	kafkaAddress := os.Getenv("KAFKA_SERVER_ADDRESS")
	consumerGroup := "weather-group"
	topic := "Weather"

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup([]string{kafkaAddress}, consumerGroup, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer group: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received interrupt signal. Shutting down...")
		cancel()
	}()

	// Consume messages from Kafka topic
	consumerHandler := &ConsumerHandler{DB: db}
	err = consumer.Consume(ctx, []string{topic}, consumerHandler)
	if err != nil {
		log.Fatalf("Error from Kafka consumer: %v", err)
	}
}

// ConsumerHandler implements sarama.ConsumerGroupHandler interface
type ConsumerHandler struct {
	DB *gorm.DB
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }


func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var weatherData WeatherData
		if err := json.Unmarshal(msg.Value, &weatherData); err != nil {
			fmt.Printf("Error unmarshalling data: %v", err)
			continue
		}
		
		parsedWeatherData := WeatherData{
			City: weatherData.City,
			CurrentCondition:[]CurrentCondition{
				{
					Humidity:         weatherData.CurrentCondition[0].Humidity,
					TempC:            weatherData.CurrentCondition[0].TempC,
					TempF:            weatherData.CurrentCondition[0].TempF,
					ObservationTime:  weatherData.CurrentCondition[0].ObservationTime,
					WindspeedKmph:    weatherData.CurrentCondition[0].WindspeedKmph,
					WindspeedMiles:   weatherData.CurrentCondition[0].WindspeedMiles,
					LocalObsDateTime:  weatherData.CurrentCondition[0].LocalObsDateTime,
					Pressure: weatherData.CurrentCondition[0].Pressure,
    			},
			},
		}
		if err := h.DB.Create(&parsedWeatherData).Error; err != nil {
			fmt.Print(weatherData.CurrentCondition[0].Humidity)
			fmt.Printf("Failed to save data to database: %v", err)
			continue
		}
		fmt.Printf("Temperature Data for %s saved successfully\n", weatherData.City)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (db *DB) SaveWeatherData(data WeatherData) error {
	return db.Create(&data).Error
}