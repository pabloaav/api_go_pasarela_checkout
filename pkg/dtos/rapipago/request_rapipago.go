package rapipago

import "errors"

type RequestRapipagoConsulta struct {
	IdClave            string            `json:"id_clave"`
	Nombre             string            `json:"nombre"`
	Apellido           string            `json:"apellido"`
	CodTrx             string            `json:"cod_trx"`
	Canal              string            `json:"canal"`
	FechaHoraOperacion string            `json:"fecha_hora_operacion"`
	CodigoRespuesta    string            `json:"codigo_respuesta"`
	Msg                string            `json:"msg"`
	DatoAdicional      string            `json:"dato_adicional"`
	Facturas           []FacturaRapipago `json:"facturas"`
	Barra              string            `json:"barra"`
	Payment            string            `json:"payment"`
	IdNumero           string            `json:"id_numero"`
}

type FacturaRapipago struct {
	IdNumero         string `json:"id_numero"`
	FechaVencimiento string `json:"fecha_vencimiento"`
	FechaEmision     string `json:"fecha_emision"`
	Importe          string `json:"importe"`
	Barra            string `json:"barra"`
	Texto1           string `json:"texto1"`
}

type RequestRapipagoConfirmacion struct {
	IdNumero           string `json:"id_numero"`
	CodTrx             string `json:"cod_trx"`
	CodTrxOriginal     string `json:"cod_trx_original"`
	Canal              string `json:"canal"`
	Importe            string `json:"importe"`
	Barra              string `json:"barra"`
	FechaHoraOperacion string `json:"fecha_hora_operacion"`
}

func (rpc RequestRapipagoConfirmacion) Validate() (erro error) {
	if len(rpc.Barra) == 0 {
		erro = errors.New("error confirmacion. codigo de barra nulo")
		return
	}
	return
}

func (rrc RequestRapipagoConfirmacion) ParseToResponse(codigo, mensage string) (response ResponseRapipagoConfirmacion) {
	response.IdNumero = rrc.IdNumero
	response.CodTrx = rrc.CodTrx
	response.Barra = rrc.Barra
	response.FechaHoraOperacion = rrc.FechaHoraOperacion
	response.CodigoRespuesta = codigo
	response.Msg = mensage
	return
}

func (rrc RequestRapipagoConsulta) ParseToResponse(codigo, mensage string) (response ResponseRapipagoImputacion) {
	response.IdNumero = rrc.IdNumero
	response.CodTrx = rrc.CodTrx
	response.Barra = rrc.Barra
	response.CodigoRespuesta = codigo
	response.Msg = mensage
	return
}
