// Pacote principal da aplicação que inicia a API.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/valdinei-santos/cpf-backend/cmd/api/routes"
	_ "github.com/valdinei-santos/cpf-backend/cmd/api/routes"
	_ "github.com/valdinei-santos/cpf-backend/cmd/api/stats"
	"github.com/valdinei-santos/cpf-backend/internal/infra/config"
	"github.com/valdinei-santos/cpf-backend/internal/infra/database"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"

	"github.com/gin-gonic/gin"
)

// @title CPF Management API
// @version 1.0
// @description Está API gerencia CPFs.
// @host localhost:8888
// @BasePath /api/v1/cpfs
func main() {
	fmt.Println("Iniciando...")
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Erro ao carregar variáveis do .env: %v", err)
		os.Exit(1)
	}

	log := logger.NewSlogILogger()
	fmt.Println("Iniciou Log...")

	// Conecta ao banco de dados
	client, err := database.ConnectDB(log)
	if err != nil {
		fmt.Printf("Erro ao conectar no database: %v", err)
		os.Exit(1)
	}
	db := client.Database("cpf_management")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer database.DisconnectDB(log, client, ctx)
	fmt.Println("Iniciou Database...")

	gin.DefaultWriter = io.Discard // Desabilita o log padrão do gin jogando para o io.Discard
	router := gin.Default()
	router.SetTrustedProxies(nil)
	routes.InitRoutes(&router.RouterGroup, log, db)

	log.Info("start cpf-mamagement", "PORT:", config.Port)
	err = router.Run(":" + config.Port)
	if err != nil {
		fmt.Printf("Erro ao iniciar a API na porta %v: %v", config.Port, err)
		log.Error("Erro ao iniciar a API na porta " + config.Port + " - " + err.Error())
	}
}
