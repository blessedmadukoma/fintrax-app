package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver    string `mapstructure:"DB_DRIVER"`
	DBSource    string `mapstructure:"DB_SOURCE"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_NAME     string `mapstructure:"DB_NAME"`

	DB_SOURCE_LIVE string `mapstructure:"DB_SOURCE_LIVE"`
}

func LoadEnvConfig(path string) (config Config) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Cannot load env:", err)
	}

	config.DBDriver = os.Getenv("DB_DRIVER")
	config.DB_USER = os.Getenv("DB_USER")
	config.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	config.DB_HOST = os.Getenv("DB_HOST")
	config.DB_PORT = os.Getenv("DB_PORT")
	config.DB_NAME = os.Getenv("DB_NAME")

	config.DB_SOURCE_LIVE = os.Getenv("DB_SOURCE_LIVE")

	return config
}
