package linkqr

type RequestAltaPos struct {
	SucursalCodigo string `json:"sucursal_codigo"`
	FechaVigencia  string `json:"fecha_vigencia"`
	PosCodigo      string `json:"pos_codigo"`
	QrTipo         int64  `json:"qr_tipo"`
}

type RequestPutPos struct {
	SucursalCodigo string `json:"sucursal_codigo"`
	PosCodigo      string `json:"pos_codigo"`
	FechaVigencia  string `json:"fecha_vigencia"`
	QrTipo         int64  `json:"qr_tipo"`
}

type RequestDeletePos struct {
	SucursalCodigo string `json:"sucursal_codigo"`
	PosCodigo      string `json:"pos_codigo"`
	FechaVigencia  string `json:"fecha_vigencia"`
}

type RequestGetPos struct {
	SucursalId     int64  `json:"sucursal_id"`     // ID que identifica a la Sucursal.
	CodigoSucursal string `json:"codigo_sucursal"` // Código que identifica a la sucursal.
	PosId          string `json:"pos_id"`
	CodigoPos      string `json:"codigo_pos"`
	Sort           string `json:"sort"`   // Campo por el cual se desea ordenar la consulta.
	Order          string `json:"order"`  // Direccion para el ordenamiento de la consulta.
	Items          string `json:"items"`  // Cantidad de ítems por página.
	Pagina         string `json:"pagina"` // Numero de página a retornar.
}
