//go:build integration
// +build integration

package routes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/valdinei-santos/cpf-backend/cmd/api/routes"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
)

// mockLogger implementa ILogger só pra não depender do logger real
type mockLogger struct{}

func (m *mockLogger) Debug(string, ...any)                           {}
func (m *mockLogger) Info(string, ...any)                            {}
func (m *mockLogger) Warn(string, ...any)                            {}
func (m *mockLogger) Error(string, ...any)                           {}
func (m *mockLogger) With(...any) logger.ILogger                     { return m }
func (m *mockLogger) WithContext(ctx context.Context) logger.ILogger { return m }
func (m *mockLogger) DebugContext(context.Context, string, ...any)   {}
func (m *mockLogger) InfoContext(context.Context, string, ...any)    {}
func (m *mockLogger) WarnContext(context.Context, string, ...any)    {}
func (m *mockLogger) ErrorContext(context.Context, string, ...any)   {}

func TestClientePost_Integration(t *testing.T) {
	ctx := context.Background()

	// 🚀 Sobe o container MongoDB
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer mongoC.Terminate(ctx)

	// Descobre a URL de conexão
	host, _ := mongoC.Host(ctx)
	port, _ := mongoC.MappedPort(ctx, "27017")
	mongoURI := "mongodb://" + host + ":" + port.Port()

	// 🔗 Conecta no Mongo real
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err)
	require.NoError(t, client.Ping(ctx, nil))

	db := client.Database("cpf_management")
	log := &mockLogger{}

	// ⚙️ Inicializa o router real
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.InitRoutes(&router.RouterGroup, log, db)

	// 🔥 Faz a requisição POST /cliente
	payload := map[string]any{
		"nome":      "Valdinei",
		"documento": "12345678900",
		"telefone":  "11999999999",
		"bloqueado": false,
	}
	body, _ := json.Marshal(payload)
	reqHTTP := httptest.NewRequest(http.MethodPost, "/api/v1/cliente", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, reqHTTP)

	// ✅ Verificações
	res := w.Result()
	defer res.Body.Close()

	bodyResp, _ := io.ReadAll(res.Body)
	t.Logf("Response body: %s", string(bodyResp))
	require.Equal(t, http.StatusCreated, res.StatusCode, "esperava status 201 Created")

	// 🔍 Verifica se o cliente foi inserido no Mongo
	col := db.Collection("cliente")
	count, err := col.CountDocuments(ctx, map[string]any{"documento": "12345678900"})
	fmt.Println("count:", count)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
}
