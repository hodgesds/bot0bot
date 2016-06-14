package plugins

import (
	"fmt"
	"github.com/hodgesds/bot0bot"
	"math/rand"
	"strings"
)

func init() {
	bot0bot.RegisterPlugin("iq", &Iq{})
}

type Iq struct {
}

func (iq *Iq) Handle(msg *bot0bot.Message) string {
	n := rand.Intn(200)

	if msg.User == msg.Bot {
		return ""
	}

	target := ""

	chunks := strings.Split(msg.Content, "PRIVMSG")
	if len(chunks) > 0 {
		for _, c := range chunks[1:] {
			if strings.Contains(c, "@") {
				i := strings.IndexRune(c, '@')
				if i >= 0 {
					target = strings.Split(c[i:], " ")[0]
					break
				}
			}
		}
	}

	if len(target) > 0 {
		return fmt.Sprintf("bot0bot says %s iq is %d", target, n)

	}

	return fmt.Sprintf("bot0bot says @%s iq is %d", msg.User, n)
}

func (iq *Iq) Match(msg *bot0bot.Message) bool {
	return strings.Contains(msg.Content, ":!iq")
}
