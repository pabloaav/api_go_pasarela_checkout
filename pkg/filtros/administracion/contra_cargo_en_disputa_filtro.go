package filtros

type ContraCargoEnDisputa struct {
	Paginacion
	IdCliente           uint
	IdCuenta            uint
	CargarCuentas       bool
	CargarTiposPago     bool
	CargarPagos         bool
	CargarPagosIntentos bool
	TransactionId       []string
}
