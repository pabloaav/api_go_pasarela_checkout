package linkqr

type ResponseSucursalQr struct {
	Status     string `json:"status"`
	ReturnCode string `json:"return_code"`
	Message    string `json:"message"`
	Data       struct {
		ComercioID int `json:"comercio_id"`
		SucursalID int `json:"sucursal_id"`
	} `json:"data"`
}

type Sucursal struct {
	AceptadorID             int    `json:"aceptador_id"`
	CuitAceptador           string `json:"cuit_aceptador"`
	ComercioID              int    `json:"comercio_id"`
	CuitComercio            string `json:"cuit_comercio"`
	SucursalID              int    `json:"sucursal_id"`
	CodigoSucursalComercio  string `json:"codigo_sucursal_comercio"`
	DescripcionSucursal     string `json:"descripcion_sucursal"`
	SucursalTipoID          int    `json:"sucursal_tipo_id"`
	SucursalTipoDescripcion string `json:"sucursal_tipo_descripcion"`
	CuentaClave             string `json:"cuenta_clave"`
	CuentaTipo              string `json:"cuenta_tipo"`
	DomicilioCalle          string `json:"domicilio_calle"`
	DomicilioNumero         string `json:"domicilio_numero"`
	DomicilioPiso           string `json:"domicilio_piso"`
	DomicilioDepartamento   string `json:"domicilio_departamento"`
	DomicilioProvinciaID    int    `json:"domicilio_provincia_id"`
	DomicilioProvincia      string `json:"domicilio_provincia"`
	DomicilioLocalidad      string `json:"domicilio_localidad"`
	DomicilioCodigoPostal   string `json:"domicilio_codigo_postal"`
	DomicilioLatitud        string `json:"domicilio_latitud"`
	DomicilioLongitud       string `json:"domicilio_longitud"`
	FechaVigencia           string `json:"fecha_vigencia"`
	FechaActivacion         string `json:"fecha_activacion"`
	FechaDesactivacion      string `json:"fecha_desactivacion"`
	UsuarioDesactivacion    string `json:"usuario_desactivacion"`
	FechaAlta               string `json:"fecha_alta"`
	UsuarioAlta             string `json:"usuario_alta"`
	FechaBaja               string `json:"fecha_baja"`
	UsuarioBaja             string `json:"usuario_baja"`
	FechaModificacion       string `json:"fecha_modificacion"`
	UsuarioModificacion     string `json:"usuario_modificacion"`
}

type ResponseSucursalesQrApilink struct {
	Status     string `json:"status"`
	ReturnCode string `json:"return_code"`
	Message    string `json:"message"`
	Data       struct {
		Sucursales []Sucursal `json:"sucursales"`
	} `json:"data"`
}

type ResponsePutDeleteSucursalQrApilink struct {
	Status     string `json:"status"`
	ReturnCode string `json:"return_code"`
	Message    string `json:"message"`
	Data       struct {
		ComercioID int `json:"comercio_id"`
		SucursalID int `json:"sucursal_id"`
	} `json:"data"`
}
