package get

import (
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/repository"
)

// UseCase - Struct do  caso de uso
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

// @Summary      Retorna um cliente pelo ID
// @Description  Retorna um cliente específico com base no ID fornecido
// @Tags         clientes
// @Accept       json
// @Produce      json
// @Param        id path string true "Cliente ID"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.OutputDefault
// @Router       /{id} [get]
// Execute - Executa a lógica de busca de um cliente
func (u *UseCase) Execute(id string) (*dto.Response, error) {
	u.log.Debug("Entrou get.Execute")

	// Pega o cliente no repositório pelo ID
	p, err := u.repo.GetClienteByID(id)
	if err != nil {
		u.log.Error(err.Error(), "mtd", "u.repo.GetClienteByID")
		return nil, err
	}

	// Transforma a entidade Cliente no DTO Response
	result := &dto.Response{
		ID:        p.ID.String(),
		Nome:      p.Nome.String(),
		Documento: p.Documento.String(),
		Telefone:  p.Telefone.String(),
		Bloqueado: p.Bloqueado.Bool(),
		CreatedAt: p.CreatedAt.String(),
		UpdatedAt: p.UpdatedAt.String(),
	}
	return result, nil
}
