package util_test

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
	util "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/test/domains/util/utilfake"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/test/mocks/mockrepository"
	"github.com/stretchr/testify/assert"
)

var (
	mockutils = new(mockrepository.MockRepositoryUtil)
	// mockFactory = new(mockservice.MockCrearMensajeServiceFactory)
	//mockServiceAdministracion  = new(mockservice.)
	service = util.NewUtilService(mockutils)
)

func TestRequestValidConsultarMovimientos(t *testing.T) {
	TableDriverTest := utilfake.EstructuraVerificarCbu()
	t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
		want := TableDriverTest.WantTable
		logs.Info(want)
		// _, got := service.ConsultarMovimientos(TableDriverTest.Request)
		// assert.Equal(t, got.Error(), want)
	})
}

// test para enviar email

func TestBuildEmailSend(t *testing.T) {
	TableDriverTest := utilfake.EstructuraValidarCbu()
	t.Run(TableDriverTest.TituloPrueba, func(t *testing.T) {
		want := TableDriverTest.WantTable
		logs.Info(want)
		got, _ := service.ValidarCBU(TableDriverTest.Cbu)
		assert.Equal(t, got, want)
	})
}

// test construir movimientos , caluclo de comisiones e impuestos
func TestBuildComisiones(t *testing.T) {
	TableDriverTest := utilfake.EstructuraBuildComisiones()
	for _, test := range TableDriverTest {
		t.Run(test.TituloPrueba, func(t *testing.T) {
			want := test.WantTable
			logs.Info(test.TituloPrueba)
			got := service.BuildComisiones(test.RequestMovimiento, test.RequestCuentaComision, test.RequestIva, test.ImporteSolicitado)
			assert.Equal(t, got, want)
		})
	}
}
