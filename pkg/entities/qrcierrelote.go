package entities

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/linkdtos"
	"gorm.io/gorm"
)

type Qrcierrelote struct {
	gorm.Model
	Uuid                 string                     `json:"uuid"`
	QrId                 string                     `json:"qr_id"`
	Qr                   string                     `json:"qr"`
	Concepto             linkdtos.EnumConceptoDebin `json:"concepto"`
	Moneda               linkdtos.EnumMoneda        `json:"moneda"`
	Importe              Monto                      `json:"importe"`
	Estado               linkdtos.EnumEstadoQr      `json:"estado"`
	Tipo                 linkdtos.EnumTipoQr        `json:"tipo"`
	FechaExpiracion      time.Time                  `json:"fecha_expiracion"`
	ContracargoId        string                     `json:"contracargo_id"`
	CompradorCuit        string                     `json:"comprador_cuit"`
	VendedorCuit         string                     `json:"vendedor_cuit"`
	ReferenciaBanco      string                     `json:"referencia_banco"`
	Match                int                        `json:"match"`
	BancoExternalId      int                        `json:"banco_external_id"`
	Pagoinformado        bool                       `json:"pagoinformado"`
	Fechaacreditacion    time.Time                  `json:"fechaacreditacion"`
	FechaCobro           time.Time                  `json:"fecha_cobro"`
	PagoestadoexternosId uint64                     `json:"pagoestadoexternos_id"`
	Pagoestadoexterno    Pagoestadoexterno          `json:"pagoestadoexterno" gorm:"foreignKey:pagoestadoexternos_id"`
}
