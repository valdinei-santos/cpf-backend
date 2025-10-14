package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/entities"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/vo"
)

// MockClienteRepository é um mock com a implementação da interface IClienteRepository
type MockClienteRepository struct {
	Clientes  []entities.Cliente
	mockError error
	//callCount int
}

// NewMockClienteRepository cria uma nova instancia de MockClienteRepository com 3 clientes padrão
func NewMockClienteRepository() *MockClienteRepository {
	return &MockClienteRepository{
		Clientes: []entities.Cliente{
			{ID: vo.FromUUID(uuid.New()), Nome: "Default Cliente1", Documento: "12345678901", Telefone: "11999999999", Bloqueado: false, CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
			{ID: vo.FromUUID(uuid.New()), Nome: "Default Cliente2", Documento: "10987654321", Telefone: "11888888888", Bloqueado: false, CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
			{ID: vo.FromUUID(uuid.New()), Nome: "Default Cliente3", Documento: "11122233344", Telefone: "11777777777", Bloqueado: true, CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
		},
	}
}

func (m *MockClienteRepository) SetMockError(err error) {
	m.mockError = err
}

// AddCliente - mock do método AddCliente
func (m *MockClienteRepository) AddCliente(p *entities.Cliente) error {
	if m.mockError != nil {
		return m.mockError
	}
	if p == nil {
		return domainerr.ErrClienteNotNil
	}
	// Cria um UUID
	p.ID = vo.FromUUID(uuid.New())
	// Adiciona o cliente ao slice
	m.Clientes = append(m.Clientes, *p)
	return nil
}

// GetClienteByID - mock do método GetClienteByID
func (m *MockClienteRepository) GetClienteByID(id string) (*entities.Cliente, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, domainerr.ErrClienteIDInvalid
	}
	for _, cliente := range m.Clientes {
		if cliente.ID == vo.FromUUID(idUUID) {
			return &cliente, nil
		}
	}
	return nil, errors.New("cliente não encontrado")
}

// GetManyClienteByIDs - busca vários clientes por ID
func (m *MockClienteRepository) GetManyClienteByIDs(ids []string) ([]*entities.Cliente, error) {
	if m.mockError != nil {
		return nil, m.mockError
	}

	var clientes []*entities.Cliente
	for _, id := range ids {
		idUUID, err := uuid.Parse(id)
		if err != nil {
			return nil, domainerr.ErrClienteIDInvalid
		}
		for _, cliente := range m.Clientes {
			if cliente.ID == vo.FromUUID(idUUID) {
				clientes = append(clientes, &cliente)
			}
		}
	}
	return clientes, nil
}

// GetAllClientes - mock do método GetAllClientes
func (m *MockClienteRepository) GetAllClientes(offset int64, limit int64) ([]*entities.Cliente, int64, error) {
	if m.mockError != nil {
		return nil, 0, m.mockError
	}

	total := int64(len(m.Clientes))

	// Aplica o offset e o limit para simular paginação
	if offset > total {
		return []*entities.Cliente{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	// Converte os clientes para um slice de ponteiros
	clientes := make([]*entities.Cliente, 0, end-offset)
	for i := offset; i < end; i++ {
		clientes = append(clientes, &m.Clientes[i])
	}

	return clientes, total, nil
}

// UpdateCliente - mock do método UpdateCliente
func (m *MockClienteRepository) UpdateCliente(id string, p *entities.Cliente) error {
	if m.mockError != nil {
		return m.mockError
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return domainerr.ErrClienteIDInvalid
	}
	for i, cliente := range m.Clientes {
		if cliente.ID == vo.FromUUID(idUUID) {
			// Atualiza o cliente existente com os novos valores
			p.ID = vo.FromUUID(idUUID) // Garante que o ID não seja alterado
			m.Clientes[i] = *p
			return nil
		}
	}
	return domainerr.ErrClienteNotFound
}

// DeleteCliente - mock do método DeleteCliente
func (m *MockClienteRepository) DeleteCliente(id string) error {
	if m.mockError != nil {
		return m.mockError
	}

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return domainerr.ErrClienteIDInvalid
	}
	for i, p := range m.Clientes {
		if p.ID == vo.FromUUID(idUUID) {
			m.Clientes = append(m.Clientes[:i], m.Clientes[i+1:]...)
			return nil
		}
	}
	return domainerr.ErrClienteNotFound
}

// Count - mock do método Count
func (r *MockClienteRepository) Count() (int64, error) {
	return int64(len(r.Clientes)), nil
}
