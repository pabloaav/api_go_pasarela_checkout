package apilink

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/linkdtos/linkqr"
)

func (s *aplinkService) CreateQrApiLinkService(ibmClienteid string, request linkqr.RequestApilinkCrearQr) (response linkqr.QRTelcoResponse, erro error) {

	response, erro = s.remoteRepository.CreateQrTelcoRemoteRepository(ibmClienteid, request)
	if erro != nil {
		return
	}
	return
}
