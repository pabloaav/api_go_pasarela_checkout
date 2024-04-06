package linkqr

import "time"

type Devolucion struct {
	OperacionID           int       `json:"operacion_id"`
	FechaAlta             time.Time `json:"fecha_alta"`
	MontoBrutoTransaccion float64   `json:"monto_bruto_transaccion"`
}

type Data struct {
	OperacionID               int          `json:"operacion_id"`
	CuitAceptador             string       `json:"cuit_aceptador"`
	CuitComercio              string       `json:"cuit_comercio"`
	CodigoSucursalComercio    string       `json:"codigo_sucursal_comercio"`
	CodigoPOSComercio         string       `json:"codigo_pos_comercio"`
	FechaAlta                 time.Time    `json:"fecha_alta"`
	OperacionTipoID           int          `json:"operacion_tipo_id"`
	OperacionTipoCodigo       string       `json:"operacion_tipo_codigo"`
	MontoBrutoTransaccion     float64      `json:"monto_bruto_transaccion"`
	MontoBrutoRetencion       float64      `json:"monto_bruto_retencion"`
	MontoBrutoComision        float64      `json:"monto_bruto_comision"`
	MontoNeto                 float64      `json:"monto_neto"`
	OperacionEstadoID         int          `json:"operacion_estado_id"`
	OperacionEstadoDesc       string       `json:"operacion_estado_desc"`
	PagadorNombre             string       `json:"pagador_nombre"`
	PagadorDocumentoTipo      string       `json:"pagador_documento_tipo"`
	PagadorDocumentoNumero    int          `json:"pagador_documento_numero"`
	BilleteraNombre           string       `json:"billetera_nombre"`
	BilleteraCuit             int          `json:"billetera_cuit"`
	RechazoCodigo             string       `json:"rechazo_codigo"`
	RechazoDescripcion        string       `json:"rechazo_descripcion"`
	OperacionIDOriginal       int          `json:"operacion_id_original"`
	FechaPagoOrigen           time.Time    `json:"fecha_pago_origen"`
	FlagDevolucionesAprobadas bool         `json:"flag_devoluciones_aprobadas"`
	DevolucionesAsociadas     []Devolucion `json:"devoluciones_asociadas"`
}

type ResponseTransaccion struct {
	Status     string `json:"status"`
	ReturnCode string `json:"return_code"`
	Message    string `json:"message"`
	Data       Data   `json:"data"`
}
