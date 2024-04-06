package services

import (
	"errors"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/pasarela/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/filtros/administracion"

	pasalera_repository "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/pasarela/repositories"
)

// Interfaz
type Service interface {
	GetPlanCuotasService(importe int) (response []dtos.PlanCuotasResponse, err error)
	HealthcheckService() (err error)

}

// Variable
var pasarela *service

// Estructura
type service struct {
	pasarelaRemoteRepository pasalera_repository.PasarelaRemoteRepository
	utilService              util.UtilService
}

// Constructor
func NewPasarelaService(rrp pasalera_repository.PasarelaRemoteRepository, us util.UtilService) Service {
	pasarela = &service{
		pasarelaRemoteRepository: rrp,
		utilService:              us,
	}
	return pasarela
}

// Funciones que implementan la interfaz
func (s *service) GetPlanCuotasService(importe int) (response []dtos.PlanCuotasResponse, err error) {
	fecha := time.Now().Format("2006-01-02")
	fechaActual, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		err = errors.New("error de conversión de datos")
		return
	}
	responseInstallments, erro := s.pasarelaRemoteRepository.GetInstallments(fechaActual)
	if erro != nil {
		return
	}
	for _, valueMedioPagoInstallment := range responseInstallments {
		var details []dtos.PlanCuotasResponseDetalle
		var installmentTemp dtos.PlanCuotasResponse
		for _, valueInstallment := range valueMedioPagoInstallment.Installments {
			// logs.Info(valueInstallment.VigenciaHasta)
			if valueInstallment.VigenciaHasta == nil {
				installmentTemp = dtos.PlanCuotasResponse{
					Id:                      valueInstallment.ID,
					Descripcion:             valueInstallment.Descripcion,
					MediopagoinstallmentsID: valueInstallment.MediopagoinstallmentsID,
				}
				for _, valueInstalmentDetail := range valueInstallment.Installmentdetail {
					details = append(details, dtos.PlanCuotasResponseDetalle{
						InstallmentsID: valueInstallment.ID,
						Cuota:          uint(valueInstalmentDetail.Cuota),
						Tna:            valueInstalmentDetail.Tna,
						Tem:            valueInstalmentDetail.Tem,
						Coeficiente:    valueInstalmentDetail.Coeficiente,
					})
				}
				break
			}
			// logs.Info("============================\n")
			// fmt.Printf("after %v - before %v \n", valueInstallment.VigenciaDesde.After(fechaActual), valueInstallment.VigenciaDesde.Before(fechaActual))
			// logs.Info("============================\n")

			// fmt.Printf("%v--%v--%v \n", valueInstallment.VigenciaDesde, fechaActual, valueInstallment.VigenciaHasta)
			// logs.Info("============================\n")
			// logs.Info("============================\n")

			if (fechaActual.After(valueInstallment.VigenciaDesde) && fechaActual.Before(*valueInstallment.VigenciaHasta)) || (fechaActual.Equal(valueInstallment.VigenciaDesde) && fechaActual.Before(*valueInstallment.VigenciaHasta)) || (fechaActual.After(valueInstallment.VigenciaDesde) && fechaActual.Equal(*valueInstallment.VigenciaHasta)) {
				installmentTemp = dtos.PlanCuotasResponse{
					Id:                      valueInstallment.ID,
					Descripcion:             valueInstallment.Descripcion,
					MediopagoinstallmentsID: valueInstallment.MediopagoinstallmentsID,
				}
				for _, valueInstalmentDetail := range valueInstallment.Installmentdetail {
					details = append(details, dtos.PlanCuotasResponseDetalle{
						InstallmentsID: valueInstallment.ID,
						Cuota:          uint(valueInstalmentDetail.Cuota),
						Tna:            valueInstalmentDetail.Tna,
						Tem:            valueInstalmentDetail.Tem,
						Coeficiente:    valueInstalmentDetail.Coeficiente,
					})
				}
				break
			}

		}
		response = append(response, dtos.PlanCuotasResponse{
			Id:                      installmentTemp.Id,
			Descripcion:             installmentTemp.Descripcion,
			MediopagoinstallmentsID: installmentTemp.MediopagoinstallmentsID,
			Installmentdetail:       details,
		})
	}


//REVIEW - Se comento este codigo por que realiza una consulta hacia la api de pasarela, y se trajo el servicio de api pasarela aca
	//planes, err = s.pasarelaRemoteRepository.GetPlanCuotasRepository()

	// if err != nil {
	// 	return response, err
	// }

	imp := entities.Monto(importe).Float64()
	// calcular el importe de cada cuota multiplicando el importe del pago por el coeficiente de cada instalment detail
	for i, installment := range response {
		for j, installmentdetail := range installment.Installmentdetail {
			importeTotal := s.utilService.ToFixed(imp*installmentdetail.Coeficiente, 2)
			response[i].Installmentdetail[j].ImporteTotal = importeTotal
			cuota := response[i].Installmentdetail[j].ImporteTotal / float64(installmentdetail.Cuota)
			response[i].Installmentdetail[j].ImporteCuota = s.utilService.ToFixed(cuota, 2)
		}
	}

	return response, err
}
func (s *service) HealthcheckService() (err error){
	filtroConfiguraciones:=filtros.ConfiguracionFiltro{
		Nombre: "ESTADO_APLICACION",
	}
	_,_,err=   s.pasarelaRemoteRepository.GetConfiguraciones(filtroConfiguraciones)
	if err != nil {
		return 
	}
	return
}