package mockservice

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	"github.com/stretchr/testify/mock"
)

type MockCreditPaymentMethod struct {
	mock.Mock
}

func (mk *MockCreditPaymentMethod) CreateResultado(request *dtos.ResultadoRequest, pago *entities.Pago, cuenta *entities.Cuenta) (*entities.Pagointento, error) {
	args := mk.Called(request, pago, cuenta)
	result := args.Get(0)
	return result.(*entities.Pagointento), args.Error(1)
}
