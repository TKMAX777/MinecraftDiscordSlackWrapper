package main

// User is a format of a user profile
type User struct {
	Name      string   `yaml:"Name"`
	Discord   string   `yaml:"Discord"`
	Slack     string   `yaml:"Slack"`
	Groups    []string `yaml:"Groups"`
	IsDefault bool     `yaml:"IsDefault"`
}
