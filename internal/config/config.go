package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hryak228pizza/wbTech.L0/internal/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DBUser  string
	DBPassw string
	DBName  string
	Dsn     string

	KafkaBroker string
	KafkaTopic  string
	KafkaGroup  string

	HttpPort  string
	CacheSize int
}

func LoadCfg() *Config {

	// load .env variables
	err := godotenv.Load()
	if err != nil {
		logger.L().Fatal("Error loading .env file")
	}

	// getting .env variables
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassw := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	httpPort := os.Getenv("HTTP_PORT")
	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		logger.L().Fatal("failed to parse cachesize from .env",
			zap.String("error", err.Error()),
		)
	}

	// create dsn string
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassw, dbHost, dbPort, dbName,
	)

	return &Config{
		DBUser:  dbUser,
		DBPassw: dbPassw,
		DBName:  dbName,
		Dsn:     dsn,

		KafkaBroker: os.Getenv("KAFKA_BROKER"),
		KafkaTopic:  os.Getenv("KAFKA_TOPIC"),
		KafkaGroup:  os.Getenv("KAFKA_GROUP"),

		HttpPort:  httpPort,
		CacheSize: cacheSize,
	}
}
