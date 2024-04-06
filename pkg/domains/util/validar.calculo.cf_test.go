package util_test

import (
	"testing"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/test/mocks/mockrepository"
	"github.com/stretchr/testify/assert"
)

type TableDriverTestValidacionCalculoCF struct {
	TituloPrueba string
	WantError    string
	Want         utildtos.ResponseValidarCF
	DataSet      utildtos.RequestValidarCF
	Cupon        entities.Monto
	Arancel      float64
	Tna          float64
	Cuotas       float64
	Dias         float64
}

func EstructuraValidacionCalculoCF() (tableDriverTest TableDriverTestValidacionCalculoCF) {
	tableDriverTest = TableDriverTestValidacionCalculoCF{
		TituloPrueba: "validar calculo coto finaciero",
		WantError:    "",
		Want: utildtos.ResponseValidarCF{
			CostoFinanciero:      87.41,
			ValorPresente:        894.59,
			CostoTotalPorcentaje: 8.74,
			ValorCoeficiente:     1.0977,
		},
		DataSet: utildtos.RequestValidarCF{
			Cupon:        691380, //1519194,
			ArancelMonto: 0.0100,
			Tna:          0,
			Cuotas:       1,
			Dias:         14,
		},
	}
	return
}

func TestCalculoCFValidar1(t *testing.T) {
	tabletDriverTest := EstructuraValidacionCalculoCF()
	//utilService := new(mockservice.MockUtilService)
	repository := new(mockrepository.MockRepositoryUtil)
	servicio := util.NewUtilService(repository)
	t.Run(tabletDriverTest.TituloPrueba, func(t *testing.T) {
		want := tabletDriverTest.Want
		// utilService.On("ValidarCalculoCF", tabletDriverTest.Cupon, tabletDriverTest.Arancel, tabletDriverTest.Tna, tabletDriverTest.Cuotas, tabletDriverTest.Dias).Return(true)
		// got := utilService.ValidarCalculoCF(tabletDriverTest.Cupon, tabletDriverTest.Arancel, tabletDriverTest.Tna, tabletDriverTest.Cuotas, tabletDriverTest.Dias)
		got := servicio.ValidarCalculoCF(tabletDriverTest.DataSet)
		assert.Equal(t, got, want)
	})
}

func EstructuraValidacionCalculoCF1() (tableDriverTest TableDriverTestValidacionCalculoCF) {
	tableDriverTest = TableDriverTestValidacionCalculoCF{
		TituloPrueba: "validar calculo coto finaciero",
		WantError:    "",
		Want: utildtos.ResponseValidarCF{
			CostoFinanciero:      0,       //2758.54,
			ValorPresente:        6844.66, //12105.33,
			CostoTotalPorcentaje: 0,       //18.22,
			ValorCoeficiente:     1,       //1.2279,
		},
		DataSet: utildtos.RequestValidarCF{
			Cupon:        691380, //100000, //1519194,
			ArancelMonto: 0.0100,
			Tna:          0, //81.0,
			Cuotas:       1,
			Dias:         14,
		},
	}
	return
}

func TestCalculoCFValidar2(t *testing.T) {
	tabletDriverTest := EstructuraValidacionCalculoCF1()
	//utilService := new(mockservice.MockUtilService)
	repository := new(mockrepository.MockRepositoryUtil)
	servicio := util.NewUtilService(repository)
	t.Run(tabletDriverTest.TituloPrueba, func(t *testing.T) {
		want := tabletDriverTest.Want
		// utilService.On("ValidarCalculoCF", tabletDriverTest.Cupon, tabletDriverTest.Arancel, tabletDriverTest.Tna, tabletDriverTest.Cuotas, tabletDriverTest.Dias).Return(true)
		// got := utilService.ValidarCalculoCF(tabletDriverTest.Cupon, tabletDriverTest.Arancel, tabletDriverTest.Tna, tabletDriverTest.Cuotas, tabletDriverTest.Dias)
		got := servicio.ValidarCalculoCF(tabletDriverTest.DataSet)
		assert.Equal(t, got, want)
	})
}
