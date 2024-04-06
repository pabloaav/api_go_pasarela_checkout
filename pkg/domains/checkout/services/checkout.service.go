package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/checkout/repositories"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/pagooffline"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/prisma"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/util"
	webhook "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/multipagosdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/rapipago"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/utildtos"
	webhookDto "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/webhook"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/filtros/administracion"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

// var (
// 	PrismaServiceVar prisma.Service = prisma.Resolve()
// )

type Service interface {
	// NewPago genera un pago en la base de datos y devuelve al url para acceder al checkout y pagarlo.
	NewPago(ctx context.Context, request *dtos.PagoRequest, apiKey string) (*dtos.PagoResponse, error)
	// GetPaid devuelve los datos de un pago que necesita el checkout para mostrar al usuario pagador.
	GetPaid(barcode string) (*dtos.CheckoutResponse, error)
	// GetPagoResultado obtiene los parámetros del checkout y ejecuta el pago devolviendo detalles del resultado.
	GetPagoResultado(ctx context.Context, request *dtos.ResultadoRequest) (*dtos.ResultadoResponse, error)
	// CheckPrisma funcionalidad que ayuda al checkout a saber si puede contar con los servicios de Prisma.
	CheckPrisma() error
	// GetBilling devuelve el recibo del pago en un archivo pdf
	GetBilling(uuid string) (*bytes.Buffer, error)
	// GetTarjetas devuelve datos de las tarjetas
	GetTarjetas() (*[]entities.Mediopago, error)
	// GetMatenimietoSistema permite verificar si el checkout esta en mantenimiento o no
	GetMatenimietoSistema() (estado bool, fecha time.Time, erro error)
	// GetPagoStatus devuelve un bool si el pago estado del pago es
	GetPagoStatus(barcode string) (status bool, erro error)
	// Notifica a los cliente el pago que se realizó
	NotificarPagos(listaPagos []webhookDto.WebhookResponse) (pagoupdate []uint)

	HashOperacionTarjeta(number string, pagointento_id int64) (status bool, erro error)
	ControlTarjetaHash(number string) (status bool, erro error)

	// GetMultiPagoStatus devuelve respuesta para Consulta Multipago con :uuid
	GetMultiPagoStatus(request multipagosdtos.RequestConsultaMultipago) (response multipagosdtos.ResponseConsultaMultipago, erro error)

	// GetMultiPagoControl devuelve respuesta para Consulta Multipago sobre transacciones hechas en un periodo
	GetMultiPagoControl(request multipagosdtos.RequestControlMultipago) (response multipagosdtos.ResponseControlMultipago, erro error)

	// PostMultiPago registra pago por medioPago Multipago
	PostMultiPago(ctx context.Context, request multipagosdtos.RequestPagoMultipago) (response multipagosdtos.ResponsePagoMultipago, erro error)

	// GetRapipagoQuery devuelve los datos de un pago hecho por rapipago
	GetRapipagoQuery(request rapipago.RequestRapipagoConsulta) (response rapipago.ResponseRapipagoConsulta, err error)

	// PostRapipagoPago actualiza los datos de un pago hecho por rapipago. Pasa a estado aprobado.
	PostRapipagoPago(ctx context.Context, request rapipago.RequestRapipagoConsulta) (response rapipago.ResponseRapipagoImputacion, err error)

	// Control apikey Adquiriente
	ControlAdquirienteApikey(apikey string, adquiriente string) (control bool, err error)

	//Confirmacion de un pago por parte de Rapipago
	RapipagoConfirmacionService(request rapipago.RequestRapipagoConfirmacion) (response rapipago.ResponseRapipagoConfirmacion, err error)

	// getEstadoAppService obtiene un registro de la tabla Configuraciones para conocer el estado de la aplicacion
	GetEstadoAppService() (err error)

	// Middleware de apikey por cuenta
	GetCuentaByApiKey(apikey string) (result bool, erro error)
}

type service struct {
	repository         repositories.Repository
	commons            commons.Commons
	payment            PaymentFactory
	prismaService      prisma.Service
	pagoOffLineService pagooffline.Service
	utilService        util.UtilService
	webhook            webhook.RemoteRepository
}

/* -------------------- Constructores ---------------------- */

func NewService(r repositories.Repository, c commons.Commons, ps prisma.Service, polS pagooffline.Service, utilService util.UtilService, webhook webhook.RemoteRepository) Service {
	return &service{
		repository:         r,
		commons:            c,
		payment:            &paymentFactory{},
		prismaService:      ps,
		pagoOffLineService: polS,
		utilService:        utilService,
		webhook:            webhook,
	}
}

func NewServiceWithPayment(r repositories.Repository, c commons.Commons, p PaymentFactory) Service { //, util util.UtilService
	return &service{
		repository: r,
		commons:    c,
		payment:    p,
		//utilService: util,
	}
}

/* ------------------ Funciones de Interfaz ----------------- */

func (s *service) NewPago(ctx context.Context, request *dtos.PagoRequest, apiKey string) (*dtos.PagoResponse, error) {
	// Valido campos obligatorios
	err := request.Validar()
	if err != nil {
		return nil, err
	}

	// objeto para verificar entre otras cosas validez de fechas
	validar := commons.NewAlgoritmoVerificacion()

	// se convierten las fechas a un formato correcto para la funciones de dias entre vencimientos

	primer_vencimiento := commons.ConvertirFechaYYYYMMDD(request.FirstDueDate)
	segundo_vencimiento := commons.ConvertirFechaYYYYMMDD(request.SecondDueDate)

	cant_dias, err := validar.CalcularDiasEntreFechas(primer_vencimiento, segundo_vencimiento)

	if err != nil {
		return nil, errors.New(err.Error())
	}
	if cant_dias > 99 {
		return nil, errors.New("la cantidad de dias entre fechas de vencimiento no puede ser mayor a dos dígitos")
	}

	// mayusculas al nombre del pagadro y la descripcion
	request.ToFormatStr()
	// Valido los montos de los items con el primer total a pagar
	err = s.validarMontos(request.Items, entities.Monto(request.FirstTotal))
	if err != nil {
		return nil, err
	}

	s.repository.BeginTx()
	defer func() {
		if err != nil {
			s.repository.RollbackTx()
		} else {
			s.repository.CommitTx()
		}
	}()

	// Busco Cuenta con apiKey
	cuenta, err := s.repository.GetCuentaByApikey(apiKey)
	if err != nil {
		return nil, err
	}
	// almacenar cliente id en context
	ctx = ctxWithClienteID(ctx, uint(cuenta.ClientesID))

	// Busco id de Tipo de pago asociado a la cuenta. Puede ser mas de uno
	var tipoPagoID int64

	for _, tp := range *cuenta.Pagotipos {
		if strings.EqualFold(request.PaymentType, tp.Pagotipo) {
			tipoPagoID = int64(tp.ID)
		}
	}

	// si tipoPagoID sigue siendo 0 significa q no hay configuracion de cuentas ni pagotipos
	if tipoPagoID <= 0 {
		return nil, fmt.Errorf("en la configuración de cuentas, no hay tipo de pago correcto para %s", request.PaymentType)
	}

	// Parseo string a fechas
	fechaVencimiento, err := time.Parse("02-01-2006", request.FirstDueDate)
	if err != nil {
		return nil, fmt.Errorf("error en fecha de vencimiento %s", err.Error())
	}
	var fechaSegundoVencimiento time.Time
	if len(request.SecondDueDate) > 0 {
		fechaSegundoVencimiento, err = time.Parse("02-01-2006", request.SecondDueDate)
		if err != nil {
			return nil, fmt.Errorf("error en fecha de segundo vencimiento: %s", err.Error())
		}
	}

	var fechaHoraExpiracion time.Time
	var time_layout = "2006-01-02 15:04:05"
	// Fecha y hora de expiracion
	if len(request.FechaHoraExpiracion) > 0 {
		fechaHoraExpiracion, err = time.Parse(time_layout, request.FechaHoraExpiracion)
		if err != nil {
			logs.Error(err)
		}
	}

	// Genero UUID unico para el pago
	// label para saltar a este punto si genera un uuid repetido
reintento:
	pagoid := s.commons.NewUUID()
	//logs.Info("codigo unico " + pagoid.String())

	// se genera registro en la tabla pasarela.pagos
	pago := entities.Pago{
		PagostipoID:         tipoPagoID,
		PagoestadosID:       1,
		Description:         request.Description,
		FirstDueDate:        fechaVencimiento,
		FirstTotal:          entities.Monto(request.FirstTotal),
		SecondDueDate:       fechaSegundoVencimiento,
		SecondTotal:         entities.Monto(request.SecondTotal),
		PayerName:           request.PayerName,
		PayerEmail:          request.PayerEmail,
		ExternalReference:   request.ExternalReference,
		Metadata:            request.Metadata,
		Uuid:                pagoid,
		PdfUrl:              config.APP_BACKGROUND_CHECKOUT_URL + "/checkout/bill/" + pagoid,
		Pagoitems:           request.Items,
		Expiration:          request.Expiration,
		FechaHoraExpiracion: fechaHoraExpiracion,
	}

	pagodb, err := s.repository.CreatePago(ctx, &pago)
	if err != nil {
		if strings.Contains(err.Error(), "uuid_UNIQUE") {
			// cuando se intenta guardar un uuid repetido, salta a la linea 119 con el label reintento
			goto reintento
		}
		return nil, fmt.Errorf("NewPago: %s", err)
	}

	if err = s.repository.CreatePagoEstadoLog(
		ctx,
		&entities.Pagoestadologs{PagosID: int64(pagodb.ID), PagoestadosID: 1},
	); err != nil {
		logs.Error(err)
	}

	// Devuelvo response adecuada
	items := make([]dtos.PagoResponseItems, 0)
	if len(request.Items) > 0 {
		for _, t := range request.Items {
			items = append(items, dtos.PagoResponseItems{
				Quantity:    int64(t.Quantity),
				Description: t.Description,
				Amount:      t.Amount.Float64(),
				Identifier:  t.Identifier,
			})
		}
	}

	response := dtos.PagoResponse{
		ID:                  int64(pagodb.ID),
		Estado:              "pending",
		Description:         pagodb.Description,
		FirstDueDate:        pagodb.FirstDueDate.Format("02-01-2006"),
		FirstTotal:          pagodb.FirstTotal.Float64(),
		SecondDueDate:       pagodb.SecondDueDate.Format("02-01-2006"),
		SecondTotal:         pagodb.SecondTotal.Float64(),
		PayerName:           pagodb.PayerName,
		PayerEmail:          pagodb.PayerEmail,
		ExternalReference:   pagodb.ExternalReference,
		Metadata:            pagodb.Metadata,
		Uuid:                pagodb.Uuid,
		CheckoutUrl:         config.APP_CHECKOUT_URL + "/checkout/" + pagodb.Uuid,
		CreatedAt:           pagodb.CreatedAt.Format("02-01-2006"),
		Items:               items,
		Expiration:          pagodb.Expiration,
		FechaHoraExpiracion: pago.FechaHoraExpiracion.Format(time_layout),
	}

	return &response, nil
}

func (s *service) GetPaid(uuid string) (*dtos.CheckoutResponse, error) {

	var secondDueDate bool
	// valida uuid
	if len(uuid) <= 0 {
		return nil, fmt.Errorf("debe enviar código único del pago, envió: %s", uuid)
	}

	if ok, err := s.commons.IsValidUUID(uuid); !ok {
		return nil, fmt.Errorf("el identificador del pago no es válido: %w", err)
	}

	filtroPago := filtros.PagoFiltro{
		Uuids:            []string{uuid},
		CargaMedioPagos:  true,
		CargarPagoEstado: true,
	}

	// busca pago por uuid en la base de datos
	pago, err := s.repository.GetPaymentByUuid(filtroPago)
	if err != nil {
		return nil, fmt.Errorf("error al obtener pago: %s", err.Error())
	}

	// si ya pagó devuelve un error
	if pago.PagoestadosID != 1 {
		if pago.PagoIntentos[len(pago.PagoIntentos)-1].Mediopagos.ChannelsID != 5 {
			filtroMedioPago := make(map[string]interface{})
			filtroMedioPago["id"] = pago.PagoIntentos[len(pago.PagoIntentos)-1].MediopagosID
			medioPago, erro := s.repository.GetMediopago(filtroMedioPago)
			if erro != nil {
				return nil, fmt.Errorf("error al obtener medio de pago: %s", erro.Error())
			}
			return nil, fmt.Errorf("el pago ya fue procesado a través del medio de pago %v", medioPago.Mediopago)
		} else { // Comportamiento del pago con QR. Si no está en un estado final, se sigue mostrando este QR.
			if pago.PagoEstados.Final {
				return nil, fmt.Errorf("el pago ya fue procesado.")
			}
			expirado, _ := verificarExpiracionQR(pago.PagoIntentos[len(pago.PagoIntentos)-1].CreatedAt)
			if expirado {
				return nil, fmt.Errorf("QR vencido, por favor genere otro.")
			}
		}
		/* Acá se debería verificar el vencimiento o no del QR en base a la hora de creación. */

	}

	// TIEMPO DE EXPIRACION DEL CHECKOUT. COMPARA CON VALOR EXPIRATION DEL PAGO
	err = _getPaymentExpiration(pago)
	if err != nil {
		return nil, err
	}

	pagotipochannels, err := s.repository.GetPagotipoChannelByPagotipoId(pago.PagostipoID)
	if err != nil {
		return nil, err
	}

	cuotas, err := s.repository.GetPagotipoIntallmentByPagotipoId(pago.PagostipoID)
	if err != nil {
		return nil, err
	}

	tipo, err := s.repository.GetPagotipoById(pago.PagostipoID)
	if err != nil {
		return nil, err
	}

	preference, err := s.repository.GetPreferencesByIdClienteRepository(uint(tipo.Cuenta.ClientesID))
	if err != nil {
		return nil, err
	}

	// array con los nombres de los channels
	var channelsArrayString []string
	for _, c := range *pagotipochannels {
		channelsArrayString = append(channelsArrayString, c.Channel.Channel)
	}

	// concatena las cuotas disponibles en un string
	var cuotasString string
	for key, value := range *cuotas {
		if len(*cuotas) == key+1 {
			cuotasString += value.Cuota
		} else {
			cuotasString += value.Cuota + ","
		}
	}

	// primer importe y fecha de vencimiento
	importe := pago.FirstTotal
	dueDate := pago.FirstDueDate
	// si se indica fecha de vencimiento, la comparamos para cobrar el segundo monto
	if !pago.FirstDueDate.IsZero() {
		hoy := time.Now().Local()
		hoyDate, err := time.Parse("2006-01-02T00:00:00Z", hoy.Format("2006-01-02T00:00:00Z"))
		if err != nil {
			return nil, fmt.Errorf("formato de fecha invalido")
		}
		// Si la fecha actual es posterior a la fecha indicada en el primer vencimiento (FirstDueDate) el importe que se considera como importe de pago es el segundo total (SecondTotal)
		if hoyDate.After(pago.FirstDueDate) {
			importe = pago.SecondTotal
			dueDate = pago.SecondDueDate
			secondDueDate = true
		}
	}

	// items de un pago
	byteItems, err := getPaymentItemsToResponse(pago)

	if err != nil {
		return nil, err
	}
	urlQr := ""

	if len(pago.PagoIntentos) > 0 {
		urlQr = pago.PagoIntentos[len(pago.PagoIntentos)-1].Qr
	}

	if tipo.SendUuid {
		tipo.BackUrlSuccess = _addUrlParam(tipo.BackUrlSuccess, "uudi", pago.Uuid)
		tipo.BackUrlPending = _addUrlParam(tipo.BackUrlPending, "uudi", pago.Uuid)
	}

	externalId := ""
	if len(pago.PagoIntentos) > 0 {
		if len(pago.PagoIntentos[len(pago.PagoIntentos)-1].ExternalID) > 0 && pago.PagoIntentos[len(pago.PagoIntentos)-1].ExternalID != "0" {
			externalId = pago.PagoIntentos[len(pago.PagoIntentos)-1].ExternalID
		}
	}

	// armo la respuesta
	response := dtos.CheckoutResponse{
		Estado:               pago.PagoEstados.Nombre,
		Description:          pago.Description,
		DueDate:              dueDate.Format("02-01-2006"),
		SecondDueDate:        secondDueDate,
		Total:                importe.Int64(), // `json:"first_total"`
		PayerName:            pago.PayerName,
		PayerEmail:           pago.PayerEmail,
		ExternalReference:    pago.ExternalReference,
		Metadata:             pago.Metadata,
		Uuid:                 pago.Uuid,
		PdfUrl:               pago.PdfUrl,
		CreatedAt:            pago.CreatedAt.String(),
		BackUrlSuccess:       tipo.BackUrlSuccess,
		BackUrlPending:       tipo.BackUrlPending,
		BackUrlRejected:      tipo.BackUrlRejected,
		IncludedChannels:     channelsArrayString, //strings.Split(tipo.IncludedChannels, ","),
		IncludedInstallments: cuotasString,        //tipo.IncludedInstallments,
		Items:                string(byteItems),   // items en formato cadena o string
		Preference: dtos.ResponsePreference{
			Client:         preference.Cliente.Cliente,
			MainColor:      preference.Maincolor,
			SecondaryColor: preference.Secondarycolor,
			Logo:           preference.Logo,
		},
		Cliente:             tipo.Cuenta.Cliente.Cliente,
		Url_qr:              urlQr,
		FechaHoraExpiracion: pago.FechaHoraExpiracion.String(),
		ExternalId:          externalId,
	}

	return &response, nil
}

func (s *service) GetPagoResultado(ctx context.Context, request *dtos.ResultadoRequest) (*dtos.ResultadoResponse, error) {
	// validaciones basicas
	if err := request.Validar(); err != nil {
		return nil, err
	}
	request.ToFormatStr()
	// valido el metodo de pago con el de la base de datos

	if request.CardNumber != "" {
		control, err := s.ControlTarjetaHash(request.CardNumber)
		if err != nil {
			return nil, err
		}

		if control {
			return nil, fmt.Errorf("Error con Tarjeta. Comunicarse a través de nuestro correo de soporte si persiste.")
		}

	}

	channel, err := s.repository.GetChannelByName(request.Channel)
	if err != nil {
		return nil, err
	}

	// valido el medio de pago
	filtro := make(map[string]interface{})
	filtro["channels_id"] = channel.ID
	if len(request.CardBrand) > 0 {
		filtro["mediopago"] = request.CardBrand
	} else {
		return nil, fmt.Errorf("error tarjeta invalida")
	}
	/* REVIEW: revisar cuando trae el medio de pago */
	medio, err := s.repository.GetMediopago(filtro)
	if err != nil {
		return nil, fmt.Errorf("error en medio de pago: %s", err.Error())
	}

	// al request le paso el id externo de medio de pago
	request.PaymentMethodID, _ = strconv.ParseInt(medio.ExternalID, 10, 64)

	filtroPago := filtros.PagoFiltro{
		Uuids: []string{request.Uuid},
	}

	// obtengo datos del pago mediante el uuid
	pago, err := s.repository.GetPaymentByUuid(filtroPago)
	if err != nil {
		return nil, err
	}

	// ver tiempo de expiracion para realizar el pago en el checkout
	err = _getPaymentExpiration(pago)
	if err != nil {
		return nil, err
	}

	// obtengo datos del tipo de pago configurado por el cliente
	tipo, err := s.repository.GetPagotipoById(pago.PagostipoID)
	if err != nil {
		return nil, err
	}

	// obtengo datos de la cuenta bancaria a la cual corresponde el pago
	cuenta, err := s.repository.GetCuentaById(int64(tipo.CuentasID))
	if err != nil {
		return nil, err
	}

	// calculo el importe a pagar segun las fechas de vencimiento
	importe := pago.FirstTotal

	// si se indica fecha de vencimiento, la comparamos para cobrar el segundo monto
	if !pago.FirstDueDate.IsZero() {
		fechaHoy := time.Now().Local() // Fecha actual
		hoyDate, err := time.Parse("2006-01-02T00:00:00Z", fechaHoy.Format("2006-01-02T00:00:00Z"))
		if err != nil {
			return nil, fmt.Errorf("formato de fecha invalido")
		}
		logs.Info("fecha actual")
		logs.Info(hoyDate)
		logs.Info("fecha primer venc.")
		logs.Info(pago.FirstDueDate)
		logs.Info(hoyDate.After(pago.FirstDueDate))

		if hoyDate.After(pago.FirstDueDate) {
			importe = pago.SecondTotal
		}
	}
	// el monto lo pasamos como integer a las apis
	request.Importe = importe.Int64()
	fechaActual, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		erro := errors.New("error convertir fecha actual")
		return nil, erro
	}
	cuotas, _ := strconv.ParseInt(request.Installments, 10, 64)
	installments, err := s.repository.GetInstallmentsByMedioPagoInstallmentsId(medio.MediopagoinstallmentsID)
	if err != nil {
		return nil, err
	}
	var installmentId int64
	for _, valueInstallment := range installments {
		if valueInstallment.VigenciaHasta == nil {
			installmentId = int64(valueInstallment.ID)
			break
		}
		if (fechaActual.After(valueInstallment.VigenciaDesde) && fechaActual.Before(*valueInstallment.VigenciaHasta)) || (fechaActual.Equal(valueInstallment.VigenciaDesde) && fechaActual.Before(*valueInstallment.VigenciaHasta)) || (fechaActual.After(valueInstallment.VigenciaDesde) && fechaActual.Equal(*valueInstallment.VigenciaHasta)) {
			installmentId = int64(valueInstallment.ID)
			break
		}
	}

	installmentsDetails, err := s.repository.GetInstallmentDetails(installmentId, cuotas) //medio.InstallmentsID int64(installments.ID)
	if err != nil {
		return nil, err
	}
	filtroConfiguracion := filtros.ConfiguracionFiltro{
		Buscar:     true,
		Nombrelike: "IMPUESTO_SOBRE_COEFICIENTE",
	}
	configuracionImpuesto, err := s.utilService.GetConfiguracionesService(filtroConfiguracion)
	if err != nil {
		return nil, fmt.Errorf("no se pudo realizar transaccion: %s", err.Error())
	}
	impuestoId, err := strconv.Atoi(configuracionImpuesto[0].Valor)
	if err != nil {
		return nil, fmt.Errorf("no se pudo realizar transaccion: %s", err.Error())
	}
	impuesto, err := s.utilService.GetImpuestoByIdService(int64(impuestoId))
	if err != nil {
		return nil, fmt.Errorf("no se pudo realizar transaccion: %s", err.Error())
	}
	installmentsDetails.Impuesto = impuesto.Porcentaje
	// obtengo constructor de pago mediante patrón factory
	paymentMethod, err := s.payment.GetPaymentMethod(int(channel.ID))
	if err != nil {
		return nil, err
	}

	// genero un transaction_id
	transactionID := s.commons.NewUUID()

	// envio los datos al constructor correspondiente para procesar el pago
	resultado, err := paymentMethod.CreateResultado(request, pago, cuenta, transactionID, installmentsDetails)
	if err != nil {
		return nil, err
	}

	// busco en la base el id de installmentdetails
	resultado.InstallmentdetailsID = int64(installmentsDetails.Id)
	//resultado.InstallmentdetailsID = s.repository.GetInstallmentDetailsID(medio.InstallmentsID, cuotas)
	resultado.MediopagosID = int64(medio.ID)

	// algunas apis devuelven el monto en entero, me aseguro q se guarde en bd el float64
	resultado.Amount = importe

	// agrego cliente id 1 al context para la auditoria
	ctx = ctxWithClienteID(ctx, 1)

	// almaceno el resultado en la base de datos
	if ok, err := s.repository.CreateResultado(ctx, resultado); !ok {
		return nil, err
	}

	if resultado.CardLastFourDigits != "" {
		if _, err := s.HashOperacionTarjeta(request.CardNumber, int64(resultado.ID)); err != nil {
			return nil, err
		}
	}

	// actualizo el estado del pago
	// cuando el pago se procesa con exito, se le coloca una fecha a PaidAt,
	// cuando hay un error, se devuelve PaidAt con fecha 0, y no se actualiza el pago.

	if !resultado.PaidAt.IsZero() {
		if channel.Channel == "DEBIN" {
			pago.PagoestadosID = 2
		} else if channel.Channel == "OFFLINE" || channel.Channel == "MULTIPAGOS" {
			pago.PagoestadosID = 2
		} else if channel.Channel == "QR" {
			pago.PagoestadosID = 2
		} else {
			pago.PagoestadosID = 4
		}
		if ok, err := s.repository.UpdatePago(ctx, pago); !ok {
			return nil, err
		}
		if err = s.repository.CreatePagoEstadoLog(
			ctx,
			&entities.Pagoestadologs{PagosID: int64(pago.ID), PagoestadosID: pago.PagoestadosID},
		); err != nil {
			logs.Error(err)
		}
	}
	var importePagado entities.Monto
	if resultado.Valorcupon > 0 {
		importePagado = entities.Monto(resultado.Valorcupon)
	} else {
		importePagado = entities.Monto(resultado.Amount)
	}

	estadoPago, erro := s.repository.GetPagoEstado(pago.PagoestadosID)
	if erro != nil {
		estadoPago.Nombre = "PENDIENTE"
		logs.Error(erro)
	}

	// items de un pago
	byteItems, _ := getPaymentItemsToResponse(pago)

	//var vencimientoQR time.Time
	if channel.Channel == "QR" {
		//vencimientoQR = pago.FechaHoraExpiracion.Add(10 * time.Minute)
	}

	// Cargo el externalId para la respuesta.
	externalId := ""
	if resultado != nil {
		if resultado.ExternalID != "0" {
			externalId = resultado.ExternalID
		}
	}

	// Esta es la respuesta que se envia al checkout como resultado del proceso de pago
	response := dtos.ResultadoResponse{
		ID:                  int64(resultado.ID),
		Estado:              resultado.StateComment,
		EstadoPago:          estadoPago.Nombre,
		Exito:               !resultado.PaidAt.IsZero(),
		Uuid:                request.Uuid, // informacion para comprobante pdf
		Channel:             channel.Channel,
		Description:         pago.Description,
		FirstDueDate:        pago.FirstDueDate.Format("02-01-2006"),
		FirstTotal:          importe.Float64(),
		SecondDueDate:       pago.SecondDueDate.Format("02-01-2006"),
		SecondTotal:         pago.SecondTotal.Float64(),
		PayerName:           pago.PayerName,
		PayerEmail:          pago.PayerEmail,
		ExternalReference:   pago.ExternalReference,
		Metadata:            pago.Metadata,
		PdfUrl:              pago.PdfUrl,
		CreatedAt:           resultado.CreatedAt.String(),
		ImportePagado:       importePagado.Float64(),
		Items:               string(byteItems),           // informacion para comprobante pdf
		ClienteName:         tipo.Cuenta.Cliente.Cliente, // informacion para comprobante pdf
		ClienteCuit:         tipo.Cuenta.Cliente.Cuit,    // informacion para comprobante pdf
		Mediopago:           medio.Mediopago,             // informacion para comprobante pdf
		Barcode:             resultado.Barcode,           // informacion para comprobante pdf
		NumeroOperacion:     pago.ID,
		UrlQr:               resultado.Qr,
		ExternalId:          externalId,
		FechaHoraExpiracion: pago.FechaHoraExpiracion.String(),
	}

	// notificacion de pago exitoso por medio del webhook
	if tipo.BackUrlNotificacionPagos != "" && response.Exito {
		if err := s.notificacionPago(response, tipo); err != nil {
			logs.Error(fmt.Sprintf("webhook:no se pudo notificar el pago %v: %s", response.NumeroOperacion, err.Error()))
		} else {
			logs.Info(fmt.Sprintf("webhook:se notifico con éxito el pago: %v", response.NumeroOperacion))
		}
	}

	// enviar email de pago exitoso
	if response.Exito && medio.Mediopago != "QR" {
		// Construir el texto html del mensaje del email
		dir_url_comprobante := config.APP_BACKGROUND_CHECKOUT_URL + "/checkout/bill/" + pago.Uuid

		url_imagen_descargaDoc := "https://img.icons8.com/?size=512&id=2mGSkp2owx0d&format=png"

		var mensaje string
		if channel.Channel == "DEBIT" || channel.Channel == "CREDIT" {
			mensaje = "<ul style='list-style: none;text-align: left;display:inline-block; line-height: 23px'><li> Fecha: <b>#4</b></li><li> Referencia: <b>#0</b></li><li> Identificador de la transacción: <b>#1</b></li><li> Medio de pago: <b>#2</b></li><li> Concepto: <b>#3</b></li><li>Nro solicitud: <b>#5</b></li> <li style='padding-top:6px;' > <a href='" + dir_url_comprobante + "'><img src='" + url_imagen_descargaDoc + "' width='16' height='16'> Información de Pago</a> </li></ul>"
		} else {
			mensaje = "<p style='display: none'> <b>#0</b> <b>#1</b> <b>#2</b> </p> <a href='" + dir_url_comprobante + "'><img src='" + url_imagen_descargaDoc + "' width='16' height='16'> Información de Pago</a>"
		}

		var descripcion utildtos.DescripcionTemplate
		var detallesPago []utildtos.DetallesPago
		var emailContacto string
		for _, det := range pago.Pagoitems {
			var identificador string
			if len(det.Identifier) > 0 {
				identificador = " - " + det.Identifier
			}
			detallesPago = append(detallesPago, utildtos.DetallesPago{
				Descripcion: det.Description + identificador,
				Cantidad:    fmt.Sprintf("%v", det.Quantity),
				Monto:       fmt.Sprintf("$%v", s.utilService.ToFixed(det.Amount.Float64(), 2)),
			})
		}
		emailContacto = cuenta.Cliente.Emailcontacto
		if cuenta.Cliente.Emailcontacto == "" {
			emailContacto = config.EMAIL_SOPORTE
		}
		descripcion = utildtos.DescripcionTemplate{
			Cliente:       cuenta.Cliente.Cliente,
			Cuit:          cuenta.Cliente.Cuit,
			Detalles:      detallesPago,
			EmailContacto: emailContacto,
			TotalPagado:   fmt.Sprintf("$%v", response.ImportePagado),
		}

		// enviar mail al usuario pagador
		var arrayEmail []string
		var email string
		email = request.HolderEmail
		if request.HolderEmail == "" {
			email = pago.PayerEmail
		}
		arrayEmail = append(arrayEmail, email)
		params := utildtos.RequestDatosMail{
			Email:            arrayEmail,
			Asunto:           "Información de Pago",
			Nombre:           pago.PayerName,
			Mensaje:          mensaje,
			CamposReemplazar: []string{response.ExternalReference, pago.Uuid, medio.Mediopago, response.Description, resultado.PaidAt.Format("02-01-2006"), fmt.Sprintf("%v", pago.ID)},
			FiltroReciboPago: true,
			Descripcion:      descripcion,
			From:             "Wee.ar",
			TipoEmail:        "template",
			CanalPago:        channel.Channel,
		}
		erro = s.utilService.EnviarMailService(params)

		// control de error al enviar email
		if erro != nil {

			message := fmt.Errorf("error al enviar email al usuario pagador: " + erro.Error())

			log := entities.Log{
				Tipo:          entities.Error,
				Mensaje:       message.Error(),
				Funcionalidad: "GetPagoResultado",
			}

			err := s.utilService.CreateLogService(log)

			if err != nil {
				mensaje := fmt.Sprintf("Crear Log: %s. GetPagoTipo: %s", err.Error(), erro.Error())
				logs.Error(mensaje)
			}
		}

		if erro != nil {
			logs.Error(erro.Error())
		}
	}

	return &response, nil
}

func (s *service) CheckPrisma() error {
	check, err := s.prismaService.CheckService()
	if err != nil {
		return err
	}
	if !check {
		return fmt.Errorf("el servicio de prisma no está disponible")
	}
	return nil
}

func (s *service) GetBilling(uuid string) (*bytes.Buffer, error) {
	filtroPago := filtros.PagoFiltro{
		Uuids: []string{uuid},
	}
	pago, err := s.repository.GetPaymentByUuid(filtroPago)
	if err != nil {
		return nil, err
	}

	intento, err := s.repository.GetValidPagointentoByPagoId(int64(pago.ID))
	if err != nil {
		return nil, err
	}

	medioPago, err := s.repository.GetMediopago(map[string]interface{}{"id": intento.MediopagosID})
	if err != nil {
		return nil, err
	}

	channel, err := s.repository.GetChannelById(uint(medioPago.ChannelsID))
	if err != nil {
		return nil, err
	}

	pagotipo, err := s.repository.GetPagotipoById(pago.PagostipoID)
	if err != nil {
		return nil, err
	}
	logs.Info(pagotipo)
	cuenta, err := s.repository.GetCuentaById(int64(pagotipo.CuentasID))
	if err != nil {
		return nil, err
	}

	cliente, err := s.repository.GetClienteByApikey(cuenta.Apikey)
	if err != nil {
		return nil, err
	}

	// generar el pdf del comprobante de pago con los datos
	file, err := _getBillingPdf(pago, cliente, channel, intento)

	if err != nil {
		message := "no se pudo generar el comprobante de pago en pdf: " + err.Error()
		logs.Error(message)
	}

	return &file, nil
}

func (s *service) GetTarjetas() (*[]entities.Mediopago, error) {
	return s.repository.GetMediosDePagos()
}

func verificarExpiracionQR(fecha time.Time) (expirado bool, err error) {
	// Obtener la hora actual
	ahora := time.Now()

	// Calcular la diferencia de tiempo en minutos
	diferencia := ahora.Sub(fecha).Minutes()

	// Verificar si han pasado 9 minutos
	expirado = diferencia > 9

	return expirado, nil
}

func (s *service) GetMatenimietoSistema() (estado bool, fecha time.Time, erro error) {
	estado, fecha, err := s.utilService.GetMatenimietoSistemaService()
	if err != nil {
		erro = fmt.Errorf("el servicio no está disponible")
		return
	}
	// filtro := filtros.ConfiguracionFiltro{
	// 	Nombre: "ESTADO_APLICACION",
	// }
	// estadoConfiguracion, err := s.utilService.GetConfiguracionService(filtro)
	// if err != nil {
	// 	estado = true
	// 	erro = fmt.Errorf("el servicio no está disponible")
	// 	return
	// }
	// if estadoConfiguracion.Valor != "sin valor" {
	// 	fecha, err = time.Parse(time.RFC3339, estadoConfiguracion.Valor)
	// 	if err != nil {
	// 		estado = true
	// 		logs.Error("error al convertir fecha de configuración")
	// 		erro = fmt.Errorf("el servicio no está disponible")
	// 		return
	// 	}
	// 	if !fecha.IsZero() {
	// 		estado = true
	// 		return
	// 	}
	// }
	// estado = false
	return
}

/* ------------------ Funciones Propias ------------------ */

/*
Autor: Jose Alarcon
Fecha: 21/06/2022
Descripción: webhook notificacion de pago al cliente
Verificar que el cliente tenga una url configurada y que ademas para notifcar el pago sea exitoso
*/
func (s *service) notificacionPago(response dtos.ResultadoResponse, tipo *entities.Pagotipo) error {
	var result []dtos.ResultadoResponse
	result = append(result, response)
	notificacionPago := dtos.ResultadoResponseWebHook{
		Url:               tipo.BackUrlNotificacionPagos,
		ResultadoResponse: result,
	}

	if err := s.webhook.NotificarPago(notificacionPago); err != nil {
		return err
	}

	peticionWebHook := dtos.RequestWebServicePeticion{
		Operacion: "NotificarPago",
		Vendor:    "WebHook",
	}

	if err := s.utilService.CrearPeticionesService(peticionWebHook); err != nil {
		logs.Error("no se pudo registrar la peticion" + err.Error())
	}

	if err := s.repository.UpdateEstadoNotificadoInicial(uint(response.NumeroOperacion)); err != nil {
		logs.Error(fmt.Sprintf("no se pudo actualizar el estado inicial notificado: %v", response.NumeroOperacion))
	} else {
		logs.Info(fmt.Sprintf("webhook:se actualizo con éxito el estado inicial notificado, del pago: %v", response.NumeroOperacion))
	}

	return nil
}

func (s *service) notificacionPagoOnline(response dtos.ResultadoResponse, tipo *entities.Pagotipo, online bool) error {
	var result []dtos.ResultadoResponse
	result = append(result, response)
	notificacionPago := dtos.ResultadoResponseWebHook{
		Url:               tipo.BackUrlNotificacionPagos,
		ResultadoResponse: result,
	}

	if err := s.webhook.NotificarPago(notificacionPago); err != nil {
		return err
	}

	peticionWebHook := dtos.RequestWebServicePeticion{
		Operacion: "NotificarPago",
		Vendor:    "WebHook",
	}

	if err := s.utilService.CrearPeticionesService(peticionWebHook); err != nil {
		logs.Error("no se pudo registrar la peticion" + err.Error())
	}

	if err := s.repository.UpdateEstadoNotificadoInicial(uint(response.NumeroOperacion)); err != nil {
		logs.Error(fmt.Sprintf("no se pudo actualizar el estado inicial notificado: %v", response.NumeroOperacion))
	} else {
		logs.Info(fmt.Sprintf("webhook:se actualizo con éxito el estado inicial notificado, del pago: %v", response.NumeroOperacion))
	}

	if err := s.repository.UpdateEstadoNotificadoOnline(uint(response.ID)); err != nil {
		logs.Error(fmt.Sprintf("no se pudo actualizar el estado inicial notificado: %v", response.NumeroOperacion))
	} else {
		logs.Info(fmt.Sprintf("webhook:se actualizo con éxito el campo notificado online, del pagointento: %v", response.ID))
	}

	return nil
}

func (s *service) validarMontos(items []entities.Pagoitems, total entities.Monto) error {
	var totalItems entities.Monto
	for _, t := range items {
		totalItems += entities.Monto(int64(t.Quantity) * int64(t.Amount))
	}
	if totalItems != total {
		return fmt.Errorf("el total de los items no coincide con el total del pago")
	}
	return nil
}

// obtener los items de un pago para devolverlos en array de bytes con el fin de asignarlos a una respuesta
func getPaymentItemsToResponse(pago *entities.Pago) ([]byte, error) {
	// convierto los items del pago para mostrarlos en el frontend
	items := make([]dtos.PagoResponseItems, 0)
	if len(pago.Pagoitems) > 0 {
		for _, pago_item := range pago.Pagoitems {
			items = append(items, dtos.PagoResponseItems{
				Quantity:    int64(pago_item.Quantity),
				Description: pago_item.Description,
				Amount:      pago_item.Amount.Float64(),
				Identifier:  pago_item.Identifier,
			})
		}
	}
	byteItems, err := json.Marshal(items)
	return byteItems, err
}

func ctxWithClienteID(ctx context.Context, id uint) context.Context {
	audit := ctx.Value(entities.AuditUserKey{}).(entities.Auditoria)
	audit.CuentaID = id
	newCtx := context.WithValue(ctx, entities.AuditUserKey{}, audit)
	return newCtx
}

func formatFechaString(fecha time.Time, formatoFecha string) string {
	fechaStr := fecha.Format(formatoFecha)
	fechaArrayStr := strings.Split(fechaStr[0:10], "-")
	fechaVto := fmt.Sprintf("%v-%v-%v", fechaArrayStr[2], fechaArrayStr[1], fechaArrayStr[0])
	return fechaVto
}
func qrExpirado(pagointentoQr *entities.Pagointento) bool {
	hoy := time.Now().Local()
	tiempoExpiracion := 10
	diferencia := hoy.Sub(pagointentoQr.CreatedAt)
	minutos := diferencia.Minutes() // en float64
	if tiempoExpiracion != 0 && minutos > float64(tiempoExpiracion) {
		return true
	}
	fmt.Println(minutos / 100)
	return false
}

func _getPaymentExpiration(pago *entities.Pago) (erro error) {
	hoy := time.Now()
	tiempoExpiracion := pago.Expiration
	cadenaFechaHora := hoy.Format("2006-01-02 15:04:05")
	layout := "2006-01-02 15:04:05"
	// Para eliminar la diferencia -3 y comparar correctamente
	fechaHoy, _ := time.Parse(layout, cadenaFechaHora)

	// cheackear fecha y hora de expiracion
	if !pago.FechaHoraExpiracion.IsZero() {
		diferencia := pago.FechaHoraExpiracion.Sub(fechaHoy)
		minutos := diferencia.Minutes()
		if math.Trunc(minutos) < 1 {
			return fmt.Errorf("el pago expiró, vuelva a generarlo")
		}
	} else {
		diferencia := hoy.Sub(pago.CreatedAt)
		minutos := diferencia.Minutes()
		if tiempoExpiracion != 0 && minutos > float64(tiempoExpiracion) {
			return fmt.Errorf("el pago expiró, vuelva a generarlo")
		}
	}

	return
}

func _addUrlParam(baseURL, key, value string) string {

	u, err := url.Parse(baseURL)
	if err != nil {
		logs.Error("_addUrlParam: error al adicionar uuid como parametro en success back url: " + err.Error())
	}

	q := u.Query()
	q.Add(key, value)
	u.RawQuery = q.Encode()

	return u.String()
}

/* -------------------------------- Funciones PDF Comprobante ------------------------------------------- */

func getHeaderAndContent(pagoItems *[]entities.Pagoitems) (header []string, contents [][]string) {
	// La cabecera de la tabla. Los nombres de las columnas
	header = []string{"Transacción", "Producto", "Cantidad", "Precio"}

	// cantidad de items del pago
	size := len(*pagoItems)
	items := make([][]string, size)

	for i, x := range *pagoItems {
		identificador := ""
		if x.Identifier != "" {
			identificador = x.Identifier
		}
		// identificador, descripcion, cantidad, monto
		items[i] = []string{
			identificador, x.Description, fmt.Sprint(x.Quantity), strconv.FormatFloat(x.Amount.Float64(), 'f', 2, 64),
		}
	}

	contents = items

	return header, contents
}

func getTelCoGreenColor() color.Color {
	return color.Color{
		Red:   195,
		Green: 216,
		Blue:  46,
	}
}

func getHeaderTextColor() color.Color {
	return color.NewBlack()
}

func getTelCoSoftBlueColor() color.Color {
	return color.Color{
		Red:   0,
		Green: 184,
		Blue:  241,
	}
}

func getDarkGrayColor() color.Color {
	return color.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

// size dinamico segun long del texto que se recibe
func _resolveColumnWidthSize(texto string) (colSize, colSpaceSize int) {
	long := len(texto)
	var columMaxSize int = 12
	switch {
	case long <= 30:
		colSize = 3
		colSpaceSize = columMaxSize - colSize
	case long > 30 && long < 43:
		colSize = 4
		colSpaceSize = columMaxSize - colSize
	case long >= 43:
		colSize = 5
		colSpaceSize = columMaxSize - colSize
	default:
		colSize = 4
		colSpaceSize = columMaxSize - colSize
	}
	return
}

func _buildHeading(m pdf.Maroto, cliente *entities.Cliente, intento *entities.Pagointento) {
	green := getTelCoGreenColor()
	negro := getHeaderTextColor()
	blanco := color.NewWhite()

	// RegisterHeader
	m.RegisterHeader(func() {
		m.Row(50, func() {
			m.Col(12, func() {

				// err = m.FileImage(filepath.FromSlash("./assets/images/cabecera_recibo.png"), props.Rect{
				// 	Top: 30,
				// })

				// En producción o en el servidor
				err := m.FileImage(filepath.Join(filepath.Base(config.DIR_BASE), "api", "assets", "images", "cabecera_recibo.png"), props.Rect{})

				if err != nil {
					logs.Error("_buildHeading: la imagen no se pudo cargar al intentar crear el comprobante de pago pdf: " + err.Error())
				}
			})
		})
	})

	texto := cliente.Cliente + " - CUIT " + cliente.Cuit

	// segun long de texto se adapta columnas
	colSize, colSpaceSize := _resolveColumnWidthSize(texto)

	m.Row(10, func() {
		m.Col(uint(colSpaceSize), func() {
			m.Text(intento.HolderName, props.Text{Size: 8, Top: 5, Left: 10})
		})
		m.SetBackgroundColor(green)
		m.Col(uint(colSize), func() {
			m.Text(texto, props.Text{
				Style: consts.Normal,
				Color: negro,
				Align: consts.Right,
				Size:  8,
				Top:   3,
				Right: 2,
			})
		})
		m.SetBackgroundColor(blanco)
	})
	m.Row(10, func() {
		m.Col(uint(colSpaceSize), func() {
			m.Text("CUIL/DNI: "+intento.HolderNumber, props.Text{Size: 8, Left: 10})
		})
		m.Col(uint(colSize), func() {
			fechaEmision := time.Now().Format("02-01-2006")

			m.Text("Fecha de emisión: "+fechaEmision, props.Text{
				Style: consts.Normal,
				Color: negro,
				Align: consts.Right,
				Size:  8,
				Top:   3,
				Right: 2,
			})
		})
		m.SetBackgroundColor(blanco)
	})
}

func _buildBodyList(m pdf.Maroto, pago *entities.Pago, channel entities.Channel, intento *entities.Pagointento) {
	// set de colores
	celeste := getTelCoSoftBlueColor()
	blanco := color.NewWhite()
	darkGrayColor := getDarkGrayColor()

	// contenido y cabeceras de la tabla
	header, contents := getHeaderAndContent(&pago.Pagoitems)

	// solo mostrar cuando el channel es debit o credit
	if channel.Channel == "DEBIT" || channel.Channel == "CREDIT" {
		m.Row(10, func() {
			m.Col(12, func() {
				m.Text("Verás este pago en tu resumen como TelCo Wee!", props.Text{
					Top:   5,
					Style: consts.Italic,
					Align: consts.Center,
					Color: color.NewBlack(),
				})
			})
		})
	}

	// Referencia del pago
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Referencia de pago: "+pago.ExternalReference, props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
	})
	// una linea de espacio o separacion
	m.Line(5, props.Line{
		Color: blanco,
	})

	m.SetBackgroundColor(celeste)

	// Tabla de items del pago
	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			GridSizes: []uint{3, 4, 2, 3},
			Size:      10,
			Style:     consts.Bold,
			Color:     blanco,
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 4, 2, 3},
		},
		Align:                consts.Center,
		HeaderContentSpace:   1,
		AlternatedBackground: &blanco,
		Line:                 true,
		LineProp: props.Line{
			Color: celeste,
		},
		VerticalContentPadding: 7, // alto de fila
	})

	// Totales y vencimientos, segun canales de pago
	if channel.Channel == "OFFLINE" {
		//primer vencimiento
		m.Row(12, func() {

			m.Col(1, func() {
				m.Text("Primer Vto.:", props.Text{
					Top:   2,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Left,
				})
				m.ColSpace(1)
				fecha := pago.FirstDueDate
				m.Text(formatFechaString(fecha, "2006-01-02T00:00:00Z"), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})

			m.ColSpace(4)
			m.Col(2, func() {
				m.Text("Total:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Right,
				})
			})
			m.Col(4, func() {
				var total float64 = 0
				if intento.Valorcupon > 0 {
					total = intento.Valorcupon.Float64()
				} else {
					total = intento.Amount.Float64()
				}
				m.Text(strconv.FormatFloat(total, 'f', 2, 64), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})
		})

		//Segundo vencimiento
		m.Row(12, func() {
			// m.ColSpace(7) // Incluye un espacio en blanco de 7 columnas a la izquierda
			// m.ColSpace(1) // Incluye un espacio en blanco de 7 columnas a la izquierda
			m.Col(1, func() {
				m.Text("Segundo Vto.:", props.Text{
					Top:   2,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Left,
				})
				m.ColSpace(1)
				fecha := pago.SecondDueDate
				m.Text(formatFechaString(fecha, "2006-01-02T00:00:00Z"), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})

			m.ColSpace(4)
			m.Col(2, func() {
				m.Text("Total:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Right,
				})
			})
			m.Col(4, func() {
				var total float64 = 0

				total = pago.SecondTotal.Float64()

				m.Text(strconv.FormatFloat(total, 'f', 2, 64), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})
		})

		// mostrar logo de rapipago y multipagos
		m.Row(20, func() {
			m.Col(2, func() {
				// LOCAL
				// err := m.FileImage(filepath.FromSlash("./assets/images/logo-rapipago.png"), props.Rect{
				// 	Top: 11,
				// })

				// SERVER
				err := m.FileImage(filepath.Join(filepath.Base(config.DIR_BASE), "api", "assets", "images", "logo-rapipago.png"), props.Rect{
					Top: 11,
				})

				if err != nil {
					logs.Error("_buildFooter: la imagen no se pudo cargar al intentar crear el comprobante de pago pdf: " + err.Error())
				}
			})
			m.Col(2, func() {})
			m.Col(2, func() {
				// LOCAL
				// err := m.FileImage(filepath.FromSlash("./assets/images/logo-multipago.png"), props.Rect{
				// 	Top: 10,
				// })

				// SERVER
				err := m.FileImage(filepath.Join(filepath.Base(config.DIR_BASE), "api", "assets", "images", "logo-multipago.png"), props.Rect{
					Top: 10,
				})

				if err != nil {
					logs.Error("_buildFooter: la imagen no se pudo cargar al intentar crear el comprobante de pago pdf: " + err.Error())
				}
			})
		})

		// se genera el codigo de barra
		m.Row(15, func() { //15
			m.Col(6, func() {

				_ = m.Barcode(intento.Barcode, props.Barcode{
					Percent: 0,
					Proportion: props.Proportion{
						Width:  20,
						Height: 2,
					},
				})
				m.Text(intento.Barcode, props.Text{
					Top:    12,
					Family: "",
					Style:  consts.Bold,
					Size:   9,
					Align:  consts.Center,
				})
			})
			m.ColSpace(6)
		})
	}

	if channel.Channel == "MULTIPAGOS" {
		//primer vencimiento
		m.Row(12, func() {

			m.Col(1, func() {
				m.Text("Primer Vto.:", props.Text{
					Top:   2,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Left,
				})
				m.ColSpace(1)
				fecha := pago.FirstDueDate
				m.Text(formatFechaString(fecha, "2006-01-02T00:00:00Z"), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})

			m.ColSpace(4)
			m.Col(2, func() {
				m.Text("Total:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Right,
				})
			})
			m.Col(4, func() {
				var total float64 = 0
				if intento.Valorcupon > 0 {
					total = intento.Valorcupon.Float64()
				} else {
					total = intento.Amount.Float64()
				}
				m.Text(strconv.FormatFloat(total, 'f', 2, 64), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})
		})

		//Segundo vencimiento
		m.Row(12, func() {
			// m.ColSpace(7) // Incluye un espacio en blanco de 7 columnas a la izquierda
			// m.ColSpace(1) // Incluye un espacio en blanco de 7 columnas a la izquierda
			m.Col(1, func() {
				m.Text("Segundo Vto.:", props.Text{
					Top:   2,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Left,
				})
				m.ColSpace(1)
				fecha := pago.SecondDueDate
				m.Text(formatFechaString(fecha, "2006-01-02T00:00:00Z"), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})

			m.ColSpace(4)
			m.Col(2, func() {
				m.Text("Total:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Right,
				})
			})
			m.Col(4, func() {
				var total float64 = 0

				total = pago.SecondTotal.Float64()

				m.Text(strconv.FormatFloat(total, 'f', 2, 64), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})
		})

		m.Row(12, func() {

			m.Col(5, func() {

				m.Text("Codigo Multipagos:", props.Text{
					Top:   5,
					Style: consts.Bold,
				})
			})
			m.ColSpace(1)
			m.Col(6, func() {
				codigo_multipago := pago.PagoIntentos[len(pago.PagoIntentos)-1].Barcode
				m.Text(codigo_multipago, props.Text{
					Top:   5,
					Style: consts.Bold,
				})
			})

		})

	}

	if channel.Channel == "CREDIT" {
		// Si existe valor cupon, y por lo tanto es un pago en cuotas, se muestra el costo financiado
		if intento.Valorcupon != 0 {
			dif := intento.Valorcupon - intento.Amount
			costo_financiero := dif.Float64()
			m.Row(7, func() {
				m.ColSpace(7) // Incluye un espacio en blanco de 7 columnas a la izquierda
				m.Col(2, func() {
					m.Text("Costo Financiero:", props.Text{
						Top:   5,
						Style: consts.Bold,
						Size:  8,
						Align: consts.Right,
					})
				})
				m.Col(3, func() {
					m.Text(strconv.FormatFloat(costo_financiero, 'f', 2, 64), props.Text{
						Top:   5,
						Style: consts.Bold,
						Size:  8,
						Align: consts.Center,
					})
				})
			})
		}
		m.Row(7, func() {
			m.ColSpace(4)
			m.Col(4, func() {
				m.Text("Total:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Right,
				})
			})
			m.Col(5, func() {
				var total float64 = 0
				if intento.Valorcupon > 0 {
					total = intento.Valorcupon.Float64()
				} else {
					total = intento.Amount.Float64()
				}
				m.Text(strconv.FormatFloat(total, 'f', 2, 64), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})
		})
	}

	if channel.Channel == "DEBIN" || channel.Channel == "DEBIT" || channel.Channel == "QR" {
		m.Row(7, func() { //15
			m.ColSpace(4)
			m.Col(4, func() {
				m.Text("Total:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Right,
				})
			})
			m.Col(5, func() {
				var total float64 = 0
				total = intento.Amount.Float64()
				m.Text(strconv.FormatFloat(total, 'f', 2, 64), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
			})
		})
	}

	// Medio de Pago, Numero de Solicitud, Numero de Operacion, Codigo Autorizacion
	m.Row(4, func() {
		var cardLastFour string
		cardLastFour = ""
		if channel.Channel == "CREDIT" || channel.Channel == "DEBIT" {
			cardLastFour = ". Los últimos 4 dígitos de su tarjeta son: " + intento.CardLastFourDigits
		}
		m.Col(15, func() {
			m.Text("Medio de Pago: "+channel.Nombre+cardLastFour, props.Text{
				Top:   16,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
				Color: darkGrayColor,
			})

		})
	})
	m.Row(4, func() {
		m.Col(15, func() {
			nroSolicitud := strconv.Itoa(int(pago.ID))
			m.Text("Nro. Solicitud: "+nroSolicitud, props.Text{
				Top:   16,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
				Color: darkGrayColor,
			})
		})
	})
	m.Row(4, func() {
		m.Col(15, func() {
			nroOperacion := strconv.Itoa(int(intento.ID))
			m.Text("Nro. Op.: "+nroOperacion, props.Text{
				Top:   16,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
				Color: darkGrayColor,
			})
		})

	})
	m.Row(4, func() {
		var codigoAutorizacion string
		if channel.Channel == "DEBIN" || channel.Channel == "OFFLINE" || channel.Channel == "QR" {
			codigoAutorizacion = intento.ExternalID
		}
		if channel.Channel == "CREDIT" || channel.Channel == "DEBIT" {
			codigoAutorizacion = intento.AuthorizationCode
		}
		m.Col(15, func() {
			m.Text("Código Autorización: "+codigoAutorizacion, props.Text{
				Top:   10,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Left,
				Color: darkGrayColor,
			})
		})
	})
}

func _buildFooter(m pdf.Maroto) {
	m.RegisterFooter(func() {
		m.Row(50, func() {
			m.Col(12, func() {

				// err = m.FileImage(filepath.FromSlash("./assets/images/footer_recibo.png"), props.Rect{
				// 	Top: 30,
				// })

				// En producción o en el servidor
				err := m.FileImage(filepath.Join(filepath.Base(config.DIR_BASE), "api", "assets", "images", "footer_recibo.png"), props.Rect{
					Top: 30,
				})

				if err != nil {
					logs.Error("_buildHeading: la imagen no se pudo cargar al intentar crear el comprobante de pago pdf: " + err.Error())
				}
			})
		})
	})
}

func _getBillingPdf(pago *entities.Pago, cliente *entities.Cliente, channel entities.Channel, intento *entities.Pagointento) (file bytes.Buffer, erro error) {

	// instancia de objeto PDF
	m := pdf.NewMaroto(consts.Portrait, consts.A4)

	m.SetPageMargins(0, 0, 0)

	// Header
	_buildHeading(m, cliente, intento)

	m.SetPageMargins(10, 0, 10)

	// Body
	_buildBodyList(m, pago, channel, intento)

	m.SetPageMargins(0, 0, 0)

	// Footer
	_buildFooter(m)

	// set de variables de retorno
	file, erro = m.Output()
	return
}

func (s *service) HashOperacionTarjeta(number string, pagointento_id int64) (status bool, erro error) {
	//Hashear card number
	textoPlano := number
	hash := md5.Sum([]byte(textoPlano))
	hashstring := hex.EncodeToString(hash[:])

	hasheado := entities.Uuid{
		Uuid: hashstring,
	}

	s.repository.SaveHasheado(&hasheado, uint(pagointento_id))

	return
}

func (s *service) ControlTarjetaHash(number string) (status bool, erro error) {
	//Hashear card number
	textoPlano := number
	hash := md5.Sum([]byte(textoPlano))
	hashstring := hex.EncodeToString(hash[:])

	status, err := s.repository.GetHasheado(hashstring)

	if err != nil {
		erro = err
		return
	}

	return
}

func (s *service) GetPagoStatus(uuid string) (status bool, erro error) {
	// valida uuid
	if len(uuid) <= 0 {
		return false, fmt.Errorf("debe enviar código único del pago, envió: %s", uuid)
	}

	if ok, err := s.commons.IsValidUUID(uuid); !ok {
		return false, fmt.Errorf("el identificador del pago no es válido: %w", err)
	}

	filtroPago := filtros.PagoFiltro{
		Uuids: []string{uuid},
	}
	// busca pago por uuid en la base de datos
	pago, err := s.repository.GetPaymentByUuid(filtroPago)
	if err != nil {
		return false, fmt.Errorf("error al obtener pago: %s", err.Error())
	}

	// si ya pagó (pagoestado_id != pending) devuelve un error
	if pago.PagoestadosID != 1 {
		filtroMedioPago := make(map[string]interface{})
		if len(pago.PagoIntentos) == 0 {
			return false, fmt.Errorf("el pago seleccionado no posee ningún intento de pago")
		}

		filtroMedioPago["id"] = pago.PagoIntentos[len(pago.PagoIntentos)-1].MediopagosID
		medioPago, erro := s.repository.GetMediopago(filtroMedioPago)
		if erro != nil {
			return false, fmt.Errorf("error al obtener medio de pago: %s", erro.Error())
		}
		return false, fmt.Errorf("el pago ya fue procesado a través del medio de pago %v", medioPago.Mediopago)
	}

	return true, nil
}

func (s *service) GetRapipagoQuery(request rapipago.RequestRapipagoConsulta) (response rapipago.ResponseRapipagoConsulta, err error) {
	// valida uuid
	if len(request.IdClave) <= 0 {
		// err = fmt.Errorf("debe enviar una clave de identificacion")
		response.CodigoRespuesta = "9" //Parámetros incorrectos o faltantes
		response.Msg = "Debe enviar una clave de identificacion"
		return
	}

	var barcodes []string
	var dni string

	switch len(request.IdClave) {
	case 7:
		dni = request.IdClave
	case 8:
		dni = request.IdClave
	case 48:
		barcodes = append(barcodes, request.IdClave)
	default:
		// err = fmt.Errorf("debe enviar una clave de identificacion con longitud válida.")
		response.CodigoRespuesta = "9" //Parámetros incorrectos o faltantes
		response.Msg = "Debe enviar una clave de identificacion con longitud válida."
		return
	}

	filtro := filtros.RapipagosFiltro{
		Barcodes:  barcodes,
		DNI:       dni,
		ConEstado: true,
	}

	// busca pago por los filtros en la base de datos
	pagos, erro := s.repository.GetPagosRapipago(filtro)

	if err != nil {
		logs.Error(erro.Error())
		return
	}

	for _, pago := range pagos {

		importeControl := strconv.FormatFloat(pago.FirstTotal.Float64(), 'f', 2, 64)
		fechaHoy := time.Now()
		if fechaHoy.After(pago.FirstDueDate) {
			importeControl = strconv.FormatFloat(pago.SecondTotal.Float64(), 'f', 2, 64)
		}

		lastPagointento := len(pago.PagoIntentos) - 1
		importe := importeControl
		importe = commons.AgregarCerosString(importe, 11, "LEFT")
		pagoResponse := rapipago.FacturaRapipago{
			IdNumero:         pago.PagoIntentos[lastPagointento].Barcode,
			Barra:            pago.PagoIntentos[lastPagointento].Barcode,
			Importe:          importe,
			FechaEmision:     pago.PagoIntentos[lastPagointento].PaidAt.Format("2006-01-02"),
			FechaVencimiento: pago.SecondDueDate.Format("2006-01-02"),
			Texto1:           pago.Description,
		}

		response.Facturas = append(response.Facturas, pagoResponse)

		response.Nombre = pago.PayerName
	}

	if len(response.Facturas) == 0 {
		response.Facturas = []rapipago.FacturaRapipago{}
		response.CodigoRespuesta = "6" //No existe registro (el cliente existe, pero no tiene deuda)
		response.Msg = "no se encontro el pago requerido."
		response.IdClave = request.IdClave
		response.CodTrx = request.CodTrx
		return
	}

	response.IdClave = request.IdClave
	response.CodigoRespuesta = "0"
	response.Msg = "Trx ok"
	response.CodTrx = request.CodTrx

	return
}

func (s *service) PostRapipagoPago(ctx context.Context, request rapipago.RequestRapipagoConsulta) (response rapipago.ResponseRapipagoImputacion, err error) {
	// valida uuid
	if len(request.Barra) <= 0 {
		// err = fmt.Errorf("debe enviar un codigo de barra.")
		response.CodigoRespuesta = "9" //Parámetros incorrectos o faltantes
		response.Msg = "Parámetros incorrectos o faltantes"
		return
	}

	var barcodes []string

	barcodes = append(barcodes, request.Barra)

	filtroRapipago := filtros.RapipagosFiltro{
		Barcodes:   barcodes,
		ConCliente: true,
	}

	// busca pago por codigo en la base de datos
	pagos, erro := s.repository.GetPagosRapipago(filtroRapipago)

	if erro != nil {
		logs.Error(erro.Error())
		response.CodigoRespuesta = "1" //Error en BD
		response.Msg = "Error en BD"
		return
	}

	if len(pagos) != 1 {
		response.CodigoRespuesta = "6" //No existe registro (el cliente existe, pero no tiene deuda)
		response.Msg = "No se encontro el pago requerido."
		return
	}

	pago := pagos[0]

	if pago.PagoestadosID == 4 {
		response.CodigoRespuesta = "14" //Pago registrado con anterioridad
		response.Msg = "Pago registrado con anterioridad."
		return
	}

	filtroMedioPago := make(map[string]interface{})
	filtroMedioPago["mediopago"] = request.Canal
	medioPago, err := s.repository.GetMediopago(filtroMedioPago)
	if err != nil {
		response.CodigoRespuesta = "1"
		response.Msg = "Error al obtener informacion de medio de pago"
		return
	}

	// Separo intento de pago con codigo
	var pagointento entities.Pagointento
	for _, PI := range pago.PagoIntentos {
		if PI.Barcode == request.Barra {
			pagointento = PI
			break
		}
	}

	// Obtiene pagoestado con id 4 ("Pagado")
	pagoestado, err := s.repository.GetPagoEstado(4)
	if err != nil {
		logs.Error(err.Error())
		response.CodigoRespuesta = "1" //Error en BD
		response.Msg = "Error en BD"
		return
	}

	pago.PagoestadosID = int64(pagoestado.ID)

	exito, err := s.repository.UpdatePagoEstado(ctx, pago)
	if err != nil {
		logs.Error(err.Error())
		response.CodigoRespuesta = "1" //Error en BD
		response.Msg = "Error en BD"
		return
	}

	// // Si pago actualizado exitoso notifico con webhook
	notifResponse := dtos.ResultadoResponse{
		ID:                int64(pagointento.ID),
		Estado:            pagointento.StateComment,
		EstadoPago:        pagoestado.Nombre,
		Exito:             exito,
		Uuid:              pagointento.Barcode,       // informacion para comprobante pdf
		Channel:           medioPago.Channel.Channel, //
		Description:       pago.Description,
		FirstDueDate:      pago.FirstDueDate.Format("02-01-2006"),
		FirstTotal:        pagointento.Amount.Float64(),
		SecondDueDate:     pago.SecondDueDate.Format("02-01-2006"),
		SecondTotal:       pago.SecondTotal.Float64(),
		PayerName:         pago.PayerName,
		PayerEmail:        pago.PayerEmail,
		ExternalReference: pago.ExternalReference,
		Metadata:          pago.Metadata,
		PdfUrl:            pago.PdfUrl,
		CreatedAt:         pagointento.CreatedAt.String(),
		ImportePagado:     pagointento.Amount.Float64(),
		ClienteName:       pago.PagosTipo.Cuenta.Cliente.Cliente, // informacion para comprobante pdf
		ClienteCuit:       pago.PagosTipo.Cuenta.Cliente.Cuit,    // informacion para comprobante pdf
		Mediopago:         medioPago.Mediopago,                   // informacion para comprobante pdf
		NumeroOperacion:   pago.ID,
	}

	// notificacion de pago exitoso por medio del webhook
	if pago.PagosTipo.BackUrlNotificacionPagos != "" && notifResponse.Exito {
		if err := s.notificacionPagoOnline(notifResponse, &pago.PagosTipo, true); err != nil {
			logs.Info(fmt.Sprintf("webhook:no se pudo notificar el pago online %v: %s", notifResponse.NumeroOperacion, err.Error()))
		}

	}

	response.IdNumero = pago.PagoIntentos[len(pago.PagoIntentos)-1].Barcode
	response.Barra = pago.PagoIntentos[len(pago.PagoIntentos)-1].Barcode
	response.CodTrx = request.CodTrx
	response.CodigoRespuesta = "0" //Transaccion aceptada
	response.Msg = "Trx ok"

	return
}

func (s *service) GetMultiPagoStatus(request multipagosdtos.RequestConsultaMultipago) (response multipagosdtos.ResponseConsultaMultipago, erro error) {
	// valida uuid

	var dni string
	var barcode string

	switch len(request.Id_clave) {
	case 7:
		dni = request.Id_clave
	case 8:
		dni = request.Id_clave
	case 48:
		barcode = request.Id_clave
	default:
		// erro = fmt.Errorf("debe enviar una clave de identificacion con longitud válida.")
		response.CodigoRespuesta = "9"
		response.Msg = "Debe enviar codigo barra con longitud válida del pago"
		return
	}

	filtroMultiPago := filtros.FiltroMultipago{
		Codigo:   barcode,
		ValorDoc: dni,
	}

	// busca pago por uuid en la base de datos
	pagos, err := s.repository.GetPaymentMultipago(filtroMultiPago)
	if err != nil {
		// erro = fmt.Errorf("error al obtener pago: %s", err.Error())
		response.CodigoRespuesta = "6" //No existe registro
		response.Msg = "No se encontraron pagos asociados pendientes."
		return
	}

	for _, pago := range pagos {

		importeControl := strconv.FormatFloat(pago.FirstTotal.Float64(), 'f', 2, 64)
		fechaHoy := time.Now()
		if fechaHoy.After(pago.FirstDueDate) {
			importeControl = strconv.FormatFloat(pago.SecondTotal.Float64(), 'f', 2, 64)
		}

		responseFactura := multipagosdtos.FacturaMultipago{
			Barra:            pago.PagoIntentos[len(pago.PagoIntentos)-1].Barcode,
			FechaEmision:     pago.PagoIntentos[len(pago.PagoIntentos)-1].PaidAt.Format("2006-01-02"),
			FechaVencimiento: pago.SecondDueDate.Format("2006-01-02"),
			Importe:          importeControl,
		}

		response.Facturas = append(response.Facturas, responseFactura)

		response.Nombre = pago.PagoIntentos[len(pago.PagoIntentos)-1].HolderName
	}

	response.CodTrx = request.CodTrx
	response.Id_clave = request.Id_clave
	response.CodigoRespuesta = "0"
	response.Msg = "Trx ok"

	return
}

func (s *service) GetMultiPagoControl(request multipagosdtos.RequestControlMultipago) (response multipagosdtos.ResponseControlMultipago, erro error) {

	fechaInicio, err := time.Parse("02-01-2006", request.FechaInicio)
	fechaFin, err := time.Parse("02-01-2006", request.FechaFin)

	filtroMultiPago := filtros.PagoFiltro{
		FechaPagoInicio:   fechaInicio,
		FechaPagoFin:      fechaFin,
		MedioPagoId:       39,
		CargaPagoIntentos: true,
	}

	pagos, err := s.repository.GetPayments(filtroMultiPago)
	if err != nil {
		// erro = fmt.Errorf("error al obtener pago: %s", err.Error())
		response.CodigoRespuesta = "6" //No existe registro
		response.Msg = "No se encontraron pagos asociados pendientes."
		return
	}

	for _, pago := range pagos {

		responseFactura := multipagosdtos.FacturaMultipagoControl{
			CodOperacion: pago.Uuid,
			Barra:        pago.PagoIntentos[len(pago.PagoIntentos)-1].Barcode,
		}

		response.Facturas = append(response.Facturas, responseFactura)

	}

	response.CodigoRespuesta = "0"
	response.Msg = "Trx ok"

	return
}

func (s *service) PostMultiPago(ctx context.Context, request multipagosdtos.RequestPagoMultipago) (response multipagosdtos.ResponsePagoMultipago, erro error) {

	// valida código de pago
	if len(request.Id_clave) != 48 {
		response.CodigoRespuesta = "9"
		response.Msg = "Debe enviar codigo barra con longitud válida del pago"
		return
	}

	// valida medio de pago no vacio
	if len(request.Canal) == 0 {
		response.CodigoRespuesta = "9"
		response.Msg = "Debe enviar un medio válido de pago"
		return
	}

	filtroMultipago := filtros.FiltroMultipago{
		Codigo: request.Id_clave,
	}

	// busca pago por uuid en la base de datos, campo codigo_multipago
	pagos, err := s.repository.GetPaymentMultipago(filtroMultipago)
	if err != nil {
		response.CodigoRespuesta = "6"
		response.Msg = "No se obtuvo la información de un pago pendiente"
		return
	}

	if len(pagos) != 1 {
		response.CodigoRespuesta = "5"
		response.Msg = "No se obtuvo información de un pago único"
		return
	}

	pago := pagos[0]

	if len(pago.PagoIntentos) == 0 {
		response.CodigoRespuesta = "4"
		response.Msg = "No se obtuvo la información del pago"
		return
	}

	if pago.PagoestadosID == 4 || pago.PagoestadosID == 7 {
		response.CodigoRespuesta = "14"
		response.Msg = "El pago ya fue registrado con anterioridad"
		return
	}

	// Separo intento de pago con codigo
	var pagointento entities.Pagointento
	for _, PI := range pago.PagoIntentos {
		if PI.Barcode == request.Id_clave {
			pagointento = PI
			break
		}
	}

	filtroMedioPago := make(map[string]interface{})
	filtroMedioPago["mediopago"] = request.Canal
	medioPago, err := s.repository.GetMediopago(filtroMedioPago)
	if err != nil {
		response.CodigoRespuesta = "1"
		response.Msg = "Error al obtener informacion de medio de pago"
		return
	}

	// actualizar ESTADO DE PAGO a PAID -> APROBADO
	pago.PagoestadosID = 4

	// actualizar la medio de pago con la recibida
	pagointento.MediopagosID = int64(medioPago.ID)
	// actualizar el pago y el pagointento
	exito, err := s.repository.UpdatePagoMP(ctx, &pago, &pagointento)
	if err != nil {
		response.CodigoRespuesta = "1"
		response.Msg = "Error al actuailizar informacion del pago"
		return
	}

	if !exito {
		response.CodigoRespuesta = "1"
		response.Msg = "Error al actuailizar informacion del pago"
		return
	}

	// // Si pago actualizado exitoso notifico con webhook
	notifResponse := dtos.ResultadoResponse{
		ID:                int64(pagointento.ID),
		Estado:            pagointento.StateComment,
		EstadoPago:        pago.PagoEstados.Nombre,
		Exito:             exito,
		Uuid:              pagointento.Barcode,       // informacion para comprobante pdf
		Channel:           medioPago.Channel.Channel, //
		Description:       pago.Description,
		FirstDueDate:      pago.FirstDueDate.Format("02-01-2006"),
		FirstTotal:        pagointento.Amount.Float64(),
		SecondDueDate:     pago.SecondDueDate.Format("02-01-2006"),
		SecondTotal:       pago.SecondTotal.Float64(),
		PayerName:         pago.PayerName,
		PayerEmail:        pago.PayerEmail,
		ExternalReference: pago.ExternalReference,
		Metadata:          pago.Metadata,
		PdfUrl:            pago.PdfUrl,
		CreatedAt:         pagointento.CreatedAt.String(),
		ImportePagado:     pagointento.Amount.Float64(),
		ClienteName:       pago.PagosTipo.Cuenta.Cliente.Cliente, // informacion para comprobante pdf
		ClienteCuit:       pago.PagosTipo.Cuenta.Cliente.Cuit,    // informacion para comprobante pdf
		Mediopago:         medioPago.Mediopago,                   // informacion para comprobante pdf
		NumeroOperacion:   pago.ID,
	}

	// // notificacion de pago exitoso por medio del webhook
	if pago.PagosTipo.BackUrlNotificacionPagos != "" && notifResponse.Exito {
		if err := s.notificacionPagoOnline(notifResponse, &pago.PagosTipo, true); err != nil {
			logs.Info(fmt.Sprintf("webhook:no se pudo notificar el pago online %v: %s", notifResponse.NumeroOperacion, err.Error()))
		}
	}

	response.CodTrx = request.CodTrx
	response.Id_clave = request.Id_clave
	response.CodigoRespuesta = "0"
	response.Msg = "Trx ok"
	response.CodOperacion = pago.Uuid
	return
}

func (s *service) ControlAdquirienteApikey(apikey string, adquiriente string) (control bool, err error) {
	return s.repository.GetApikeyAdquiriente(apikey, adquiriente)
}

func (s *service) RapipagoConfirmacionService(request rapipago.RequestRapipagoConfirmacion) (response rapipago.ResponseRapipagoConfirmacion, err error) {

	if err = request.Validate(); err != nil {
		response = request.ParseToResponse("9", "Parámetros incorrectos o faltantes")
		return
	}

	var barcodes []string
	barcodes = append(barcodes, request.Barra)
	filtroRapipago := filtros.RapipagosFiltro{
		Barcodes:  barcodes,
		ConEstado: true,
	}

	// busca pago por codigo en la base de datos
	pagos, err := s.repository.GetPagosRapipagoAprobados(filtroRapipago)

	if err != nil {
		logs.Error(err.Error())
		response = request.ParseToResponse("1", "Error en BD")
		return
	}

	if len(pagos) < 1 {
		err = errors.New("error al consultar pagos con datos de confirmacion")
		response = request.ParseToResponse("6", "No existe registro")
		return
	}

	response = request.ParseToResponse("0", "Trx ok")
	lastPagointentoIndex := len(pagos[0].PagoIntentos) - 1
	response.FechaHoraOperacion = pagos[0].PagoIntentos[lastPagointentoIndex].PaidAt.Format("2006-01-02 15:04:05")

	return
}

func (s *service) NotificarPagos(listaPagos []webhookDto.WebhookResponse) (pagoupdate []uint) {
	// var pagosupdate []uint
	for _, webhook := range listaPagos {
		var pagosinformados []uint
		erro := s.webhook.NotificarPagos(webhook)
		if erro != nil {
			logs.Info(erro) //solo informar el error continuar enviando los pagos a los demas clientes
			log := entities.Log{
				Tipo:          entities.Error,
				Funcionalidad: "NotificarPagos",
				Mensaje:       fmt.Sprintf("webhook: no se pudo notificar pagos al cliente .: %s%s", erro, webhook.Url),
			}

			err := s.utilService.CreateLogService(log)

			if err != nil {
				logs.Info(utildtos.ERROR_LOG + "NotificarPagos." + erro.Error())
			}
		} else {
			logs.Info(fmt.Sprintf("webhook: se notifico con exito al cliente:%s", webhook.Url))
			for _, pago := range webhook.ResultadoResponseWebHook {
				pagoupdate = append(pagoupdate, uint(pago.Id))
				pagosinformados = append(pagosinformados, uint(pago.Id))
			}
			logs.Info(fmt.Sprintf("webhook: se notifico con exito al cliente los siguientes pagos:%v", pagosinformados))

			// crear logs
			log := entities.Log{
				Tipo:          entities.Info,
				Funcionalidad: "NotificarPagos",
				Mensaje:       fmt.Sprintf("webhook: se notifico con exito al cliente:%s", webhook.Url),
			}

			err := s.utilService.CreateLogService(log)

			if err != nil {
				logs.Info(utildtos.ERROR_LOG + "NotificarPagos." + erro.Error())
			}
		}
	}
	return
}

func (s *service) GetEstadoAppService() (err error) {
	filtroConfiguraciones := filtros.ConfiguracionFiltro{
		Nombre: "ESTADO_APLICACION",
	}

	_, erro := s.repository.GetEstadoAplicacionRepository(filtroConfiguraciones)

	if erro != nil {
		return
	}

	return
}

func (s *service) GetCuentaByApiKey(apikey string) (result bool, erro error) {
	cuenta, err := s.repository.GetCuentaByApiKey(apikey)
	if err != nil {
		log := entities.Log{
			Tipo:          "info",
			Funcionalidad: "GetCuentaByApiKey",
			Mensaje:       err.Error(),
		}
		err = s.utilService.CreateLogService(log)
		if err != nil {
			logs.Error("error al intentar registrar logs de erro en GetCuentaByApiKey")
		}
		erro = errors.New("api-key invalido")
		return
	}
	result = false
	if len(cuenta.Apikey) > 0 {
		result = true
	}
	return
}
