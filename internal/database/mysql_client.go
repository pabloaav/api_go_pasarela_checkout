package database

import (
	"fmt"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLClient contiene la instancia de base de datos
type MySQLClient struct {
	*gorm.DB
	TX *gorm.DB
}

// NewMySQLClient cliente de la base de datos en MySql
func NewMySQLClient() *MySQLClient {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.DB_USER, config.DB_PASSW, config.DB_HOST, config.DB_PORT, config.DB_NAME)
	// logs.Info(dsn)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
			return time.Now().In(loc)
		},
	})

	// NUEVA CONFIGURACION PARA CONTROL DE CONECCIONES
	database, _ := gormDB.DB()
	database.SetMaxIdleConns(20)           // maximo de conexiones inactivas que se mantendra en el grupo de conexiones
	database.SetMaxOpenConns(50)           // maximo de conexiones abiertas
	database.SetConnMaxLifetime(time.Hour) //establece el numero maximo de conexiones abiertas permitidas
	//establece la duración máxima que una conexión puede permanecer abierta en el pool antes de ser cerrada y reemplazada por una nueva.

	if err != nil {
		logs.Error("cannot connect database" + err.Error())
		panic(err)
	}

	return &MySQLClient{gormDB, nil}
}
