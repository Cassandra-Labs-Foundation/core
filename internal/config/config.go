package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	JWT    JWTConfig
}

// ServerConfig holds server related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// JWTConfig holds JWT related configuration
type JWTConfig struct {
	Secret       string
	ExpiryMinutes int
}

// Load returns configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  time.Duration(getEnvAsInt("SERVER_READ_TIMEOUT", 10)) * time.Second,
			WriteTimeout: time.Duration(getEnvAsInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second,
		},
		JWT: JWTConfig{
			Secret:       getEnv("JWT_SECRET", "your-secret-key"),
			ExpiryMinutes: getEnvAsInt("JWT_EXPIRY_MINUTES", 60),
		},
	}
}

// Simple helper function to read environment variables with a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// Helper function to read an environment variable as an integer with a default value
func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}