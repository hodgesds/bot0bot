package bot0bot

import (
	"time"
)

type Message struct {
	Content   string `json:"content,`
	Timestamp string `json:"timestamp,`
	User      string `json:"user,"`
	Bot       string `json:"bot,"`
}

func NewMessage(user, bot, content string) *Message {
	return &Message{
		Bot:       bot,
		Content:   content,
		Timestamp: time.Now().String(),
		User:      user,
	}
}
