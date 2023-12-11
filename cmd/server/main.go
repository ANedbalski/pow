package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"pow/config"
	"pow/internal/app"
	"pow/internal/book"
	"pow/internal/pow"
	"pow/internal/repository"
	"pow/internal/server"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.NewServerConfig("./config", "server.yml")
	if err != nil {
		panic(fmt.Errorf("Unable to load config: %w \n", err))
	}

	bookRepo := repository.NewInMemoBookRepo()
	powRepo := repository.NewInMemoPOWRepo()

	alg := pow.New(cfg.POWDifficulty, time.Now)
	service := book.New(bookRepo)
	application := app.NewFactory(alg, service, powRepo, cfg.POWTTL)
	srv := server.NewTCP(cfg.Addr, logger, application)

	err = srv.Start(ctx)
	if err != nil {
		panic(fmt.Errorf("Unable to start server: %w \n", err))
	}

	go func() {
		graceful := true
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		for {
			sig := <-sig

			if !graceful {
				logger.Info("caught another exit signal, now hard dying", "signal", sig)
				os.Exit(1)
			}
			graceful = false

			go func() {
				logger.Info("starting graceful shutdown", "signal", sig)
				if err := srv.Stop(); err != nil {
					logger.Error("failed to stop server gracefully", "error", err)
				}
			}()
		}
	}()

	srv.Running()
}
