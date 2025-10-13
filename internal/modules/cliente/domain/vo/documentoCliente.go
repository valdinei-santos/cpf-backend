package vo

import "github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"

type DocumentoCliente string

func NewDocumentoCliente(desc string) (DocumentoCliente, error) {
	if len(desc) < 5 || len(desc) > 100 {
		return "", domainerr.ErrClienteDocumentoInvalid
	}
	return DocumentoCliente(desc), nil
}

func (d DocumentoCliente) String() string {
	return string(d)
}
