package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/prisma"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	prismadtos "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/prismadtos/prismatransacciones"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
)

type debitPayment struct {
	service prisma.Service
}

func NewDebitPayment(s prisma.Service) PaymentMethod {
	return &debitPayment{service: s}
}

func (c *debitPayment) CreateResultado(request *dtos.ResultadoRequest, pago *entities.Pago, cuenta *entities.Cuenta, transactionID string, installmentsDetails *dtos.InstallmentDetailsResponse) (*entities.Pagointento, error) {
	//serv := prisma.Resolve()

	// Validaciones básicas
	err := request.Validar()
	if err != nil {
		return nil, err
	}

	// armo request del token
	tokenReq := prismadtos.StructToken{
		Card: prismadtos.Card{
			CardNumber:          request.CardNumber,
			CardExpirationMonth: request.CardMonth,
			CardExpirationYear:  request.CardYear[2:4],
			SecurityCode:        request.CardCode,
			CardHolderName:      request.HolderName,
			CardHolderIdentification: prismadtos.CardHolderIdentification{
				TypeDni:   prismadtos.EnumTipoDocumento(request.HolderDocType),
				NumberDni: request.HolderDocNum,
			},
		},
		TypePay: "simple",
	}

	// llamo servicio de obtener el token
	tokenInterface, err := c.service.SolicitarToken(tokenReq)
	if err != nil {
		return nil, err
	}
	// casteo la respuesta a estructura de token
	token := tokenInterface.(prismadtos.PagoToken)
	logs.Info(token)

	//siteid, _ := strconv.Atoi(config.SITEID)
	payReq := prismadtos.StructPayments{
		PagoSimple: prismadtos.PaymentsSimpleRequest{
			Customerid: prismadtos.Customerid{
				ID: fmt.Sprint(cuenta.ID),
			},
			SiteTransactionID: transactionID,
			SiteID:            config.SITEID, //int64(siteid),
			Token:             token.Id,
			PaymentMethodID:   request.PaymentMethodID,
			//PaymentMethodID:   31,
			Bin:          token.Bin,
			Amount:       request.Importe,
			Currency:     "ARS",
			Installments: 1,
			Description:  pago.Description,
			PaymentType:  "single",
			//EstablishmentName: cuenta.Cuenta,
			EstablishmentName: "TelCo Wee",
			//Email:             request.HolderEmail,
			Customeremail: prismadtos.Customeremail{
				Email: request.HolderEmail,
			},
			SubPayments: make([]interface{}, 0),
		},
		TypePay: "simple",
	}

	// llamo al servicio del payment
	var payment prismadtos.PaymentsSimpleResponse
	paymentInterface, err := c.service.Payments(payReq)
	// casteo la respuesta a estructura de payment
	if err != nil {
		payment.Status = err.Error()
	} else {
		payment = paymentInterface.(prismadtos.PaymentsSimpleResponse)
	}

	logs.Info(payment)
	// armo la respuesta del metodo
	paidAt, _ := time.Parse("2006-01-02T15:04Z", payment.Date)
	if err != nil {
		paidAt = time.Time{}
	}
	// convierto site id de la respuesta de la api
	api_site_id, _ := strconv.Atoi(payment.SiteID)

	response := entities.Pagointento{
		PagosID:      int64(pago.ID),
		MediopagosID: 2,
		ExternalID:   fmt.Sprint(payment.ID),
		PaidAt:       paidAt,
		ReportAt:     time.Now().Local(),
		IsAvailable:  false,
		Amount:       entities.Monto(request.Importe),
		//Valorcupon:           entities.Monto(request.Importe),
		StateComment:         payment.Status,
		InstallmentdetailsID: 1,
		HolderType:           request.HolderDocType,
		HolderNumber:         request.HolderDocNum,
		HolderName:           request.HolderName,
		HolderEmail:          request.HolderEmail,
		TicketNumber:         payment.StatusDetails.Ticket,
		AuthorizationCode:    payment.StatusDetails.CardAuthorizationCode,
		CardLastFourDigits:   token.LastFourDigits,
		TransactionID:        transactionID,
		SiteID:               int64(api_site_id),
	}
	fmt.Printf("CheckOut - site_transacction_id : %v - id : %v - descripcion %v\n", response.TransactionID, response.ExternalID, pago.Description)
	return &response, nil
}
