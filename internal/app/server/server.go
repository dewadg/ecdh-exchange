package server

import (
	"os"

	"github.com/dewadg/ecdh-exchange/internal/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Run(cfg *config.Config) {
	app := fiber.New()
	app.Post("/exchange", handleKeyExchange())
	app.Get("/exchange/:sessionID", handleGetSecret())

	addr := "127.0.0.1:8000"
	if os.Getenv("APP_ENV") == "production" {
		addr = "0.0.0.0:8000"
	}

	if err := app.Listen(addr); err != nil {
		logrus.WithError(err).Fatal("http server error")
	}
}
