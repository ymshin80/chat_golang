package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

//환경변수 관리용
type Config struct {
	DB struct {
		Database string
		URL string
	}

	Kafka struct {
		URL string
		ClientID string
	}
}


func NewConfig(path string) *Config {
	c := new(Config)

	if f, err := os.Open(path); err != nil {
		panic(err)
	
	} else if err := toml.NewDecoder(f).Decode(c); err != nil {
		panic(err)
	} else {
		return c
	}
}