package main_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"testing"

// 	main "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/api"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/database"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
// 	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
// 	"github.com/gofiber/fiber/v2"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func newDatabaseTest() *database.MySQLClient {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.USER_TEST, config.PASSW_TEST, config.HOST_TEST, config.PORT_TEST, config.DB_NAME)
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		logs.Error("cannot create mysql client")
// 		panic(err)
// 	}
// 	// cada test va a crear una conexion asi que es importante que la libere cuando deje de usarla
// 	// mysqldb, err := db.DB()
// 	// if err != nil {
// 	// 	logs.Error("cannot get database specific driver")
// 	// 	panic(err)
// 	// }
// 	// defer mysqldb.Close()
// 	// nos aseguramos que las tablas esten creadas de acuerdo a las estructuras definidas en las entities
// 	// db.AutoMigrate(
// 	// 	entities.Adquiriente{},
// 	// 	entities.Apilinkcierrelote{},
// 	// 	entities.Channel{},
// 	// 	entities.Cliente{},
// 	// 	entities.Cuenta{},
// 	// 	entities.Cuentacomision{},
// 	// 	entities.Impuesto{},
// 	// 	entities.Installment{},
// 	// 	entities.Installmentdetail{},
// 	// 	entities.Log{},
// 	// 	entities.Mediopago{},
// 	// 	entities.Movimiento{},
// 	// 	entities.Movimientocomisiones{},
// 	// 	entities.Notificacione{},
// 	// 	entities.Pago{},
// 	// 	entities.Pagoestado{},
// 	// 	entities.Pagoestadoexterno{},
// 	// 	entities.Pagoestadologs{},
// 	// 	entities.Pagointento{},
// 	// 	entities.Pagoitems{},
// 	// 	entities.Pagotipo{},
// 	// 	entities.Prismacierrelote{},
// 	// 	entities.Transferencia{},
// 	// )

// 	return &database.MySQLClient{DB: db}
// }

// var app *fiber.App

// func TestMainApp(t *testing.T) {
// 	httpClient := http.DefaultClient
// 	sqlClient := newDatabaseTest()
// 	osFile := os.File{}

// 	app = main.InicializarApp(httpClient, sqlClient, &osFile)

// 	// limpiamos la tabla
// 	// insertamos datos necesarios para la prueba

// 	request := dtos.PagoRequest{
// 		PayerName:         "Fer",
// 		Description:       "Test de integracion",
// 		FirstTotal:        1000,
// 		FirstDueDate:      "10-07-2021",
// 		ExternalReference: "15685",
// 		SecondDueDate:     "10-08-2021",
// 		SecondTotal:       1010,
// 		PayerEmail:        "fernando.castro@telco.com.ar",
// 		PaymentType:       "sellos",
// 	}
// 	requestJson, _ := json.Marshal(request)
// 	req, _ := http.NewRequest("POST", "/checkout", bytes.NewBuffer(requestJson))
// 	req.Header.Add("apikey", "123123123123123")
// 	req.Header.Add("content-type", "application/json")
// 	req.Header.Add("Cache-Control", "no-cache")
// 	resp, err := app.Test(req, -1)
// 	if err != nil {
// 		logs.Error("error al ejecutar el pago: " + err.Error())
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != 201 {
// 		bytresp, _ := io.ReadAll(resp.Body)
// 		fmt.Println(string(bytresp))
// 	}
// 	var response *dtos.CheckoutResponse
// 	json.NewDecoder(resp.Body).Decode(&response)

// }
