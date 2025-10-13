package vo

import "github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"

type TelefoneCliente string

func NewTelefoneCliente(desc string) (TelefoneCliente, error) {
	if len(desc) < 3 || len(desc) > 11 {
		return "", domainerr.ErrClienteTelefoneInvalid
	}
	return TelefoneCliente(desc), nil
}

func (t TelefoneCliente) String() string {
	return string(t)
}
