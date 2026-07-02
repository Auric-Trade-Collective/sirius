package config

type Config struct {
	Isolated map[string]Entry `toml:"isolated"`
	Host map[string]Entry `toml:"host"`
}

type Entry struct {
	Name string `toml:"name"`
	Args string `toml:"args"`
}
