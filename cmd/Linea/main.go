package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/alizarazot/2026-i-dbm/internal/server"
	"github.com/charmbracelet/log"
)

const (
	ENV_BASE = "LINEA_"
	ENV_ADDR = ENV_BASE + "ADDR"
)

func main() {
	if err := run(os.Getenv, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "Server exited with error:", err)
		os.Exit(1)
	}
}

func run(getenv func(string) string, stderr io.Writer) error {
	loggerHandler := log.New(stderr)
	logger := slog.New(loggerHandler)

	server := server.NewServer(logger)

	addr := getenv(ENV_ADDR)
	if addr == "" {
		return fmt.Errorf("an address for http needs to be specified on %q environment variable", ENV_ADDR)
	}

	if err := server.Start(addr); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
