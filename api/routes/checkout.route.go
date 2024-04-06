package routes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/api/middlewares/middlewareapikey"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
	checkout_service "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/checkout/services"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/multipagosdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/rapipago"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	"github.com/gofiber/fiber/v2"
)

func CheckoutRoutes(app fiber.Router, middlewaresApiKey middlewareapikey.MiddlewareApiKey, service checkout_service.Service) {
	app.Post("/", createPaymentRequest(service))
	app.Get("/:uuid", getPaymentRequest(service))
	app.Post("/prisma", checkPrisma(service))
	app.Post("/pagar", postPagar(service))
	app.Get("/bill/:barcode", getBill(service)) // pdf comprobante del pago
	app.Get("/tarjetas/all", getTarjetas(service))
	app.Get("/pago/estado/:barcode", getVerificarPagoEstado(service))
	app.Get("/consulta/estado", middlewaresApiKey.ValidarApikeyCuenta(), getEstadoApp(service))

	// app.Get("/control/:hash", controlHash(service))

	app.Post("/multipago/consulta", postVerificarMultipagoEstado(service))
	app.Post("/multipago/pago", PagoMultipago(service))
	app.Post("/multipago/control", postControlMultipagoPagos(service))

	app.Post("/rapipago/pago", postRapipagoPago(service))
	app.Post("/rapipago/consulta", postRapipagoQuery(service))
	app.Post("/rapipago/confirmacion", postRapipagoConfirmacion(service))
}

// crear una solicitud de pago
func createPaymentRequest(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		api := c.Get("apiKey")
		if len(api) <= 0 {
			return fiber.NewError(400, "Debe enviar una api key válida")
		}

		var request dtos.PagoRequest

		err := c.BodyParser(&request)
		if err != nil {
			apiShort := api[0:4]
			mensajeError := fmt.Sprintf("Error en el formato de los parametros enviados al crear la solicitud de pago: APIKEY:%s... || ERROR:%s. El body sin parsear es:%s", apiShort, err.Error(), c.Body())
			logs.Error(mensajeError)
			return fiber.NewError(400, "Parámetros incorrectos: "+err.Error())
		}

		estado, fecha, err := service.GetMatenimietoSistema()
		if err != nil {
			return c.Status(503).JSON(&fiber.Map{
				"status":  estado,
				"message": "el sistema estara en mantenimiento hasta " + time.Now().Format(time.RFC822Z),
			})
			// return fiber.NewError(503, "Parámetros incorrectos: "+err.Error())
		}
		if estado {
			fechaString := fmt.Sprintf("%v-%v-%v %v:%v:%v", fecha.Day(), fecha.Month(), fecha.Year(), fecha.Hour(), fecha.Minute(), fecha.Second())

			return c.Status(503).JSON(&fiber.Map{
				"status":  estado,
				"message": "el sistema estara en mantenimiento hasta " + fechaString,
			})
		}

		ctx := getCheckoutContext(c)

		// data es de tipo dtos.PagoResponse
		data, err := service.NewPago(ctx, &request, api)
		if err != nil {
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"data":    data,
			"message": "Solicitud de pago generada",
		})
	}
}

// obtener un pago generado mediante una solicitud de pago
func getPaymentRequest(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		code := c.Params("uuid")

		if len(code) <= 0 {
			return fiber.NewError(400, "debe enviar un código de pago válido")
		}

		// data es tipo dtos.CheckoutResponse
		data, err := service.GetPaid(code)

		if err != nil {
			return c.JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		data.BaseUrl = config.APP_CHECKOUT_URL //config.APP_HOST

		return c.JSON(&fiber.Map{
			"status":  true,
			"data":    data,
			"message": "Solicitud de Pago pendiente recuperada con exito",
		})
	}
}

// checkear si el servicio de prisma esta disponible
func checkPrisma(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		err := service.CheckPrisma()

		if err != nil {
			return c.JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(), // "el servicio de prisma no está disponible"
			})
		}

		return c.Status(200).JSON(&fiber.Map{
			"status":  true,
			"message": "el servicio de prisma está funcionando correctamente.",
		})
	}
}

// efectuar un pago
func postPagar(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		logs.Info(fmt.Sprint("postPagar: ", c.IP()))
		var request dtos.ResultadoRequest
		err := c.BodyParser(&request)
		if err != nil {
			return fiber.NewError(400, "Parámetros incorrectos: "+err.Error())
		}
		logs.Info(request.Channel)
		logs.Info(request.Uuid)
		ctx := getCheckoutContext(c)
		res, err := service.GetPagoResultado(ctx, &request)
		if err != nil {
			return fiber.NewError(400, err.Error())
		}
		return c.Status(200).JSON(&fiber.Map{
			"status":  res.Exito,
			"data":    res,
			"message": res.Estado,
		})
	}
}

// Retorna un stream para mostrar un pdf con los datos del pago
func getBill(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		code := c.Params("barcode")

		file, err := service.GetBilling(code)
		if err != nil {
			return c.JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		c.Set("Content-Disposition", "filename=recibo.pdf")
		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Length", fmt.Sprint(file.Len()))

		return c.SendStream(file) // c.SendStream(file)
	}
}

// obtiene datos de la tabla mediopagos
func getTarjetas(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		tar, err := service.GetTarjetas()
		if err != nil {
			return c.JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
		return c.JSON(&fiber.Map{
			"status":  true,
			"data":    tar,
			"message": "tarjetas enviadas",
		})
	}
}

func getCheckoutContext(c *fiber.Ctx) context.Context {
	userctx := entities.Auditoria{
		IP: c.IP(),
	}
	ctx := context.WithValue(c.Context(), entities.AuditUserKey{}, userctx)
	return ctx
}

func getVerificarPagoEstado(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		code := c.Params("barcode")

		response, err := service.GetPagoStatus(code) // servicio que verifica el estado del pago
		if err != nil {
			return c.Status(200).JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.JSON(&fiber.Map{
			"status":  true,
			"data":    response,
			"message": "Control hash",
		})
	}
}

func controlHash(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		hash := c.Params("hash")

		response, err := service.ControlTarjetaHash(hash)
		if err != nil {
			return c.Status(200).JSON(&fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.JSON(&fiber.Map{
			"status":  true,
			"data":    response,
			"message": "Control tarjeta bloqueada",
		})
	}
}

func postVerificarMultipagoEstado(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		apikey := c.Get("ApiKey")
		if len(apikey) <= 0 {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Apikey invalida",
			}
			return c.Status(200).JSON(errorFormato)
		}

		controlApikey, err := service.ControlAdquirienteApikey(apikey, "multipagos")
		if err != nil {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Apikey invalida",
			}
			return c.Status(200).JSON(errorFormato)
		}

		if controlApikey == false {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Apikey invalida",
			}
			return c.Status(200).JSON(errorFormato)
		}

		var request multipagosdtos.RequestConsultaMultipago
		err = c.BodyParser(&request)
		if err != nil {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Valores con formatos incorrectos",
			}
			logs.Error(err.Error())
			return c.Status(200).JSON(errorFormato)
		}

		response, err := service.GetMultiPagoStatus(request) // servicio que verifica el estado del pago
		if err != nil {

			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "10",
				Msg:             "Error en el servicio",
			}
			logs.Error(err.Error())
			return c.Status(200).JSON(errorFormato)
		}

		// logs.Info(response)
		return c.Status(200).JSON(response)
	}
}

func postControlMultipagoPagos(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		apikey := c.Get("ApiKey")
		if len(apikey) <= 0 {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Apikey invalida",
			}
			return c.Status(200).JSON(errorFormato)
		}

		controlApikey, err := service.ControlAdquirienteApikey(apikey, "multipagos")
		if err != nil {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Apikey invalida",
			}
			return c.Status(200).JSON(errorFormato)
		}

		if controlApikey == false {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Apikey invalida",
			}
			return c.Status(200).JSON(errorFormato)
		}

		var request multipagosdtos.RequestControlMultipago
		err = c.BodyParser(&request)
		if err != nil {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Valores con formatos incorrectos",
			}
			logs.Error(err.Error())
			return c.Status(200).JSON(errorFormato)
		}

		response, err := service.GetMultiPagoControl(request) // servicio que verifica el estado del pago
		if err != nil {

			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "10",
				Msg:             "Error en el servicio",
			}
			logs.Error(err.Error())
			return c.Status(200).JSON(errorFormato)
		}

		// logs.Info(response)
		return c.Status(200).JSON(response)
	}
}

func PagoMultipago(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")
		apikey := c.Get("ApiKey")
		if len(apikey) <= 0 {
			// return fmt.Errorf("acceso no autorizado, debe enviar una apikey para autenticación")
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Apikey invalida",
			}
			return c.Status(200).JSON(errorFormato)
		}

		controlApikey, err := service.ControlAdquirienteApikey(apikey, "multipagos")
		if err != nil {
			// return fiber.NewError(400, err.Error())
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "1",
				Msg:             "Error en la DB.",
			}
			logs.Error(err.Error())
			return c.Status(200).JSON(errorFormato)
		}

		if controlApikey == false {
			// return fiber.NewError(400, "apikey invalida")
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "5",
				Msg:             "Apikey invalida",
			}
			logs.Error(err.Error())
			return c.Status(200).JSON(errorFormato)

		}

		var request multipagosdtos.RequestPagoMultipago

		err = c.BodyParser(&request)
		if err != nil {
			// return fiber.NewError(400, "Parámetros incorrectos: "+err.Error())
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Valores con formatos incorrectos",
			}
			logs.Error(err.Error())
			return c.Status(200).JSON(errorFormato)

		}

		ctx := getCheckoutContext(c)
		response, err := service.PostMultiPago(ctx, request) // servicio que registra pago por Multipago
		if err != nil {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "1",
				Msg:             "Error en la DB.",
			}
			logs.Error(err.Error())
			return c.Status(200).JSON(errorFormato)
		}

		// logs.Info(response)
		return c.Status(200).JSON(response)
	}
}

func postRapipagoPago(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request rapipago.RequestRapipagoConsulta

		err := c.BodyParser(&request)
		if err != nil {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Valores con formatos incorrectos",
			}
			response := request.ParseToResponse(errorFormato.CodigoRespuesta, errorFormato.Msg)
			logs.Error(err.Error())
			return c.Status(200).JSON(response)
		}

		// Control por apikey
		errorFormato, err := ControlApikey(c, service, "postRapipagoPago")
		if errorFormato.CodigoRespuesta != "0" || err != nil {
			response := request.ParseToResponse(errorFormato.CodigoRespuesta, errorFormato.Msg)

			return c.Status(200).JSON(response)
		}

		ctx := getCheckoutContext(c)
		response, err := service.PostRapipagoPago(ctx, request) // servicio que actualiza el estado del pago rapipago
		if err != nil {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "10",
				Msg:             "Error en el servicio",
			}
			logs.Error(err.Error())
			response := request.ParseToResponse(errorFormato.CodigoRespuesta, errorFormato.Msg)
			return c.Status(200).JSON(response)
		}
		if response.CodigoRespuesta != "0" {
			response = request.ParseToResponse(response.CodigoRespuesta, response.Msg)
		}
		return c.Status(200).JSON(response)
	}
}

func postRapipagoQuery(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request rapipago.RequestRapipagoConsulta

		err := c.BodyParser(&request)
		if err != nil {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "9",
				Msg:             "Valores con formatos incorrectos",
			}
			logs.Error(err.Error())
			return c.Status(200).JSON(errorFormato)
		}

		// Control por apikey
		errorFormato, err := ControlApikey(c, service, "postRapipagoQuery")
		if errorFormato.CodigoRespuesta != "0" || err != nil {
			return c.Status(200).JSON(errorFormato)
		}

		response, err := service.GetRapipagoQuery(request) // servicio que verifica el estado del pago rapipago
		if err != nil {
			errorFormato := rapipago.ResponseRapipagoImputacion{
				CodigoRespuesta: "10",
				Msg:             "Error en el servicio",
			}
			logs.Error(err.Error())
			return c.Status(200).JSON(errorFormato)
		}

		return c.Status(200).JSON(response)
	}
}

func postRapipagoConfirmacion(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		var request rapipago.RequestRapipagoConfirmacion

		err := c.BodyParser(&request)
		if err != nil {
			response := request.ParseToResponse("9", "Parámetros incorrectos o faltantes")
			logs.Error(err.Error())
			return c.Status(200).JSON(response)
		}

		// Control por apikey
		errorFormato, err := ControlApikey(c, service, "postRapipagoQuery")
		if errorFormato.CodigoRespuesta != "0" || err != nil {
			return c.Status(200).JSON(errorFormato)
		}

		response, err := service.RapipagoConfirmacionService(request) // servicio que verifica el estado del pago rapipago
		if err != nil {
			logs.Error(err.Error())
			return c.Status(200).JSON(response)
		}

		return c.Status(200).JSON(response)
	}
}

func ControlApikey(c *fiber.Ctx, s checkout_service.Service, funcionalidad string) (errorFormato rapipago.ResponseRapipagoImputacion, err error) {

	apikey := c.Get("ApiKey")
	if len(apikey) == 0 || len(apikey) < 36 {
		errorFormato = rapipago.ResponseRapipagoImputacion{
			CodigoRespuesta: "9",
			Msg:             "Apikey invalida",
		}
		err = errors.New("apikey invalida en " + funcionalidad)
		logs.Error(err.Error())
		return
	}

	controlApikey, err := s.ControlAdquirienteApikey(apikey, "rapipago")
	if err != nil {
		errorFormato = rapipago.ResponseRapipagoImputacion{
			CodigoRespuesta: "1",
			Msg:             "Error en la DB.",
		}
		logs.Error(err.Error())
		return
	}

	if !controlApikey {
		errorFormato = rapipago.ResponseRapipagoImputacion{
			CodigoRespuesta: "5",
			Msg:             "Operacion invalida. Apikey no valida",
		}
		err = errors.New("Operacion invalida. Apikey no valida en " + funcionalidad)
		logs.Error(err.Error())
		return
	}
	// si encuentra la apikey responde con el mensaje de exito
	errorFormato = rapipago.ResponseRapipagoImputacion{
		CodigoRespuesta: "0",
		Msg:             "Trx ok",
	}
	return
}

func getEstadoApp(service checkout_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		erro := service.GetEstadoAppService()

		if erro != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&fiber.Map{
				"code":          400,
				"status":        false,
				"statusMessage": "El sistema no se encuentra disponible",
			})
		}

		return c.JSON(&fiber.Map{
			"code":          200,
			"status":        true,
			"statusMessage": "El sistema se encuentra disponible",
		})
	}
}
