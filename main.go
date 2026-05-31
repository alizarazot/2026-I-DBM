package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	handler := log.New(os.Stderr)
	logger := slog.New(handler)

	e := echo.NewWithConfig(echo.Config{Logger: logger})

	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	addr := os.Getenv("LINEA_ADDR")
	if addr == "" {
		logger.Error("an address for http needs to be specified on 'LINEA_ADDR' environment variable")
		os.Exit(2)
	}
	if err := e.Start(addr); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
