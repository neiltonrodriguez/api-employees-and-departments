package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port          string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	AppEnv        string
	LogLevel      string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
	CacheTTL      string
}

func Load() (*Config, error) {
	c := &Config{
		Port:          getenv("APP_PORT", "8080"),
		DBHost:        getenv("DATABASE_HOST", "localhost"),
		DBPort:        getenv("DATABASE_PORT", "5432"),
		DBUser:        getenv("DATABASE_USER", "postgres"),
		DBPassword:    getenv("DATABASE_PASSWORD", "postgres"),
		DBName:        getenv("DATABASE_NAME", "companydb"),
		DBSSLMode:     getenv("DATABASE_SSLMODE", "disable"),
		AppEnv:        getenv("APP_ENV", "development"),
		LogLevel:      getenv("LOG_LEVEL", "info"),
		RedisHost:     getenv("REDIS_HOST", "localhost"),
		RedisPort:     getenv("REDIS_PORT", "6379"),
		RedisPassword: getenv("REDIS_PASSWORD", ""),
		RedisDB:       getenv("REDIS_DB", "0"),
		CacheTTL:      getenv("CACHE_TTL_SECONDS", "300"),
	}
	return c, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBName, c.DBPassword, c.DBSSLMode)
}
