package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/j-gc/plantpal-backend/internal/modules/auth/application"
	usershttp "github.com/j-gc/plantpal-backend/internal/modules/authinfrastructure/http"
	"github.com/j-gc/plantpal-backend/internal/modules/authinfrastructure/persistence"
	"github.com/j-gc/plantpal-backend/internal/modules/authinfrastructure/security"
	"github.com/j-gc/plantpal-backend/internal/shared/config"
	"github.com/j-gc/plantpal-backend/internal/shared/jwt"
	"github.com/j-gc/plantpal-backend/internal/shared/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	log := logger.New(cfg.LogLevel)

	// Initialize database
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Verify database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}
	log.Info("Database connected successfully")

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize dependencies
	userRepo := persistence.NewUserRepository(db)
	hasher := security.NewBcryptHasher()
	tokenIssuer := jwt.NewHS256Issuer(cfg.JWTSecret, "plantpal-backend")
	userService := application.NewService(userRepo, hasher, tokenIssuer)

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Setup API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		authHandlers := usershttp.NewAuthHandlers(userService)
		authHandlers.RegisterRoutes(api.Group("/auth"))
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Info("Starting server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("Server exited")
}
