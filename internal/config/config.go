package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost       string
	DBPort       int
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	Port         string
	Environment  string
	RedisURL     string
	AWSRegion    string
	AWSBucket    string
	AWSAccessKey string
	AWSSecretKey string

	// Elasticsearch configuration
	ElasticsearchURL      string `mapstructure:"ELASTICSEARCH_URL"`
	ElasticsearchUsername string `mapstructure:"ELASTICSEARCH_USERNAME"`
	ElasticsearchPassword string `mapstructure:"ELASTICSEARCH_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load("/home/administrator/test/irm_backend/.env"); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %v", err)
	}

	viper.SetDefault("ELASTICSEARCH_URL", "http://localhost:9200")
	viper.AutomaticEnv()

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	config.DBHost = getEnv("DB_HOST", "localhost")
	config.DBPort = dbPort
	config.DBUser = getEnv("DB_USER", "postgres")
	config.DBPassword = getEnv("DB_PASSWORD", "")
	config.DBName = getEnv("DB_NAME", "irm_db")
	config.JWTSecret = getEnv("JWT_SECRET", "your-secret-key")
	config.Port = getEnv("PORT", "8080")
	config.Environment = getEnv("ENV", "development")
	config.RedisURL = getEnv("REDIS_URL", "redis://localhost:6379")
	config.AWSRegion = getEnv("AWS_REGION", "")
	config.AWSBucket = getEnv("AWS_BUCKET", "")
	config.AWSAccessKey = getEnv("AWS_ACCESS_KEY", "")
	config.AWSSecretKey = getEnv("AWS_SECRET_KEY", "")
	config.ElasticsearchURL = getEnv("ELASTICSEARCH_URL", "http://localhost:9200")
	config.ElasticsearchUsername = getEnv("ELASTICSEARCH_USERNAME", "")
	config.ElasticsearchPassword = getEnv("ELASTICSEARCH_PASSWORD", "")

	return &config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
