package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config - ...
type Config struct {
	Port string
	//ArqLog string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Erro ao carregar .env local.")
	}

	config := &Config{
		Port: os.Getenv("PORT"),
		//ArqLog: os.Getenv("ARQ_LOG"),
	}

	// Validação
	if config.Port == "" {
		return nil, fmt.Errorf("PORT não encontrada no .env")
	}
	//if config.ArqLog == "" {
	//	return nil, fmt.Errorf("ARQ_LOG não encontrada no .env")
	//}

	return config, nil
}
