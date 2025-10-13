package create_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdinei-santos/cpf-backend/internal/domain/globalerr"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/repository"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/create"
)

func TestExecute(t *testing.T) {
	// Pega um ID válido do mock de repositório
	//mockRepoWithCliente := repository.NewMockClienteRepository()
	//validID := mockRepoWithCliente.Clientes[0].ID.String()

	// Tabela de teste para cenários
	tests := []struct {
		name         string
		repo         *repository.MockClienteRepository
		logger       *logger.MockILogger
		input        *dto.Request
		expectedResp *dto.Response
		expectedErr  error
		expectDebug  bool
		expectError  bool
	}{
		{
			name:   "Deve retornar sucesso com dados válidos",
			repo:   repository.NewMockClienteRepository(), // Cria uma nova instância para não compartilhar estado. Dava erro
			logger: logger.NewMockILogger(),
			input: &dto.Request{
				Nome:      "Cliente Teste",
				Documento: "12345678900",
				Telefone:  "11999999999",
				Bloqueado: false,
			},
			expectedErr: nil,
			expectDebug: true,
			expectError: false,
		},
		{
			name:        "Deve retornar error quando dados inválidos são fornecidos",
			repo:        repository.NewMockClienteRepository(),
			logger:      logger.NewMockILogger(),
			input:       &dto.Request{Nome: ""}, // Nome vazio causa erro
			expectedErr: domainerr.ErrClienteNomeInvalid,
			expectDebug: true,
			expectError: true,
		},
		{
			name: "Deve retornar error quando repositório falha",
			repo: func() *repository.MockClienteRepository {
				r := repository.NewMockClienteRepository()
				r.SetMockError(globalerr.ErrInternal)
				return r
			}(),
			logger: logger.NewMockILogger(),
			input: &dto.Request{
				Nome:      "Cliente Teste 2",
				Documento: "12345678900",
				Telefone:  "11999999999",
				Bloqueado: false,
			},
			expectedErr: globalerr.ErrInternal,
			expectDebug: true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := create.NewUseCase(tt.repo, tt.logger)

			resp, err := uc.Execute(tt.input)

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
