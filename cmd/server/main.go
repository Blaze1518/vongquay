package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/Blaze1518/sinhnhatf168/docs"
	"github.com/Blaze1518/sinhnhatf168/internal/auth"
	"github.com/Blaze1518/sinhnhatf168/internal/config"
	"github.com/Blaze1518/sinhnhatf168/internal/db"
	"github.com/Blaze1518/sinhnhatf168/internal/game/campaign"
	"github.com/Blaze1518/sinhnhatf168/internal/game/prize"
	"github.com/Blaze1518/sinhnhatf168/internal/game/ticket"
	"github.com/Blaze1518/sinhnhatf168/internal/game/winner"
	"github.com/Blaze1518/sinhnhatf168/internal/migrate"
	"github.com/Blaze1518/sinhnhatf168/internal/server"
)

// @title HappyBirthDayF168
// @version 1.0
// @description API documentation for HappyBirthDayF168
// @termsOfService http://swagger.io/terms/

// @host api-sinhnhatf168.attservice.org
// @BasePath /api/v1

func main() {
	if err := run(); err != nil {
		slog.Error("Server crashed on startup", "error", err)
		os.Exit(1)
	}
}

func run() error {
	logger := slog.Default()
	logger.Info("Starting server...")

	cfg, err := config.LoadConfig(logger) 

	if err != nil {
		logger.Error("Lỗi khi tải cấu hình", "error", err)
		return err
	}

	
	database, err := db.NewPostgresDBFromDatabaseConfig(cfg.Database, logger)
	if err != nil {
		logger.Error("Lỗi khi kết nối đến cơ sở dữ liệu", "error", err)
		return err
	}

	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("Có lỗi khi lấy DB: %w", err)
	}

	timeout := time.Duration(cfg.Migrations.Timeout) * time.Second
	migrator, err := migrate.New(sqlDB, migrate.Config{
		MigrationsDir: cfg.Migrations.Directory,
		Timeout:       time.Duration(cfg.Migrations.Timeout) * time.Second,
		LockTimeout:   time.Duration(cfg.Migrations.LockTimeout) * time.Second,
	})

	if err != nil {
		return fmt.Errorf("Có lỗi khi tạo migrator: %w", err)
	}
    defer migrator.Close() // Nhớ close để giải phóng bộ nhớ

    // 1. Tạo context kiểm soát timeout cho việc chạy migration
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // 2. THAY VÌ CHỈ XEM VERSION, HÃY RA LỆNH UP LUÔN
    // Hàm Up() bên trong của bạn đã tự xử lý nếu không có file mới (ErrNoChange) rồi.
    if err := migrator.Up(ctx); err != nil {
        return fmt.Errorf("Tự động chạy migration thất bại: %w", err)
    }

	campaignRepo := campaign.NewRepository(database)
	prizeRepo := prize.NewRepository(database)
	ticketRepo := ticket.NewRepository(database)
	winnerRepo := winner.NewRepository(database)

	campaignService := campaign.NewService(campaignRepo)
	prizeService := prize.NewService(prizeRepo)
	ticketService := ticket.NewService(ticketRepo)
	importService := ticket.NewImportService(ticketRepo)
	winnerService := winner.NewService(winnerRepo, ticketRepo, prizeRepo, database)

	campaignHandler := campaign.NewHandler(campaignService)
	prizeHandler := prize.NewHandler(prizeService)
	ticketHandler := ticket.NewHandler(ticketService, importService)
	winnerHandler := winner.NewHandler(winnerService)

	authHandler := auth.NewHandler()

	router := server.SetupRouter(authHandler, campaignHandler, prizeHandler, ticketHandler, winnerHandler, cfg)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}
	addr := fmt.Sprintf(":%s", port)
	logger.Info("Server đang lắng nghe tại", "địa chỉ", addr)

	if err := router.Run(addr); err != nil {
		logger.Error("Lỗi khi chạy server", "error", err)
		return err
	}

	return nil
}