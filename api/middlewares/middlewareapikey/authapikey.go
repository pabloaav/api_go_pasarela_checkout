package middlewareapikey

import (
	"errors"
	"fmt"

	checkout_service "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/checkout/services"

	"github.com/gofiber/fiber/v2"
)

type MiddlewareApiKey struct {
	Service checkout_service.Service
}

func (m *MiddlewareApiKey) ValidarApikeyCuenta() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		apikey := c.Get("ApiKey")
		if len(apikey) <= 0 {
			return errors.New("Debe enviar una api key válida")
		}
		isApiKeyValid, err := m.Service.GetCuentaByApiKey(apikey)
		if err != nil {
			return fmt.Errorf("error " + err.Error())
		}
		if !isApiKeyValid {
			return fmt.Errorf("api-key no es valida")
		}

		return c.Next()
	}
}
