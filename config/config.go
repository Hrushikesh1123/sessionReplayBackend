package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config struct holds configuration
type Config struct {
	Host          string
	Port          int
	User          string
	Password      string
	DBName        string
	KafkaBrokers  []string
	KafkaTopics   []string
	KafkaGroupIDs []string
}

// LoadConfig loads configuration from environment variables or `.env` file
func LoadConfig() (*Config, error) {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:          os.Getenv("DB_HOST"),
		Port:          port,
		User:          os.Getenv("DB_USER"),
		Password:      os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		KafkaBrokers:  strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		KafkaTopics:   strings.Split(os.Getenv("KAFKA_TOPICS"), ","),
		KafkaGroupIDs: strings.Split(os.Getenv("KAFKA_GROUP_IDS"), ","),
	}, nil
}
