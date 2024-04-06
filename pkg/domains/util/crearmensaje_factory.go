package util

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/utildtos"
)

const (
	EMAIL_TEMPLATE = 1
	EMAIL_ADJUNTO  = 2
)

type CrearMensajeFactory interface {
	GetCrearMensajeMethod(m int) (CrearMensajeMethod, error)
}
type crearMensajeFactory struct{}

func NewCrearMensajeFactory() CrearMensajeFactory {
	return &crearMensajeFactory{}
}
func (cmf *crearMensajeFactory) GetCrearMensajeMethod(m int) (CrearMensajeMethod, error) {
	switch m {
	case EMAIL_TEMPLATE:
		return NewEmailTemplateCrearMensaje(), nil
	// case EMAIL_ADJUNTO:
	// 	return NewEmailAdjuntoCrearMensaje(), nil
	default:
		return nil, fmt.Errorf("no se reconoce el metodo de creacion de mensaje %d", m)
	}
}

type CrearMensajeMethod interface {
	MensajeResultado(subject string, to []string, params utildtos.RequestDatosMail) (mensaje string, erro error)
}
