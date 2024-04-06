package entities

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/enumsdtos"
	"gorm.io/gorm"
)

type Webservicespeticione struct {
	gorm.Model
	Operacion string
	Vendor    enumsdtos.EnumVendor
}
