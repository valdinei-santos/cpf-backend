package vo

import "github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"

type NomeCliente string

func NewNomeCliente(nome string) (NomeCliente, error) {
	if len(nome) < 3 || len(nome) > 50 {
		return "", domainerr.ErrClienteNomeInvalid
	}
	return NomeCliente(nome), nil
}

func (n NomeCliente) String() string {
	return string(n)
}
