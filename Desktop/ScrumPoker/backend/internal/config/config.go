package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          string
	AppOrigin     string
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPass     string
	MySQLDB       string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	PresenceTTL   time.Duration
}

func Load() Config {
	return Config{
		Port:          env("PORT", "8080"),
		AppOrigin:     env("APP_ORIGIN", "http://localhost:5173"),
		MySQLHost:     env("MYSQL_HOST", "localhost"),
		MySQLPort:     env("MYSQL_PORT", "3306"),
		MySQLUser:     env("MYSQL_USER", "root"),
		MySQLPass:     env("MYSQL_PASSWORD", ""),
		MySQLDB:       env("MYSQL_DATABASE", "scrum_poker"),
		RedisAddr:     env("REDIS_ADDR", "localhost:6379"),
		RedisPassword: env("REDIS_PASSWORD", ""),
		RedisDB:       envInt("REDIS_DB", 0),
		PresenceTTL:   envDuration("PRESENCE_TTL", 45*time.Second),
	}
}

func (c Config) MySQLDSN() string {
	return c.MySQLUser + ":" + c.MySQLPass + "@tcp(" + c.MySQLHost + ":" + c.MySQLPort + ")/" + c.MySQLDB + "?parseTime=true&multiStatements=true&charset=utf8mb4"
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
