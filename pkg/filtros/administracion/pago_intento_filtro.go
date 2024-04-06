package filtros

type PagoIntentoFiltro struct {
	ExternalId              bool
	ExternalIds             []string
	TicketNumber            []string
	CodigoAutorizacion      []string
	TransaccionesId         []string
	Barcode                 []string
	PagosId                 []uint
	Channel                 bool
	CargarPago              bool
	CargarPagoItems         bool
	CargarPagoTipo          bool
	CargarPagoEstado        bool
	CargarMovimientos       bool
	CargarCuenta            bool
	CargarCliente           bool
	CargarCuentaComision    bool
	CargarImpuestos         bool
	CargarInstallmentdetail bool
}
