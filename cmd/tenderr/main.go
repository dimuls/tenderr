package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"tednerr"
	tenderrClickhouse "tednerr/clickhouse"
	"tednerr/postgres"
)

const defaultConfigPath = "config.yaml"

func main() {
	level := zap.NewAtomicLevel()

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		level,
	)

	logger := zap.New(core)

	err := run(logger, level)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func run(logger *zap.Logger, logLevel zap.AtomicLevel) error {

	var configPath string

	flag.StringVar(&configPath, "config", defaultConfigPath, "path to config")

	flag.Parse()

	var config tenderr.Config

	err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("config load: %w", err)
	}

	switch config.LogLevel {
	case "debug":
		logLevel.SetLevel(zap.DebugLevel)
	case "info":
		logLevel.SetLevel(zap.InfoLevel)
	case "warn":
		logLevel.SetLevel(zap.WarnLevel)
	case "error":
		logLevel.SetLevel(zap.ErrorLevel)
	default:
		logLevel.SetLevel(zap.InfoLevel)
	}

	db, err := sqlx.Open("postgres", config.PostgresURL)
	if err != nil {
		return fmt.Errorf("sqlx open: %w", err)
	}

	defer db.Close()

	storage := &postgres.Storage{
		DB: db,
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.Clickhouse.Addr},
		Auth: clickhouse.Auth{
			Database: config.Clickhouse.Database,
			Username: config.Clickhouse.Username,
			Password: config.Clickhouse.Password,
		},
	})
	if err != nil {
		return fmt.Errorf("clickhouse open: %w", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("clickhouse conn ping: %w", err)
	}

	err = storage.Migrate()
	if err != nil {
		return fmt.Errorf("postgres storage migrate: %w", err)
	}

	server := &tenderr.Server{
		Addr:       config.Addr,
		Storage:    storage,
		LogStorage: &tenderrClickhouse.LogStorage{Conn: conn},
		CORS:       config.CORS,
		Logger:     logger,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = server.ListenAndServe(ctx)
	if err != nil {
		return fmt.Errorf("server listen and serve: %w", err)
	}

	server.Logger.Info("service stopped, exiting")

	return nil
}
