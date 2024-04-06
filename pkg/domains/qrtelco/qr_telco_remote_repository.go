package qrtelco

import (
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/linkdtos/linkqr"
)

type RepositoryRemotePrisma interface {
	PostCreateQrTelcoRemoteRepository(request linkqr.RequestApilinkCrearQr) (err error)
}

/* type remoteRepository struct {
	HTTPClient  *http.Client
	UtilService util.UtilService
}

func NewRemote(http *http.Client, u util.UtilService) RepositoryRemotePrisma {
	return &remoteRepository{
		HTTPClient:  http,
		UtilService: u,
	}
}

//
func (r *remoteRepository) PostCreateQrTelcoRemoteRepository() (request linkqr.RequestApilinkCrearQr, err error) {
	var erro ErrorEstructura
	payload := strings.NewReader("{}")
	base, err := url.Parse(config.URL_PRISMA)
	if err != nil {
		logs.Error("Error al crear base url" + err.Error())
	}
	base.Path += config.URI_HEALTHCHECK
	logs.Info(base.String())
	req, _ := http.NewRequest("GET", base.String(), payload)
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		logs.Error("error al verificar servico de prisma: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		json.NewDecoder(resp.Body).Decode(&erro)
		fmt.Println(erro)
		return response, &erro
	}
	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil
}
*/
