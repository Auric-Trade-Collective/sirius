package config

import "errors"

type Config struct {
	Isolated map[string]Entry `toml:"isolated"`
	Host map[string]*Entry `toml:"host"`
}

type Entry struct {
	Name string `toml:"name"`
	Args string `toml:"args"`
	NeedsTTY *bool `toml:"needs_tty"`
	NeedsDev []string `toml:"needs_dev"`
	NeedsDep []string `toml:"needs_dep"`
	OnExit *string  `toml:"on_exit"`
}

func (c Config) FindEntryByName(name string) (*Entry, error){
	if _, exists := c.Host[name]; exists {
		return c.Host[name], nil
	}

	return nil, errors.New("")
}
