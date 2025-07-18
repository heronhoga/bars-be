package utils

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

func LoadEnv() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found or error loading it:", err)
		} else {
			log.Println(".env file successfully loaded")
		}
	})
}
