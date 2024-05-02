package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/config"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/handlers"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/kafka"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/repository"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/usecase"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("loading envs", err)
	}

	c, err := config.Read()
	if err != nil {
		slog.Error("reading config:", err)
		return
	}

	dbPool, err := repository.Connect(c)
	if err != nil {
		slog.Error("database connect", err)
	}

	defer func() {
		if dbPool != nil {
			dbPool.Close()
		}
	}()

	repo := repository.NewStorageMessage(dbPool)
	service := usecase.NewService(&repo)
	handler := handlers.NewServiceMessage(service)

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}

	wg.Add(1)
	consumer, err := kafka.RunConsumer(ctx, wg, c, handler.Handler())
	if err != nil {
		slog.Error("kafka consumer", err)
	}
	slog.Info("kafka consumer started")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	keepRunning := true
	for keepRunning {
		select {
		case <-ctx.Done():
			slog.Info("terminate: context done")
			keepRunning = false
		case <-ch:
			slog.Info("terminate: signal")
			keepRunning = false
		}
	}

	cancel()

	wg.Wait()

	if err = consumer.Close(); err != nil {
		slog.Error("closing consumer", err)
		return
	}
}
