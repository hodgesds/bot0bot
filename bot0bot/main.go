package main

import (
	// import bot0bot before importing plugins
	"github.com/hodgesds/bot0bot"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"

	// XXX: Register plugins here by importing:
	_ "github.com/hodgesds/bot0bot/plugins"
)

func main() {
	app := cli.NewApp()
	app.Action = runBot

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "log, l",
			Value:  "-",
			Usage:  "log file `STDOUT`",
			EnvVar: "BOT0BOT_LOG",
		},
		cli.StringFlag{
			Name:  "log-level",
			Value: "info",
			Usage: "log level (debug,info,warn,error,fatal)",
		},
		cli.StringFlag{
			Name:  "log-format",
			Value: "json",
			Usage: "log level (json,text)",
		},
		cli.StringFlag{
			Name:  "channel",
			Value: "",
			Usage: "Channel should be #<twitch user>",
		},
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "twitch account name",
		},
		cli.StringFlag{
			Name:  "name",
			Value: "",
			Usage: "bot name, should be same as user",
		},
		cli.StringFlag{
			Name:  "host",
			Value: "irc.twitch.tv",
			Usage: "IRC Host",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "Oauth password see http://www.twitchapps.com/tmi/",
		},
		cli.IntFlag{
			Name: "port,p",
			//Value: 6697,
			//Value: 443,
			Value: 6667,
			Usage: "tcp port",
		},
		cli.StringSliceFlag{
			Name:  "plugins",
			Value: &cli.StringSlice{"hello", "iq"},
			Usage: "Breadbot plugins",
		},
	}

	app.Run(os.Args)
}

func runBot(c *cli.Context) error {
	setupLogging(c)
	bot := bot0bot.NewBot(
		c.String("channel"),
		c.String("user"),
		c.String("name"),
		c.String("password"),
		c.String("host"),
		c.Int("port"),
		c.StringSlice("plugins"),
	)
	bot.Start()
	defer bot.Stop()

	return nil
}

func setupLogging(c *cli.Context) {
	bot0bot.Log = logrus.New()

	if c.String("log") == "-" {
		bot0bot.Log.Out = os.Stdout
	}

	switch c.String("log-level") {
	case "debug":
		bot0bot.Log.Level = logrus.DebugLevel

	case "info":
		bot0bot.Log.Level = logrus.InfoLevel

	case "warn":
		bot0bot.Log.Level = logrus.WarnLevel

	case "error":
		bot0bot.Log.Level = logrus.ErrorLevel

	case "fatal":
		bot0bot.Log.Level = logrus.FatalLevel

	default:
		bot0bot.Log.Level = logrus.InfoLevel
	}

	switch c.String("log-format") {
	case "text":
		bot0bot.Log.Formatter = &logrus.TextFormatter{}

	case "json":
		bot0bot.Log.Formatter = &logrus.JSONFormatter{}

	default:
		bot0bot.Log.Formatter = &logrus.JSONFormatter{}
	}
}
