package app

import (
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

	// Setup routes
	routes := NewRoutes(a.App, container)
	routes.Setup()

	// Start server
	if err := a.Listen(container.GetServerAddress()); err != nil {
		container.Dependencies.Logger.LogPanic(err.Error())
	}
}
