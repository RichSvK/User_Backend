package config

import (
	"github.com/joho/godotenv"
)

func LoadEnv(filename string) {
	_ = godotenv.Load(filename)
}