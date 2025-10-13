package update_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/repository"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/update"
)

func TestExecute(t *testing.T) {
	mockRepo := repository.NewMockClienteRepository()
	validClienteID := mockRepo.Clientes[0].ID.String()

	tests := []struct {
		name         string
		id           string
		repo         *repository.MockClienteRepository
		logger       *logger.MockILogger
		input        *dto.Request
		expectedResp *dto.Response
		expectedErr  error
		expectDebug  bool
		expectError  bool
	}{
		{
			name:   "Deve atualizar um cliente com sucesso",
			id:     validClienteID,
			repo:   mockRepo,
			logger: logger.NewMockILogger(),
			input: &dto.Request{
				Nome:      "Cliente Atualizado",
				Documento: "12345678900",
				Telefone:  "11999999999",
				Bloqueado: false,
			},
			expectedErr: nil,
			expectDebug: true,
			expectError: false,
		},
		{
			name:   "Deve retornar erro ao tentar atualizar cliente com ID inválido",
			id:     "id-invalido",
			repo:   mockRepo,
			logger: logger.NewMockILogger(),
			input: &dto.Request{
				Nome:      "Cliente Atualizado",
				Documento: "12345678900",
				Telefone:  "11999999999",
				Bloqueado: false,
			},
			expectedErr: domainerr.ErrClienteNotFound,
			expectDebug: true,
			expectError: true,
		},
		{
			name: "Deve retornar erro quando o repositório falha ao atualizar",
			id:   validClienteID,
			repo: func() *repository.MockClienteRepository {
				r := repository.NewMockClienteRepository()
				r.SetMockError(domainerr.ErrClienteNotFound)
				return r
			}(),
			logger: logger.NewMockILogger(),
			input: &dto.Request{
				Nome:      "Cliente Atualizado",
				Documento: "12345678900",
				Telefone:  "11999999999",
				Bloqueado: false,
			},
			expectedErr: domainerr.ErrClienteNotFound,
			expectDebug: true,
			expectError: true,
		},
		{
			name:   "Deve retornar erro ao tentar atualizar com dados de entrada inválidos",
			id:     validClienteID,
			repo:   mockRepo,
			logger: logger.NewMockILogger(),
			input: &dto.Request{
				Nome:      "", // Nome vazio
				Documento: "12345678900",
				Telefone:  "11999999999",
				Bloqueado: false,
			},
			expectedErr: domainerr.ErrClienteNomeInvalid,
			expectDebug: true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := update.NewUseCase(tt.repo, tt.logger)

			resp, err := uc.Execute(tt.id, tt.input)

			//Verifique se não houve erro
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.Nil(t, err)
				//Verifique somente os campos que não mudam
				assert.Equal(t, tt.input.Nome, resp.Nome)
				assert.Equal(t, tt.input.Documento, resp.Documento)
				assert.Equal(t, tt.input.Telefone, resp.Telefone)
				assert.Equal(t, tt.input.Bloqueado, resp.Bloqueado)

				// Verifique se o ID foi gerado
				assert.NotEmpty(t, resp.ID)

				// Verifique se as datas não estão vazias e estão no formato correto
				assert.NotEmpty(t, resp.CreatedAt)
				assert.NotEmpty(t, resp.UpdatedAt)
			}
			assert.Equal(t, tt.expectDebug, tt.logger.DebugCalled)
			assert.Equal(t, tt.expectError, tt.logger.ErrorCalled)
		})
	}
}
