package administraciondtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
)

type ResponsePagos struct {
	Pagos           []ResponsePago `json:"data"`
	SaldoPendiente  entities.Monto `json:"saldo_pendiente"`
	SaldoDisponible entities.Monto `json:"saldo_disponible"`
	Meta            dtos.Meta      `json:"meta"`
}

type ResponsePago struct {
	Identificador       uint           `json:"identificador"`
	Cuenta              string         `json:"cuenta"`
	Pagotipo            string         `json:"pagotipo"`
	Fecha               time.Time      `json:"fecha"`
	ExternalReference   string         `json:"external_reference"`
	PayerName           string         `json:"payer_name"`
	Estado              string         `json:"estado"`
	NombreEstado        string         `json:"nombre_estado"`
	Amount              entities.Monto `json:"amount"`
	FechaPago           time.Time      `json:"fecha_pago"`
	Channel             string         `json:"channel"`
	NombreChannel       string         `json:"nombre_channel"`
	UltimoPagoIntentoId uint64         `json:"ultimo_pago_intento_id"`
	TransferenciaId     uint64         `json:"transferencia_id"`
	FechaTransferencia  string         `json:"fecha_transferencia"`
	ReferenciaBancaria  string         `json:"referencia_bancaria"`
	PagoItems           []PagoItems    `json:"pago_items"`
}

type PagoItems struct {
	Descripcion   string
	Identificador string
	Cantidad      int64
	Monto         float64
}
