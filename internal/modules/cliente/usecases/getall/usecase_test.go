package getall_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/repository"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/getall"
)

func TestExecute(t *testing.T) {
	// Cria uma instância do mock de repositório para ser usada em todos os cenários
	mockRepo := repository.NewMockClienteRepository()

	tests := []struct {
		name         string
		repo         *repository.MockClienteRepository
		logger       *logger.MockILogger
		page         int
		size         int
		expectedResp *dto.ResponseManyPaginated
		expectedErr  error
		expectDebug  bool
		expectError  bool
	}{
		{
			name:   "Deve retornar a primeira página de clientes com sucesso",
			repo:   mockRepo,
			logger: logger.NewMockILogger(),
			page:   1,
			size:   2,
			expectedResp: &dto.ResponseManyPaginated{
				Clientes: []dto.Response{
					{
						ID:        mockRepo.Clientes[0].ID.String(),
						Nome:      mockRepo.Clientes[0].Nome.String(),
						Documento: mockRepo.Clientes[0].Documento.String(),
						Telefone:  mockRepo.Clientes[0].Telefone.String(),
						Bloqueado: mockRepo.Clientes[0].Bloqueado.Bool(),
					},
					{
						ID:        mockRepo.Clientes[1].ID.String(),
						Nome:      mockRepo.Clientes[1].Nome.String(),
						Documento: mockRepo.Clientes[1].Documento.String(),
						Telefone:  mockRepo.Clientes[1].Telefone.String(),
						Bloqueado: mockRepo.Clientes[1].Bloqueado.Bool(),
					},
				},
				TotalItems:   int64(len(mockRepo.Clientes)),
				TotalPages:   2, // 3 clientes com 2 por página = 2 páginas
				CurrentPage:  1,
				ItemsPerPage: 2,
			},
			expectedErr: nil,
			expectDebug: true,
			expectError: false,
		},
		{
			name:   "Deve retornar a segunda página de clientes com sucesso",
			repo:   mockRepo,
			logger: logger.NewMockILogger(),
			page:   2,
			size:   2,
			expectedResp: &dto.ResponseManyPaginated{
				Clientes: []dto.Response{
					{
						ID:        mockRepo.Clientes[2].ID.String(),
						Nome:      mockRepo.Clientes[2].Nome.String(),
						Documento: mockRepo.Clientes[2].Documento.String(),
						Telefone:  mockRepo.Clientes[2].Telefone.String(),
						Bloqueado: mockRepo.Clientes[2].Bloqueado.Bool(),
					},
				},
				TotalItems:   int64(len(mockRepo.Clientes)),
				TotalPages:   2,
				CurrentPage:  2,
				ItemsPerPage: 2,
			},
			expectedErr: nil,
			expectDebug: true,
			expectError: false,
		},
		{
			name:   "Deve retornar página vazia para página que não existe",
			repo:   mockRepo,
			logger: logger.NewMockILogger(),
			page:   10,
			size:   2,
			expectedResp: &dto.ResponseManyPaginated{
				Clientes:     []dto.Response{},
				TotalItems:   int64(len(mockRepo.Clientes)),
				TotalPages:   2,
				CurrentPage:  10,
				ItemsPerPage: 2,
			},
			expectedErr: nil,
			expectDebug: true,
			expectError: false,
		},
		{
			name:   "Deve retornar todos os clientes quando o limite é maior que o total",
			repo:   mockRepo,
			logger: logger.NewMockILogger(),
			page:   1,
			size:   10, // Maior que o número de clientes
			expectedResp: &dto.ResponseManyPaginated{
				Clientes: []dto.Response{
					{
						ID:        mockRepo.Clientes[0].ID.String(),
						Nome:      mockRepo.Clientes[0].Nome.String(),
						Documento: mockRepo.Clientes[0].Documento.String(),
						Telefone:  mockRepo.Clientes[0].Telefone.String(),
						Bloqueado: mockRepo.Clientes[0].Bloqueado.Bool(),
					},
					{
						ID:        mockRepo.Clientes[1].ID.String(),
						Nome:      mockRepo.Clientes[1].Nome.String(),
						Documento: mockRepo.Clientes[1].Documento.String(),
						Telefone:  mockRepo.Clientes[1].Telefone.String(),
						Bloqueado: mockRepo.Clientes[1].Bloqueado.Bool(),
					},
					{
						ID:        mockRepo.Clientes[2].ID.String(),
						Nome:      mockRepo.Clientes[2].Nome.String(),
						Documento: mockRepo.Clientes[2].Documento.String(),
						Telefone:  mockRepo.Clientes[2].Telefone.String(),
						Bloqueado: mockRepo.Clientes[2].Bloqueado.Bool(),
					},
				},
				TotalItems:   int64(len(mockRepo.Clientes)),
				TotalPages:   1,
				CurrentPage:  1,
				ItemsPerPage: 10,
			},
			expectedErr: nil,
			expectDebug: true,
			expectError: false,
		},
		{
			name: "Deve retornar erro quando o repositório falha",
			repo: func() *repository.MockClienteRepository {
				r := repository.NewMockClienteRepository()
				r.SetMockError(errors.New("erro de conexão com o banco de dados"))
				return r
			}(),
			logger:       logger.NewMockILogger(),
			page:         1,
			size:         2,
			expectedResp: nil,
			expectedErr:  errors.New("erro de conexão com o banco de dados"),
			expectDebug:  true,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := getall.NewUseCase(tt.repo, tt.logger)

			resp, err := uc.Execute(int64(tt.page), int64(tt.size))

			assert.Equal(t, tt.expectedResp, resp)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectDebug, tt.logger.DebugCalled)
			assert.Equal(t, tt.expectError, tt.logger.ErrorCalled)
		})
	}
}
