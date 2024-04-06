package linkqr

type RequestSucursalQr struct {
	Calle          string `json:"calle"`
	CodigoPostal   string `json:"codigo_postal"`
	CuentaClave    string `json:"cuenta_clave"`
	CuentaTipo     string `json:"cuenta_tipo"`
	Departamento   string `json:"departamento"`
	Descripcion    string `json:"descripcion"`
	FechaVigencia  string `json:"fecha_vigencia"`
	Latitud        string `json:"latitud"`
	Localidad      string `json:"localidad"`
	Longitud       string `json:"longitud"`
	Numero         string `json:"numero"`
	Piso           string `json:"piso"`
	Provincia      string `json:"provincia"`
	SucursalCodigo string `json:"sucursal_codigo"`
	TipoSucursal   string `json:"tipo_sucursal"`
}

// Request para obtener sucursales de un comercio ().
type RequestGetSucursales struct {
	SucursalId     int64  `json:"sucursal_id"`     // ID que identifica a la Sucursal.
	CodigoSucursal string `json:"codigo_sucursal"` // Código que identifica a la sucursal.
	Sort           string `json:"sort"`            // Campo por el cual se desea ordenar la consulta. Valores permitidos --> sucursal_id ; fecha_alta ; fecha_baja ; domiclio_codigo_postal. Valor default --> sucursal_id
	Order          string `json:"order"`           // Direccion para el ordenamiento de la consulta. Valores permitidos --> ASC (ascendiente) ; DESC (descendiente). Valor default --> DESC
	Items          string `json:"items"`           // Cantidad de ítems por página.
	Pagina         string `json:"pagina"`          // Numero de página a retornar.
}

type RequestDeleteSucursal struct {
	FechaVigencia string `json:"fecha_vigencia"` // Fecha de vigencia para la baja. Debe ser posterior a la fecha actual. Formato --> YYYY-MM-DD
}
