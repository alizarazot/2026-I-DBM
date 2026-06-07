package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/alizarazot/2026-i-dbm/internal/constants"
	"github.com/alizarazot/2026-i-dbm/internal/database"
	"github.com/alizarazot/2026-i-dbm/internal/server"

	"charm.land/log/v2"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	if err := run(os.Getenv, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "Server exited with error:", err)
		os.Exit(1)
	}
}

func run(getenv func(string) string, stderr io.Writer) error {
	ctx := context.Background()

	loggerHandler := log.New(stderr)
	logger := slog.New(loggerHandler)

	mongoClient, err := mongo.Connect(mongoOptions.Client().ApplyURI(getenv(constants.ENV_MONGODB_URI)))
	if err != nil {
		return fmt.Errorf("unable to connect to MongoDB: %w", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	mongoDatabase := getenv(constants.ENV_MONGODB_DATABASE)

	var pingResult bson.M
	pingCtx, pingCtxCancel := context.WithTimeout(ctx, time.Second)
	defer pingCtxCancel()
	if err := mongoClient.Database(mongoDatabase).RunCommand(pingCtx, bson.D{{Key: "ping", Value: 1}}).Decode(&pingResult); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}
	logger.Info("Successful ping to database", "pingResult", pingResult)

	server := server.NewServer(
		logger,
		[]byte(getenv(constants.ENV_JWT_SECRET)),
		database.NewAuthStore(mongoClient, mongoDatabase, constants.DB_COLLECTION_USER))

	addr := getenv(constants.ENV_ADDR)
	if addr == "" {
		return fmt.Errorf("an address for http needs to be specified on %q environment variable", constants.ENV_ADDR)
	}

	if err := server.Start(addr); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
