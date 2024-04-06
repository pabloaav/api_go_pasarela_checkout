package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/pagooffline"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/pagoofflinedtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	uuid "github.com/satori/go.uuid"
)

type multipagoPayment struct {
	service pagooffline.Service
}

func NewMultipagoPayment(s pagooffline.Service) PaymentMethod {
	return &multipagoPayment{service: s}
}

func (rapi *multipagoPayment) CreateResultado(request *dtos.ResultadoRequest, pago *entities.Pago, cuenta *entities.Cuenta, transactionID string, installmentsDetails *dtos.InstallmentDetailsResponse) (*entities.Pagointento, error) {
	// se arma el dto request DE RAPIPAGO
	recargo := pago.SecondTotal - pago.FirstTotal
	dtoRequestRapipago := pagoofflinedtos.OffLineRequestResponse{
		CodigoEmpresa:     config.CODIGO_EMPRESA_RAPIPAGO,
		NumeroCliente:     strconv.Itoa(int(cuenta.ClientesID)),
		NumeroComprobante: strconv.Itoa(int(pago.ID)),
		Importe:           request.Importe,
		FechaPrimerVto:    pago.FirstDueDate,
		ImporteRecargo:    entities.Monto.Int64(recargo),
		FechaSegundoVto:   pago.SecondDueDate,
	}
	rapipagoRequest_ov := pagoofflinedtos.New(dtoRequestRapipago)
	if len(rapipagoRequest_ov.GetErrors()) > 0 {
		listErrores := rapipagoRequest_ov.GetErrors()
		logs.Error(listErrores)
		var errorsString string
		for _, valueErrores := range listErrores {
			errorsString = " - " + valueErrores
		}
		return nil, errors.New(errorsString) //rapipagoRequest_ov.GetErrors()

	}

	codigoMultipago := uuid.NewV4().String()

	response := entities.Pagointento{
		PagosID:              int64(pago.ID),
		MediopagosID:         39,
		InstallmentdetailsID: 1,
		ExternalID:           pago.ExternalReference,
		PaidAt:               time.Now(),
		ReportAt:             time.Now().Local(),
		StateComment:         "approved",
		HolderType:           request.HolderDocType,
		HolderNumber:         request.HolderDocNum,
		HolderName:           request.HolderName,
		HolderEmail:          request.HolderEmail,
		TransactionID:        transactionID,
		Barcode:              codigoMultipago,
	}

	return &response, nil
}
