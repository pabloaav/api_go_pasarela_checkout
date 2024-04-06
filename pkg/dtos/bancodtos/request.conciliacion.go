package bancodtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
)

type RequestConciliacion struct {
	Transferencias   administraciondtos.TransferenciaResponsePaginado
	ListaApilink     []*entities.Apilinkcierrelote
	ListaRapipago    []*entities.Rapipagocierrelote
	TipoConciliacion int64
}
