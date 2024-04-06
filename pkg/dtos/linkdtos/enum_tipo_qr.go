package linkdtos

import "errors"

type EnumTipoQr string

const (
	QrDefault  EnumTipoQr = "QR"
	QrEstatico EnumTipoQr = "QR-ESTATICO"
)

func (e EnumTipoQr) IsValid() error {
	switch e {
	case QrDefault, QrEstatico:
		return nil
	}
	return errors.New("tipo EnumTipoQr con formato inválido")
}
