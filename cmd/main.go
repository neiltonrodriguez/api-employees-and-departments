package main

import (
	"api-employees-and-departaments/internal/config"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // optional
	_, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
}
