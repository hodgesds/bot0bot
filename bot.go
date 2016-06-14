package bot0bot

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net"
	"net/textproto"
	"strings"
	"time"
)

type Bot struct {
	Channel  string `json:"channel,"`
	Host     string `json:"host,"`
	Name     string `json:"name,"`
	Password string `json:"-,"` // i assume this is what you want...
	Port     int    `json:"port,"`
	User     string `json:"user,"`

	plugins map[string]Plugin

	room map[string][]string

	conn net.Conn

	lastMsg int64
	msgChan chan string

	// channels
	stopChan chan bool
}

func NewBot(
	channel, user, name, password, host string,
	port int,
	pluginNames []string,
) *Bot {
	// setup plugins
	plugins := map[string]Plugin{}
	for _, name := range pluginNames {
		plugin, err := getPlugin(name)
		if err != nil {
			Log.Errorf("unknown plugin %s", name)
			continue
		}

		plugins[name] = plugin
	}

	if password == "" {
		fmt.Printf("Enter password for %s@%s:%d\n", user, host, port)
		p, err := terminal.ReadPassword(0)
		if err != nil {
			Log.Panic(err)
		}
		password = string(p)
	}

	return &Bot{
		Channel:  channel,
		Host:     host,
		Name:     name,
		Password: password,
		Port:     port,
		User:     user,

		room: map[string][]string{},

		plugins: plugins,

		msgChan:  make(chan string),
		stopChan: make(chan bool),
	}
}

func (b *Bot) Start() {
	var err error
	if b.Port == 6697 || b.Port == 443 {
		b.conn, err = tls.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", b.Host, b.Port),
			&tls.Config{},
		)
		if err != nil {
			Log.Panic(err)
			return
		}
	} else {
		b.conn, err = net.Dial("tcp",
			fmt.Sprintf("%s:%d", b.Host, b.Port),
		)
		if err != nil {
			Log.Panic(err)
			return
		}
	}

	fmt.Fprintf(b.conn, "USER %s 8 * :%s\r\n", b.User, b.User)
	fmt.Fprintf(b.conn, "PASS %s\r\n", b.Password)
	fmt.Fprintf(b.conn, "NICK %s\r\n", b.User)
	fmt.Fprintf(b.conn, "JOIN %s\r\n", b.Channel)

	reader := bufio.NewReader(b.conn)
	tp := textproto.NewReader(reader)

	for {
		select {
		case <-b.stopChan:
			err := b.conn.Close()
			if err != nil {
				Log.Error(err)
			}
			return

		case msg := <-b.msgChan:
			b.message(msg)

		default:
			msg, err := tp.ReadLine()
			if err != nil {
				Log.Error(err)
			}

			b.onReadline(msg)
		}
	}
}

func (b *Bot) onReadline(msg string) {
	Log.Infoln(msg)
	if strings.Contains(msg, "PING") {
		b.onPing(msg)
	} else if strings.Contains(msg, ".tmi.twitch.tv PRIVMSG "+b.Channel) {
		b.onMessage(msg)
	} else if strings.Contains(msg, ".tmi.twitch.tv JOIN "+b.Channel) {
		b.onJoin(msg)
	} else if strings.Contains(msg, ".tmi.twitch.tv PART "+b.Channel) {
		b.onPart(msg)
	} else if strings.Contains(msg, ":jtv MODE "+b.Channel+" +o ") {
		b.onMod(msg)
	} else if strings.Contains(msg, ":jtv MODE "+b.Channel+" -o ") {
	}
}

func (b *Bot) onMessage(msg string) {
	userdata := strings.Split(msg, ".tmi.twitch.tv PRIVMSG "+b.Channel)
	username := strings.Split(userdata[0], "@")
	//usermessage := strings.Replace(userdata[1], " :", "", 1)

	m := NewMessage(username[1], b.Name, msg)

	for _, plugin := range b.plugins {
		if plugin.Match(m) {
			response := plugin.Handle(m)
			if len(response) > 0 {
				b.message(response)
			}
		}
	}
	//if b.userLastMsg[username[1]]+int64(ircbot.userMaxLastMsg) >= time.Now().Unix() {
	//	b.timeout(username[1], "spam")
	//}
	//b.userLastMsg[username[1]] = time.Now().Unix()
	//go b.CmdInterpreter(username[1], usermessage)
}

func (b *Bot) onUnMod(msg string) {
	usermod := strings.Split(msg, ":jtv MODE "+b.Channel+" -o ")
	username := usermod[1]
	Log.Infof("%s is no longer a mod!\n", username)
}

func (b *Bot) onMod(msg string) {
	usermod := strings.Split(msg, ":jtv MODE "+b.Channel+" +o ")
	username := usermod[1]
	Log.Infof("%s is now a mod!\n", username)
}

func (b *Bot) onPing(msg string) {
	Log.Info("Got ping from server...")
	pongdata := strings.Split(msg, "PING ")
	fmt.Fprintf(b.conn, "PONG %s\r\n", pongdata[1])
}

func (b *Bot) onJoin(msg string) {
	userjoindata := strings.Split(msg, ".tmi.twitch.tv JOIN "+b.Channel)
	userjoined := strings.Split(userjoindata[0], "@")
	username := userjoined[1]

	if username == b.User || username == b.Name {
		Log.Info("Not saying hello to myself")
		return
	}

	Log.Infof("%s has joined!\n", username)

	_, ok := b.room[username]
	if !ok {
		Log.Infof("Adding %s to room", username)
		b.room[username] = []string{}
		fmt.Fprintf(b.conn, "PRIVMSG "+b.Channel+" :Welcome "+username+"!\r\n")
	}

}

func (b *Bot) onPart(msg string) {
	userjoindata := strings.Split(msg, ".tmi.twitch.tv PART "+b.Channel)
	userjoined := strings.Split(userjoindata[0], "@")
	username := userjoined[1]
	Log.Infof("%s has left!\n", username)
}

func (b *Bot) message(msg string) {
	now := time.Now().Unix()
	if now-b.lastMsg < 3 {
		Log.Debugf("Requeuing message: %s", msg)
		go func() { b.msgChan <- msg }()
		return
	}

	fmt.Fprintf(b.conn, "PRIVMSG "+b.Channel+" :"+msg+"\r\n")

	b.lastMsg = now
}

func (b *Bot) Stop() {
	Log.Info("Stopping bot")
	close(b.stopChan)
}
