package config

import (
	"log"
	"os"

	"github.com/joho/godotenv" // 用于本地开发加载.env文件
)

type AppConfig struct {
	MySQLDSN    string
	RedisAddr   string
	KafkaBroker string
	KafkaTopic  string
	ServerPort  string
}

func LoadConfig() *AppConfig {
	// 加载 .env 文件，仅在本地开发时使用
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, assuming environment variables are set.")
	}

	cfg := &AppConfig{
		MySQLDSN:    getEnv("MYSQL_DSN", "root:qq1595215609@tcp(mysql:3306)/todo_app?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:  getEnv("KAFKA_TOPIC", "todo_events"),
		ServerPort:  getEnv("SERVER_PORT", "8000"),
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
