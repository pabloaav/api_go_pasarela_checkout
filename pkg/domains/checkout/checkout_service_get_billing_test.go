package checkout

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/checkout/services"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/test/mocks/mockrepository"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/test/mocks/mockservice"
	"github.com/stretchr/testify/assert"
)

func TestGetBilling(t *testing.T) {
	assertions := assert.New(t)
	repo := mockrepository.MockRepository{}
	mockUtilService := new(mockservice.MockUtilService)
	mockRepositoryWebHook := new(mockrepository.MockRepositoryWebHook)
	serv := services.NewService(&repo, &mockservice.MockCommonsService{}, &mockservice.MockPrismaService{}, &mockservice.MockPagoOffLineService{}, mockUtilService, mockRepositoryWebHook)

	t.Run("Cuando no se encuentra el pago devuelvo un error", func(t *testing.T) {
		uuidNoValido := "Uuid-de-pago-inexistente"
		repo.On("GetPagoByUuid", uuidNoValido).Return(&entities.Pago{}, fmt.Errorf("no existe pago con identificador Uuid-de-pago-inexistente")).Once()

		_, err := serv.GetBilling(uuidNoValido)

		assertions.EqualError(err, "no existe pago con identificador Uuid-de-pago-inexistente")
	})

	t.Run("Cuando no se encuentran los datos del pago intento correcto se devuelve un error", func(t *testing.T) {
		pagoSinIntentoValido := pagoCreated
		pagoSinIntentoValido.ID = 55
		repo.On("GetPagoByUuid", pagoSinIntentoValido.Uuid).Return(&pagoSinIntentoValido, nil).Once()
		repo.On("GetValidPagointentoByPagoId", int64(pagoSinIntentoValido.ID)).Return(&entities.Pagointento{}, fmt.Errorf("no se encontró intento con el id de pago: 55")).Once()

		_, err := serv.GetBilling(pagoSinIntentoValido.Uuid)

		assertions.EqualError(err, "no se encontró intento con el id de pago: 55")
	})

	t.Run("Cuando se genera correctamente el comprobante devuelvo un buffer de bytes correspondiente al PDF", func(t *testing.T) {
		repo.On("GetPagoByUuid", pagoCreated.Uuid).Return(&pagoCreated, nil).Once()
		repo.On("GetValidPagointentoByPagoId", int64(pagoCreated.ID)).Return(&pagoIntentoValido, nil).Once()

		pdf, err := serv.GetBilling(pagoCreated.Uuid)

		assertions.NoError(err)
		assertions.NotZero(pdf)
		assertions.IsType(pdf, &bytes.Buffer{})
	})
}
