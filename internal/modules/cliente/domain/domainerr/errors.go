// Pacote com erros locais do domínio.
package domainerr

import "errors"

var (
	ErrClienteNomeInvalid      = errors.New("o nome deve ter entre 3 e 50 caracteres")
	ErrClienteDocumentoInvalid = errors.New("o documento deve ter entre 11 e 14 digitos")
	ErrClienteTelefoneInvalid  = errors.New("o telefone deve ter entre 3 e 11 digitos")
	ErrClienteBloqueadoInvalid = errors.New("o bloqueio deve ser true ou false")
	ErrClienteIDInvalid        = errors.New("ID inválido")
	ErrClienteUUIDInvalid      = errors.New("UUID inválido")
	ErrClienteNotFound         = errors.New("cliente não encontrado")
	ErrClienteNotNil           = errors.New("cliente não pode ser nil")
	ErrClienteNotFoundMany     = errors.New("nenhum cliente encontrado")
	ErrDuplicatekey            = errors.New("registro já existe")
)
