package main

import (
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/api/handler"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/config"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/repository"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/usecase"
	"fmt"
	"log"
	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	c, err := config.Read()
	if err != nil {
		log.Println("failed to read config:", err.Error())
		return
	}

	dbPool, err := repository.Connect(c)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		if dbPool != nil {
			dbPool.Close()
		}
	}()
	repo := repository.NewStorageMessage(dbPool)
	service := usecase.NewMessageService(&repo)
	handlers := handler.NewMessageHandler(service)
	router := handlers.RegisterHandlers()
	slog.Info("starting listening port: ")
	err = handler.Serve(c, router)

	if err != nil {
		slog.Error("running server", err)
	}
}
