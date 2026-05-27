package server

import (
	"time"

	"github.com/Blaze1518/sinhnhatf168/internal/auth"
	"github.com/Blaze1518/sinhnhatf168/internal/config"
	"github.com/Blaze1518/sinhnhatf168/internal/errors"
	"github.com/Blaze1518/sinhnhatf168/internal/game/campaign"
	"github.com/Blaze1518/sinhnhatf168/internal/game/prize"
	"github.com/Blaze1518/sinhnhatf168/internal/game/ticket"
	"github.com/Blaze1518/sinhnhatf168/internal/game/winner"
	"github.com/Blaze1518/sinhnhatf168/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(authHandler *auth.Handler, campaignHandler *campaign.Handler, prizeHandler *prize.Handler, ticketHandler *ticket.Handler, winnerHandler *winner.Handler, cfg *config.Config) *gin.Engine {
	router := gin.New()

	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	skipPaths := config.GetSkipPaths(cfg.App.Environment)
	loggerConfig := middleware.NewLoggerConfig(
		cfg.Logging.GetLogLevel(),
		skipPaths,
	)
	router.Use(middleware.Logger(loggerConfig))
	router.Use(errors.ErrorHandler())
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	router.Use(cors.New(corsConfig))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var defaultLimit gin.HandlerFunc
	rlCfg := cfg.Ratelimit
	if rlCfg.Enabled {
		defaultLimit = middleware.NewRateLimitMiddleware(
			rlCfg.Window,
			rlCfg.Requests,
			middleware.ResolveClientIP,
			nil,
		)
	}

	var authLimit, drawLimit, importLimit gin.HandlerFunc
	if rlCfg.Enabled {
		authLimit = middleware.NewRateLimitMiddleware(time.Minute, 5, middleware.ResolveClientIP, nil)
		drawLimit = middleware.NewRateLimitMiddleware(time.Minute, 10, middleware.ResolveClientIP, nil)
		importLimit = middleware.NewRateLimitMiddleware(time.Minute, 2, middleware.ResolveClientIP, nil)
	}

	// router.Use(whitelistip.WhitelistIPMiddleware(
	// 	middleware.ResolveClientIP,
	// 	whitelistIPService,
	// ))

	v1 := router.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			if rlCfg.Enabled {
				authGroup.GET("/login", authLimit, authHandler.Login)
			} else {
				authGroup.GET("/login", authHandler.Login)
			}
		}
		gameGroup := v1.Group("/game")
		if rlCfg.Enabled {
			gameGroup.Use(defaultLimit) 
		}
		{
			campaignGroup := gameGroup.Group("campaign")
			{
				campaignGroup.POST("/", campaignHandler.Create)
			}
			prizeGroup := gameGroup.Group("prize")
			{
				prizeGroup.POST("/", prizeHandler.Create)
			}
			ticketGroup := gameGroup.Group("ticket")
			{
				ticketGroup.POST("/", ticketHandler.Create)
				if rlCfg.Enabled {
					ticketGroup.POST("/import", importLimit, ticketHandler.ImportExcel)
				} else {
					ticketGroup.POST("/import", ticketHandler.ImportExcel)
				}
			}
			winnerGroup := gameGroup.Group("winner")
			{
				if rlCfg.Enabled {
					winnerGroup.POST("/draw", drawLimit, winnerHandler.Draw)
				} else {
					winnerGroup.POST("/draw", winnerHandler.Draw)
				}
			}
		}
	}


	return router
}