package get_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/repository"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/get"
)

func TestExecute(t *testing.T) {
	// Pega um ID válido do mock de repositório
	mockRepoWithCliente := repository.NewMockClienteRepository()
	validID := mockRepoWithCliente.Clientes[0].ID.String()

	tests := []struct {
		name         string
		repo         *repository.MockClienteRepository
		logger       *logger.MockILogger
		inputID      string
		expectedResp *dto.Response
		expectedErr  error
		expectDebug  bool
		expectError  bool
	}{
		{
			name:    "Deve retornar sucesso quando o cliente é encontrado",
			repo:    mockRepoWithCliente, // Usa o mock que tem o cliente
			logger:  logger.NewMockILogger(),
			inputID: validID,
			expectedResp: &dto.Response{
				ID:        mockRepoWithCliente.Clientes[0].ID.String(),
				Nome:      mockRepoWithCliente.Clientes[0].Nome.String(),
				Documento: mockRepoWithCliente.Clientes[0].Documento.String(),
				Telefone:  mockRepoWithCliente.Clientes[0].Telefone.String(),
				Bloqueado: mockRepoWithCliente.Clientes[0].Bloqueado.Bool(),
			},
			expectedErr: nil,
			expectDebug: true,
			expectError: false,
		},
		{
			name:         "Deve retornar erro quando o ID não é encontrado",
			repo:         repository.NewMockClienteRepository(), // Usa uma nova instância sem o cliente
			logger:       logger.NewMockILogger(),
			inputID:      uuid.New().String(), // Gera um ID que não existe
			expectedResp: nil,
			expectedErr:  domainerr.ErrClienteNotFound,
			expectDebug:  true,
			expectError:  true,
		},
		{
			name:         "Deve retornar erro quando o ID é inválido",
			repo:         repository.NewMockClienteRepository(),
			logger:       logger.NewMockILogger(),
			inputID:      "id-invalido", // ID que não é um UUID
			expectedResp: nil,
			expectedErr:  domainerr.ErrClienteIDInvalid, //errors.New("ID inválido: invalid UUID length: 11"),
			expectDebug:  true,
			expectError:  true,
		},
		{
			name: "Deve retornar erro quando o repositório falha",
			repo: func() *repository.MockClienteRepository {
				r := repository.NewMockClienteRepository()
				r.SetMockError(errors.New("erro de conexão com o banco de dados"))
				return r
			}(),
			logger:       logger.NewMockILogger(),
			inputID:      validID,
			expectedResp: nil,
			expectedErr:  errors.New("erro de conexão com o banco de dados"),
			expectDebug:  true,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := get.NewUseCase(tt.repo, tt.logger)

			resp, err := uc.Execute(tt.inputID)

			assert.Equal(t, tt.expectedResp, resp)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expectDebug, tt.logger.DebugCalled)
			assert.Equal(t, tt.expectError, tt.logger.ErrorCalled)
		})
	}
}
