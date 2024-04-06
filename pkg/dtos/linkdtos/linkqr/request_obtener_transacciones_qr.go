package linkqr

import "time"

type RequestGetTransaccionesQr struct {
	FechaDesde      time.Time `json:"fecha_desde"`
	FechaHasta      time.Time `json:"fecha_hasta"`
	CodigoSucursal  string    `json:"codigo_sucursal"`
	CodigoPos       string    `json:"codigo_pos"`
	EstadoOperacion string    `json:"estado_operacion"`
	Sort            string    `json:"sort"`
	Order           string    `json:"order"`
	Items           int64     `json:"items"`
	Pagina          int64     `json:"pagina"`
}
