//go:build integration
// +build integration

package routes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/valdinei-santos/cpf-backend/cmd/api/routes"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
)

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

type testEnv struct {
	ctx    context.Context
	client *mongo.Client
	db     *mongo.Database
	router *gin.Engine
}

func setupIntegrationTest(t *testing.T) *testEnv {
	ctx := context.Background()

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
	t.Cleanup(func() { mongoC.Terminate(ctx) })

	host, _ := mongoC.Host(ctx)
	port, _ := mongoC.MappedPort(ctx, "27017")
	mongoURI := "mongodb://" + host + ":" + port.Port()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err)
	require.NoError(t, client.Ping(ctx, nil))

	db := client.Database("cpf_management")
	log := &mockLogger{}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.InitRoutes(&router.RouterGroup, log, db)

	return &testEnv{ctx, client, db, router}
}

// -----------------------------------------------------------------------------
// POST /api/v1/cliente
// -----------------------------------------------------------------------------
func TestClientePost_Integration(t *testing.T) {
	env := setupIntegrationTest(t)

	payload := map[string]any{
		"nome":      "Valdinei",
		"documento": "12345678900",
		"telefone":  "11999999999",
		"bloqueado": false,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cliente", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	col := env.db.Collection("cliente")
	count, err := col.CountDocuments(env.ctx, bson.M{"documento": "12345678900"})
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
}

// -----------------------------------------------------------------------------
// GET /api/v1/cliente/:id
// -----------------------------------------------------------------------------
func TestClienteGetByID_Integration(t *testing.T) {
	env := setupIntegrationTest(t)

	// Cria um cliente
	body := []byte(`{"nome":"João","documento":"98765432100","telefone":"11988888888","bloqueado":false}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cliente", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	id := created["id"].(string)

	// Faz GET /cliente/:id
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/cliente/"+id, nil)
	wGet := httptest.NewRecorder()
	env.router.ServeHTTP(wGet, reqGet)

	require.Equal(t, http.StatusOK, wGet.Code)
	require.Contains(t, wGet.Body.String(), "João")
}

// -----------------------------------------------------------------------------
// GET /api/v1/cliente
// -----------------------------------------------------------------------------
func TestClienteList_Integration(t *testing.T) {
	env := setupIntegrationTest(t)

	// Insere dois clientes
	for _, nome := range []string{"Maria", "Pedro"} {
		body := map[string]any{"nome": nome, "documento": nome + "123", "telefone": "1111", "bloqueado": false}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/cliente", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
	}

	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/cliente", nil)
	wList := httptest.NewRecorder()
	env.router.ServeHTTP(wList, reqList)

	require.Equal(t, http.StatusOK, wList.Code)
	require.Contains(t, wList.Body.String(), "Maria")
	require.Contains(t, wList.Body.String(), "Pedro")
}

// -----------------------------------------------------------------------------
// PUT /api/v1/cliente/:id
// -----------------------------------------------------------------------------
func TestClientePut_Integration(t *testing.T) {
	env := setupIntegrationTest(t)

	// Cria cliente
	body := []byte(`{"nome":"Carlos","documento":"99999999999","telefone":"11999999999","bloqueado":false}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cliente", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	id := created["id"].(string)

	// Atualiza cliente
	update := []byte(`{"nome":"Carlos Atualizado","documento":"12345678900","telefone":"11988888888","bloqueado":false}`)
	reqPut := httptest.NewRequest(http.MethodPut, "/api/v1/cliente/"+id, bytes.NewBuffer(update))
	reqPut.Header.Set("Content-Type", "application/json")
	wPut := httptest.NewRecorder()
	env.router.ServeHTTP(wPut, reqPut)

	require.Equal(t, http.StatusOK, wPut.Code)
	require.Contains(t, wPut.Body.String(), "Atualizado")
}

// -----------------------------------------------------------------------------
// DELETE /api/v1/cliente/:id
// -----------------------------------------------------------------------------
func TestClienteDelete_Integration(t *testing.T) {
	env := setupIntegrationTest(t)

	// Cria cliente
	body := []byte(`{"nome":"Ana","documento":"77777777777","telefone":"11777777777","bloqueado":false}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/cliente", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	id := created["id"].(string)

	// Deleta cliente
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/v1/cliente/"+id, nil)
	wDel := httptest.NewRecorder()
	env.router.ServeHTTP(wDel, reqDel)
	require.Equal(t, http.StatusNoContent, wDel.Code)

	// Verifica se foi removido
	col := env.db.Collection("cliente")
	count, err := col.CountDocuments(env.ctx, bson.M{"id": id})
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}
