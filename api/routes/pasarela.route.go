package routes

import (
	"fmt"
	"strconv"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
	pasalera_service "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/pasarela/services"
	"github.com/gofiber/fiber/v2"
)

func PasarelaRoutes(app fiber.Router, service pasalera_service.Service, middlewares middlewares.MiddlewareManager) {
	/* 	obtiene la tabla de intereses por cuotas, utilizado para ver costos desde el checkout */
	app.Get("/plan-cuotas", getPlanCuotas(service))
	app.Get("/healthcheck",    middlewares.ValidarApikey(),   getHealthcheck(service))
	
}

func getPlanCuotas(service pasalera_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logs.Info("accedio al contolador planes de cuotas")
		importe := c.Query("importe")
		if len(importe) <= 0 {
			return fiber.NewError(400, "el importe enviado no es válido")
		}
		logs.Info(fmt.Sprintf("parametro enviado %v", importe))
		monto, err := strconv.Atoi(importe)

		if err != nil {
			return fmt.Errorf("no se puede procesar el id de pago")
		}

		response, err := service.GetPlanCuotasService(monto)

		if err != nil {
			logs.Info(fmt.Sprintf("error en response"))
			return fiber.NewError(400, "Error: "+err.Error())
		}
		return c.JSON(response)
	}
}
func getHealthcheck(service pasalera_service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		err := service.HealthcheckService()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(&fiber.Map{
				"status":        false,
				"statusMessage": "El sistema no se encuentra disponible",
	
			})	
		}
		return c.JSON(&fiber.Map{
			"status":        true,
			"statusMessage": "El sistema se encuentra disponible",

		})

	}
}