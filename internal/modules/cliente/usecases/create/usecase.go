package create

import (
	"github.com/valdinei-santos/cpf-backend/internal/domain/globalerr"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/entities"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/repository"
)

// UseCase - Estrutura para o caso de uso de criação de cliente
type UseCase struct {
	repo repository.IClienteRepository
	log  logger.ILogger
}

// NewUseCase - Construtor do caso de uso
func NewUseCase(r repository.IClienteRepository, l logger.ILogger) *UseCase {
	return &UseCase{
		repo: r,
		log:  l,
	}
}

// @Summary      Cria um novo cliente
// @Description  Cria um novo cliente com os dados fornecidos
// @Tags         clientes
// @Accept       json
// @Produce      json
// @Param        cliente body dto.Request true "Dados do cliente a ser criado"
// @Success      201 {object} dto.Response
// @Failure      400 {object} string "Erro na requisição"
// @Router       / [post]
// Execute - Executa a lógica de criação de um cliente
func (u *UseCase) Execute(in *dto.Request) (*dto.Response, error) {
	u.log.Debug("Entrou create.Execute")

	// Cria o objeto Cliente a partir do DTO de entrada
	p, err := entities.NewCliente(in.Nome, in.Documento, in.Telefone, in.Bloqueado)
	if err != nil {
		u.log.Error(err.Error(), "mtd", "entities.NewCliente")
		return nil, err
	}

	// Salva o cliente no repositório
	err = u.repo.AddCliente(p)
	if err != nil {
		u.log.Error(err.Error(), "mtd", "u.repo.AddCliente")
		return nil, globalerr.ErrInternal
	}
	u.log.Debug("Cliente criado com sucesso", "id", p.ID)
	// Retorna o DTO de saída
	resp := &dto.Response{
		ID:        p.ID.String(),
		Nome:      p.Nome.String(),
		Documento: p.Documento.String(),
		Telefone:  p.Telefone.String(),
		Bloqueado: p.Bloqueado.Bool(),
		CreatedAt: p.CreatedAt.String(),
		UpdatedAt: p.UpdatedAt.String(),
	}
	return resp, nil
}
