package framework

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscToken string `json:"disc_token"`
}

func LoadEnv() *Config {
	// Required to have .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	config := &Config{
		DiscToken: os.Getenv("DISCORD_TOKEN"),
	}

	log.Print("Successfully loaded secrets.")
	return config
}
