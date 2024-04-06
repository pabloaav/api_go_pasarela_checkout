package apilink

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/linkdtos/linkqr"
)

func (r *remoteRepository) CreateQrTelcoRemoteRepository(ibmClienteid string, request linkqr.RequestApilinkCrearQr) (response linkqr.QRTelcoResponse, erro error) {
	base, erro := url.Parse(config.HOSTTELCOQR)

	if erro != nil {
		return
	}

	base.Path += "qr/create"

	json_data, _ := json.Marshal(request)

	req, _ := http.NewRequest("POST", base.String(), bytes.NewBuffer(json_data))

	erro = executeRequest(r, req, ERROR_CREATE_QR, &response)

	return
}
