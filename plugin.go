package bot0bot

import (
	"fmt"
	"sync"
)

var (
	pluginMu sync.RWMutex
	plugins  = make(map[string]Plugin)
)

type Plugin interface {
	Handle(*Message) string
	Match(*Message) bool
}

func RegisterPlugin(name string, plugin Plugin) {
	pluginMu.Lock()
	defer pluginMu.Unlock()
	if plugin == nil {
		panic("Registered plugin is nil")
	}
	if _, dup := plugins[name]; dup {
		panic("Plugin already registered with name " + name)
	}
	plugins[name] = plugin
}

func unregisterAllPlugins() {
	pluginMu.Lock()
	defer pluginMu.Unlock()
	// For tests.
	plugins = make(map[string]Plugin)
}

func hasPlugin(plugin string) bool {
	pluginMu.RLock()
	defer pluginMu.RUnlock()

	for p, _ := range plugins {
		if plugin == p {
			return true
		}
	}

	return false
}

func getPlugin(name string) (Plugin, error) {
	pluginMu.RLock()
	defer pluginMu.RUnlock()

	for pluginName, plugin := range plugins {
		if name == pluginName {
			return plugin, nil
		}
	}

	return nil, fmt.Errorf("plugin %s not registered", name)
}
