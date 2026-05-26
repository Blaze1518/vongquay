package db

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Blaze1518/sinhnhatf168/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDBFromDatabaseConfig(cfg config.DatabaseConfig, logger *slog.Logger) (*gorm.DB, error) {
	dsnMaster := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
	cfg.Host, cfg.User, cfg.Password, "postgres", cfg.Port, cfg.SSLMode)
	
	logger.Info("Đang kết nối tới PostgreSQL tổng để kiểm tra database...", "host", cfg.Host, "port", cfg.Port)
	dbMaster, err := gorm.Open(postgres.Open(dsnMaster), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		return nil, fmt.Errorf("Không thể kết nối đến cơ sở dữ liệu: %w", err)
	}

	var count int
	checkQuery := "SELECT COUNT(*) FROM pg_database WHERE datname = ?"
	dbMaster.Raw(checkQuery, cfg.Name).Scan(&count)

	if count == 0 {
		logger.Info("Cơ sở dữ liệu chưa tồn tại. Đang tiến hành tạo mới...", "database_name", cfg.Name)
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.Name)
		if err := dbMaster.Exec(createDBQuery).Error; err != nil {
			return nil, fmt.Errorf("Thất bại khi tự động tạo database: %w", err)
		}
		logger.Info("Tạo cơ sở dữ liệu thành công!", "database_name", cfg.Name)
	}

	sqlMaster, err := dbMaster.DB()

	if err == nil {
		sqlMaster.Close()
	}

	dsnTarget := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsnTarget), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		return nil, fmt.Errorf("không thể kết nối đến cơ sở dữ liệu ứng dụng: %w", err)
	}


	sqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("không thể lấy DB từ gorm: %w", err)
	}

	sqlDB.SetConnMaxLifetime(time.Minute * 30)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}