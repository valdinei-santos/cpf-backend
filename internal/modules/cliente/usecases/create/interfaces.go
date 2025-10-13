package create

import "github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"

// IUsecase - ...
type IUsecase interface {
	Execute(p *dto.Request) (*dto.Response, error)
}
