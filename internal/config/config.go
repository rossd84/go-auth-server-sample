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

type AppConfig struct {
    Env        string
    LogLevel   string
    LogFile    string
	Port string
	Version string
    DB         DBConfig
}

func LoadEnv() {
    env := os.Getenv("ENV")
    if env == "" {
        env = "development"
        log.Println("⚠️  ENV not set, defaulting to development")
    }

    filename := fmt.Sprintf(".env.%s", env)
    if err := godotenv.Load(filename); err != nil {
        log.Printf("⚠️  Failed to load %s, falling back to system env", filename)
    } else {
        log.Printf("✅ Loaded environment config from %s", filename)
    }
}

func LoadDBConfig() DBConfig {
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

func LoadAppConfig() AppConfig {
    LoadEnv()
    return AppConfig{
        Env:      os.Getenv("ENV"),
		Port: os.Getenv("PORT"),
        LogLevel: os.Getenv("LOG_LEVEL"),
        LogFile:  os.Getenv("LOG_FILE_PATH"),
		Version: os.Getenv("VERSION"),
        DB:       LoadDBConfig(),
    }
}

func (c DBConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
    )
}

