package cierrelotedtos

import (
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
)

type ResponseMovimientoTotales struct {
	Id                 int64
	Empresa            string
	FechaPresentacion  time.Time
	TipoRegistro       string
	ComercioNro        string
	EstablecimientoNro string
	Codop              string
	TipoAplicacion     string
	FechaPago          time.Time
	ImporteTotal       entities.Monto
	SignoImporteTotal  string
	DetalleMovimientos []ResponseMoviminetoDetalles
}

func (rmt *ResponseMovimientoTotales) EntityToDtos(entityMovimientoTotales entities.Prismamovimientototale) {
	rmt.Id = 0
	if entityMovimientoTotales.ID > 0 {
		rmt.Id = int64(entityMovimientoTotales.ID)
	}
	rmt.Empresa = entityMovimientoTotales.Empresa
	rmt.FechaPresentacion = entityMovimientoTotales.FechaPresentacion
	rmt.TipoRegistro = entityMovimientoTotales.TipoRegistro
	rmt.ComercioNro = entityMovimientoTotales.ComercioNro
	rmt.EstablecimientoNro = entityMovimientoTotales.EstablecimientoNro
	rmt.Codop = entityMovimientoTotales.Codop
	rmt.TipoAplicacion = entityMovimientoTotales.TipoAplicacion
	rmt.FechaPago = entityMovimientoTotales.FechaPago
	rmt.ImporteTotal = entityMovimientoTotales.ImporteTotal
	rmt.SignoImporteTotal = entityMovimientoTotales.SignoImporteTotal
}
