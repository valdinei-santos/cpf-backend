package repository

import (
	"encoding/json"

	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/entities"
)

// ClienteRepo é um repositório para gerenciar produtos
type ClienteRepo struct {
	filePath string
	products []*entities.Cliente
	mutex    sync.Mutex
}

// NewClienteRepo - cria uma nova instância do repositório
func NewClienteRepo(filePath string) (*ClienteRepo, error) {
	repo := &ClienteRepo{
		filePath: filePath,
	}
	err := repo.Load()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// Load - carrega os dados do arquivo JSON para o repositório
func (r *ClienteRepo) Load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		// Se o arquivo não existe, inicializa a lista vazia
		if os.IsNotExist(err) {
			r.products = []*entities.Cliente{}
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &r.products)
}

// Save - salva os dados do repositório no arquivo JSON
func (r *ClienteRepo) Save() error {
	data, err := json.MarshalIndent(r.products, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

// AddCliente - adiciona um novo produto ao repositório
func (r *ClienteRepo) AddCliente(p *entities.Cliente) error {
	r.mutex.Lock()         // Bloqueia o acesso ao repositório/Arquivo JSON
	defer r.mutex.Unlock() // Libera o acesso ao repositório/Arquivo JSON
	r.products = append(r.products, p)
	return r.Save()
}

// GetClienteByID - busca um produto por ID
func (r *ClienteRepo) GetClienteByID(id string) (*entities.Cliente, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, domainerr.ErrClienteIDInvalid
	}
	for _, product := range r.products {
		if product.ID == idUUID {
			return product, nil
		}
	}
	return nil, domainerr.ErrClienteNotFound
}

// GetManyClienteByIDs - busca vários produtos por ID
func (r *ClienteRepo) GetManyClienteByIDs(ids []string) ([]*entities.Cliente, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var products []*entities.Cliente
	for _, id := range ids {
		idUUID, err := uuid.Parse(id)
		if err != nil {
			return nil, domainerr.ErrClienteIDInvalid
		}
		for _, product := range r.products {
			if product.ID == idUUID {
				products = append(products, product)
			}
		}
	}
	if len(products) == 0 {
		return nil, domainerr.ErrClienteNotFoundMany
	}
	return products, nil
}

// GetAllClientes - retorna todos os produtos
func (r *ClienteRepo) GetAllClientes(offset int, limit int) ([]*entities.Cliente, int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Garante que o offset e o limit não extrapolem o tamanho do array
	total := len(r.products)
	if offset >= total {
		return []*entities.Cliente{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}

	return r.products[offset:end], total, nil
}

// UpdateCliente - atualiza os dados de um produto existente
func (r *ClienteRepo) UpdateCliente(id string, p *entities.Cliente) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return domainerr.ErrClienteIDInvalid
	}
	for i, product := range r.products {
		if product.ID == idUUID {
			r.products[i] = p
			return r.Save()
		}
	}
	return domainerr.ErrClienteNotFound
}

// DeleteCliente - remove um produto do repositório
func (r *ClienteRepo) DeleteCliente(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return domainerr.ErrClienteIDInvalid
	}
	for i, product := range r.products {
		if product.ID == idUUID {
			// Remove o produto da slice
			r.products = append(r.products[:i], r.products[i+1:]...)
			return r.Save()
		}
	}
	return domainerr.ErrClienteNotFound
}

func (r *ClienteRepo) Count() (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return len(r.products), nil
}
