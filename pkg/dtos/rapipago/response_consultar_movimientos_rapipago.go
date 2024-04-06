package rapipago

type ResponseConsultarMovimientosRapipago struct {
	FechaCobro      string
	Importe         int64
	CodigoBarras    string
	BancoExternalId int64
	Match           bool
}
