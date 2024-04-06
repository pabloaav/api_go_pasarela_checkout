package linkqr

import "time"

type ResponseOrdenPagoQr struct {
	Status     string        `json:"status"`
	ReturnCode string        `json:"return_code"`
	Message    string        `json:"message"`
	Data       RespuestaData `json:"data"`
}

type RespuestaData struct {
	OperacionID string `json:"operacion_id"`
}

type ResponseTransaccionesQr struct {
	Status     string         `json:"status"`
	ReturnCode string         `json:"return_code"`
	Message    string         `json:"message"`
	Data       RespuestaDatos `json:"data"`
}

type RespuestaDatos struct {
	Operaciones []Operacion `json:"operaciones"`
	Paginado    Paginado    `json:"paginado"`
}

type Operacion struct {
	OperacionID               int       `json:"operacion_id"`
	CuitAceptador             string    `json:"cuit_aceptador"`
	CuitComercio              string    `json:"cuit_comercio"`
	CodigoSucursalComercio    string    `json:"codigo_sucursal_comercio"`
	CodigoPosComercio         string    `json:"codigo_pos_comercio"`
	FechaAlta                 time.Time `json:"fecha_alta"`
	OperacionTipoID           int       `json:"operacion_tipo_id"`
	OperacionTipoCodigo       string    `json:"operacion_tipo_codigo"`
	MontoBrutoTransaccion     float64   `json:"monto_bruto_transaccion"`
	MontoBrutoRetencion       float64   `json:"monto_bruto_retencion"`
	MontoBrutoComision        float64   `json:"monto_bruto_comision"`
	MontoNeto                 float64   `json:"monto_neto"`
	OperacionEstadoID         int       `json:"operacion_estado_id"`
	OperacionEstadoDesc       string    `json:"operacion_estado_desc"`
	PagadorNombre             string    `json:"pagador_nombre"`
	PagadorDocumentoTipo      string    `json:"pagador_documento_tipo"`
	PagadorDocumentoNumero    int       `json:"pagador_documento_numero"`
	BilleteraNombre           string    `json:"billetera_nombre"`
	BilleteraCuit             int       `json:"billetera_cuit"`
	RechazoCodigo             string    `json:"rechazo_codigo"`
	RechazoDescripcion        string    `json:"rechazo_descripcion"`
	OperacionIDOriginal       int       `json:"operacion_id_original"`
	FechaPagoOrigen           time.Time `json:"fecha_pago_origen"`
	FlagDevolucionesAprobadas bool      `json:"flag_devoluciones_aprobadas"`
}

type Paginado struct {
	Pagina       int `json:"pagina"`
	TotalPaginas int `json:"total_paginas"`
	TotalItems   int `json:"total_items"`
}
