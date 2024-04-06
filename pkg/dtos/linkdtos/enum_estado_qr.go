package linkdtos

import (
	"errors"
)

type EnumEstadoQr string

const (
	Aprobada EnumEstadoQr = "APROBADA" // ESTADO provisorio para el QR (hasta tener documentacion de estados)
)

func (e EnumEstadoQr) IsValid() error {
	switch e {
	case Aprobada:
		return nil
	}
	return errors.New("tipo EnumEstadoQr con formato inválido")
}
