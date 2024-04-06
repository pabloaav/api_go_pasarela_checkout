package rapipago

type ResponseRapipagoConsulta struct {
	IdClave         string            `json:"id_clave"`
	Nombre          string            `json:"nombre"`
	Apellido        string            `json:"apellido"`
	CodTrx          string            `json:"cod_trx"`
	CodigoRespuesta string            `json:"codigo_respuesta"`
	Msg             string            `json:"msg"`
	DatoAdicional   string            `json:"dato_adicional"`
	Facturas        []FacturaRapipago `json:"facturas"`
}

type ResponseRapipagoImputacion struct {
	IdNumero        string `json:"id_numero"`
	CodTrx          string `json:"cod_trx"`
	Barra           string `json:"barra"`
	CodigoRespuesta string `json:"codigo_respuesta"`
	Msg             string `json:"msg"`
}

type ResponseRapipagoConfirmacion struct {
	IdNumero           string `json:"id_numero"`
	CodTrx             string `json:"cod_trx"`
	Barra              string `json:"barra"`
	FechaHoraOperacion string `json:"fecha_hora_operacion"`
	CodigoRespuesta    string `json:"codigo_respuesta"`
	Msg                string `json:"msg"`
}
