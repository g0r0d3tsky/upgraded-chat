package handler

import (
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[*websocket.Conn]struct{})
)

// TODO: куда выносить имя топика?
const topic = "general-chat"

type mes struct {
	UserNickname string
	Content      string
}

type MessageService interface {
	Push(topic string, message *domain.Message) error
	GetMessages(ctx context.Context, amount int) ([]*domain.Message, error)
}

type MessageHandler struct {
	service MessageService
}

func NewMessageHandler(service MessageService) *MessageHandler {
	return &MessageHandler{
		service: service,
	}
}

func (h *MessageHandler) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer func(connection *websocket.Conn) {
		err := connection.Close()
		if err != nil {
			slog.Error("connection closing", err)
		}
	}(connection)
	h.sendLastMessages(context.TODO(), connection)

	clients[connection] = struct{}{}
	defer delete(clients, connection)

	for {
		mt, messageBytes, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break
		}

		var message mes
		err = json.Unmarshal(messageBytes, &message)
		if err != nil {
			slog.Error("error decoding message:", err)
			continue
		}

		var mess domain.Message
		mess.UserNickname = message.UserNickname
		mess.Content = message.Content

		err = h.service.Push(topic, &mess)
		if err != nil {
			slog.Error("saving message:", err)
			return
		}

		go writeMessage(message)

		go messageHandler(message)
	}
}

func writeMessage(message mes) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("Error encoding message:", err)
		return
	}

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			return
		}
	}
}

func messageHandler(message mes) {
	fmt.Printf("%v : %v \n", message.UserNickname, message.Content)
}

func (h *MessageHandler) sendLastMessages(ctx context.Context, connection *websocket.Conn) {
	//TODO: здесь константу 10 тоже хотелось бы куда-то вынести?
	messages, err := h.service.GetMessages(ctx, 10)
	if err != nil {
		slog.Error("getting messages: ", err)
		return
	}
	for _, msg := range messages {
		messageBytes, err := json.Marshal(msg)
		if err != nil {
			slog.Error("encoding message:", err)
			continue
		}

		err = connection.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			slog.Error("sending message:", err)
			continue
		}
	}
}

func (h *MessageHandler) RegisterHandlers() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/chat", h.echo).Methods("GET")
	return router
}
