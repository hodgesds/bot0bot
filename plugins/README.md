# Plugins
A Plugin must implement the `Plugin` interface:

```go
type Plugin interface {
	Handle(*Message) string
	Match(*Message) bool
}
```

In addition a plugin should register itself with bot0bot:

```go
import "github.com/hodgesds/bot0bot"

func init() {
	bot0bot.RegisterPlugin("myPlugin", MyInterface{})
}
```
