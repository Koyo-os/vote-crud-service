package config

type Config struct{
	Port int
	Host string
}

func NewConfig() *Config {
	return &Config{
		Port: 8051,
		Host: "localhost",
	}
}