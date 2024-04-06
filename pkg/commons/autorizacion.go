package commons

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func VerificarAutorizacion(c *fiber.Ctx, username string, password string) (autorizado bool, err error) {
	// Obtener el encabezado "Authorization" de la solicitud
	authHeader := c.Get(fiber.HeaderAuthorization)

	// Verificar que el encabezado "Authorization" está presente
	if authHeader == "" {
		err := fmt.Errorf("no se encuentra el encabezado.")
		return false, err
	}

	// Decodificar las credenciales Basic Auth del encabezado
	auth := strings.SplitN(authHeader, " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		err := fmt.Errorf("no se puede decodificar basic auth.")
		return false, err
	}

	payload, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		err := fmt.Errorf("no estás autorizado para realizar esta acción.")
		return false, err
	}

	credentials := strings.SplitN(string(payload), ":", 2)
	if len(credentials) != 2 {
		err := fmt.Errorf("no estás autorizado para realizar esta acción.")
		return false, err
	}

	usernameAuth, passwordAuth := credentials[0], credentials[1]

	// Verificar si las credenciales coinciden con las esperadas
	if usernameAuth != username || passwordAuth != password {
		err := fmt.Errorf("las credenciales no coinciden.")
		return false, err
	}
	autorizado = true
	return autorizado, err
}
