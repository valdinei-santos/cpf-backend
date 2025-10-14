package delete

import (
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/repository"
)

// UseCase - Estrutura para o caso de uso de delete do cliente
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

// @Summary      Deleta um cliente pelo ID
// @Description  Deleta um cliente específico com base no ID fornecido
// @Tags         clientes
// @Produce      json
// @Param        id path string true "ID do cliente a ser deletado"
// @Success      200 {object} dto.OutputDefault "Cliente deletado com sucesso"}
// @Failure      404 {string} string "cliente não encontrado"
// @Router       /{id} [delete]
// Execute - Executa a lógica para deletar um cliente
func (u *UseCase) Execute(id string) error {
	u.log.Debug("Entrou delete.Execute")

	err := u.repo.DeleteCliente(id)
	if err != nil {
		u.log.Error(err.Error(), "mtd", "u.repo.Delete")
		return err
	}

	return nil
}
