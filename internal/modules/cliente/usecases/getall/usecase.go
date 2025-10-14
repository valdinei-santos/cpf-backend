package getall

import (
	"math"

	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
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

// @Summary        Lista todos os clientes
// @Description    Retorna uma lista de clientes, paginada
// @Tags           clientes
// @Produce        json
// @Param          page query int64 false "Numero da página a ser retornada"
// @Param          size query int64 false "Quantidade de itens na página a ser retornada"
// @Success        200 {array} dto.ResponseManyPaginated
// @Failure        500 {string} string "Erro interno do servidor"
// @Router         / [get]
// Execute - Executa a lógica para buscar todos os clientes
func (u *UseCase) Execute(page int64, size int64) (*dto.ResponseManyPaginated, error) {
	u.log.Debug("Entrou getall.Execute")

	// Calcula o offset para o repositório
	offset := (page - 1) * size

	// Busca o subconjunto de clientes e o total de itens
	paginatedClientes, totalItems, err := u.repo.GetAllClientes(offset, size)
	if err != nil {
		u.log.Error("Erro ao buscar clientes: ", err)
		return nil, err
	}

	// Converte as entidades para DTOs
	clienteList := make([]dto.Response, len(paginatedClientes))
	for i, p := range paginatedClientes {
		clienteList[i] = dto.Response{
			ID:        p.ID.String(),
			Nome:      p.Nome.String(),
			Documento: p.Documento.String(),
			Telefone:  p.Telefone.String(),
			Bloqueado: p.Bloqueado.Bool(),
			CreatedAt: p.CreatedAt.String(),
			UpdatedAt: p.UpdatedAt.String(),
		}
	}

	// Calcula o total de páginas
	totalPages := int(math.Ceil(float64(totalItems) / float64(size)))
	if totalPages == 0 && totalItems > 0 { // Lida com o caso de 1 página.
		totalPages = 1
	}

	// Constrói a resposta paginada
	result := &dto.ResponseManyPaginated{
		Clientes:     clienteList,
		TotalItems:   totalItems,
		TotalPages:   int64(totalPages),
		CurrentPage:  page,
		ItemsPerPage: size,
	}

	return result, nil
}
