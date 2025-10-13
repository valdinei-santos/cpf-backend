package vo

import "fmt"

type BloqueadoCliente bool

func NewBloqueadoCliente(bloq bool) (BloqueadoCliente, error) {
	return BloqueadoCliente(bloq), nil
}

func (b BloqueadoCliente) Bool() bool {
	return b.Bool()
}

func (b BloqueadoCliente) String() string {
	return fmt.Sprintf("%t", b)
}
