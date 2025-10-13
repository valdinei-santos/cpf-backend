package routes_test

import (
	"testing"
)

func TestInitRoutes(t *testing.T) {
	// Mocks
	/* mockILogger := logger.NewMockILogger()
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
	} */
}
