package util

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

func GenerateQRImage(url string) (string, error) {
	// Crear un código QR con la URL proporcionada.
	qrCode, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return "", err
	}

	// Generar el código QR como bytes en formato PNG.
	qrImage, err := qrCode.PNG(256)
	if err != nil {
		return "", err
	}

	// Convertir los bytes de la imagen en una cadena base64.
	qrBase64 := base64.StdEncoding.EncodeToString(qrImage)

	return qrBase64, nil
}
