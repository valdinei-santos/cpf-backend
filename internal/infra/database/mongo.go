package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Constantes para as variáveis de ambiente
const (
	MONGO_HOST = "MONGO_HOST"
	MONGO_PORT = "MONGO_PORT"
	MONGO_USER = "MONGO_USER"
	MONGO_PASS = "MONGO_PASS"
)

// Conexao é a variável exportada que qualquer pacote pode usar para interagir com o DB.
//var Conexao *mongo.Client

// ConnectDB estabelece a conexão com o MongoDB.
func ConnectDB(log logger.ILogger) (*mongo.Client, error) {
	// Cria o Connection URI (string de conexão)
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/admin?authSource=admin",
		getEnv(MONGO_USER, "useradmin"),
		getEnv(MONGO_PASS, "useradmin"),
		getEnv(MONGO_HOST, "localhost"),
		getEnv(MONGO_PORT, "27017"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second)

	// Conecta ao MongoDB
	var err error
	Client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.ErrorContext(ctx, "Erro ao conectar ao MongoDB: %v", err)
	}

	// Tentativa de ping para confirmar a conexão imediatamente
	if err = Client.Ping(ctx, nil); err != nil {
		Client.Disconnect(context.Background())
		return nil, fmt.Errorf("falha ao pingar o MongoDB: %w", err)
	}

	log.Info("Conexão com MongoDB estabelecida com sucesso!")
	return Client, nil
}

// DisconnectDB fecha a conexão com o MongoDB.
func DisconnectDB(log logger.ILogger, client *mongo.Client, ctx context.Context) {
	if client != nil {
		if err := client.Disconnect(ctx); err != nil {
			log.ErrorContext(ctx, "Erro ao desconectar do MongoDB: %v", err)
		}
		log.Info("Conexão com MongoDB fechada.")
	}
}

// getEnv retorna a variável de ambiente ou um valor padrão.
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
