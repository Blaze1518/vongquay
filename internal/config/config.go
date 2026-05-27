package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// executablePath is a package-level variable so tests can override it.
var executablePath = os.Executable

type Config struct {
	App        AppConfig        `mapstructure:",squash"`
	Database   DatabaseConfig   `mapstructure:",squash"`
	JWT        JWTConfig        `mapstructure:",squash"`
	Server     ServerConfig     `mapstructure:",squash"`
	Logging    LoggingConfig    `mapstructure:",squash"`
	Ratelimit  RateLimitConfig  `mapstructure:",squash"`
	Migrations MigrationsConfig `mapstructure:",squash"`
	Health     HealthConfig     `mapstructure:",squash"`
	PgAdmin    PgAdminConfig    `mapstructure:",squash"`
}

type AppConfig struct {
	Name        string `mapstructure:"app_name" yaml:"name"`
	Version     string `mapstructure:"app_version" yaml:"version"`
	Environment string `mapstructure:"app_environment" yaml:"environment"`
	Debug       bool   `mapstructure:"app_debug" yaml:"debug"`
}

type DatabaseConfig struct {
	Host          string `mapstructure:"database_host" yaml:"host"`
	Port          int    `mapstructure:"database_port" yaml:"port"`
	User          string `mapstructure:"database_user" yaml:"user"`
	Password      string `mapstructure:"database_password" yaml:"password"`
	Name          string `mapstructure:"database_name" yaml:"name"`
	SSLMode       string `mapstructure:"database_ssl_mode" yaml:"sslmode"`
	PostgresUser  string `mapstructure:"postgres_user" yaml:"postgres_user"`
	PostgresPass  string `mapstructure:"postgres_password" yaml:"postgres_password"`
	PostgresDB    string `mapstructure:"postgres_db" yaml:"postgres_db"`
	URL           string `mapstructure:"db_url" yaml:"db_url"`
	ContainerName string `mapstructure:"database_container_name" yaml:"database_container_name"`
}

type JWTConfig struct {
	Secret          string        `mapstructure:"jwt_secret" yaml:"secret"`
	AccessTokenTTL  time.Duration `mapstructure:"jwt_access_token_ttl" yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"jwt_refresh_token_ttl" yaml:"refresh_token_ttl"`
	TTLHours        int           `mapstructure:"jwt_ttlhours" yaml:"ttlhours"`
}

type ServerConfig struct {
	Port            string `mapstructure:"app_port" yaml:"port"` 
	ReadTimeout     int    `mapstructure:"server_readtimeout" yaml:"readtimeout"`
	WriteTimeout    int    `mapstructure:"server_writetimeout" yaml:"writetimeout"`
	IdleTimeout     int    `mapstructure:"server_idletimeout" yaml:"idletimeout"`
	ShutdownTimeout int    `mapstructure:"server_shutdowntimeout" yaml:"shutdowntimeout"`
	MaxHeaderBytes  int    `mapstructure:"server_maxheaderbytes" yaml:"maxheaderbytes"`
}

type LoggingConfig struct {
	Level string `mapstructure:"logging_level" yaml:"level"`
}

type RateLimitConfig struct {
	Enabled  bool          `mapstructure:"ratelimit_enabled" yaml:"enabled"`
	Requests int           `mapstructure:"ratelimit_requests" yaml:"requests"`
	Window   time.Duration `mapstructure:"ratelimit_window" yaml:"window"`
}

type MigrationsConfig struct {
	Directory   string `mapstructure:"migrations_directory" yaml:"directory"`
	Timeout     int    `mapstructure:"migrations_timeout" yaml:"timeout"`
	LockTimeout int    `mapstructure:"migrations_locktimeout" yaml:"locktimeout"`
}

type HealthConfig struct {
	Interval             time.Duration `mapstructure:"health_check_interval" yaml:"interval"`
	Timeout              time.Duration `mapstructure:"health_check_timeout" yaml:"timeout"` 
	Retries              int           `mapstructure:"health_check_retries" yaml:"retries"`
	StartPeriod          time.Duration `mapstructure:"health_check_start_period" yaml:"start_period"`
	DatabaseCheckEnabled bool          `mapstructure:"database_check_enabled" yaml:"database_check_enabled"`
}

type PgAdminConfig struct {
	DefaultEmail    string `mapstructure:"pgadmin_default_email" yaml:"default_email"`
	DefaultPassword string `mapstructure:"pgadmin_default_password" yaml:"default_password"`
}

func LoadConfig(logger *slog.Logger) (*Config, error) {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	// Đồng bộ hóa cách map key từ Env sang Struct (Viper tự động chuyển uppercase thành lowercase)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// --- Đặt giá trị mặc định (Defaults) dựa theo file .env của bạn ---
	
	// Application & Server
	v.SetDefault("app_name", "sinhnhatf168")
	v.SetDefault("app_version", "1.0.0")
	v.SetDefault("app_environment", "production")
	v.SetDefault("app_debug", false)
	v.SetDefault("app_port", "8080")
	
	// Database & Docker Postgres Container
	v.SetDefault("database_host", "postgres-db")
	v.SetDefault("database_port", 5432)
	v.SetDefault("database_user", "admin")
	v.SetDefault("database_password", "your_password_here")
	v.SetDefault("database_name", "sinhnhatf168")
	v.SetDefault("database_ssl_mode", "disable")
	v.SetDefault("postgres_user", "admin")
	v.SetDefault("postgres_password", "your_password_here")
	v.SetDefault("postgres_db", "sinhnhatf168")
	v.SetDefault("db_url", "postgres://admin:your_password_here@postgres-db:5432/sinhnhatf168")
	v.SetDefault("database_container_name", "go_api_db")
	
	// JWT
	v.SetDefault("jwt_secret", "super-secret-key-change-me-in-production")
	v.SetDefault("jwt_access_token_ttl", 15 * time.Minute)
	v.SetDefault("jwt_refresh_token_ttl", 168 * time.Hour)
	v.SetDefault("jwt_ttlhours", 0)
	
	// Logging
	v.SetDefault("logging_level", "info")
	
	// Rate Limiting
	v.SetDefault("ratelimit_enabled", true)
	v.SetDefault("ratelimit_requests", 100)
	v.SetDefault("ratelimit_window", 1 * time.Minute)
	
	// Migrations
	v.SetDefault("migrations_directory", "migrations")
	v.SetDefault("migrations_timeout", 30)
	v.SetDefault("migrations_locktimeout", 10)
	
	// Docker Health Check
	v.SetDefault("health_check_interval", 15 * time.Second)
	v.SetDefault("health_check_timeout", 5 * time.Second)
	v.SetDefault("health_check_retries", 5)
	v.SetDefault("health_check_start_period", 15 * time.Second)
	v.SetDefault("database_check_enabled", false) // Giữ thuộc tính bổ sung cũ của bạn

	// PgAdmin
	v.SetDefault("pgadmin_default_email", "blazer@attcloud.org")
	v.SetDefault("pgadmin_default_password", "blazer@attcloud.org")
	
	// Đọc từ file .env nếu có (ghi đè lên các Default ở trên)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("không thể tải file cấu hình: %w", err)
		}
	}

	logger.Info("Môi trường hiện tại được xác định", "env", v.AllSettings())

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("không thể ép kiểu vào struct Config: %w", err)
	}

	logger.Info("Dữ liệu cấu hình đã được Unmarshal vào Struct thành công", "config", cfg)

	return &cfg, nil
}

func GetSkipPaths(env string) []string {
	switch env {
	case "production":
		return []string{"/health", "/health/live", "/health/ready", "/metrics", "/debug", "/pprof"}
	case "development", "test":
		return []string{"/health", "/health/live", "/health/ready"}
	default:
		return []string{"/health", "/health/live", "/health/ready"}
	}
}

func (c *LoggingConfig) GetLogLevel() slog.Level {
	switch strings.ToLower(c.Level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}