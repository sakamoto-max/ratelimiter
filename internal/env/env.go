package env

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv(filePath string) {
	err := godotenv.Load(filePath)
	if err != nil {
		log.Fatalf("failed to load env file : %v", err)
	}
}
