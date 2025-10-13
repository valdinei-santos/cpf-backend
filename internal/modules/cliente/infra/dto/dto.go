package dto

// Request -
type Request struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Documento string `json:"documento"`
	Telefone  string `json:"telefone"`
	Bloqueado bool   `json:"bloquaeado"`
}

// Response -
type Response struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Documento string `json:"documento"`
	Telefone  string `json:"telefone"`
	Bloqueado bool   `json:"bloquaeado"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ResponseManyPaginated struct {
	Clientes     []Response `json:"clientes"`
	TotalItems   int64      `json:"totalItems"`
	TotalPages   int64      `json:"totalPages"`
	CurrentPage  int64      `json:"currentPage"`
	ItemsPerPage int64      `json:"itemsPerPage"`
}

// ResponseMany -
type ResponseMany struct {
	Products []Response `json:"clientes"`
}

// OutputDefault - Struct com a resposta da API
type OutputDefault struct {
	Title    string  `json:"title"`
	Detail   string  `json:"detail"`
	Instance *string `json:"instance,omitempty"`
}
