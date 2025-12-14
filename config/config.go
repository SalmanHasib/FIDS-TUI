package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	APIKey               string
	AirportCode          string
	UpdateInterval       time.Duration
	LookaheadHours       int
	TotalFlights         int
	FlightsPerPage       int
	MaxPages             int
	PageRotationInterval time.Duration
	CharAnimationSpeed   time.Duration
}

// LoadConfig loads configuration from environment variables and sets defaults
func LoadConfig() *Config {
	cfg := &Config{
		APIKey:               getEnv("FLIGHTAWARE_API_KEY", ""),
		AirportCode:          getEnv("AIRPORT_CODE", ""),
		UpdateInterval:       10 * time.Minute,
		LookaheadHours:       6,
		TotalFlights:         50,
		FlightsPerPage:       15,
		MaxPages:             3,
		PageRotationInterval: 15 * time.Second,
		CharAnimationSpeed:   250 * time.Millisecond,
	}

	// Override with environment variables if set
	if val := os.Getenv("UPDATE_INTERVAL"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			cfg.UpdateInterval = d
		}
	}

	if val := os.Getenv("PAGE_ROTATION_INTERVAL"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			cfg.PageRotationInterval = d
		}
	}

	if val := os.Getenv("MAX_PAGES"); val != "" {
		if pages, err := strconv.Atoi(val); err == nil && pages > 0 {
			cfg.MaxPages = pages
		}
	}

	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
