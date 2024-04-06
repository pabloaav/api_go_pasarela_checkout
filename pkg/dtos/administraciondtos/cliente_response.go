package administraciondtos

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	"gorm.io/gorm"
)

type ResponseFacturacionPaginado struct {
	Clientes []ResponseFacturacion `json:"data"`
	Meta     dtos.Meta             `json:"meta"`
}

type ResponseFacturacion struct {
	Id             uint               `json:"id"`      // Id cliente
	Cliente        string             `json:"cliente"` // nombre Cliente abreviado
	RazonSocial    string             `json:"razon_social"`
	NombreFantasia string             `json:"nombre_fantasia"`
	Email          string             `json:"email"`
	Cuit           string             `json:"cuit"`
	Personeria     string             `json:"personeria"`
	Impuestos      []ResponseImpuesto `json:"impuestos"`
	// Iva              ResponseFacturacionIva      `json:"iva"`
	// Iibb             ResponseFacturacionIibb     `json:"iibb"`
	RetiroAutomatico bool                        `json:"retiro_automatico"`
	ReporteBatch     bool                        `json:"reporte_batch"`
	NombreReporte    string                      `json:"nombre_reporte"`
	Cuenta           []ResponseFacturacionCuenta `json:"cuenta"`
}

func (r *ResponseFacturacion) FromEntity(c entities.Cliente) {
	r.Id = c.ID
	r.Cliente = c.Cliente
	r.RazonSocial = c.Razonsocial
	r.NombreFantasia = c.Nombrefantasia
	r.Email = c.Email
	r.Cuit = c.Cuit
	r.Personeria = c.Personeria
	r.RetiroAutomatico = c.RetiroAutomatico
	r.ReporteBatch = c.ReporteBatch
	r.NombreReporte = c.NombreReporte
	if c.Iva != nil {
		// var iva ResponseFacturacionIva
		// iva.FromEntity(*c.Iva)
		// r.Iva = iva
		var impuestoiva ResponseImpuesto
		iva := entities.Impuesto{
			Model:      gorm.Model{ID: c.Iva.ID},
			Impuesto:   c.Iva.Impuesto,
			Porcentaje: c.Iva.Porcentaje,
			Tipo:       c.Iva.Tipo,
			Fechadesde: c.Iva.Fechadesde,
		}
		impuestoiva.FromImpuesto(iva)
		r.Impuestos = append(r.Impuestos, impuestoiva)
	}
	if c.Iibb != nil {
		// var iibb ResponseFacturacionIibb
		// iibb.FromEntity(*c.Iibb)
		// r.Iibb = iibb
		var impuestoIibb ResponseImpuesto
		iibb := entities.Impuesto{
			Model:      gorm.Model{ID: c.Iibb.ID},
			Impuesto:   c.Iibb.Impuesto,
			Porcentaje: c.Iibb.Porcentaje,
			Tipo:       c.Iibb.Tipo,
			Fechadesde: c.Iibb.Fechadesde,
		}
		impuestoIibb.FromImpuesto(iibb)
		r.Impuestos = append(r.Impuestos, impuestoIibb)
	}
	if c.Cuentas != nil {
		for _, c := range *c.Cuentas {
			var cuenta ResponseFacturacionCuenta
			cuenta.FromEntity(c)
			r.Cuenta = append(r.Cuenta, cuenta)
		}
	}
}

type ResponseFacturacionIva struct {
	Id         uint    `json:"id"`
	Impuesto   string  `json:"impuesto"`
	Tipo       string  `json:"tipo"`
	Porcentaje float64 `json:"porcentaje"`
}

func (r *ResponseFacturacionIva) FromEntity(c entities.Impuesto) {
	r.Id = c.ID
	r.Impuesto = c.Impuesto
	r.Tipo = c.Tipo
	r.Porcentaje = c.Porcentaje
}

type ResponseFacturacionIibb struct {
	Id         uint    `json:"id"`
	Impuesto   string  `json:"impuesto"`
	Tipo       string  `json:"tipo"`
	Porcentaje float64 `json:"porcentaje"`
}

func (r *ResponseFacturacionIibb) FromEntity(c entities.Impuesto) {
	r.Id = c.ID
	r.Impuesto = c.Impuesto
	r.Tipo = c.Tipo
	r.Porcentaje = c.Porcentaje
}

type ResponseFacturacionCuenta struct {
	Id                   uint                            `json:"id"`
	Cuenta               string                          `json:"cuenta"`
	Cbu                  string                          `json:"cbu"`
	Cvu                  string                          `json:"cvu"`
	DiasRetiroAutomatico int64                           `json:"dias_retiro_automatico"`
	Rubro                ResponseFacturacionRubro        `json:"rubro"`
	Comisiones           []ResponseFacturacionComisiones `json:"comisiones"`
	TiposPago            []ResponseFacturacionTiposPago  `json:"tipos_pago"`
}

func (r *ResponseFacturacionCuenta) FromEntity(c entities.Cuenta) {
	r.Id = c.ID
	r.Cuenta = c.Cuenta
	r.Cbu = c.Cbu
	r.Cvu = c.Cvu
	r.DiasRetiroAutomatico = c.DiasRetiroAutomatico
	if c.Rubro != nil {
		var rubro ResponseFacturacionRubro
		rubro.FromEntity(*c.Rubro)
		r.Rubro = rubro
	}
	if c.Cuentacomisions != nil {
		for _, c := range *c.Cuentacomisions {
			var comisiones ResponseFacturacionComisiones
			comisiones.FromEntity(c)
			r.Comisiones = append(r.Comisiones, comisiones)
		}
	}
	if c.Pagotipos != nil {
		for _, pt := range *c.Pagotipos {
			var pagotipos ResponseFacturacionTiposPago
			pagotipos.FromEntity(pt)
			r.TiposPago = append(r.TiposPago, pagotipos)
		}
	}
}

type ResponseFacturacionRubro struct {
	Id    uint   `json:"id"`
	Rubro string `json:"rubro"`
}

func (r *ResponseFacturacionRubro) FromEntity(c entities.Rubro) {
	r.Id = c.ID
	r.Rubro = c.Rubro
}

type ResponseFacturacionComisiones struct {
	Nombre        string  `json:"nombre"`
	Comision      float64 `json:"comision"`
	VigenciaDesde string  `json:"vigencia_desde"`
	Canal         string  `json:"canal"`
}

func (r *ResponseFacturacionComisiones) FromEntity(c entities.Cuentacomision) {
	r.Nombre = c.Cuentacomision
	if c.ChannelArancel.Tipocalculo == "FIJO" {
		r.Comision = c.Comision
	} else {
		r.Comision = c.Comision + c.ChannelArancel.Importe
	}
	r.VigenciaDesde = c.VigenciaDesde.Format("2006-01-02")
	r.Canal = c.Channel.Nombre
}

type ResponseFacturacionTiposPago struct {
	PagoTipo                 string
	BackUrlSuccess           string
	BackUrlPending           string
	BackUrlRejected          string
	BackUrlNotificacionPagos string
}

func (r *ResponseFacturacionTiposPago) FromEntity(c entities.Pagotipo) {
	r.PagoTipo = c.Pagotipo
	r.BackUrlSuccess = c.BackUrlSuccess
	r.BackUrlPending = c.BackUrlPending
	r.BackUrlRejected = c.BackUrlRejected
	r.BackUrlNotificacionPagos = c.BackUrlNotificacionPagos
}
