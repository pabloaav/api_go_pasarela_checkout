package ribcradtos

import (
	"fmt"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/administraciondtos"
)

type RICuentasClienteRequest struct {
	FechaInicio    time.Time
	FechaFin       time.Time
	Ruta           string
	CbuCuentaTelco string
}

func (r *RICuentasClienteRequest) IsValid() (erro error) {

	if r.FechaInicio.IsZero() {
		erro = fmt.Errorf(administraciondtos.ERROR_FECHA_INICIO_INVALIDA)
		return
	}

	if !r.FechaFin.IsZero() && r.FechaInicio.After(r.FechaFin) {
		erro = fmt.Errorf(administraciondtos.ERROR_FECHA_INICIO_INVALIDA)
		return
	}
	if r.FechaFin.IsZero() {
		erro = fmt.Errorf(administraciondtos.ERROR_FECHA_FIN_INVALIDA)
		return
	}

	if len(strings.TrimSpace(r.Ruta)) < 1 {
		erro = fmt.Errorf(administraciondtos.ERROR_RUTA_INVALIDA)
		return
	}

	return
}
