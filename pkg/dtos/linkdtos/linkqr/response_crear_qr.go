package linkqr

import (
	"math/rand"
	"time"
)

type QrDato struct {
	QRData        string `json:"qr_data"`
	FechaCreacion string `json:"fecha_creacion"`
	OperacionID   int    `json:"operacion_id"`
}

type ResponseApilinkCrearQr struct {
	Status     string `json:"status"`
	ReturnCode string `json:"return_code"`
	Message    string `json:"message"`
	Data       QrDato `json:"data"`
}

type QRTelcoResponse struct {
	Data struct {
		Status     string `json:"status"`
		ReturnCode string `json:"return_code"`
		Message    string `json:"message"`
		QrDato     QrDato `json:"data"`
	} `json:"data"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomQRData() string {
	// Genera datos aleatorios para qr_data (en este caso, se mantiene el mismo valor constante)
	return "00020101021143270023ar.com.redlink.columbia5015001120013048372512600220000174400000000083090520497005303032540536.405802AR5914EmpanadasJose6004CABA6104102862130709QRTERMDEF800800044700630491BD"
}

func GenerateRandomDate() string {
	// Genera una fecha aleatoria
	return time.Now().Format(time.RFC3339)
}

func GenerateRandomOperacionID() int {
	// Genera un ID de operación aleatorio
	return rand.Intn(9000) + 1000
}

func GenerateRandomResponse() ResponseApilinkCrearQr {
	qrData := QrDato{
		QRData:        GenerateRandomQRData(),
		FechaCreacion: GenerateRandomDate(),
		OperacionID:   GenerateRandomOperacionID(),
	}

	response := ResponseApilinkCrearQr{
		Status:     "SUCCESS",
		ReturnCode: "solicitud_qr_ok",
		Message:    "La solicitud de qr fue ejecutada correctamente",
		Data:       qrData,
	}

	return response
}
