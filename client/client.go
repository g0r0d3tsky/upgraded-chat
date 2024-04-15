package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
)

type Message struct {
	UserNickname string
	Content      string
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your nickname: ")
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8085/chat", nil)
	if err != nil {
		log.Fatal("Unable to connect to WebSocket server:", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Unable to close connection", err)
		}
	}(conn)

	done := make(chan struct{})
	go readMessages(conn, done)

	go writeMessages(conn, reader, nickname)

	<-interrupt
	fmt.Println("Interrupt signal received. Exiting...")
}

func readMessages(conn *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		var receivedMessage Message
		err = json.Unmarshal(messageBytes, &receivedMessage)
		if err != nil {
			log.Println("Error decoding message:", err)
			continue
		}

		fmt.Println(receivedMessage.UserNickname + ": " + receivedMessage.Content)
	}
}

func writeMessages(conn *websocket.Conn, reader *bufio.Reader, nickname string) {
	for {
		content, _ := reader.ReadString('\n')
		content = strings.TrimSpace(content)

		message := Message{
			UserNickname: nickname,
			Content:      content,
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Println("Error encoding message:", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}
}
