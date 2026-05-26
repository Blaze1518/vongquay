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
}

type AppConfig struct {
	Name        string `mapstructure:"app_name" yaml:"name"`
	Version     string `mapstructure:"app_version" yaml:"version"`
	Environment string `mapstructure:"app_environment" yaml:"environment"`
	Debug       bool   `mapstructure:"app_debug" yaml:"debug"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"database_host" yaml:"host"`
	Port     int    `mapstructure:"database_port" yaml:"port"`
	User     string `mapstructure:"database_user" yaml:"user"`
	Password string `mapstructure:"database_password" yaml:"password"`
	Name     string `mapstructure:"database_name" yaml:"name"`
	SSLMode  string `mapstructure:"database_ssl_mode" yaml:"sslmode"`
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
	Timeout              time.Duration `mapstructure:"health_check_timeout" yaml:"timeout"` 
	DatabaseCheckEnabled bool          `mapstructure:"database_check_enabled" yaml:"database_check_enabled"`
}

func LoadConfig(logger *slog.Logger) (*Config, error) {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer("_", "."))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("không thể tải file cấu hình: %w", err)
		}
	}

	logger.Info("Môi trường hiện tại được xác định", "env", v.AllSettings())

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Không thể ép kiểu vào struct Config: %w", err)
	}

	logger.Info("Dữ liệu cấu hình đã được Unmarshal vào Struct thành công", "config", cfg)

	return &cfg, nil
}

func GetSkipPaths(env string) []string {
	switch env {
	case "production":
		return []string{"/health", "/health/live", "/health/ready", "/metrics", "/debug", "/pprof"}
	case "development":
		return []string{"/health", "/health/live", "/health/ready"}
	case "test":
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