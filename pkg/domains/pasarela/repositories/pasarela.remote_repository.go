package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
	dto "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/pasarela/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/filtros/administracion"
)

// Interfaz
type PasarelaRemoteRepository interface {
	GetPlanCuotasRepository() (response []dto.PlanCuotasResponse, err error)
	GetConfiguraciones(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, totalFilas int64, erro error) 
	GetInstallments(fechaDesde time.Time) (medioPagoInstallments []entities.Mediopagoinstallment, erro error) 

}

// Estructura
type remoteRepository struct {
	HTTPClient *http.Client
	SQLClient        *database.MySQLClient
	utilService      util.UtilService

}

// Constructor
func NewPasarelaRepository(http *http.Client,sqlClient *database.MySQLClient, t util.UtilService) PasarelaRemoteRepository {
	return &remoteRepository{
		HTTPClient: http,
		SQLClient:        sqlClient,
		utilService:      t,
	}
}

func (r *remoteRepository) GetPlanCuotasRepository() (response []dto.PlanCuotasResponse, err error) {

	payload := strings.NewReader("{}")
	base, err := url.Parse(config.API_PASARELA_URL)
	if err != nil {
		logs.Error("Error al crear base url" + err.Error())
	}
	base.Path += PLAN_CUOTAS

	logs.Info(base.String())

	request, _ := http.NewRequest("GET", base.String(), payload)

	resp, err := r.HTTPClient.Do(request)

	if err != nil {
		logs.Error("error al comunicarse con el servicio de pasarela: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return response, err
	}

	json.NewDecoder(resp.Body).Decode(&response)
	return response, nil

}
func (r *remoteRepository) GetConfiguraciones(filtro filtros.ConfiguracionFiltro) (configuraciones []entities.Configuracione, totalFilas int64, erro error) {

	resp := r.SQLClient.Model(entities.Configuracione{})
	
	if filtro.Id > 0 {
		resp.Where("id = ?", filtro.Id)
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("nombre like ?", fmt.Sprintf("%%%s%%", filtro.Nombre))
	}

	if filtro.Number > 0 && filtro.Size > 0 {

		resp.Count(&totalFilas)

		if resp.Error != nil {
			erro = fmt.Errorf("error al cargar el total de filas de la consulta")
		}

		offset := (filtro.Number - 1) * filtro.Size
		resp.Limit(int(filtro.Size))
		resp.Offset(int(offset))

	}

	resp.Find(&configuraciones)

	if resp.Error != nil {

		erro = fmt.Errorf("error al cargar las configuraciones")

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConfiguraciones",
		}

		err := r.utilService.CreateLogService(log)

		if err != nil {
			mensaje := fmt.Sprintf("Crear Log: %s. %s", err.Error(), resp.Error.Error())
			logs.Error(mensaje)
		}
	}

	return
}
func (r *remoteRepository) GetInstallments(fechaDesde time.Time) (medioPagoInstallments []entities.Mediopagoinstallment, erro error) {
	res := r.SQLClient.Table("mediopagoinstallments as mpi")
	res.Preload("Installments")
	res.Preload("Installments.Installmentdetail")
	res.Find(&medioPagoInstallments)
	if res.Error != nil {
		logs.Info(res.Error)
		erro = errors.New("error al crear detalle plan de cuota")
		return
	}
	return
}