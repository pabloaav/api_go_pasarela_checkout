package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/linkdtos/linktransferencia"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
)

type RequestTransferenciaAutomatica struct {
	CuentaId      uint64
	Cuenta        string
	DatosClientes DatosClientes
	Request       RequestTransferenicaCliente
}
type RequestTransferenicaCliente struct {
	Transferencia         linktransferencia.RequestTransferenciaCreateLink `json:"transferencia,omitempty"`
	ListaMovimientosId    []uint64                                         `json:"lista_movimientos_id,omitempty"`
	ListaMovimientosIdNeg []uint64                                         `json:"lista_movimientos_id_neg,omitempty"`
}

type ResponseTransferenciaAutomatica struct {
	CuentaId uint64         `json:"cuentaid"`
	Cuenta   string         `json:"cuenta"`
	Origen   string         `json:"origen"`
	Destino  string         `json:"destino"`
	Importe  entities.Monto `json:"importe"`
	Error    string         `json:"error"`
}

type DatosClientes struct {
	NombreCliente string
	EmailCliente  string
}
