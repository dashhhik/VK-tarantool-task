package http

import (
	"VK-test/infrastructure/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type AuthHandler interface {
	Login(ctx *fiber.Ctx) error
}

type ReadHandler interface {
	Read(ctx *fiber.Ctx) error
}

type WriteHandler interface {
	Write(ctx *fiber.Ctx) error
}

type HandlerContainer struct {
	Auth  AuthHandler
	Read  ReadHandler
	Write WriteHandler
}

func NewHandlerContainer(authHandler AuthHandler, writeHandler WriteHandler, readHandler ReadHandler) *HandlerContainer {
	return &HandlerContainer{
		Auth:  authHandler,
		Write: writeHandler,
		Read:  readHandler,
	}
}

func registerRoutes(container HandlerContainer, app *fiber.App) {
	app.Post("/api/login", container.Auth.Login)
	app.Post("/api/write", container.Write.Write, middleware.Protected())
	app.Post("/api/read", container.Read.Read, middleware.Protected())
}

func NewHTTPServer(handlerContainer HandlerContainer, cfg fiber.Config, logger *zap.Logger) {
	app := fiber.New(cfg)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, User-Agent",
		AllowMethods: "POST",
	}))

	registerRoutes(handlerContainer, app)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + strconv.Itoa(8000)); err != nil {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()
	<-signalChannel

	_ = app.Shutdown()
}
