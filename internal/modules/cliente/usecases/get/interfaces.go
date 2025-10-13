package get

import "github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"

// IUsecase - ...
type IUsecase interface {
	Execute(id string) (*dto.Response, error)
}
