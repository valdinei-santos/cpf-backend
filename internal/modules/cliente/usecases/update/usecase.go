package update

import (
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"
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

// @Summary      Atualiza um cliente pelo ID
// @Description  Atualiza um cliente existente com base no ID
// @Tags         clientes
// @Accept       json
// @Produce      json
// @Param        id path string true "Cliente ID"
// @Param        cliente body dto.Request  true  "Dados do cliente para atualização"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.OutputDefault
// @Router       /{id} [put]
// Execute - Executa a lógica de criação de um cliente
func (u *UseCase) Execute(id string, in *dto.Request) (*dto.Response, error) {
	u.log.Debug("Entrou update.Execute")

	// Pega o cliente no repositório pelo ID
	c, err := u.repo.GetClienteByID(id)
	if err != nil {
		u.log.Error(err.Error(), "mtd", "u.repo.GetClienteByID")
		return nil, domainerr.ErrClienteNotFound
	}

	// Altera o objeto Cliente a partir do Cliente enviado no DTO de entrada
	pNew, err := entities.UpdateCliente(id, in.Nome, in.Documento, in.Telefone, in.Bloqueado, c.CreatedAt)
	if err != nil {
		u.log.Error(err.Error(), "mtd", "entities.UpdateCliente")
		return nil, err
	}

	// Altera o cliente no repositório
	err = u.repo.UpdateCliente(id, pNew)
	if err != nil {
		u.log.Error(err.Error(), "mtd", "u.repo.UpdateCliente")
		return nil, err
	}

	// Retorna o DTO de saída
	resp := &dto.Response{
		ID:        pNew.ID.String(),
		Nome:      pNew.Nome.String(),
		Documento: pNew.Documento.String(),
		Telefone:  pNew.Telefone.String(),
		Bloqueado: pNew.Bloqueado.Bool(),
		CreatedAt: pNew.CreatedAt.String(),
		UpdatedAt: pNew.UpdatedAt.String(),
	}
	return resp, nil
}
