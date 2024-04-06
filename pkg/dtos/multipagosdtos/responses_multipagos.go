package multipagosdtos

type ResponseConsultaMultipago struct {
	CodTrx          string             `json:"cod_trx"`
	Id_clave        string             `json:"id_clave"`
	Nombre          string             `json:"nombre"`
	CodigoRespuesta string             `json:"codigo_respuesta"`
	Msg             string             `json:"msg"`
	Facturas        []FacturaMultipago `json:"facturas"`
}

type FacturaMultipago struct {
	FechaEmision     string `json:"fecha_emision"`
	FechaVencimiento string `json:"fecha_vencimiento"`
	Importe          string `json:"importe"`
	Barra            string `json:"barra"`
}

type ResponsePagoMultipago struct {
	CodTrx          string `json:"cod_trx"`
	Id_clave        string `json:"id_clave"`
	CodigoRespuesta string `json:"codigo_respuesta"`
	Msg             string `json:"msg"`
	CodOperacion    string `json:"cod_operacion"`
}

type ResponseControlMultipago struct {
	CodigoRespuesta string                    `json:"codigo_respuesta"`
	Msg             string                    `json:"msg"`
	Facturas        []FacturaMultipagoControl `json:"facturas"`
}

type FacturaMultipagoControl struct {
	CodOperacion string `json:"cod_operacion"`
	Barra        string `json:"barra"`
}
