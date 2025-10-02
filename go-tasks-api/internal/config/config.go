package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	// Server config
	Port         string        `env:"PORT" default:"8000"`
	Host         string        `env:"HOST" default:"0.0.0.0"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" default:"30s"`
	Debug        bool          `env:"DEBUG" default:"false"`
	// App config
	AppName     string `env:"APP_NAME" default:"Tasks API"`
	AppVersion  string `env:"APP_VERSION" default:"1.0.0"`
	Environment string `env:"ENVIRONMENT" default:"development"`

	// Database config (futuro)
	DatabaseURL string `env:"DATABASE_URL" default:""`

	// Auth config (futuro)
	JWTSecret string `env:"JWT_SECRET" default:""`
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	cfg := &Config{}

	// Server
	cfg.Port = getEnvWithDefault("PORT", "8000")
	cfg.Host = getEnvWithDefault("HOST", "0.0.0.0")
	cfg.ReadTimeout = parseDuration("READ_TIMEOUT", "30s")
	cfg.WriteTimeout = parseDuration("WRITE_TIMEOUT", "30s")
	cfg.Debug = parseBool("DEBUG", "false")

	// App
	cfg.AppName = getEnvWithDefault("APP_NAME", "Tasks API")
	cfg.AppVersion = getEnvWithDefault("APP_VERSION", "1.0.0")
	cfg.Environment = getEnvWithDefault("ENVIRONMENT", "development")

	// Database
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")

	// Auth
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	// Validar configuración crítica en producción
	if cfg.Environment == "production" {
		if cfg.JWTSecret == "" {
			return nil, fmt.Errorf("JWT_SECRET is required in production")
		}
		if cfg.DatabaseURL == "" {
			return nil, fmt.Errorf("DATABASE_URL is required in production")
		}
	}

	return cfg, nil
}

// getEnvWithDefault obtiene una variable de entorno con fallback
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parsea una duración desde ENV con fallback
func parseDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

// parseBool parsea un booleano desde ENV con fallback
func parseBool(key, defaultValue string) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	boolValue, _ := strconv.ParseBool(defaultValue)
	return boolValue
}

// IsDevelopment verifica si estamos en desarrollo
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction verifica si estamos en producción
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
