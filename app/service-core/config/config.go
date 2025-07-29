package config

import (
	"os"
	"strings"
	"time"
)

func isRunningTest() bool {
	for _, arg := range os.Args {
		if strings.HasSuffix(arg, ".test") {
			return true
		}
	}

	return false
}

func MustSetEnv(active bool, key string) string {
	value := os.Getenv(key)
	if active && value == "" {
		if isRunningTest() {
			return "test"
		}
		panic("Missing environment variable: " + key)
	}

	return os.Getenv(key)
}

type Config struct {
	// General
	LogLevel  string
	HTTPPort  string
	GRPCPort  string
	Domain    string
	CoreURL   string
	AdminURL  string
	ClientURL string
	CronToken string

	// Constants
	MaxFileSize     int64
	HTTPTimeout     time.Duration
	ContextTimeout  time.Duration
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration

	// Postgres
	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
}

func LoadConfig() *Config {
	const (
		HTTPTimeout     = 10 * time.Second
		ContextTimeout  = 10 * time.Second
		AccessTokenExp  = 15 * time.Minute
		RefreshTokenExp = 30 * 24 * time.Hour
		MaxFileSize     = 10 << 20
	)
	return &Config{
		LogLevel:              MustSetEnv(true, "LOG_LEVEL"),
		HTTPPort:              MustSetEnv(true, "HTTP_PORT"),
		GRPCPort:              MustSetEnv(true, "GRPC_PORT"),
		Domain:                MustSetEnv(true, "DOMAIN"),
		CoreURL:               MustSetEnv(true, "CORE_URL"),
		AdminURL:              MustSetEnv(true, "ADMIN_URL"),
		ClientURL:             MustSetEnv(true, "CLIENT_URL"),
		CronToken:             MustSetEnv(true, "CRON_TOKEN"),
		HTTPTimeout:           HTTPTimeout,
		ContextTimeout:        ContextTimeout,
		AccessTokenExp:        AccessTokenExp,
		RefreshTokenExp:       RefreshTokenExp,
		MaxFileSize:           MaxFileSize,
		PostgresHost:          MustSetEnv(true, "POSTGRES_HOST"),
		PostgresPort:          MustSetEnv(true, "POSTGRES_PORT"),
		PostgresDB:            MustSetEnv(true, "POSTGRES_DB"),
		PostgresUser:          MustSetEnv(true, "POSTGRES_USER"),
		PostgresPassword:      MustSetEnv(true, "POSTGRES_PASSWORD"),
	}
}

func LoadTestConfig() *Config {
	const (
		HTTPTimeout                = 10 * time.Second
		ContextTimeout             = 10 * time.Second
		AccessTokenExp             = 5 * time.Minute
		RefreshTokenExp            = 30 * 24 * time.Hour
		MaxFileSize                = 10 << 20
	)
	return &Config{
		LogLevel:              "debug",
		HTTPPort:              "8080",
		GRPCPort:              "50051",
		Domain:                "localhost",
		CoreURL:               "http://localhost:8080",
		AdminURL:              "http://localhost:8080",
		ClientURL:             "http://localhost:3000",
		CronToken:             "test",
		HTTPTimeout:           HTTPTimeout,
		ContextTimeout:        ContextTimeout,
		AccessTokenExp:        AccessTokenExp,
		RefreshTokenExp:       RefreshTokenExp,
		MaxFileSize:           MaxFileSize,
		PostgresHost:          "localhost",
		PostgresPort:          "5432",
		PostgresDB:            "test",
		PostgresUser:          "test",
		PostgresPassword:      "test",
	}
}
