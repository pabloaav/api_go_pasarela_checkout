package utildtos

import "errors"

type EnumTipoEmail string

const (
	Template EnumTipoEmail = "template"
	Adjunto  EnumTipoEmail = "adjunto"
)

func (e EnumTipoEmail) IsValid() (int, error) {
	switch e {
	case Template:
		return 1, nil
	case Adjunto:
		return 2, nil
	}
	return 0, errors.New("el tipo de correo no es valido")
}
