package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	ServerPort string
	ServerHost string

	// Database
	DatabaseURL string

	// Kafka
	KafkaBrokers     string
	KafkaTopicEvents string
	KafkaUsername    string
	KafkaPassword    string

	// Game settings
	MatchmakingTimeout time.Duration // Time before bot is assigned
	ReconnectTimeout   time.Duration // Time allowed for reconnection
	BotMoveDelay       time.Duration // Artificial delay for bot moves

	// Feature flags
	KafkaEnabled bool
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Debug().Msg("No .env file found, using environment variables")
	}

	cfg := &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		ServerHost:         getEnv("SERVER_HOST", "0.0.0.0"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/connect_four?sslmode=disable"),
		KafkaBrokers:       getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopicEvents:   getEnv("KAFKA_TOPIC_EVENTS", "game-events"),
		KafkaUsername:      getEnv("KAFKA_USERNAME", ""),
		KafkaPassword:      getEnv("KAFKA_PASSWORD", ""),
		MatchmakingTimeout: getDurationEnv("MATCHMAKING_TIMEOUT_SECONDS", 10) * time.Second,
		ReconnectTimeout:   getDurationEnv("RECONNECT_TIMEOUT_SECONDS", 30) * time.Second,
		BotMoveDelay:       getDurationEnv("BOT_MOVE_DELAY_MS", 300) * time.Millisecond,
		KafkaEnabled:       getBoolEnv("KAFKA_ENABLED", false),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue int) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return time.Duration(intVal)
		}
	}
	return time.Duration(defaultValue)
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
