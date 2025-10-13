package getall

import "github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"

// IUsecase - ...
type IUsecase interface {
	Execute(page int64, size int64) (*dto.ResponseManyPaginated, error)
}
