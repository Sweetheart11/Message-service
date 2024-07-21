package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	Kafka       Kafka      `yaml:"kafka"`
	Database    Database   `yaml:"database"`
}

type HTTPServer struct {
	Addr        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type Kafka struct {
	Broker string `yaml:"broker" env-required:"true"`
	Topic  string `yaml:"topic" env-required:"true"`
}

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
	Name     string
}

func MustLoad() *Config {
	configPath := FetchConfigPath()
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	cfg := &Config{
		Env: getEnv("ENV", "local"),
		HTTPServer: HTTPServer{
			Addr:        getEnv("HTTP_SERVER_ADDRESS", "localhost:8080"),
			Timeout:     getDurationEnv("HTTP_SERVER_TIMEOUT", 4*time.Second),
			IdleTimeout: getDurationEnv("HTTP_SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Kafka: Kafka{
			Broker: getEnv("KAFKA_BROKER", "localhost:9092"),
			Topic:  getEnv("KAFKA_TOPIC", "messages"),
		},
		Database: Database{
			User:     getEnvOrPanic("DB_USER"),
			Password: getEnvOrPanic("DB_PASSWORD"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Name:     getEnv("DB_NAME", "messages"),
		},
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvOrPanic(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("environment variable %s is required", key)
	}
	return value
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		parsedValue, err := time.ParseDuration(value)
		if err != nil {
			log.Fatalf("invalid duration format for %s: %s", key, value)
		}
		return parsedValue
	}
	return defaultValue
}

func FetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", ".env", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
