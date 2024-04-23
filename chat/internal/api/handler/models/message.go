package models

type Message struct {
	UserNickname string `json:"user_nickname"`
	Content      string `json:"content"`
}
