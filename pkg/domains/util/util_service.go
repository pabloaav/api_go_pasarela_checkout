package util

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/administraciondtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/filtros/administracion"
)

type UtilService interface {
	CreateNotificacionService(notificacion entities.Notificacione) (erro error)
	CreateLogService(log entities.Log) (erro error)
	LogError(erro string, funcionalidad string)
	CrearPeticionesService(peticiones dtos.RequestWebServicePeticion) (erro error)

	//CONFIGURACIONES
	GetConfiguracionService(filtro filtros.ConfiguracionFiltro) (configuracion administraciondtos.ResponseConfiguracion, erro error)
	GetConfiguracionesService(filtro filtros.ConfiguracionFiltro) (configuraciones []administraciondtos.ResponseConfiguracion, erro error)
	CreateConfiguracionService(config administraciondtos.RequestConfiguracion) (id uint, erro error)
	FirstOrCreateConfiguracionService(nombre string, descripcion string, valor string) (key string, erro error)

	//Redondeo
	ToFixed(num float64, precision int) float64

	//Calcular Comisiones
	BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto, importeSolicitado entities.Monto) (erro error)

	// impuestos
	GetImpuestoByIdService(id int64) (impuesto entities.Impuesto, erro error)
	CalcularValorCuponService(importe, coeficiente, impuesto float64) (valorCupon int64)
	// CalcularCostoFinancieroIvaService(valorCupon, porcentajeArancel, coeficiente, porcentajeIVA float64) (importeCFMasIva int64)

	// manejo cadena
	RightStr(cadenaStr string, valueStr int) (rightStr string)
	LeftStr(cadenaStr string, valueStr int) (LeftStr string)
	BuildStr(cadenaStr string, valueStr int) (center string)

	GetMatenimietoSistemaService() (estado bool, fecha time.Time, erro error)

	EnviarMailService(params utildtos.RequestDatosMail) (erro error)

	CsvCreate(name string, data [][]string) error

	ValidarCBU(cbu string) (res bool, erro error)

	ValidarCalculoCF(RequestValidarCF utildtos.RequestValidarCF) (responseValidarCF utildtos.ResponseValidarCF)
}

var util *utilService

func NewUtilService(r UtilRepository) UtilService {
	util := utilService{
		repository: r,
		factory:    &crearMensajeFactory{},
	}

	return &util
}
func NewUtilWithService(r UtilRepository, f CrearMensajeFactory) UtilService {
	util := utilService{
		repository: r,
		factory:    f,
	}

	return &util
}

// Resolve devuelve la instancia antes creada
func Resolve() *utilService {
	return util
}

type utilService struct {
	repository UtilRepository
	factory    CrearMensajeFactory
}

func (r *utilService) CreateNotificacionService(notificacion entities.Notificacione) (erro error) {
	return r.repository.CreateNotificacion(notificacion)

}

func (r *utilService) CreateLogService(log entities.Log) (erro error) {
	return r.repository.CreateLog(log)
}

func (r *utilService) LogError(erro string, funcionalidad string) {

	log := entities.Log{
		Tipo:          entities.Error,
		Mensaje:       erro,
		Funcionalidad: funcionalidad,
	}

	err := r.CreateLogService(log)

	if err != nil {
		mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), erro)
		logs.Error(mensaje)
	}
}

func (r *utilService) CrearPeticionesService(peticiones dtos.RequestWebServicePeticion) (erro error) {
	peticionEntity := entities.Webservicespeticione{
		Operacion: peticiones.Operacion,
		Vendor:    peticiones.Vendor,
	}
	err := r.repository.CrearPeticionesRepository(peticionEntity)
	if err != nil {
		mensaje := fmt.Sprintf("error al registrar peticion: %s", err.Error())
		erro = errors.New(mensaje)
		logs.Error(mensaje)
		return
	}
	return nil

}

func (s *utilService) GetConfiguracionService(filtro filtros.ConfiguracionFiltro) (configuracion administraciondtos.ResponseConfiguracion, erro error) {

	response, erro := s.repository.GetConfiguracion(filtro)

	if erro != nil {
		return
	}

	configuracion.FromEntity(response)

	return
}

func (s *utilService) GetConfiguracionesService(filtro filtros.ConfiguracionFiltro) (configuraciones []administraciondtos.ResponseConfiguracion, erro error) {
	response, erro := s.repository.GetConfiguracionesRepository(filtro)
	if erro != nil {
		return
	}
	for _, valueResponse := range response {
		configuracion := administraciondtos.ResponseConfiguracion{}
		configuracion.FromEntity(valueResponse)
		configuraciones = append(configuraciones, configuracion)
	}
	return
}

func (s *utilService) CreateConfiguracionService(config administraciondtos.RequestConfiguracion) (id uint, erro error) {

	erro = config.IsValid(false)

	if erro != nil {
		return
	}

	request := config.ToEntity(false)

	return s.repository.CreateConfiguracion(request)

}

func (s *utilService) FirstOrCreateConfiguracionService(nombre string, descripcion string, valor string) (key string, erro error) {

	if len(strings.TrimSpace(nombre)) < 1 || len(strings.TrimSpace(valor)) < 1 {
		erro = fmt.Errorf("el campo nombre o el campo valor es inválido")
		return
	}

	filtro := filtros.ConfiguracionFiltro{
		Nombre: nombre,
	}

	response, erro := s.GetConfiguracionService(filtro)

	if erro != nil || response.Id == 0 {

		configuracion := administraciondtos.RequestConfiguracion{
			Nombre:      nombre,
			Descripcion: descripcion,
			Valor:       valor,
		}
		_, erro = s.CreateConfiguracionService(configuracion)

		if erro != nil {
			return
		}

		response.Valor = valor
	}

	key = response.Valor

	return

}

func (c *utilService) BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto, importeSolicitado entities.Monto) (erro error) {

	var descuentos entities.Monto

	/*
	   TODO: se debe recorrer cuentacomisiones y registrar las comisiones dependiendo del tipo de pago
	*/
	for _, cc := range *cuentacomisiones {
		if cc.Comision > -1 {
			var totalMin float64
			var totalMax float64
			// analizar
			var importeComisiones float64
			var importeComisionTelco float64
			var importeComisionProveedor float64
			var resultImporte float64
			switch cc.ChannelArancel.Tipocalculo {
			case "FIJO":

				importeComisiones = c.ToFixed(((float64(movimiento.Monto) / 100 * cc.Comision) + cc.ChannelArancel.Importe), 4) //+ cc.ChannelArancel.Importe // 100
				importeComisionTelco = c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 4)                            // 98
				importeComisionProveedor = cc.ChannelArancel.Importe                                                            //2
				resultImporte = c.ToFixed((importeComisionTelco + importeComisionProveedor), 4)                                 // 100
			case "PORCENTAJE":
				// valores porcentual del proveedor
				totalMin = cc.Importeminimo + cc.ChannelArancel.Importeminimo
				totalMax = cc.Importemaximo + cc.ChannelArancel.Importemaximo
				// calcular comision general
				comisionGeneral := cc.Comision + cc.ChannelArancel.Importe
				importeComisiones = c.ToFixed((float64(movimiento.Monto) / 100 * comisionGeneral), 4)
				// calcular comision telco y proveedor
				importeComisionTelco = c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 4)
				importeComisionProveedor = c.ToFixed((float64(movimiento.Monto) / 100 * cc.ChannelArancel.Importe), 4)
				// controlar calculos
				resultImporte = c.ToFixed((float64(movimiento.Monto)/100*cc.Comision)+(float64(movimiento.Monto)/100*cc.ChannelArancel.Importe), 4) // c.ToFixed((importeComisionTelco + importeComisionProveedor), 4)

			}
			if importeComisiones == resultImporte {

				var valorPorcentajeTelco float64
				var valorPorcentajeProveedor float64

				if totalMax > 0 {
					notificacion := entities.Notificacione{
						Tipo:        entities.NotificacionComisionConMaximo,
						Descripcion: fmt.Sprintf("comision con maximo. %s", erro.Error()),
					}
					c.CreateNotificacionService(notificacion)
				}

				min, max := VerificarCalculo(importeComisiones, totalMin, totalMax)

				RequestComision := RequestComision{
					ImporteComisionTelco:     importeComisionTelco,
					ImporteComisionProveedor: importeComisionProveedor,
					MinTelco:                 cc.Importeminimo,
					MaxTelco:                 cc.Importemaximo,
					MinProveedor:             cc.ChannelArancel.Importeminimo,
					MaxProveedor:             cc.ChannelArancel.Importemaximo,
					MinBool:                  min,
					MaxBool:                  max,
				}

				valorComisionTelcoVerif, valorComisionProveedorVerif := VerificarMinimoMaximo(RequestComision)

				// calculo de comision de telco
				valorPorcentajeTelco = cc.Comision
				if importeComisionTelco != valorComisionTelcoVerif {
					if cc.ChannelArancel.Tipocalculo == "FIJO" {
						valorComisionTelcoVerif = valorComisionTelcoVerif - cc.ChannelArancel.Importe
					}
					importeComisionTelco = valorComisionTelcoVerif
					if movimiento.Monto > 0 {
						valorPorcentajeTelco = valorComisionTelcoVerif
					} else {
						valorPorcentajeTelco = -1.00 * valorComisionTelcoVerif
					}
				}

				// calculo de comision de proveedor
				valorPorcentajeProveedor = cc.ChannelArancel.Importe
				if importeComisionProveedor != valorComisionProveedorVerif {
					importeComisionProveedor = valorComisionProveedorVerif
					if movimiento.Monto > 0 {
						valorPorcentajeProveedor = valorComisionProveedorVerif
					} else {
						valorPorcentajeProveedor = -1.00 * valorComisionProveedorVerif
					}
				}

				movimientoComision := entities.Movimientocomisiones{
					CuentacomisionsID:   cc.ID,
					Monto:               entities.Monto(int64(importeComisionTelco * 100)),
					Porcentaje:          valorPorcentajeTelco, //cc.Comision,
					Montoproveedor:      entities.Monto(int64(importeComisionProveedor * 100)),
					Porcentajeproveedor: valorPorcentajeProveedor,
				}
				descuentos += movimientoComision.Monto + movimientoComision.Montoproveedor
				movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

				//Calculo de Iva sobre comision
				if iva != nil && iva.Porcentaje > -1 {
					impuestoIva := c.ToFixed((importeComisionTelco * iva.Porcentaje), 2)
					impuestoIvaProveedor := c.ToFixed((importeComisionProveedor * iva.Porcentaje), 2)
					movimientoImpuesto := entities.Movimientoimpuestos{
						ImpuestosID:    uint64(iva.ID),
						Monto:          entities.Monto(int64(impuestoIva * 100)),
						Montoproveedor: entities.Monto(int64(impuestoIvaProveedor * 100)),
						Porcentaje:     iva.Porcentaje,
					}
					descuentos += movimientoImpuesto.Monto + movimientoImpuesto.Montoproveedor
					movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
				}
			}
		}
	}
	if importeSolicitado < movimiento.Monto {
		movimiento.Monto = importeSolicitado
	}
	movimiento.Monto -= descuentos

	//Si es un movimiento normal le resto si es una devolucion le sumo.
	// if movimiento.Monto > 0 {
	// 	movimiento.Monto -= descuentos
	// } else {
	// 	movimiento.Monto += descuentos
	// }

	return
}

type RequestComision struct {
	ImporteComisionTelco     float64
	ImporteComisionProveedor float64
	MinTelco                 float64
	MaxTelco                 float64
	MinProveedor             float64
	MaxProveedor             float64
	MinBool                  bool
	MaxBool                  bool
}

func (s *utilService) ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func (s *utilService) GetImpuestoByIdService(id int64) (impuesto entities.Impuesto, erro error) {

	impuesto, err := s.repository.GetImpuestoByIdRepository(id)
	if err != nil {
		erro = errors.New(ERROR_CONFIGURACIONES)
		return
	}
	return
}

func (s *utilService) CalcularValorCuponService(importe, coeficiente, impuesto float64) (valorCupon int64) {
	logs.Info("realizando calculo valor cupon")
	logs.Info(impuesto)
	// importeCupon := importe * coeficiente
	// costoFinancieroNeto := importeCupon - importe
	// valorACobrar := costoFinancieroNeto + (costoFinancieroNeto * impuesto)
	// valorCFmasIVA := s.ToFixed(valorACobrar, 2)
	// //inporteMasCfIva := importe + valorCFmasIVA
	// //finalValor := s.ToFixed(inporteMasCfIva, 2)
	// valorCupon = entities.Monto(s.ToFixed(importe+valorCFmasIVA, 2) * 100).Int64()
	valorCupon = entities.Monto(s.ToFixed(importe*coeficiente, 2) * 100).Int64()
	return
}

// func (s *utilService) CalcularCostoFinancieroIvaService(valorCupon, porcentajeArancel, coeficiente, porcentajeIVA float64) (importeCFMasIva int64) {
// 	importeSinArancel := valorCupon - (valorCupon * porcentajeArancel)
// 	importeSinCostoFinanciero := importeSinArancel / coeficiente
// 	importeTemporalCostoFinanciero := importeSinArancel - importeSinCostoFinanciero
// 	importeTemporalIva := importeTemporalCostoFinanciero * porcentajeIVA
// 	importeTemporalCFMasIva := (importeTemporalCostoFinanciero + importeTemporalIva)
// 	importeredondeo := s.ToFixed(importeTemporalCFMasIva, 2)
// 	importeCFMasIva = entities.Monto(importeredondeo * 100).Int64()
// 	return
// }

func (s *utilService) RightStr(cadenaStr string, valueStr int) (rightStr string) {
	totalStr := len(cadenaStr) - valueStr
	rightStr = cadenaStr[totalStr:]
	return
}

func (s *utilService) LeftStr(cadenaStr string, valueStr int) (LeftStr string) {
	LeftStr = cadenaStr[0:valueStr]
	return
}

func (s *utilService) BuildStr(cadenaStr string, valueStr int) (center string) {
	totalStr := len(cadenaStr) - (valueStr * 2)
	for i := 0; i < totalStr; i++ {
		center += "0"
	}
	return
}

func (s *utilService) GetMatenimietoSistemaService() (estado bool, fecha time.Time, erro error) {
	filtro := filtros.ConfiguracionFiltro{
		Nombre: "ESTADO_APLICACION",
	}
	//estadoConfiguracion, err := s.GetConfiguracionService(filtro)
	response, err := s.repository.GetConfiguracion(filtro)
	logs.Info(response.CreatedAt)
	if err != nil {
		estado = true
		erro = fmt.Errorf("el servicio no está disponible")
		return
	}
	if response.Valor != "sin valor" {
		fecha, err = time.Parse("2006-01-02T15:04:00Z", response.Valor)
		if err != nil {
			estado = true

			logs.Error("error al convertir fecha de configuración")
			erro = fmt.Errorf("el servicio no está disponible")
			return
		}
		if !fecha.IsZero() {
			estado = true

			return
		}
	}
	estado = false
	return
}

func (s *utilService) EnviarMailService(params utildtos.RequestDatosMail) (erro error) {
	err := params.IsValid()
	if err != nil {
		erro = errors.New(err.Error())
		return
	}
	typoEmail, err := params.TipoEmail.IsValid()
	if err != nil {
		return err
	}
	var mensajeCompleto string
	if len(params.CamposReemplazar) != 0 {
		mensajeCompleto = construirMensaje(params.Mensaje, params.CamposReemplazar)
	} else {
		mensajeCompleto = params.Mensaje
	}
	// paramsEmail := utildtos.ParamsEmail{
	// 	Email:   params.Email,
	// 	Nombre:  params.Nombre,
	// 	Mensaje: mensajeCompleto,
	// }
	// t, err := template.ParseFiles("../api/views/templateemail/send_mail.html")
	// //t, err := template.ParseFiles("./api/views/templateemail/send_mail.html")
	// if err != nil {
	// 	logs.Error(err.Error())
	// 	erro = errors.New("error al obtener template" + err.Error())
	// 	return
	// }
	// buf := new(bytes.Buffer)
	// erro = t.Execute(buf, paramsEmail)
	// if erro != nil {
	// 	return erro
	// }

	// Datos del remitente para el email.
	smtpUsername := config.SMTP_USERNAME
	smtpPassowrd := config.SMTP_PASSWORD

	// password := config.EMAIL_PASS

	// Crear mensaje
	to := params.Email
	from := config.EMAIL_FROM

	smtpHost := config.SMTPHOST
	smtpPort := config.SMTPPORT
	address := smtpHost + ":" + smtpPort
	params.Mensaje = mensajeCompleto
	params.From = from

	factoryCrearMensaje, err := s.factory.GetCrearMensajeMethod(typoEmail)
	if err != nil {
		return err
	}
	message, err := factoryCrearMensaje.MensajeResultado(params.Asunto, to, params)
	if err != nil {
		return err
	}

	//message := createMessage(params.Asunto, to, from, buf.String())

	//message := createMessage(params.Asunto, to, params.From+" "+from, buf.String())

	// Authentication.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassowrd, smtpHost)

	// Sending email.

	err = smtp.SendMail(address, auth, params.From, to, []byte(message))

	//err = smtp.SendMail(address, auth, params.From+" "+from, to, []byte(message))

	if err != nil {
		logs.Error("error enviar email - " + err.Error())
		return errors.New(ERROR_ENVIAR_EMAIL)
	}
	logs.Error("email enviado con exito.")

	return
}

// //Función para crear el mensage que se enviará al cliente
// func createMessage(subject string, to []string, from, value string) string {

// 	body := value
// 	header := make(map[string]string)
// 	header["From"] = from
// 	for _, valueTo := range to {
// 		header["To"] = valueTo
// 	}
// 	//header["To"] = to[0]
// 	header["Subject"] = subject
// 	header["MIME-Version"] = "1.0"
// 	header["Content-Type"] = "text/html; charset=\"utf-8\""
// 	header["Content-Transfer-Encoding"] = "base64"

// 	message := ""
// 	for k, v := range header {
// 		message += fmt.Sprintf("%s: %s\r\n", k, v)
// 	}
// 	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
// 	return message

// }

func construirMensaje(mensaje string, campos []string) (mensajeCompleto string) {
	for key, value := range campos {
		mensaje = strings.Replace(mensaje, "#"+strconv.Itoa(key), value, 1)
	}
	mensajeCompleto = mensaje
	return
}

func VerificarMinimoMaximo(request RequestComision) (importeComisionTelco, importeComisionProveedor float64) {

	if request.MinBool && !request.MaxBool {
		if request.MinTelco != 0 && request.MinProveedor != 0 {
			importeComisionTelco = request.MinTelco
			importeComisionProveedor = request.MinProveedor
			return
		}

		if request.MinTelco != 0 && request.MinProveedor == 0 {
			importeComisionTelco = request.MinTelco - request.ImporteComisionProveedor
			importeComisionProveedor = request.ImporteComisionProveedor
			return
		}

		// if request.MinTelco == 0 && request.MinProveedor != 0 {
		// 	importeComisionTelco = request.ImporteComisionTelco
		// 	importeComisionProveedor = request.MinProveedor
		// 	return
		// }
	}

	if request.MaxBool && !request.MinBool {
		// if request.MaxTelco != 0 && request.MaxProveedor != 0 {
		// 	importeComisionTelco = request.MaxTelco
		// 	importeComisionProveedor = request.MaxProveedor
		// 	return
		// }

		if request.MaxTelco != 0 && request.MaxProveedor == 0 {
			importeComisionTelco = request.MaxTelco
			importeComisionProveedor = request.ImporteComisionProveedor
			return
		}
	}

	// if !request.MinBool && !request.MaxBool {
	// 	if request.MinTelco != 0 && request.MinProveedor != 0 {
	// 		importeComisionTelco = request.MinTelco
	// 		importeComisionProveedor = request.MinProveedor
	// 		return
	// 	}
	// 	if request.MinTelco != 0 && request.MinProveedor == 0 {
	// 		importeComisionTelco = request.MinTelco - request.ImporteComisionProveedor
	// 		importeComisionProveedor = request.ImporteComisionProveedor
	// 		return
	// 	}
	// 	if request.MaxTelco != 0 && request.MaxProveedor != 0 {
	// 		importeComisionTelco = request.MaxTelco
	// 		importeComisionProveedor = request.MaxProveedor
	// 		return
	// 	}

	// 	if request.MaxTelco != 0 && request.MaxProveedor == 0 {
	// 		importeComisionTelco = request.MaxTelco
	// 		importeComisionProveedor = request.ImporteComisionProveedor
	// 		return
	// 	}
	// 	importeComisionTelco = request.ImporteComisionTelco
	// 	importeComisionProveedor = request.ImporteComisionProveedor
	// 	return
	// }
	importeComisionTelco = request.ImporteComisionTelco
	importeComisionProveedor = request.ImporteComisionProveedor
	return
}

func VerificarCalculo(comision, minimo, maximo float64) (resultMin, resultMax bool) {

	if minimo != 0 && maximo == 0 {
		if math.Abs(comision) < minimo {
			resultMin = true
		}
		return
	}
	// if minimo == 0 && maximo != 0 {
	// 	if math.Abs(comision) > maximo {
	// 		resultMax = true
	// 	}
	// 	return
	// }
	// if minimo != 0 && maximo != 0 {
	// 	if math.Abs(comision) < minimo {
	// 		resultMin = true
	// 	}
	// 	if math.Abs(comision) > maximo {
	// 		resultMax = true
	// 	}

	// 	return
	// }
	return
}

// convertir datos  a excel // utilizado para enviar reportes a clientes
func (s *utilService) CsvCreate(name string, data [][]string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = ';'
	defer w.Flush()

	for _, d := range data {
		err := w.Write(d)
		if err != nil {
			return err
		}
	}

	return nil
}

// validar un CBU de acuerdo a sus diferentes posibilidades
func (s *utilService) ValidarCBU(cbu string) (res bool, erro error) {
	erro = validarLargoCbu(cbu)
	if erro != nil {
		return
	}
	erro = validarCodigoBanco(cbu[0:8])
	if erro != nil {
		return
	}
	erro = validarCuenta(cbu[8:22])
	if erro != nil {
		return
	}
	res = true
	return
}

func (s *utilService) ValidarCalculoCF(RequestValidarCF utildtos.RequestValidarCF) (responseValidarCF utildtos.ResponseValidarCF) {

	// se calcula el valor neto del cupon "Net"
	cuponNeto := RequestValidarCF.Cupon.Float64() * (1 - RequestValidarCF.ArancelMonto)
	// se calcula el valor de la cuota
	montoCuota := cuponNeto / RequestValidarCF.Cuotas

	air := RequestValidarCF.Tna / 100
	air_12 := air / 12

	//calculos intermedios
	intermedio := 1 + (air_12)
	exponente := RequestValidarCF.Cuotas - 1
	intermedio2 := (1 + (air)*((30-RequestValidarCF.Dias)/360))
	intermedio3 := 1 - 1/intermedio
	intermedio4 := 1 - 1/math.Pow(intermedio, exponente)
	// calcular terminso formula CF
	primerTermino := 1 / intermedio2

	segundoTermino := 1 / intermedio * (intermedio4 / intermedio3 / intermedio2)
	if math.IsNaN(segundoTermino) {
		segundoTermino = 0
	}

	// formula calculo CF
	responseValidarCF.CostoFinanciero = s.ToFixed(cuponNeto-montoCuota*(primerTermino+segundoTermino), 2)
	// calculo valor presente
	valorPresenteTemporal := cuponNeto - responseValidarCF.CostoFinanciero
	responseValidarCF.ValorPresente = s.ToFixed(valorPresenteTemporal, 2)
	// calculo ceficiente formula1
	coeficienteUno := s.ToFixed(cuponNeto/responseValidarCF.ValorPresente, 4)
	coeficienteDos := s.ToFixed(1/(1-responseValidarCF.CostoFinanciero/RequestValidarCF.Cupon.Float64()/(1-RequestValidarCF.ArancelMonto)), 4)
	//calculo costo total
	responseValidarCF.CostoTotalPorcentaje = s.ToFixed((responseValidarCF.CostoFinanciero/RequestValidarCF.Cupon.Float64())*100, 2)
	responseValidarCF.ValorCoeficiente = coeficienteUno
	logs.Info(coeficienteDos)

	return
}

// validar la longitud del CBU
func validarLargoCbu(cbu string) error {
	if commons.StringIsEmpity(cbu) {
		return fmt.Errorf("cbu está en blanco")
	}
	if len(cbu) != 22 {
		return fmt.Errorf("longitud de cbu no es válido: %d", len(cbu))
	}
	return nil
}

// validar la parte del Codigo de Banco del CBU que son los 8 primeros digitos
func validarCodigoBanco(codigo string) error {
	if len(codigo) != 8 {
		return fmt.Errorf("el código de banco es incorrecto")
	}
	banco := codigo[0:3] // numero de entidad

	digitoVerificador := codigo[3:4]

	sucursal := codigo[4:7] // numero de sucursal
	// fmt.Println("sucursal: " + sucursal)
	digitoVerificador2 := codigo[7:8]
	// fmt.Println("digito verificador 2: " + digitoVerificador2)

	var suma int
	var x int

	x, _ = strconv.Atoi(banco[0:1])
	suma = x * 7
	x, _ = strconv.Atoi(banco[1:2])
	suma = suma + x
	x, _ = strconv.Atoi(banco[2:3])
	suma = suma + (x * 3)
	x, _ = strconv.Atoi(digitoVerificador)
	suma = suma + (x * 9)
	x, _ = strconv.Atoi(sucursal[0:1])
	suma = suma + (x * 7)
	x, _ = strconv.Atoi(sucursal[1:2])
	suma = suma + x
	x, _ = strconv.Atoi(sucursal[2:3])
	suma = suma + (x * 3)

	var diferencia = (10 - (suma % 10)) % 10
	digito, _ := strconv.Atoi(digitoVerificador2)
	if diferencia != digito {
		return fmt.Errorf("codigo de banco inválido")
	}
	return nil
}

// validar la parte del Cuenta del CBU que son los 14 ultimos digitos
func validarCuenta(cuenta string) error {
	if len(cuenta) != 14 {
		return fmt.Errorf("logitud de cuenta inválido: %d", len(cuenta))
	}
	digitoVerificador, _ := strconv.Atoi(cuenta[13:14])

	var suma int
	var x int

	x, _ = strconv.Atoi(cuenta[0:1])
	suma = x * 3
	x, _ = strconv.Atoi(cuenta[1:2])
	suma = suma + (x * 9)
	x, _ = strconv.Atoi(cuenta[2:3])
	suma = suma + (x * 7)
	x, _ = strconv.Atoi(cuenta[3:4])
	suma = suma + x
	x, _ = strconv.Atoi(cuenta[4:5])
	suma = suma + (x * 3)
	x, _ = strconv.Atoi(cuenta[5:6])
	suma = suma + (x * 9)
	x, _ = strconv.Atoi(cuenta[6:7])
	suma = suma + (x * 7)
	x, _ = strconv.Atoi(cuenta[7:8])
	suma = suma + (x * 1)
	x, _ = strconv.Atoi(cuenta[8:9])
	suma = suma + (x * 3)
	x, _ = strconv.Atoi(cuenta[9:10])
	suma = suma + (x * 9)
	x, _ = strconv.Atoi(cuenta[10:11])
	suma = suma + (x * 7)
	x, _ = strconv.Atoi(cuenta[11:12])
	suma = suma + (x * 1)
	x, _ = strconv.Atoi(cuenta[12:13])
	suma = suma + (x * 3)

	var diferencia = (10 - (suma % 10)) % 10

	if diferencia != digitoVerificador {
		return fmt.Errorf("error en cuenta bancaria")
	}

	return nil
}

// mes de junio
// func (c *utilService) BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto) (erro error) {

// 	var descuentos entities.Monto
// 	/*
// 	   TODO: se debe recorrer cuentacomisiones y registrar las comisiones dependiendo del tipo de pago
// 	*/
// 	for _, cc := range *cuentacomisiones {
// 		if cc.Comision > -1 {

// 			comision := c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 2)

// 			movimientoComision := entities.Movimientocomisiones{
// 				CuentacomisionsID: cc.ID,
// 				Monto:             entities.Monto(int64(comision * 100)),
// 				Porcentaje:        cc.Comision,
// 			}
// 			descuentos += movimientoComision.Monto
// 			movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

// 			//Calculo de Iva sobre comision
// 			if iva != nil && iva.Porcentaje > -1 {
// 				impuestoIva := c.ToFixed((comision * iva.Porcentaje), 2)
// 				movimientoImpuesto := entities.Movimientoimpuestos{
// 					ImpuestosID: uint64(iva.ID),
// 					Monto:       entities.Monto(int64(impuestoIva * 100)),
// 					Porcentaje:  iva.Porcentaje,
// 				}
// 				descuentos += movimientoImpuesto.Monto
// 				movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
// 			}
// 		}
// 	}
// 	//Si es un movimiento normal le resto si es una devolucion le sumo.
// 	if movimiento.Monto > 0 {
// 		movimiento.Monto -= descuentos
// 	} else {
// 		movimiento.Monto += descuentos
// 	}

// 	return
// }

/*
	// if comision < cc.Importeminimo && cc.Importeminimo != 0 {
		// 	importeComision = cc.Importeminimo
		// } else if comision > cc.Importemaximo && cc.Importemaximo != 0 {
		// 	importeComision = cc.Importemaximo
		// } else {
		// 	importeComision = comision
		// }
*/

//06-10-2022
// func (c *utilService) BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto) (erro error) {

// 	var descuentos entities.Monto
// 	/*
// 	   TODO: se debe recorrer cuentacomisiones y registrar las comisiones dependiendo del tipo de pago
// 	*/
// 	for _, cc := range *cuentacomisiones {
// 		if cc.Comision > -1 {
// 			// calcular comision general

// 			// calcular comision telco y proveedor

// 			// controlar calculos

// 			// verificar si el importe de la comision no supera el minimo,
// 			// supera el maximo o se encuetra entre el minimo y el maximo
// 			// de telco y proveedor

// 			var valorPorcentaje float64
// 			comision := c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 2)
// 			importeComision := VerificarMinimoMaximo(comision, cc.Importeminimo, cc.Importemaximo)

// 			valorPorcentaje = cc.Comision
// 			if comision != importeComision {
// 				if movimiento.Monto > 0 {
// 					valorPorcentaje = importeComision
// 				} else {
// 					valorPorcentaje = -1.00 * importeComision
// 				}
// 			}

// 			movimientoComision := entities.Movimientocomisiones{
// 				CuentacomisionsID: cc.ID,
// 				Monto:             entities.Monto(int64(importeComision * 100)),
// 				Porcentaje:        valorPorcentaje, //cc.Comision,
// 			}
// 			descuentos += movimientoComision.Monto
// 			movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

// 			//Calculo de Iva sobre comision
// 			if iva != nil && iva.Porcentaje > -1 {
// 				impuestoIva := c.ToFixed((importeComision * iva.Porcentaje), 2)
// 				movimientoImpuesto := entities.Movimientoimpuestos{
// 					ImpuestosID: uint64(iva.ID),
// 					Monto:       entities.Monto(int64(impuestoIva * 100)),
// 					Porcentaje:  iva.Porcentaje,
// 				}
// 				descuentos += movimientoImpuesto.Monto
// 				movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
// 			}
// 		}
// 	}
// 	movimiento.Monto -= descuentos

// 	//Si es un movimiento normal le resto si es una devolucion le sumo.
// 	// if movimiento.Monto > 0 {
// 	// 	movimiento.Monto -= descuentos
// 	// } else {
// 	// 	movimiento.Monto += descuentos
// 	// }

// 	return
// }

// 14-10-2022
// // valores porcentual del proveedor
// 			// calcular comision general
// 			comisionGeneral := cc.Comision + cc.ChannelArancel.Importe
// 			importeComisiones := c.ToFixed((float64(movimiento.Monto) / 100 * comisionGeneral), 4)
// 			// calcular comision telco y proveedor
// 			importeComisionTelco := c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 4)
// 			importeComisionProveedor := c.ToFixed((float64(movimiento.Monto) / 100 * cc.ChannelArancel.Importe), 4)
// 			// controlar calculos
// 			resultImporte := c.ToFixed((importeComisionTelco + importeComisionProveedor), 4)
// 			if importeComisiones == resultImporte {
// 				//var valorPorcentaje float64
// 				var valorPorcentajeTelco float64
// 				var valorPorcentajeProveedor float64
// 				// comision := c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 2)
// 				// importeComision := VerificarMinimoMaximo(comisionGeneral, cc.Importeminimo, cc.Importemaximo)

// 				// valorPorcentaje = cc.Comision
// 				// if importeComisiones != importeComision {
// 				// 	if movimiento.Monto > 0 {
// 				// 		valorPorcentaje = importeComision
// 				// 	} else {
// 				// 		valorPorcentaje = -1.00 * importeComision
// 				// 	}
// 				// }
// 				//verifico valor minimo y maximo de telco
// 				importeComisionTelcoVerif := VerificarMinimoMaximo(cc.Comision, cc.Importeminimo, cc.Importemaximo)
// 				valorPorcentajeTelco = cc.Comision
// 				if importeComisionTelco != importeComisionTelcoVerif {
// 					if movimiento.Monto > 0 {
// 						valorPorcentajeTelco = importeComisionTelcoVerif
// 					} else {
// 						valorPorcentajeTelco = -1.00 * importeComisionTelcoVerif
// 					}
// 				}
// 				//verifico valor minimo y maximo del proveedor
// 				// cc.ChannelArancel.Importe es igual a la comision del proveedor
// 				importeComisionProveedorVerif := VerificarMinimoMaximo(cc.ChannelArancel.Importe, cc.ChannelArancel.Importeminimo, cc.ChannelArancel.Importemaximo)
// 				valorPorcentajeProveedor = cc.ChannelArancel.Importe
// 				if importeComisionProveedor != importeComisionProveedorVerif {
// 					if movimiento.Monto > 0 {
// 						valorPorcentajeProveedor = importeComisionProveedorVerif
// 					} else {
// 						valorPorcentajeProveedor = -1.00 * importeComisionProveedorVerif
// 					}
// 				}

// 				movimientoComision := entities.Movimientocomisiones{
// 					CuentacomisionsID:   cc.ID,
// 					Monto:               entities.Monto(int64(importeComisionTelco * 100)),
// 					Porcentaje:          valorPorcentajeTelco, //cc.Comision,
// 					Montoproveedor:      entities.Monto(int64(importeComisionProveedor * 100)),
// 					Porcentajeproveedor: valorPorcentajeProveedor,
// 				}
// 				descuentos += movimientoComision.Monto + movimientoComision.Montoproveedor
// 				movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

// 				//Calculo de Iva sobre comision
// 				if iva != nil && iva.Porcentaje > -1 {
// 					impuestoIva := c.ToFixed((importeComisionTelco * iva.Porcentaje), 2)
// 					impuestoIvaProveedor := c.ToFixed((importeComisionProveedor * iva.Porcentaje), 2)
// 					movimientoImpuesto := entities.Movimientoimpuestos{
// 						ImpuestosID:    uint64(iva.ID),
// 						Monto:          entities.Monto(int64(impuestoIva * 100)),
// 						Montoproveedor: entities.Monto(int64(impuestoIvaProveedor * 100)),
// 						Porcentaje:     iva.Porcentaje,
// 					}
// 					descuentos += movimientoImpuesto.Monto + movimientoImpuesto.Montoproveedor
// 					movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
// 				}
// 			}

// 		}

// 27/12/22
// func (c *utilService) BuildComisiones(movimiento *entities.Movimiento, cuentacomisiones *[]entities.Cuentacomision, iva *entities.Impuesto, importeSolicitado entities.Monto) (erro error) {

// 	var descuentos entities.Monto

// 	/*
// 	   TODO: se debe recorrer cuentacomisiones y registrar las comisiones dependiendo del tipo de pago
// 	*/
// 	for _, cc := range *cuentacomisiones {
// 		if cc.Comision > -1 {
// 			// analizar
// 			var importeComisiones float64
// 			var importeComisionTelco float64
// 			var importeComisionProveedor float64
// 			var resultImporte float64
// 			switch cc.ChannelArancel.Tipocalculo {
// 			case "FIJO":

// 				importeComisiones = c.ToFixed(((float64(movimiento.Monto) / 100 * cc.Comision) + cc.ChannelArancel.Importe), 4) //+ cc.ChannelArancel.Importe // 100
// 				importeComisionTelco = c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 4)                            // 98
// 				importeComisionProveedor = cc.ChannelArancel.Importe                                                            //2
// 				resultImporte = c.ToFixed((importeComisionTelco + importeComisionProveedor), 4)                                 // 100
// 			case "PORCENTAJE":
// 				// valores porcentual del proveedor
// 				// calcular comision general

// 				comisionGeneral := cc.Comision + cc.ChannelArancel.Importe
// 				importeComisiones = c.ToFixed((float64(movimiento.Monto) / 100 * comisionGeneral), 4)
// 				// calcular comision telco y proveedor
// 				importeComisionTelco = c.ToFixed((float64(movimiento.Monto) / 100 * cc.Comision), 4)
// 				importeComisionProveedor = c.ToFixed((float64(movimiento.Monto) / 100 * cc.ChannelArancel.Importe), 4)
// 				// controlar calculos
// 				resultImporte = c.ToFixed((float64(movimiento.Monto)/100*cc.Comision)+(float64(movimiento.Monto)/100*cc.ChannelArancel.Importe), 4) // c.ToFixed((importeComisionTelco + importeComisionProveedor), 4)

// 			}
// 			if importeComisiones == resultImporte {

// 				var valorPorcentajeTelco float64
// 				var valorPorcentajeProveedor float64

// 				valorComisionTelcoVerif := VerificarMinimoMaximo(importeComisionTelco, cc.Importeminimo, cc.Importemaximo, cc.ChannelArancel.Importeminimo, cc.ChannelArancel.Importemaximo)
// 				valorPorcentajeTelco = cc.Comision
// 				if importeComisionTelco != valorComisionTelcoVerif {
// 					if cc.ChannelArancel.Tipocalculo == "FIJO" {
// 						valorComisionTelcoVerif = valorComisionTelcoVerif - cc.ChannelArancel.Importe
// 					}
// 					importeComisionTelco = valorComisionTelcoVerif
// 					if movimiento.Monto > 0 {
// 						valorPorcentajeTelco = valorComisionTelcoVerif
// 					} else {
// 						valorPorcentajeTelco = -1.00 * valorComisionTelcoVerif
// 					}
// 				}
// 				//verifico valor minimo y maximo del proveedor
// 				// cc.ChannelArancel.Importe es igual a la comision del proveedor
// 				// importeComisionProveedorVerif puede contener un importe o un coeficiente

// 				valorComisionProveedorVerif := VerificarMinimoMaximo(importeComisionProveedor, cc.ChannelArancel.Importeminimo, cc.ChannelArancel.Importemaximo, 0, 0)
// 				valorPorcentajeProveedor = cc.ChannelArancel.Importe
// 				if importeComisionProveedor != valorComisionProveedorVerif {
// 					importeComisionProveedor = valorComisionProveedorVerif
// 					if movimiento.Monto > 0 {
// 						valorPorcentajeProveedor = valorComisionProveedorVerif
// 					} else {
// 						valorPorcentajeProveedor = -1.00 * valorComisionProveedorVerif
// 					}
// 				}

// 				movimientoComision := entities.Movimientocomisiones{
// 					CuentacomisionsID:   cc.ID,
// 					Monto:               entities.Monto(int64(importeComisionTelco * 100)),
// 					Porcentaje:          valorPorcentajeTelco, //cc.Comision,
// 					Montoproveedor:      entities.Monto(int64(importeComisionProveedor * 100)),
// 					Porcentajeproveedor: valorPorcentajeProveedor,
// 				}
// 				descuentos += movimientoComision.Monto + movimientoComision.Montoproveedor
// 				movimiento.Movimientocomisions = append(movimiento.Movimientocomisions, movimientoComision)

// 				//Calculo de Iva sobre comision
// 				if iva != nil && iva.Porcentaje > -1 {
// 					impuestoIva := c.ToFixed((importeComisionTelco * iva.Porcentaje), 2)
// 					impuestoIvaProveedor := c.ToFixed((importeComisionProveedor * iva.Porcentaje), 2)
// 					movimientoImpuesto := entities.Movimientoimpuestos{
// 						ImpuestosID:    uint64(iva.ID),
// 						Monto:          entities.Monto(int64(impuestoIva * 100)),
// 						Montoproveedor: entities.Monto(int64(impuestoIvaProveedor * 100)),
// 						Porcentaje:     iva.Porcentaje,
// 					}
// 					descuentos += movimientoImpuesto.Monto + movimientoImpuesto.Montoproveedor
// 					movimiento.Movimientoimpuestos = append(movimiento.Movimientoimpuestos, movimientoImpuesto)
// 				}
// 			}
// 		}
// 	}
// 	if importeSolicitado < movimiento.Monto {
// 		movimiento.Monto = importeSolicitado
// 	}
// 	movimiento.Monto -= descuentos

// 	//Si es un movimiento normal le resto si es una devolucion le sumo.
// 	// if movimiento.Monto > 0 {
// 	// 	movimiento.Monto -= descuentos
// 	// } else {
// 	// 	movimiento.Monto += descuentos
// 	// }

// 	return
// }

// func VerificarMinimoMaximo(comision, minimo, maximo float64) (importeComision float64) {

// 	if minimo != 0 && maximo == 0 {
// 		if math.Abs(comision) < minimo {
// 			importeComision = minimo
// 		} else {
// 			importeComision = comision
// 		}
// 		return
// 	}
// 	if minimo == 0 && maximo != 0 {
// 		if math.Abs(comision) > maximo {
// 			importeComision = maximo
// 		} else {
// 			importeComision = comision
// 		}
// 		return
// 	}
// 	if minimo != 0 && maximo != 0 {
// 		if math.Abs(comision) < minimo {
// 			importeComision = minimo
// 		}
// 		if math.Abs(comision) > maximo {
// 			importeComision = maximo
// 		}
// 		if math.Abs(comision) > minimo && math.Abs(comision) < maximo {
// 			importeComision = comision
// 		}
// 		return
// 	}
// 	importeComision = comision
// 	return
// }

// func VerificarMinimoMaximo(request RequestComision) (importeComisionTelco, importeComisionProveedor float64) {

// 	if request.MinBool && !request.MaxBool {
// 		if request.MinTelco != 0 && request.MinProveedor != 0 {
// 			importeComisionTelco = request.MinTelco
// 			importeComisionProveedor = request.MinProveedor
// 			return
// 		}

// 		if request.MinTelco != 0 && request.MinProveedor == 0 {
// 			importeComisionTelco = request.MinTelco - request.ImporteComisionProveedor
// 			importeComisionProveedor = request.ImporteComisionProveedor
// 			return
// 		}

// 		// if request.MinTelco == 0 && request.MinProveedor != 0 {
// 		// 	importeComisionTelco = request.ImporteComisionTelco
// 		// 	importeComisionProveedor = request.MinProveedor
// 		// 	return
// 		// }
// 	}

// 	// if request.MaxBool && !request.MinBool {
// 	// 	if request.MaxTelco != 0 && request.MaxProveedor != 0 {
// 	// 		importeComisionTelco = request.MaxTelco
// 	// 		importeComisionProveedor = request.MaxProveedor
// 	// 		return
// 	// 	}

// 	// 	if request.MaxTelco != 0 && request.MaxProveedor == 0 {
// 	// 		importeComisionTelco = request.MaxTelco
// 	// 		importeComisionProveedor = request.ImporteComisionProveedor
// 	// 		return
// 	// 	}
// 	// }

// 	// if !request.MinBool && !request.MaxBool {
// 	// 	if request.MinTelco != 0 && request.MinProveedor != 0 {
// 	// 		importeComisionTelco = request.MinTelco
// 	// 		importeComisionProveedor = request.MinProveedor
// 	// 		return
// 	// 	}
// 	// 	if request.MinTelco != 0 && request.MinProveedor == 0 {
// 	// 		importeComisionTelco = request.MinTelco - request.ImporteComisionProveedor
// 	// 		importeComisionProveedor = request.ImporteComisionProveedor
// 	// 		return
// 	// 	}
// 	// 	if request.MaxTelco != 0 && request.MaxProveedor != 0 {
// 	// 		importeComisionTelco = request.MaxTelco
// 	// 		importeComisionProveedor = request.MaxProveedor
// 	// 		return
// 	// 	}

// 	// 	if request.MaxTelco != 0 && request.MaxProveedor == 0 {
// 	// 		importeComisionTelco = request.MaxTelco
// 	// 		importeComisionProveedor = request.ImporteComisionProveedor
// 	// 		return
// 	// 	}
// 	// 	importeComisionTelco = request.ImporteComisionTelco
// 	// 	importeComisionProveedor = request.ImporteComisionProveedor
// 	// 	return
// 	// }
// 	importeComisionTelco = request.ImporteComisionTelco
// 	importeComisionProveedor = request.ImporteComisionProveedor
// 	return
// }
