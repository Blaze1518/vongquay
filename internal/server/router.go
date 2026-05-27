package server

import (
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

	rlCfg := cfg.Ratelimit
	if rlCfg.Enabled {
		router.Use(
			middleware.NewRateLimitMiddleware(
				rlCfg.Window,
				rlCfg.Requests,
				middleware.ResolveClientIP,
				nil,
			),
		)
	}

	// router.Use(whitelistip.WhitelistIPMiddleware(
	// 	middleware.ResolveClientIP,
	// 	whitelistIPService,
	// ))

	v1 := router.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.GET("/login", authHandler.Login)
		}
		gameGroup := v1.Group("/game")
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
				ticketGroup.POST("/import", ticketHandler.ImportExcel)
			}
			winnerGroup := gameGroup.Group("winner")
			{
				winnerGroup.POST("/draw", winnerHandler.Draw)
			}
		}
	}


	return router
}