package administraciondtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/linkdtos/linktransferencia"

type RequestTransferenciasComisiones struct {
	Transferencia           linktransferencia.RequestTransferenciaCreateLink `json:"transferencia,omitempty"`
	MovimientosIdComisiones []uint64                                         `json:"movimientos_id_comisiones"`
}

type RequestMovimientosId struct {
	MovimientosId []uint64 `json:"movimientos_id"`
}
