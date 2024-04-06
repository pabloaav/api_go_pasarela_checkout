package mockrepository

import (
	"context"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/qrcierrelotesdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/filtros/administracion"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) BeginTx()    {}
func (m *MockRepository) CommitTx()   {}
func (m *MockRepository) RollbackTx() {}

func (m *MockRepository) CreatePago(ctx context.Context, pago *entities.Pago) (*entities.Pago, error) {
	args := m.Called(ctx, pago)
	result := args.Get(0)
	return result.(*entities.Pago), args.Error(1)
}
func (m *MockRepository) UpdatePago(ctx context.Context, pago *entities.Pago) (bool, error) {
	args := m.Called(pago)
	result := args.Bool(0)
	return result, args.Error(1)
}
func (m *MockRepository) GetPaymentByUuid(filtroPago filtros.PagoFiltro) (*entities.Pago, error) {
	args := m.Called(filtroPago.Uuids[0])
	result := args.Get(0)
	return result.(*entities.Pago), args.Error(1)
}
func (m *MockRepository) GetClienteByApikey(apikey string) (*entities.Cliente, error) {
	args := m.Called(apikey)
	result := args.Get(0)
	return result.(*entities.Cliente), args.Error(1)
}
func (m *MockRepository) GetCuentaByApikey(apikey string) (*entities.Cuenta, error) {
	args := m.Called(apikey)
	result := args.Get(0)
	return result.(*entities.Cuenta), args.Error(1)
}
func (m *MockRepository) GetPagotipoById(id int64) (*entities.Pagotipo, error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(*entities.Pagotipo), args.Error(1)
}
func (m *MockRepository) GetPagotipoChannelByPagotipoId(id int64) (*[]entities.Pagotipochannel, error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(*[]entities.Pagotipochannel), args.Error(1)
}
func (m *MockRepository) GetPagotipoIntallmentByPagotipoId(id int64) (*[]entities.Pagotipointallment, error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(*[]entities.Pagotipointallment), args.Error(1)
}
func (m *MockRepository) GetChannelByName(nombre string) (*entities.Channel, error) {
	args := m.Called(nombre)
	result := args.Get(0)
	return result.(*entities.Channel), args.Error(1)
}
func (m *MockRepository) GetCuentaById(id int64) (*entities.Cuenta, error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(*entities.Cuenta), args.Error(1)
}
func (m *MockRepository) CreateResultado(ctx context.Context, resultado *entities.Pagointento) (bool, error) {
	args := m.Called(ctx, resultado)
	result := args.Bool(0)
	return result, args.Error(1)
}
func (m *MockRepository) GetValidPagointentoByPagoId(pagoId int64) (*entities.Pagointento, error) {
	args := m.Called(pagoId)
	result := args.Get(0)
	return result.(*entities.Pagointento), args.Error(1)
}
func (m *MockRepository) GetMediosDePagos() (*[]entities.Mediopago, error) {
	args := m.Called()
	result := args.Get(0)
	return result.(*[]entities.Mediopago), args.Error(1)
}
func (m *MockRepository) GetMediopago(filtro map[string]interface{}) (*entities.Mediopago, error) {
	args := m.Called(filtro)
	result := args.Get(0)
	return result.(*entities.Mediopago), args.Error(1)
}
func (m *MockRepository) GetInstallmentDetailsID(installmentID, numeroCuota int64) int64 {
	args := m.Called(installmentID, numeroCuota)
	result := args.Get(0)
	return result.(int64)
}
func (m *MockRepository) GetInstallmentDetails(installmentID, numeroCuota int64) (installmentDetails *dtos.InstallmentDetailsResponse, erro error) {
	args := m.Called(installmentID, numeroCuota)
	result := args.Get(0)
	return result.(*dtos.InstallmentDetailsResponse), args.Error(1)

}
func (m *MockRepository) CreatePagoEstadoLog(ctx context.Context, pel *entities.Pagoestadologs) error {
	args := m.Called(ctx, pel)
	return args.Error(0)
}

func (m *MockRepository) GetInstallmentsByMedioPagoInstallmentsId(id int64) (installments []entities.Installment, erro error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.([]entities.Installment), args.Error(1)
}

func (m *MockRepository) SaveHasheado(hasheado *entities.Uuid, pagointento_id uint) (erro error) {
	args := m.Called(hasheado)
	return args.Error(1)
}

func (m *MockRepository) GetHasheado(hash string) (control bool, erro error) {
	args := m.Called(hash)
	result := args.Get(0)
	return result.(bool), args.Error(1)
}

func (m *MockRepository) GetPagoEstado(id int64) (*entities.Pagoestado, error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(*entities.Pagoestado), args.Error(1)
}
func (m *MockRepository) GetPreferencesByIdClienteRepository(id uint) (preferencia entities.Preference, erro error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(entities.Preference), args.Error(1)

}

func (m *MockRepository) GetChannelById(id uint) (channel entities.Channel, erro error) {
	args := m.Called(id)
	result := args.Get(0)
	return result.(entities.Channel), args.Error(1)

}

func (m *MockRepository) UpdateEstadoNotificadoInicial(id uint) error {
	return nil
}

func (m *MockRepository) CreateQrcierrelotesRepository(ctx context.Context, requestQrcierrelote *qrcierrelotesdtos.RequestCreateQrCierrelotes) error {
	return nil
}

func (m *MockRepository) GetPagosRepository(filtro filtros.PagoFiltro) (pagos []entities.Pago, erro error) {
	return
}

func (m *MockRepository) GetPagointentosRepository(filtro filtros.PagoIntentoFiltro) (pagointentos []entities.Pagointento, err error) {
	return
}

func (m *MockRepository) GetPagosEstadosExternos(filtro filtros.PagoEstadoExternoFiltro) (estados []entities.Pagoestadoexterno, erro error) {
	return
}

func (m *MockRepository) GetPagosRapipago(filtro filtros.RapipagosFiltro) (pagos []entities.Pago, err error) {
	args := m.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pago), args.Error(1)
}

func (m *MockRepository) GetPagosRapipagoAprobados(filtro filtros.RapipagosFiltro) (pagos []entities.Pago, err error) {
	args := m.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pago), args.Error(1)
}

func (m *MockRepository) UpdatePagoEstado(ctx context.Context, pago entities.Pago) (bool, error) {
	args := m.Called(ctx, pago)
	result := args.Get(0)
	return result.(bool), args.Error(1)
}

func (m *MockRepository) UpdatePagoIntento(ctx context.Context, pagointento entities.Pagointento) (bool, error) {
	args := m.Called(ctx, pagointento)
	result := args.Get(0)
	return result.(bool), args.Error(1)
}

func (m *MockRepository) GetPaymentMultipago(filtro filtros.FiltroMultipago) (pagos []entities.Pago, erro error) {
	args := m.Called(filtro)
	result := args.Get(0)
	return result.([]entities.Pago), args.Error(1)
}

func (m *MockRepository) UpdatePagoMP(ctx context.Context, pago *entities.Pago, pi *entities.Pagointento) (bool, error) {
	args := m.Called(ctx, pago, pi)
	result := args.Get(0)
	return result.(bool), args.Error(1)
}

func (m *MockRepository) GetApikeyAdquiriente(apikey string, adquiriente string) (control bool, err error) {
	args := m.Called(apikey, adquiriente)
	result := args.Get(0)
	return result.(bool), args.Error(1)
}

func (m *MockRepository) UpdateEstadoNotificadoOnline(id uint) error {
	return nil
}

func (m *MockRepository) GetPayments(request filtros.PagoFiltro) (pagos []entities.Pago, erro error) {
	args := m.Called(request)
	result := args.Get(0)
	return result.([]entities.Pago), args.Error(1)
}

func (m *MockRepository) GetEstadoAplicacionRepository(filtro filtros.ConfiguracionFiltro) (configuracion entities.Configuracione, err error) {
	args := m.Called(filtro)
	result := args.Get(0)
	return result.(entities.Configuracione), args.Error(1)
}

func (m *MockRepository) GetCuentaByApiKey(apikey string) (cuenta *entities.Cuenta, erro error) {
	args := m.Called(apikey)
	result := args.Get(0)
	return result.(*entities.Cuenta), args.Error(1)
}
