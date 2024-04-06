package main

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/api/middlewares/middlewareapikey"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/api/routes"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/auditoria"
	checkout_repository "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/checkout/repositories"
	checkout_service "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/checkout/services"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/pagooffline"
	pasarela_repository "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/pasarela/repositories"

	pasarela_service "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/pasarela/services"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/prisma"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/webhook"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
)

var WebSocketConnection *websocket.Conn

func InicializarApp(clienteHttp *http.Client, clienteSql *database.MySQLClient, clienteFile *os.File) *fiber.App {
	//Servicios comunes
	fileRepository := commons.NewFileRepository(clienteFile) // que es esto
	commonsService := commons.NewCommons(fileRepository)
	algoritmoVerificacionService := commons.NewAlgoritmoVerificacion()
	middlewares := middlewares.MiddlewareManager{HTTPClient: clienteHttp}

	utilRepository := util.NewUtilRepository(clienteSql)
	utilService := util.NewUtilService(utilRepository)

	//ApiLink
	apiLinkRemoteRepository := apilink.NewRemote(clienteHttp, utilService)
	apiLinkRepository := apilink.NewRepository(clienteSql, utilService)
	apilink.NewService(apiLinkRemoteRepository, apiLinkRepository)

	auditoriaRespository := auditoria.NewAuditoriaRepository(clienteSql)
	auditoriaService := auditoria.AuditoriaService(auditoriaRespository)

	/* Pasarela Backend. Se crea el repositorio y el servicio */
	pasarelaRepository := pasarela_repository.NewPasarelaRepository(clienteHttp, clienteSql, utilService)
	pasarelaService := pasarela_service.NewPasarelaService(pasarelaRepository, utilService)

	prismaRepository := prisma.NewRepository(clienteSql)
	remoteRepository := prisma.NewRepoasitory(clienteHttp, prismaRepository)
	prismaService := prisma.NewService(remoteRepository, prismaRepository, commonsService)
	pagoOffLineService := pagooffline.NewService(algoritmoVerificacionService)

	// webhooks
	webhooksRepository := webhook.NewRemote(clienteHttp)

	checkoutRepository := checkout_repository.NewRepository(clienteSql, auditoriaService, utilService)
	checkoutService := checkout_service.NewService(checkoutRepository, commonsService, prismaService, pagoOffLineService, utilService, webhooksRepository)

	// Middleware para validar apikey por cuenta
	middlewaresApiKey := middlewareapikey.MiddlewareApiKey{Service: checkoutService}

	//descomentar esto en servidor
	engine := html.New(filepath.Join(filepath.Base(config.DIR_BASE), "api", "views"), ".html")
	//descomentar esto en local
	// engine := html.New("views", ".html")
	engine.Delims("${", "}")
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var msg string
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				msg = e.Message
			}

			if msg == "" {
				msg = "No se pudo procesar el llamado a la api: " + err.Error()
			}

			_ = ctx.Status(code).JSON(internalError{
				Message: msg,
			})

			return nil
		},
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.ALLOW_ORIGIN + ", http://127.0.0.1:3300, " + config.AUTH,
		AllowHeaders: "",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	/* RUTAS */

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Corrientes Telecomunicaciones API Checkout"))
	})
	checkout := app.Group(config.API_VERSION + "/checkout")
	routes.CheckoutRoutes(checkout, middlewaresApiKey, checkoutService)

	/* 	Crear grupo de rutas Pasarela. Inyectar servicio Pasarela */
	pasarela := app.Group(config.API_VERSION + "/pasarela")
	routes.PasarelaRoutes(pasarela, pasarelaService, middlewares)

	prisma := app.Group(config.API_VERSION + "/prisma")
	routes.PrismaRoutes(prisma, prismaService, middlewares)

	//descomentar esto en local
	// app.Static("/", "./views")
	//descomentar esto en servidor
	app.Static("/", filepath.Join(filepath.Base(config.DIR_BASE), "api", "views"))

	return app
}

func main() {
	var HTTPTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     false, // <- this is my adjustment
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	var HTTPClient = &http.Client{
		Transport: HTTPTransport,
	}

	//HTTPClient.Timeout = time.Second * 120 //Todo validar si este tiempo está bien
	clienteSQL := database.NewMySQLClient()
	osFile := os.File{}

	app := InicializarApp(HTTPClient, clienteSQL, &osFile)

	_ = app.Listen(":3303")
}

type internalError struct {
	Message string `json:"message"`
}
