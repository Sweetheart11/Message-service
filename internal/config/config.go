package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
	User        string
	Password    string
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
			Address:     getEnv("HTTP_SERVER_ADDRESS", "localhost:8080"),
			Timeout:     getDurationEnv("HTTP_SERVER_TIMEOUT", 4*time.Second),
			IdleTimeout: getDurationEnv("HTTP_SERVER_IDLE_TIMEOUT", 60*time.Second),
			User:        getEnvOrPanic("HTTP_SERVER_USER"),
			Password:    getEnvOrPanic("HTTP_SERVER_PASSWORD"),
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
