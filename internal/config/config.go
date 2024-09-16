package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerAddress    string `env:"SERVER_ADDRESS"`
	PostgresUsername string `env:"POSTGRES_USERNAME"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresDatabase string `env:"POSTGRES_DATABASE"`
}

func MustLoadConfig() *Config {
	var cnf Config

	msg := "Error while init config file"

	if err := cleanenv.ReadConfig("config/.env", &cnf); err != nil {
		panic(msg + err.Error())
	}

	return &cnf
}
