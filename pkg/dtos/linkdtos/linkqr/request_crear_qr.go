package linkqr

type Cuenta struct {
	Tipo   string `json:"tipo"`
	Numero string `json:"numero"`
}

type RequestApilinkCrearQr struct {
	/* Cuenta                 Cuenta  `json:"cuenta"` */
	CodigoActividadMCC     string  `json:"codigo_actividad_mcc"`
	CodigoSucursalComercio string  `json:"codigo_sucursal_comercio"`
	CodigoPOSComercio      string  `json:"codigo_pos_comercio"`
	Monto                  float64 `json:"monto"`
}
