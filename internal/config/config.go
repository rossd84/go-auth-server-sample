package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
)

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
}

func LoadDBConfig() DBConfig {
    // Load .env.api first, fall back to .env
    if err := godotenv.Load(".env.api"); err != nil {
        log.Println("⚠️  .env.api not found, falling back to .env")
        _ = godotenv.Load(".env")
    }

	required := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			log.Fatalf("❌ Missing required environment variable: %s", key)
		}
	}

    sslMode := os.Getenv("SSL_MODE")
    if sslMode == "" {
        sslMode = "disable"
    }

    return DBConfig{
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        Name:     os.Getenv("DB_NAME"),
        SSLMode:  sslMode,
    }
}

func (c DBConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
    )
}

