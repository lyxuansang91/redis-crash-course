package config

type Config struct {
	Protocol string
	Port string
	MaxConnections int
}


const (
	Protocol = "tcp"
	Port = ":3000"
	MaxConnections = 20000
)

var defaultConfig = &Config{
	Protocol: Protocol,
	Port: Port,
	MaxConnections: MaxConnections,
}


func NewConfig() *Config {
	return defaultConfig
}