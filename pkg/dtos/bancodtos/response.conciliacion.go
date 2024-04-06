package bancodtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
)

type ResponseConciliacion struct {
	Transferencias []MovimientosTransferenciasResponse
	ListaApilink   []*entities.Apilinkcierrelote
	ListaRapipago  []*entities.Rapipagocierrelote
}

type MovimientosTransferenciasResponse struct {
	Id              uint `json:"id"`
	Match           int  `json:"match"`
	BancoExternalId int  `json:"banco_external_id"`
}
