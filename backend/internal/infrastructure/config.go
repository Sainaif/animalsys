package infrastructure

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Security SecurityConfig
	RateLimit RateLimitConfig
	Upload   UploadConfig
	Log      LogConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
	URL         string
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	URI      string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type SecurityConfig struct {
	JWTSecret              string
	JWTAccessTokenExpiry   time.Duration
	JWTRefreshTokenExpiry  time.Duration
	BcryptCost             int
	CORSAllowedOrigins     []string
}

type RateLimitConfig struct {
	RequestsPerIP     int
	WindowMinutes     int
	RequestsPerUser   int
	WindowHours       int
}

type UploadConfig struct {
	MaxSize           int64
	AllowedFileTypes  []string
	Path              string
}

type LogConfig struct {
	Level  string
	Format string
	File   string
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Read config file (optional, env vars take precedence)
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist, we can use env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Parse JWT token expiry
	accessTokenExpiry, err := time.ParseDuration(viper.GetString("JWT_ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		accessTokenExpiry = 15 * time.Minute // default
	}

	refreshTokenExpiry, err := time.ParseDuration(viper.GetString("JWT_REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		refreshTokenExpiry = 168 * time.Hour // default 7 days
	}

	config := &Config{
		App: AppConfig{
			Name:        viper.GetString("APP_NAME"),
			Environment: viper.GetString("APP_ENV"),
			Debug:       viper.GetBool("APP_DEBUG"),
			URL:         viper.GetString("APP_URL"),
		},
		Server: ServerConfig{
			Host: viper.GetString("API_HOST"),
			Port: viper.GetString("API_PORT"),
		},
		Database: DatabaseConfig{
			URI:      viper.GetString("MONGO_URI"),
			Host:     viper.GetString("MONGO_HOST"),
			Port:     viper.GetString("MONGO_PORT"),
			Database: viper.GetString("MONGO_DATABASE"),
			Username: viper.GetString("MONGO_USERNAME"),
			Password: viper.GetString("MONGO_PASSWORD"),
		},
		Security: SecurityConfig{
			JWTSecret:             viper.GetString("JWT_SECRET"),
			JWTAccessTokenExpiry:  accessTokenExpiry,
			JWTRefreshTokenExpiry: refreshTokenExpiry,
			BcryptCost:            viper.GetInt("BCRYPT_COST"),
			CORSAllowedOrigins:    viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
		},
		RateLimit: RateLimitConfig{
			RequestsPerIP:   viper.GetInt("RATE_LIMIT_REQUESTS_PER_IP"),
			WindowMinutes:   viper.GetInt("RATE_LIMIT_WINDOW_MINUTES"),
			RequestsPerUser: viper.GetInt("RATE_LIMIT_REQUESTS_PER_USER"),
			WindowHours:     viper.GetInt("RATE_LIMIT_WINDOW_HOURS"),
		},
		Upload: UploadConfig{
			MaxSize:          viper.GetInt64("MAX_UPLOAD_SIZE"),
			AllowedFileTypes: viper.GetStringSlice("ALLOWED_FILE_TYPES"),
			Path:             viper.GetString("UPLOAD_PATH"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
			File:   viper.GetString("LOG_FILE"),
		},
	}

	// Validate config
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("APP_NAME", "AnimalSys")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_DEBUG", true)
	viper.SetDefault("APP_URL", "http://localhost:3000")

	// Server defaults
	viper.SetDefault("API_HOST", "0.0.0.0")
	viper.SetDefault("API_PORT", "8080")

	// Database defaults
	viper.SetDefault("MONGO_HOST", "localhost")
	viper.SetDefault("MONGO_PORT", "27017")
	viper.SetDefault("MONGO_DATABASE", "animalsys")

	// Security defaults
	viper.SetDefault("JWT_ACCESS_TOKEN_EXPIRY", "15m")
	viper.SetDefault("JWT_REFRESH_TOKEN_EXPIRY", "168h")
	viper.SetDefault("BCRYPT_COST", 12)
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	// Rate limit defaults
	viper.SetDefault("RATE_LIMIT_REQUESTS_PER_IP", 100)
	viper.SetDefault("RATE_LIMIT_WINDOW_MINUTES", 15)
	viper.SetDefault("RATE_LIMIT_REQUESTS_PER_USER", 1000)
	viper.SetDefault("RATE_LIMIT_WINDOW_HOURS", 1)

	// Upload defaults
	viper.SetDefault("MAX_UPLOAD_SIZE", 10485760) // 10MB
	viper.SetDefault("UPLOAD_PATH", "/uploads")

	// Log defaults
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("LOG_FILE", "logs/app.log")
}

func validateConfig(config *Config) error {
	if config.Security.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(config.Security.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	if config.Database.URI == "" && config.Database.Host == "" {
		return fmt.Errorf("database configuration is required")
	}

	if config.Server.Port == "" {
		return fmt.Errorf("API_PORT is required")
	}

	return nil
}
