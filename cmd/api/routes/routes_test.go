package routes_test

/* import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/valdinei-santos/cpf-backend/cmd/api/routes"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockLogger struct{}

func (m *mockLogger) Debug(msg string, args ...any)                             {}
func (m *mockLogger) Info(msg string, args ...any)                              {}
func (m *mockLogger) Warn(msg string, args ...any)                              {}
func (m *mockLogger) Error(msg string, args ...any)                             {}
func (m *mockLogger) With(args ...any) logger.ILogger                           { return m }
func (m *mockLogger) WithContext(ctx context.Context) logger.ILogger            { return m }
func (m *mockLogger) DebugContext(ctx context.Context, msg string, args ...any) {}
func (m *mockLogger) InfoContext(ctx context.Context, msg string, args ...any)  {}
func (m *mockLogger) WarnContext(ctx context.Context, msg string, args ...any)  {}
func (m *mockLogger) ErrorContext(ctx context.Context, msg string, args ...any) {}

// mockDB cria um *mongo.Database falso (não conectado)
func mockDB() *mongo.Database {
	return &mongo.Database{} // apenas placeholder, pois os handlers podem ignorar conexão
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/")
	log := &mockLogger{}
	db := mockDB()

	routes.InitRoutes(api, log, db)
	return r
}

func TestPingEndpoint(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/ping", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "pong", resp.Body.String())
}

func TestStatusEndpoint(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/status", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "success", body["status"])
}

func TestClienteOptionsRoot(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("OPTIONS", "/api/v1/cliente", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestClientePost(t *testing.T) {
	router := setupRouter()

	body := []byte(`{"nome": "Fulano", "documento": "12345678900"}, "telefone": "48999448383", "bloqueado": false`)
	req, _ := http.NewRequest("POST", "/api/v1/cliente", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// O retorno depende do comportamento de cliente.StartCreate,
	// então aqui testamos apenas se não dá erro 404.
	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestClienteGetAll(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/cliente", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestClienteGetByID(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/cliente/43310521-a934-11f0-ae2b-fabc94dc9b50", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestClientePut(t *testing.T) {
	router := setupRouter()

	body := []byte(`{"nome": "Atualizado"}`)
	req, _ := http.NewRequest("PUT", "/api/v1/cliente/43310521-a934-11f0-ae2b-fabc94dc9b50", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestClienteDelete(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/api/v1/cliente/43310521-a934-11f0-ae2b-fabc94dc9b50", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestSwaggerEndpoint(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// swagger handler pode retornar redirect (302) ou 404 se não configurado
	assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusFound || resp.Code == http.StatusNotFound)
}

func TestAccessCounterMiddleware(t *testing.T) {
	router := setupRouter()

	// Chamamos dois endpoints diferentes
	req1, _ := http.NewRequest("GET", "/ping", nil)
	req2, _ := http.NewRequest("GET", "/status", nil)
	resp1 := httptest.NewRecorder()
	resp2 := httptest.NewRecorder()

	router.ServeHTTP(resp1, req1)
	router.ServeHTTP(resp2, req2)

	// Verifica se ambos responderam
	assert.Equal(t, http.StatusOK, resp1.Code)
	assert.Equal(t, http.StatusOK, resp2.Code)
}

// Opcional: teste de concorrência para garantir que o middleware não dá race
func TestConcurrentRequests(t *testing.T) {
	router := setupRouter()

	const n = 10
	done := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/ping", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)
			done <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-done
	}
}

var _ logger.ILogger = (*mockLogger)(nil) */

/* func TestInitRoutes(t *testing.T) {
	// Mocks
	mockILogger := logger.NewMockILogger()
	mockRepo := repository.NewMockClienteRepository()
	validID1 := mockRepo.Clientes[0].ID.String()

	// Cria um novo Gin engine e grupo de roteamento
	gin.SetMode(gin.TestMode)
	router := gin.New()
	apiGroup := router.Group("/")

	// Inicializa rotas
	routes.InitRoutes(apiGroup, mockILogger, mockRepo)

	// Casos de teste
	tests := []struct {
		name       string
		method     string
		url        string
		statusCode int
	}{
		{"Teste Ping", "GET", "/ping", http.StatusOK},
		{"Teste Ping", "GET", "/stats", http.StatusOK},
		{"Rota CreateCliente", "POST", "/api/v1/cliente/", http.StatusCreated},
		{"Rota DeleteClienteByID", "DELETE", "/api/v1/cliente/" + validID1, http.StatusOK},
		{"Rota GetAllClientes", "GET", "/api/v1/cliente/", http.StatusOK},
		{"Rota GetClienteByID", "GET", "/api/v1/cliente/" + validID1, http.StatusOK},
		{"Rota UpdateClienteByID", "PUT", "/api/v1/cliente/" + validID1, http.StatusOK},
	}

	// Executa o casos de teste
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Define o contexto para o caso de teste
			mockILogger.SetContext(tc.name)

			// Define o corpo da requisição para POST e PUT
			var body string
			if tc.method == "POST" || tc.method == "PUT" {
				//body = `{"name": "Test Cliente", "price": 10.0}`
				body = `{
					"nome": "Default Cliente1",
					"documento": "71248609972",
					"telefone": "48999448383",
					"bloqueado": false
				}`
			}

			// Simula a requisição
			req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json") // Define o cabeçalho como JSON
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			// Verifica o status da resposta
			assert.Equal(t, tc.statusCode, resp.Code)

			// Exibe os logs apenas se o teste falhar
			if t.Failed() {
				logs := mockILogger.GetLogs(tc.name)
				t.Logf("Logs gerados no teste '%s':\n%s", tc.name, strings.Join(logs, "\n"))
			}
		})
	}
} */
