package main

import (
	"VK-test/infrastructure/http"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type App struct {
	HTTPHandlers *http.HandlerContainer
	Logger       *zap.Logger
}

func NewApp(container *http.HandlerContainer, logger *zap.Logger) *App {
	return &App{
		HTTPHandlers: container,
		Logger:       logger,
	}
}

func main() {
	app, err := InitializeApp()
	if err != nil {
		app.Logger.Error("error initializing app", zap.Error(err))
		return
	}
	http.NewHTTPServer(*app.HTTPHandlers, fiber.Config{}, app.Logger)
}
