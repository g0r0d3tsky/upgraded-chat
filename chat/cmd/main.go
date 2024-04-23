package main

import (
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/api/handler"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/cache"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/config"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/kafka"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/repository"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/usecase"
	"log/slog"

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

	producer, err := kafka.New(c)

	repo := repository.NewStorageMessage(dbPool)
	redis, err := cache.Connect(c)
	if err != nil {
		slog.Error("redis connect", err)
	}
	cacheService := cache.NewCache(redis, &repo)
	service := usecase.NewMessageService(producer, cacheService, c.KafkaTopic)
	handlers := handler.NewMessageHandler(service, c.AmountLastMessages)
	router := handlers.RegisterHandlers()
	slog.Info("starting listening port: ")
	err = handler.Serve(c, router)

	if err != nil {
		slog.Error("running server", err)
	}
}
