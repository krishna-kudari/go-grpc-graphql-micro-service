package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/krishna-kudari/go-grpc-graphql-micro-service/order"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
}

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		logger.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	var r order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewOrderRepository(config.DatabaseURL)
		if err != nil {
			logger.Warn("failed to connect to database, retrying",
				slog.String("error", err.Error()))
			return err
		}
		logger.Info("connected to database")
		return nil
	})

	defer func() {
		if err := r.Close(); err != nil {
			logger.Error("failed to close repository", slog.String("error", err.Error()))
		}
	}()

	service, err := order.NewOrderService(r)
	if err != nil {
		logger.Error("failed to create order service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("starting gRPC server", slog.Int("port", 8080))
		if err := order.ListenGRPC(ctx, service, config.AccountURL, config.CatalogURL, 8080, logger); err != nil {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Info("received signal, shutting down", slog.String("signal", sig.String()))
	case err := <-errCh:
		logger.Error("server error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
