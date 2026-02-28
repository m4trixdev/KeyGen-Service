package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m4trixdev/keygen-service/config"
	"github.com/m4trixdev/keygen-service/internal/handlers"
	"github.com/m4trixdev/keygen-service/internal/middleware"
	"github.com/m4trixdev/keygen-service/internal/models"
	"github.com/m4trixdev/keygen-service/internal/repository"
	"github.com/m4trixdev/keygen-service/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config.Load()

	db, err := gorm.Open(postgres.Open(config.C.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DB] Failed to connect: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Key{}, &models.KeyUsageLog{}); err != nil {
		log.Fatalf("[DB] Migration failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	keyRepo := repository.NewKeyRepository(db)

	authSvc := services.NewAuthService(userRepo)
	keySvc := services.NewKeyService(keyRepo)

	authHandler := handlers.NewAuthHandler(authSvc)
	keyHandler := handlers.NewKeyHandler(keySvc)

	if config.C.DatabaseURL != "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.RateLimit(config.C.RateLimitPerMin))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", middleware.Authenticate(), middleware.RequireAdmin(), authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		keys := v1.Group("/keys", middleware.Authenticate())
		{
			keys.GET("", middleware.RequireAdmin(), keyHandler.List)
			keys.POST("", middleware.RequireAdmin(), keyHandler.Generate)
			keys.POST("/validate", keyHandler.Validate)
			keys.DELETE("/:id", middleware.RequireAdmin(), keyHandler.Revoke)
		}
	}

	log.Printf("[Server] Listening on :%s", config.C.Port)
	if err := r.Run(":" + config.C.Port); err != nil {
		log.Fatalf("[Server] Failed to start: %v", err)
	}
}
