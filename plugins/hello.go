package plugins

import (
	"github.com/hodgesds/bot0bot"
	"math/rand"
	"strings"
)

func init() {
	bot0bot.RegisterPlugin("hello", &Hello{})
}

type Hello struct {
}

func (h *Hello) Handle(msg *bot0bot.Message) string {
	user := "@" + msg.User
	n := rand.Intn(200)
	if (n >= 0) && (n <= 10) {
		return "Welcome " + user + "!!!"
	}
	if (n > 10) && (n <= 20) {
		return "Howdy " + user + "."
	}
	if (n > 20) && (n <= 30) {
		return "Aren't I much better than Nightbot " + user + "?"
	}
	if (n > 30) && (n <= 40) {
		return "I'm still learning things around here."
	}
	if (n > 40) && (n <= 50) {
		return "Do you speak 01100010 01101111 01110100?"
	}
	if (n > 50) && (n <= 60) {
		return "Ask me later."
	}
	if (n > 60) && (n <= 75) {
		return "meh..."
	}
	if n > 75 {
		return ""
	}
	return ""
}

func (h *Hello) Match(msg *bot0bot.Message) bool {
	if msg.User == msg.Bot {
		return false
	}

	chunks := strings.Split(msg.Content, "PRIVMSG")
	if len(chunks) > 0 {
		for _, c := range chunks[1:] {
			if strings.Contains(c, "@"+msg.Bot) {
				return true
			}
		}
	}

	return false
}
