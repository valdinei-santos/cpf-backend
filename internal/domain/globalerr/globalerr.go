package globalerr

import "errors"

var (
	ErrDuplicatekey         = errors.New("chave duplicada")
	ErrNotFound             = errors.New("não encontrado")
	ErrNotNil               = errors.New("não pode ser nil")
	ErrBadRequest           = errors.New("erro nos dados do request")
	ErrSaveInDatabase       = errors.New("erro ao salvar no banco de dados")
	ErrConnectionInDatabase = errors.New("erro de conexão com o banco de dados")
	ErrInternal             = errors.New("erro interno no servidor")
	ErrHttp400              = errors.New("dados de requisição inválidos")
	ErrHttp401              = errors.New("não autenticado")
	ErrHttp403              = errors.New("acesso não permitido")
	ErrHttp404              = errors.New("recurso não encontrado")
	ErrHttp405              = errors.New("método não permitido")
	ErrHttp409              = errors.New("conflito de recurso")
	ErrHttp422              = errors.New("entidade inutilizável")
	ErrHttp429              = errors.New("muitas requisições")
	ErrHttp500              = errors.New("erro inesperado no servidor")
	ErrHttp503              = errors.New("serviço indisponível")
)
