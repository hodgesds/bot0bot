# bot0bot
bot0bot is a twitch bot, it supports various [plugins](plugins/README.md). It
runs via command line or windows command prompt. To compile bot0bot with extra
plugins you need to clone/fork the repo and import any plugins in
[main.go](bot0bot/main.go#L11).

```sh
NAME:
   bot0bot

USAGE:
   bot0bot [global options] command [command options] [arguments...]
   
VERSION:
   0.0.1
   
COMMANDS:
GLOBAL OPTIONS:
   --log STDOUT, -l STDOUT  log file STDOUT (default: "-") [$BOT0BOT_LOG]
   --log-level value        log level (debug,info,warn,error,fatal) (default: "info")
   --log-format value       log level (json,text) (default: "json")
   --channel value          Channel should be #<twitch user>
   --user value             twitch account name
   --name value             bot name, should be same as user
   --host value             IRC Host (default: "irc.twitch.tv")
   --password value         Oauth password see http://www.twitchapps.com/tmi/
   --port value, -p value   tcp port (default: 6667)
   --plugins value          bot0bot plugins (default: "hello", "iq")
   --help, -h               show help
   --version, -v            print the version
   
```

# FAQ
Why doesn't my twitch password doesn't work?
* Try using a twitch [oauth password](http://www.twitchapps.com/tmi/)

Do I have to compile it to run?
* Nope, download binaries from [releases](https://github.com/hodgesds/bot0bot/releases)

Can I run this on Windows?
* Yep, see [releases](https://github.com/hodgesds/bot0bot/releases)
