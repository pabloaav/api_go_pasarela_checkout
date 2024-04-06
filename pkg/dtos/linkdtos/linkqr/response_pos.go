package linkqr

type ResponseGetPos struct {
	Status     string `json:"status"`
	ReturnCode string `json:"return_code"`
	Message    string `json:"message"`
	Data       struct {
		Pos []Pos `json:"pos"`
	} `json:"data"`
}

type ResponseCreatePutDeletePos struct {
	Status     string `json:"status"`
	ReturnCode string `json:"return_code"`
	Message    string `json:"message"`
	Data       struct {
		ComercioID int `json:"comercio_id"`
		SucursalID int `json:"sucursal_id"`
		PosID      int `json:"pos_id"`
	} `json:"data"`
}

type Pos struct {
	AceptadorID            int    `json:"aceptador_id"`
	CuitAceptador          string `json:"cuit_aceptador"`
	ComercioID             int    `json:"comercio_id"`
	CuitComercio           string `json:"cuit_comercio"`
	SucursalID             int    `json:"sucursal_id"`
	CodigoSucursalComercio string `json:"codigo_sucursal_comercio"`
	PosID                  int    `json:"pos_id"`
	CodigoPosComercio      string `json:"codigo_pos_comercio"`
	QrTipoID               int    `json:"qr_tipo_id"`
	QrTipo                 string `json:"qr_tipo"`
	FechaActivacion        string `json:"fecha_activacion"`
	FechaDesactivacion     string `json:"fecha_desactivacion"`
	UsuarioDesactivacion   string `json:"usuario_desactivacion"`
	FechaAlta              string `json:"fecha_alta"`
	UsuarioAlta            string `json:"usuario_alta"`
	FechaBaja              string `json:"fecha_baja"`
	UsuarioBaja            string `json:"usuario_baja"`
}
