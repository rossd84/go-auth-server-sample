package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	Port       string
	Version    string
	JWTSecret  string `mapstructure:"jwt_secret"`
	JWTRefresh string `mapstructure:"jwt_refresh"`
	JWTIssuer  string `mapstructure:"jwt_issuer"`
	DB         DBConfig
}

func BindFlags() {
	pflag.String("env", "", "Environment to run the app in (e.g. development, production)")
	pflag.String("port", "", "Port the server should run on")
	pflag.String("log_level", "", "Log level (debug, info, warn, error)")
	pflag.String("log_file_path", "", "Path to log file")
	pflag.Parse()

	if err := viper.BindPFlag("ENV", pflag.Lookup("env")); err != nil {
		log.Fatalf("failed to bind ENV flag: %v", err)
	}
	if err := viper.BindPFlag("PORT", pflag.Lookup("port")); err != nil {
		log.Fatalf("failed to bind PORT flag: %v", err)
	}
	if err := viper.BindPFlag("LOG_LEVEL", pflag.Lookup("log_level")); err != nil {
		log.Fatalf("failed to bind LOG_LEVEL flag: %v", err)
	}
	if err := viper.BindPFlag("LOG_FILE_PATH", pflag.Lookup("log_file_path")); err != nil {
		log.Fatalf("failed to bind LOG_FILE_PATH flag: %v", err)
	}
}

func LoadAppConfig() AppConfig {
	BindFlags()

	env := viper.GetString("ENV")
	if env == "" {
		env = "development"
		log.Println("⚠️  ENV not set, defaulting to development")
	}

	configFile := fmt.Sprintf("environments/.env.%s", env)
	viper.SetConfigFile(configFile)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️  Failed to load %s, falling back to environment only", configFile)
	} else {
		log.Printf("✅ Loaded config from %s", configFile)
	}

	// Required DB values
	required := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, key := range required {
		if viper.GetString(key) == "" {
			log.Fatalf("❌ Missing required config value: %s", key)
		}
	}

	return AppConfig{
		Env:        env,
		LogLevel:   viper.GetString("LOG_LEVEL"),
		LogFile:    viper.GetString("LOG_FILE_PATH"),
		Port:       viper.GetString("PORT"),
		Version:    viper.GetString("VERSION"),
		JWTSecret:  viper.GetString("JWT_SECRET"),
		JWTRefresh: viper.GetString("JWT_REFRESH"),
		JWTIssuer:  viper.GetString("JWT_ISSUER"),
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  defaultOrValue(viper.GetString("SSL_MODE"), "disable"),
		},
	}
}

func defaultOrValue(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func (c AppConfig) IsDev() bool {
	return c.Env == "development"
}
