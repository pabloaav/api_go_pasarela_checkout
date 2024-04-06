package reportedtos

import (
	"time"

	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/filtros/administracion"
)

type RequestPagosPeriodo struct {
	Paginacion                      filtros.Paginacion
	ClienteId                       uint64
	FechaInicio                     time.Time
	FechaFin                        time.Time
	PagoIntento                     uint64
	PagoIntentos                    []uint64
	TipoMovimiento                  string
	CargarComisionImpuesto          bool
	CargarMovimientosTransferencias bool
	CargarPagoIntentos              bool
	CargarCuenta                    bool
	Number                          uint32
	Size                            uint32
}
