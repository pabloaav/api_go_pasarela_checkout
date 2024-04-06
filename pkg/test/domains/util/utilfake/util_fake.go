package utilfake

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/utildtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
)

func EstructuraVerificarCbu() (tableDriverTestPeyment TableDriverTestConsultarMoviento) {
	tableDriverTestPeyment = TableDriverTestConsultarMoviento{
		TituloPrueba: "el tipo de movimiento no es valido, los valores correctos son debin, prisma, transferencia",
		WantTable:    true,
		Cbu:          "56477491421212212121212",
	}
	return
}

// // Construir el texto html del mensaje del email
// mensaje := "<p style='box-sizing:border-box;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif,'Apple Color Emoji','Segoe UI Emoji','Segoe UI Symbol';font-size:16px;line-height:1.5em;margin-top:0;text-align:center'><h2 style='text-align:center'>Operación de pago exitosa</h2> El pago de la referencia <b>#4</b> fue aprobado. <ul><li> Importe: <b>#0</b></li><li> Identificador de la transacción: <b>#1</b></li><li> Medio de pago: <b>#2</b></li><li> Concepto: <b>#3</b></li></ul></p>"
// /* enviar mail al usuario pagador */
// var arrayEmail []string
// var email string
// email = request.HolderEmail
// if request.HolderEmail == "" {
// 	email = pago.PayerEmail
// }
// arrayEmail = append(arrayEmail, email)
// params := utildtos.RequestDatosMail{
// 	Email:            arrayEmail,
// 	Asunto:           "Información de Pago",
// 	Nombre:           pago.PayerName,
// 	Mensaje:          mensaje,
// 	CamposReemplazar: []string{fmt.Sprintf("$%v", response.ImportePagado), pago.Uuid, medio.Mediopago, response.Description, response.ExternalReference},
// 	From:             "Wee.ar!",
// 	TipoEmail:        "template",
// }

func EstructuraEmail() (tableDriverTestPeyment TableDriverTestEmailSend) {
	tableDriverTestPeyment = TableDriverTestEmailSend{
		TituloPrueba: "envios de email en pagos exitosos",
		WantTable:    "",
		Request: utildtos.RequestDatosMail{
			Email:            []string{"jose.alarcon@telco.com.ar"},
			Asunto:           "Información de Pago",
			Nombre:           "jose",
			Mensaje:          "<p style='box-sizing:border-box;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif,'Apple Color Emoji','Segoe UI Emoji','Segoe UI Symbol';font-size:16px;line-height:1.5em;margin-top:0;text-align:center'><h2 style='text-align:center'>Operación de pago exitosa</h2> El pago de la referencia <b>#4</b> fue aprobado. <ul><li> Importe: <b>#0</b></li><li> Identificador de la transacción: <b>#1</b></li><li> Medio de pago: <b>#2</b></li><li> Concepto: <b>#3</b></li></ul></p>",
			CamposReemplazar: []string{"djasjds", "dsadsad", "dasdasdasdasd"},
			From:             "Wee.ar",
			TipoEmail:        "template",
		},
	}
	return
}

func EstructuraValidarCbu() (tableDriverTest TableDriverTestConsultarMoviento) {
	tableDriverTest = TableDriverTestConsultarMoviento{
		TituloPrueba: "validar cbu",
		WantTable:    true,
		Cbu:          "0940099372007393130021",
	}
	return
}

func EstructuraBuildComisiones() (tableDriverTest []TableDriverBuildComisiones) {
	// var cuentaComision *[]entities.Cuentacomision
	var tableDriver []TableDriverBuildComisiones
	var err error
	cuentaComisionOffline := append([]entities.Cuentacomision{}, entities.Cuentacomision{
		CuentasID:          4,
		ChannelsId:         3,
		ChannelarancelesId: 15,
		Cuentacomision:     "VIKING OFFLINE",
		Comision:           0.001,
		Mediopagoid:        0,
		Importeminimo:      10,
		Importemaximo:      0,
		ChannelArancel: entities.Channelarancele{
			ChannelsId:    3,
			RubrosId:      2,
			Importe:       0.0275,
			Tipocalculo:   "PORCENTAJE",
			Importeminimo: 30,
			Importemaximo: 0,
			Mediopagoid:   0,
			Pagocuota:     false,
		},
	})

	cuentaComisionDebito := append([]entities.Cuentacomision{}, entities.Cuentacomision{
		CuentasID:          4,
		ChannelsId:         3,
		ChannelarancelesId: 15,
		Cuentacomision:     "VIKING DEBITO",
		Comision:           0.0060,
		Mediopagoid:        0,
		Importeminimo:      10,
		Importemaximo:      0,
		ChannelArancel: entities.Channelarancele{
			ChannelsId:    3,
			RubrosId:      2,
			Importe:       0.0035,
			Tipocalculo:   "PORCENTAJE",
			Importeminimo: 0,
			Importemaximo: 0,
			Mediopagoid:   0,
			Pagocuota:     false,
		},
	})

	cuentaComisionCredito := append([]entities.Cuentacomision{}, entities.Cuentacomision{
		CuentasID:          4,
		ChannelsId:         3,
		ChannelarancelesId: 15,
		Cuentacomision:     "VIKING CREDITO",
		Comision:           0.0060,
		Mediopagoid:        0,
		Importeminimo:      0,
		Importemaximo:      0,
		ChannelArancel: entities.Channelarancele{
			ChannelsId:    3,
			RubrosId:      2,
			Importe:       0.01,
			Tipocalculo:   "PORCENTAJE",
			Importeminimo: 0,
			Importemaximo: 0,
			Mediopagoid:   0,
			Pagocuota:     false,
		},
	})

	tableDriverTest = append(tableDriver, TableDriverBuildComisiones{
		// &PRUEBA 1 - OFFLINE
		TituloPrueba: "Calculo de comisioones OFFLINE : Telco minimo y Proveedor minimo",
		WantTable:    err,
		RequestMovimiento: &entities.Movimiento{
			CuentasId:      4,
			PagointentosId: 187,
			Tipo:           "C",
			Monto:          268710,
			MotivoBaja:     "",
			Reversion:      false,
			Enobservacion:  false,
		},
		RequestCuentaComision: &cuentaComisionOffline,
		RequestIva: &entities.Impuesto{
			Impuesto:   "IVA",
			Porcentaje: 0.21,
		},
		ImporteSolicitado: 268710},
		// &PRUEBA 2 - DEBITO
		TableDriverBuildComisiones{
			TituloPrueba: "Calculo de comisiones DEBITO: Telco minimo y Proveedor sin minimo",
			WantTable:    err,
			RequestMovimiento: &entities.Movimiento{
				CuentasId:      4,
				PagointentosId: 187,
				Tipo:           "C",
				Monto:          50000,
				MotivoBaja:     "",
				Reversion:      false,
				Enobservacion:  false,
			},
			RequestCuentaComision: &cuentaComisionDebito,
			RequestIva: &entities.Impuesto{
				Impuesto:   "IVA",
				Porcentaje: 0.21,
			},
			ImporteSolicitado: 50000},
		// &PRUEBA 3 - CREDITO
		TableDriverBuildComisiones{
			TituloPrueba: "Calculo de comisiones Credito: Telco sin minimo y Proveedor sin minimo",
			WantTable:    err,
			RequestMovimiento: &entities.Movimiento{
				CuentasId:      4,
				PagointentosId: 187,
				Tipo:           "C",
				Monto:          1120000,
				MotivoBaja:     "",
				Reversion:      false,
				Enobservacion:  false,
			},
			RequestCuentaComision: &cuentaComisionCredito,
			RequestIva: &entities.Impuesto{
				Impuesto:   "IVA",
				Porcentaje: 0.21,
			},
			ImporteSolicitado: 1120000},
	)

	return
}
