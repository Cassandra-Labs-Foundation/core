package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
	Supabase SupabaseConfig
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

// DatabaseConfig holds database related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// SupabaseConfig holds Supabase related configuration
type SupabaseConfig struct {
	URL    string
	APIKey string
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
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "bankingcore"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Supabase: SupabaseConfig{
			URL:    getEnv("SUPABASE_URL", ""),
			APIKey: getEnv("SUPABASE_API_KEY", ""),
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