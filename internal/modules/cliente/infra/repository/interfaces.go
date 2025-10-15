package repository

import "github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/entities"

// IClienteRepository define a interface para as operações de Cliente
type IClienteRepository interface {
	AddCliente(p *entities.Cliente) error
	GetClienteByID(id string) (*entities.Cliente, error)
	GetAllClientes(offset int64, limit int64) ([]*entities.Cliente, int64, error)
	UpdateCliente(id string, p *entities.Cliente) error
	DeleteCliente(id string) error
	Count() (int64, error)
}
