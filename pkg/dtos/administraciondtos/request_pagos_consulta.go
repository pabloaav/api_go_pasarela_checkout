package administraciondtos

import (
	"fmt"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
)

type RequestPagosConsulta struct {
	Uuid              string   `json:"uuid"`
	ExternalReference string   `json:"external_reference"`
	FechaDesde        string   `json:"fecha_desde"`
	FechaHasta        string   `json:"fecha_hasta"`
	Uuids             []string `json:"uuids"`
}

func (r *RequestPagosConsulta) ToPago() entities.Pago {
	return entities.Pago{
		Uuid:              r.Uuid,
		ExternalReference: r.ExternalReference,
	}
}

func (r *RequestPagosConsulta) IsValid() error {
	if len(r.Uuid)+len(r.ExternalReference)+len(r.FechaDesde)+len(r.Uuids) <= 0 {
		return fmt.Errorf("parámetros de búsqueda insuficientes")
	}
	return nil
}
