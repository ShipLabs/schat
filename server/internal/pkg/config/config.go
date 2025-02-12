package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var Configs *Config

type Database struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Name     string `env:"NAME,required"`
}

type Config struct {
	Port string   `env:"PORT,required"`
	DB   Database `env:"" envPrefix:"DB_"`
}

func Load() {
	if loadErr := godotenv.Load(); loadErr != nil {
		fmt.Println("Error loading config file: ", loadErr)
		panic(loadErr)
	}

	config := Config{}

	if parseErr := env.Parse(&config); parseErr != nil {
		fmt.Println("Error parsing config: ", parseErr)
		panic(parseErr)
	}

	Configs = &config
}
