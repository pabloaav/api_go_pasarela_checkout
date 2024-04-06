package reportedtos

type ResponseClientesReportes struct {
	Clientes       string
	Email          string
	Fecha          string
	Pagos          []PagosReportes
	Rendiciones    []Rendiciones
	Reversiones    []Reversiones
	TotalCobrado   string
	RendicionTotal string
	TipoReporte    string
}

// type FactoryEmail struct {
// 	Clientes    string
// 	Email       string
// 	Fecha       string
// 	Pagos       []PagosReportes
// 	Rendiciones []Rendiciones
// }

type PagosReportes struct {
	Cuenta    string
	Id        string
	MedioPago string
	Estado    string
	Cuotas    string
	Monto     string
}

type Rendiciones struct {
	Cuenta                  string // Nombre de la cuenta del cliente
	Id                      string // external_reference enviada por el cliente
	FechaCobro              string // fecha que el pagador realizo el pago
	FechaDeposito           string // fecha que se le envio el dinero al cliente(transferencia)
	ImporteCobrado          string // importe solicitud de pago
	ImporteDepositado       string // importe depositado al cliente
	CantidadBoletasCobradas string // pago items
	// ComisionPorcentaje      string // comision de telco cobrada al cliente
	// ComisionIva             string // iva Cobrado al cliente
	Comision string // comision de telco cobrada al cliente
	Iva      string // iva Cobrado al cliente
}

type Reversiones struct {
	Cuenta    string
	Id        string
	MedioPago string
	Monto     string
}
