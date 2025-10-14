package delete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdinei-santos/cpf-backend/internal/domain/globalerr"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/repository"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/delete"
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
		expectedResp *dto.OutputDefault
		expectedErr  error
		expectDebug  bool
		expectError  bool
	}{
		{
			name:         "Deve retornar sucesso ao excluir um cliente",
			repo:         mockRepoWithCliente,
			logger:       logger.NewMockILogger(),
			inputID:      validID,
			expectedResp: &dto.OutputDefault{},
			expectedErr:  nil,
			expectDebug:  true,
			expectError:  false,
		},
		{
			name: "Deve retornar erro se o repositório falhar",
			repo: func() *repository.MockClienteRepository {
				r := repository.NewMockClienteRepository()
				r.SetMockError(globalerr.ErrSaveInDatabase)
				return r
			}(),
			logger:       logger.NewMockILogger(),
			inputID:      validID,
			expectedResp: nil,
			expectedErr:  globalerr.ErrSaveInDatabase,
			expectDebug:  true,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := delete.NewUseCase(tt.repo, tt.logger)

			err := uc.Execute(tt.inputID)

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
