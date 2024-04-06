package reportedtos

type ResponseReportesRendiciones struct {
	PagoIntentoId           uint64
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

type Totales struct {
	TotalCobrado string
	TotalRendido string
}

type ResponseTotales struct {
	Totales  Totales
	Detalles []ResponseReportesRendiciones
}
