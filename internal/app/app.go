package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	*fiber.App
	container *Container
}

func NewApp() *App {
	return &App{
		App: fiber.New(),
	}
}

func (a *App) Start() {
	// Initialize container
	container := NewContainer()
	a.container = container

	// Setup dependencies
	deps := container.Dependencies
	a.setupGracefulShutdown(deps)

	// Setup routes
	routes := NewRoutes(a.App, container)
	routes.Setup()

	// Start server
	deps.Logger.LogInfo(fmt.Sprintf("Starting server on %s", container.GetServerAddress()))
	if err := a.Listen(container.GetServerAddress()); err != nil {
		container.Dependencies.Logger.LogPanic(err.Error())
	}
}

func (a *App) setupGracefulShutdown(deps *Dependencies) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		deps.Logger.LogInfo("Received shutdown signal, shutting down gracefully...")

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Channel to signal completion of cleanup
		done := make(chan bool, 1)

		// Perform cleanup in a separate goroutine
		go func() {
			a.performCleanup(deps)
			done <- true
		}()

		// Wait for cleanup to complete or timeout
		select {
		case <-done:
			deps.Logger.LogInfo("Application shut down successfully")
		case <-ctx.Done():
			deps.Logger.LogError("Shutdown timeout exceeded, forcing exit")
		}

		os.Exit(0)
	}()
}

func (a *App) performCleanup(deps *Dependencies) {
	// 1. Stop accepting new requests first
	deps.Logger.LogInfo("Stopping HTTP server...")
	if err := a.App.ShutdownWithTimeout(15 * time.Second); err != nil {
		deps.Logger.LogError(fmt.Sprintf("failed to shutdown fiber app gracefully: %v", err))
	}

	// 2. Close database connections
	deps.Logger.LogInfo("Closing database connections...")
	if deps.Postgres != nil {
		if err := deps.Postgres.CloseConnection(); err != nil {
			deps.Logger.LogError(fmt.Sprintf("failed to close postgres connection: %v", err))
		} else {
			deps.Logger.LogInfo("PostgreSQL connection closed successfully")
		}
	}

	// 3. Close Redis connection
	deps.Logger.LogInfo("Closing Redis connection...")
	if deps.Redis != nil {
		if err := deps.Redis.Close(); err != nil {
			deps.Logger.LogError(fmt.Sprintf("failed to close redis connection: %v", err))
		} else {
			deps.Logger.LogInfo("Redis connection closed successfully")
		}
	}

	// 4. Close other resources if any (MinIO, etc.)
	if deps.MinioClient != nil {
		deps.Logger.LogInfo("MinIO client cleanup completed")
		// MinIO client usually doesn't need explicit closing
	}

	deps.Logger.LogInfo("All cleanup operations completed")
}
