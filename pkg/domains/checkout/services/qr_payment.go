package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/linkdtos/linkqr"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
)

type qrPayment struct {
	service     apilink.AplinkService
	utilService util.UtilService
}

func NewQrPayment(s apilink.AplinkService, util util.UtilService) PaymentMethod {
	return &qrPayment{
		service:     s,
		utilService: util,
	}
}

func (c *qrPayment) CreateResultado(request *dtos.ResultadoRequest, pago *entities.Pago, cuenta *entities.Cuenta, transactionID string, installmentsDetails *dtos.InstallmentDetailsResponse) (*entities.Pagointento, error) {

	// Id del medio de pago QR
	medioPagoQR := int64(40)

	err := commons.EsCuilValido(request.HolderCuit)
	if err != nil {
		return nil, err
	}

	importe := entities.Monto(request.Importe)
	qrRequest := linkqr.RequestApilinkCrearQr{
		CodigoActividadMCC:     "",
		CodigoSucursalComercio: config.APILINK_COD_SUCURSAL,
		CodigoPOSComercio:      config.APILINK_COD_POS,
		Monto:                  importe.Float64(),
	}

	// Genero una respuesta aleatoria simulando la de la API Qr
	//qrrr := linkqr.GenerateRandomResponse()
	//
	// Respuesta VERDADERA
	qrResult, err := c.service.CreateQrApiLinkService(config.IBM_CLIENT_ID, qrRequest)

	pagoResultado := entities.Pagointento{
		PagosID:              int64(pago.ID),
		MediopagosID:         medioPagoQR,
		InstallmentdetailsID: 1,
		ExternalID:           strconv.Itoa(qrResult.Data.QrDato.OperacionID),
		ReportAt:             time.Now().Local(),
		IsAvailable:          false,
		Amount:               entities.Monto(request.Importe),
		TransactionID:        transactionID,
		HolderName:           request.HolderName,
		HolderType:           "CUIL/CUIT",
		HolderNumber:         request.HolderCuit,
		HolderEmail:          request.HolderEmail,
		HolderCbu:            "",
	}

	if err != nil {
		logs.Error(err)
		pagoResultado.StateComment = err.Error()
	}

	var iniciados = []interface{}{"SUCCESS"}

	if commons.Include(iniciados, string(qrResult.Data.Status)) {
		formato := "2006-01-02T15:04:05.999-07:00"
		fecha, err := time.Parse(formato, qrResult.Data.QrDato.FechaCreacion)
		if err != nil {
			fmt.Println("Error al analizar la fecha:", err)
		}
		pagoResultado.PaidAt = fecha
		pago.FechaHoraExpiracion = fecha.Add(9 * time.Minute)
		pagoResultado.StateComment = "APROBADO"
		pagoResultado.Qr = qrResult.Data.QrDato.QRData

	}

	return &pagoResultado, nil
}
