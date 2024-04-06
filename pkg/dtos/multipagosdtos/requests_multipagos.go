package multipagosdtos

type RequestPagoMultipago struct {
	FechaHoraOperacion string `json:"fecha_hora_operacion"`
	Id_clave           string `json:"id_clave"`
	Importe            string `json:"importe"`
	CodTrx             string `json:"cod_trx"`
	Canal              string `json:"canal"`
}

type RequestConsultaMultipago struct {
	Id_clave string `json:"id_clave"`
	CodTrx   string `json:"cod_trx"`
}

type RequestControlMultipago struct {
	FechaInicio string `json:"fecha_inicio"`
	FechaFin    string `json:"fecha_fin"`
}
