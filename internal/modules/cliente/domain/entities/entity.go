package entities

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/vo"
)

// Cliente representa a estrutura de um cliente
type Cliente struct {
	ID        uuid.UUID
	Nome      vo.NomeCliente
	Documento vo.DocumentoCliente
	Telefone  vo.TelefoneCliente
	Bloqueado vo.BloqueadoCliente
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCliente - cria uma nova instância de Cliente
func NewCliente(nome, documento, telefone string, bloqueado bool) (*Cliente, error) {
	uuidVO, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	nomeVO, err := vo.NewNomeCliente(nome)
	if err != nil {
		return nil, err
	}
	documentoVO, err := vo.NewDocumentoCliente(documento)
	if err != nil {
		return nil, err
	}
	telefoneVO, err := vo.NewTelefoneCliente(telefone)
	if err != nil {
		return nil, err
	}
	bloqueadoVO, err := vo.NewBloqueadoCliente(bloqueado)
	if err != nil {
		return nil, err
	}

	c := &Cliente{
		ID:        uuidVO,
		Nome:      nomeVO,
		Documento: documentoVO,
		Telefone:  telefoneVO,
		Bloqueado: bloqueadoVO,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = c.validate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// UpdateCliente - altera uma instância de Cliente
func UpdateCliente(id, nome, documento, telefone string, bloqueado bool, createdAt time.Time) (*Cliente, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, domainerr.ErrClienteIDInvalid
	}
	nomeVO, err := vo.NewNomeCliente(nome)
	if err != nil {
		return nil, err
	}
	documentoVO, err := vo.NewDocumentoCliente(documento)
	if err != nil {
		return nil, err
	}
	telefoneVO, err := vo.NewTelefoneCliente(telefone)
	if err != nil {
		return nil, err
	}
	bloqueadoVO, err := vo.NewBloqueadoCliente(bloqueado)
	if err != nil {
		return nil, err
	}
	c := &Cliente{
		ID:        idUUID,
		Nome:      nomeVO,
		Documento: documentoVO,
		Telefone:  telefoneVO,
		Bloqueado: bloqueadoVO,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}
	err = c.validate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Validate - Valida os campos do Cliente
func (c *Cliente) validate() error {
	return validator.New().Struct(c)
}
