package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Environment string
	Version     string
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	Storage     StorageConfig
	CORS        CORSConfig
	Log         LogConfig
	Email       EmailConfig
	SMS         SMSConfig
	Payment     PaymentConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds MongoDB configuration
type DatabaseConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// StorageConfig holds file storage configuration
type StorageConfig struct {
	Type        string // "local" or "s3"
	LocalPath   string
	BaseURL     string
	MaxFileSize int64  // Max file size in bytes
	S3Bucket    string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Endpoint  string // For MinIO
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string // "debug", "info", "warn", "error"
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	Provider string // "sendgrid" or "mailgun"
	APIKey   string
	From     string
	FromName string
}

// SMSConfig holds SMS service configuration
type SMSConfig struct {
	Provider    string // "twilio"
	AccountSID  string
	AuthToken   string
	PhoneNumber string
}

// PaymentConfig holds payment processing configuration
type PaymentConfig struct {
	Provider       string // "stripe"
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
}

// Load reads configuration from environment variables and files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, will use env vars and defaults
	}

	cfg := &Config{
		Environment: viper.GetString("ENVIRONMENT"),
		Version:     viper.GetString("VERSION"),
		Server: ServerConfig{
			Port:         viper.GetString("SERVER_PORT"),
			ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			URI:      viper.GetString("DB_URI"),
			Database: viper.GetString("DB_NAME"),
			Timeout:  viper.GetDuration("DB_TIMEOUT"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:               viper.GetString("JWT_SECRET"),
			AccessTokenDuration:  viper.GetDuration("JWT_ACCESS_DURATION"),
			RefreshTokenDuration: viper.GetDuration("JWT_REFRESH_DURATION"),
		},
		Storage: StorageConfig{
			Type:        viper.GetString("STORAGE_TYPE"),
			LocalPath:   viper.GetString("STORAGE_LOCAL_PATH"),
			S3Bucket:    viper.GetString("STORAGE_S3_BUCKET"),
			S3Region:    viper.GetString("STORAGE_S3_REGION"),
			S3AccessKey: viper.GetString("STORAGE_S3_ACCESS_KEY"),
			S3SecretKey: viper.GetString("STORAGE_S3_SECRET_KEY"),
			S3Endpoint:  viper.GetString("STORAGE_S3_ENDPOINT"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetString("CORS_ALLOWED_ORIGINS"),
		},
		Log: LogConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
		Email: EmailConfig{
			Provider: viper.GetString("EMAIL_PROVIDER"),
			APIKey:   viper.GetString("EMAIL_API_KEY"),
			From:     viper.GetString("EMAIL_FROM"),
			FromName: viper.GetString("EMAIL_FROM_NAME"),
		},
		SMS: SMSConfig{
			Provider:    viper.GetString("SMS_PROVIDER"),
			AccountSID:  viper.GetString("SMS_ACCOUNT_SID"),
			AuthToken:   viper.GetString("SMS_AUTH_TOKEN"),
			PhoneNumber: viper.GetString("SMS_PHONE_NUMBER"),
		},
		Payment: PaymentConfig{
			Provider:       viper.GetString("PAYMENT_PROVIDER"),
			SecretKey:      viper.GetString("PAYMENT_SECRET_KEY"),
			PublishableKey: viper.GetString("PAYMENT_PUBLISHABLE_KEY"),
			WebhookSecret:  viper.GetString("PAYMENT_WEBHOOK_SECRET"),
		},
	}

	// Validate required fields
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("VERSION", "0.1.0")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_READ_TIMEOUT", 15*time.Second)
	viper.SetDefault("SERVER_WRITE_TIMEOUT", 15*time.Second)
	viper.SetDefault("DB_URI", "mongodb://localhost:27017")
	viper.SetDefault("DB_NAME", "animalsys")
	viper.SetDefault("DB_TIMEOUT", 10*time.Second)
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("JWT_SECRET", "change-me-in-production")
	viper.SetDefault("JWT_ACCESS_DURATION", 15*time.Minute)
	viper.SetDefault("JWT_REFRESH_DURATION", 168*time.Hour) // 7 days
	viper.SetDefault("STORAGE_TYPE", "local")
	viper.SetDefault("STORAGE_LOCAL_PATH", "./uploads")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "*")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("EMAIL_PROVIDER", "sendgrid")
	viper.SetDefault("SMS_PROVIDER", "twilio")
	viper.SetDefault("PAYMENT_PROVIDER", "stripe")
}

// validate checks required configuration fields
func validate(cfg *Config) error {
	if cfg.Database.URI == "" {
		return fmt.Errorf("DB_URI is required")
	}
	if cfg.Database.Database == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if cfg.JWT.Secret == "change-me-in-production" && cfg.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}
	return nil
}
