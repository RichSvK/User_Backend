package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv(filename string) {
	if err := godotenv.Load(filename); err != nil {
		log.Printf("[ERROR] error: %v", err)
	}
}
