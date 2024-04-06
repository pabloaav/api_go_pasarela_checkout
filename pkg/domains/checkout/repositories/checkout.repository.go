package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/internal/logs"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/auditoria"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/dtos/qrcierrelotesdtos"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"
	filtros "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/filtros/administracion"
	"gorm.io/gorm"
)

type Repository interface {
	BeginTx()
	CommitTx()
	RollbackTx()
	CreatePago(ctx context.Context, pago *entities.Pago) (*entities.Pago, error)
	UpdatePago(ctx context.Context, pago *entities.Pago) (bool, error)
	UpdatePagoEstado(ctx context.Context, pago entities.Pago) (bool, error)
	UpdatePagoIntento(ctx context.Context, pagointento entities.Pagointento) (bool, error)
	UpdatePagoMP(ctx context.Context, pago *entities.Pago, pi *entities.Pagointento) (bool, error)

	GetPaymentByUuid(filtroPago filtros.PagoFiltro) (*entities.Pago, error)
	GetPaymentMultipago(request filtros.FiltroMultipago) (pagos []entities.Pago, erro error)
	GetPayments(request filtros.PagoFiltro) (pagos []entities.Pago, erro error)
	GetClienteByApikey(apikey string) (*entities.Cliente, error)
	GetCuentaByApikey(apikey string) (*entities.Cuenta, error)
	GetPagotipoById(id int64) (*entities.Pagotipo, error)
	GetPagotipoChannelByPagotipoId(id int64) (*[]entities.Pagotipochannel, error)
	GetPagotipoIntallmentByPagotipoId(id int64) (*[]entities.Pagotipointallment, error)
	GetChannelByName(nombre string) (*entities.Channel, error)
	GetCuentaById(id int64) (*entities.Cuenta, error)
	CreateResultado(ctx context.Context, resultado *entities.Pagointento) (bool, error)
	GetValidPagointentoByPagoId(pagoId int64) (*entities.Pagointento, error)
	GetMediosDePagos() (*[]entities.Mediopago, error)
	GetMediopago(filtro map[string]interface{}) (*entities.Mediopago, error)
	GetInstallmentDetailsID(installmentID, numeroCuota int64) int64
	GetInstallmentDetails(installmentID, numeroCuota int64) (installmentDetails *dtos.InstallmentDetailsResponse, erro error)
	GetInstallmentsByMedioPagoInstallmentsId(id int64) (installments []entities.Installment, erro error)
	CreatePagoEstadoLog(ctx context.Context, pel *entities.Pagoestadologs) error
	// GetCuentaTelco() (*entities.Configuracione, error)
	GetPagoEstado(id int64) (*entities.Pagoestado, error)

	GetPreferencesByIdClienteRepository(id uint) (preferencia entities.Preference, erro error)
	GetChannelById(id uint) (channel entities.Channel, erro error)

	UpdateEstadoNotificadoInicial(id uint) error

	// QRCIERRELOTES
	GetPagosRepository(filtro filtros.PagoFiltro) (pagos []entities.Pago, erro error)
	CreateQrcierrelotesRepository(ctx context.Context, requestQrcierrelote *qrcierrelotesdtos.RequestCreateQrCierrelotes) error

	// PAGOINTENTOS
	GetPagointentosRepository(filtro filtros.PagoIntentoFiltro) (pagointentos []entities.Pagointento, err error)

	// PAGOESTADOS
	GetPagosEstadosExternos(filtro filtros.PagoEstadoExternoFiltro) (estados []entities.Pagoestadoexterno, erro error)

	UpdateEstadoNotificadoOnline(id uint) error

	//Hash
	SaveHasheado(hasheado *entities.Uuid, pagointento_id uint) (erro error)
	GetHasheado(hash string) (control bool, erro error)

	//Rapipago
	GetPagosRapipago(filtro filtros.RapipagosFiltro) (pagos []entities.Pago, err error)
	GetPagosRapipagoAprobados(filtro filtros.RapipagosFiltro) (pagos []entities.Pago, err error)

	//Apikey Adquiriente
	GetApikeyAdquiriente(apikey string, adquiriente string) (control bool, err error)

	// Estado Aplicacion
	GetEstadoAplicacionRepository(filtro filtros.ConfiguracionFiltro) (configuracion entities.Configuracione, err error)

	// Middleware de apikey por cuenta
	GetCuentaByApiKey(apikey string) (cuenta *entities.Cuenta, err error)
}

type repository struct {
	SQLClient        *database.MySQLClient
	auditoriaService auditoria.AuditoriaService
	utilService      util.UtilService
}

func NewRepository(sqlClient *database.MySQLClient, a auditoria.AuditoriaService, t util.UtilService) Repository {
	return &repository{
		SQLClient:        sqlClient,
		auditoriaService: a,
		utilService:      t,
	}
}

func (r *repository) BeginTx() {
	r.SQLClient.TX = r.SQLClient.DB
	r.SQLClient.DB = r.SQLClient.Begin()
}

func (r *repository) CommitTx() {
	r.SQLClient.Commit()
	r.SQLClient.DB = r.SQLClient.TX
}

func (r *repository) RollbackTx() {
	r.SQLClient.Rollback()
	r.SQLClient.DB = r.SQLClient.TX
}

func (r *repository) GetPagosRepository(filtro filtros.PagoFiltro) (pagos []entities.Pago, erro error) {
	return
}

func (r *repository) auditarCheckout(ctx context.Context, resultado interface{}) error {
	audit := ctx.Value(entities.AuditUserKey{}).(entities.Auditoria)

	audit.Operacion = strings.ToLower(audit.Query[:6])

	audit.Origen = "pasarela.checkout"

	res, _ := json.Marshal(resultado)
	audit.Resultado = string(res)

	err := r.auditoriaService.Create(&audit)

	if err != nil {
		return fmt.Errorf("auditoria: %w", err)
	}

	return nil
}

func (r *repository) CreatePago(ctx context.Context, pago *entities.Pago) (*entities.Pago, error) {
	res := r.SQLClient.WithContext(ctx).Create(&pago)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se pudo generar registro pago: " + res.Error.Error())
	}
	err := r.auditarCheckout(res.Statement.Context, res.RowsAffected)
	if err != nil {
		return nil, err
	}
	return pago, nil
}

func (r *repository) UpdatePago(ctx context.Context, pago *entities.Pago) (bool, error) {
	res := r.SQLClient.WithContext(ctx).Model(&pago).Updates(&pago)
	if res.Error != nil {
		return false, fmt.Errorf("al actualizar el estado del pago: %s", res.Error.Error())
	}
	err := r.auditarCheckout(res.Statement.Context, res.RowsAffected)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) UpdatePagoMP(ctx context.Context, pago *entities.Pago, pi *entities.Pagointento) (bool, error) {

	// se actualiza el estado del pago
	res := r.SQLClient.Table("pagos").Where("id = ?", pago.ID).Update("pagoestados_id", 4)

	if res.Error != nil {
		return false, fmt.Errorf("al actualizar el estado del pago: %s", res.Error.Error())
	}

	// se actualiza la fecha de pago del pagointento
	res2 := r.SQLClient.Table("pagointentos").Where("id = ?", pi.ID).Update("mediopagos_id", pi.MediopagosID)
	if res2.Error != nil {
		return false, fmt.Errorf("al actualizar el medio de pago del pagointento: %s", res2.Error.Error())
	}

	return true, nil
}

func (r *repository) UpdatePagoEstado(ctx context.Context, pago entities.Pago) (bool, error) {
	res := r.SQLClient.WithContext(ctx).Model(&pago).Where("id = ?", pago.ID).Update("pagoestados_id", pago.PagoestadosID)
	if res.Error != nil {
		return false, fmt.Errorf("al actualizar el estado del pago: %s", res.Error.Error())
	}

	return true, nil
}

func (r *repository) UpdatePagoIntento(ctx context.Context, pagointento entities.Pagointento) (bool, error) {
	res := r.SQLClient.WithContext(ctx).Model(&pagointento).Where("id = ?", pagointento.ID).Update("paid_at", pagointento.PaidAt)
	if res.Error != nil {
		return false, fmt.Errorf("al actualizar el estado del pago: %s", res.Error.Error())
	}

	return true, nil
}

// obtiene un pago por uuid. preload de pagointentos y pagoitems
func (r *repository) GetPaymentByUuid(filtroPago filtros.PagoFiltro) (*entities.Pago, error) {

	var pago entities.Pago

	res := r.SQLClient.Model(entities.Pago{}).Preload("Pagoitems").Preload("PagoIntentos")

	if filtroPago.CargaMedioPagos {
		res = res.Preload("PagoIntentos.Mediopagos")
	}

	if filtroPago.CargarPagoEstado {
		res = res.Preload("PagoEstados")
	}

	res = res.Where("uuid = ?", filtroPago.Uuids[0]).Find(&pago)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no existe pago con identificador %s", filtroPago.Uuids[0])
	}

	return &pago, nil
}

func (r *repository) GetPaymentMultipago(filtro filtros.FiltroMultipago) (pagos []entities.Pago, erro error) {

	res := r.SQLClient.Model(entities.Pago{})
	res.Preload("PagoIntentos").Preload("PagoEstados")
	res.Preload("PagosTipo.Cuenta.Cliente")

	res.Joins("JOIN pagointentos ON pagointentos.pagos_id = pagos.id")

	res.Where("pagointentos.barcode != '' ").Where("pagos.pagoestados_id = 2 ")

	if filtro.Codigo != "" {
		res.Where("pagointentos.barcode = ?", filtro.Codigo).Find(&pagos)
	}

	if filtro.ValorDoc != "" {
		res.Where("pagointentos.holder_number = ?", filtro.ValorDoc).Find(&pagos)
	}

	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no existe pago para la consulta")
	}

	return
}

func (r *repository) GetPayments(request filtros.PagoFiltro) (pagos []entities.Pago, erro error) {

	res := r.SQLClient.Model(entities.Pago{})

	if request.CargaPagoIntentos {
		res.Preload("PagoIntentos")
	}

	if request.CargarPagoEstado {
		res.Preload("PagoEstados")
	}

	if request.CargarCuenta {
		res.Preload("PagosTipo.Cuenta.Cliente")
	}

	if request.PagoEstadosIds != nil {
		res.Where("pagoestados_id in ?", request.PagoEstadosIds)
	}

	if !request.FechaPagoInicio.IsZero() {

		res.Where("paid_at >= ?", request.FechaPagoInicio.Format("2006-01-02"))
	}
	if !request.FechaPagoFin.IsZero() {
		res.Where("paid_at <= ?", request.FechaPagoFin.Format("2006-01-02"))

	}

	if request.PagoEstadosIds != nil {
		res.Where("pagoestados_id in ?", request.PagoEstadosIds)
	}

	res.Joins("JOIN pagointentos ON pagointentos.pagos_id = pagos.id")

	if request.MedioPagoId != 0 {
		res.Where("pagointentos.mediopagos_id = ?", request.MedioPagoId)
	}

	res.Find(&pagos)

	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no existe pago para la consulta")
	}

	return
}

func (r *repository) GetClienteByApikey(apikey string) (*entities.Cliente, error) {

	var cliente entities.Cliente

	res := r.SQLClient.Preload("Cuentas.Pagotipos").Joins("JOIN cuentas on cuentas.clientes_id = clientes.id and cuentas.apikey = ?", apikey).Find(&cliente)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró cliente con apikey: %s", apikey)
	}

	return &cliente, nil
}

func (r *repository) GetCuentaByApikey(apikey string) (*entities.Cuenta, error) {
	var cuenta entities.Cuenta
	res := r.SQLClient.Preload("Pagotipos").Where("apikey = ?", apikey).Find(&cuenta)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró cuenta con apikey: %s", apikey)
	}
	return &cuenta, nil
}

// obtiene un pagotipo por su id
func (r *repository) GetPagotipoById(id int64) (*entities.Pagotipo, error) {
	var tipo entities.Pagotipo
	// res := r.SQLClient.Find(&tipo, id)
	res := r.SQLClient.Table("pagotipos").Where("id=?", id)
	res.Preload("Cuenta.Cliente")
	res.Find(&tipo)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró tipo de pago con el id: %d", id)
	}
	return &tipo, nil
}

// obtiene pagotipochannels mediante el id del pagotipo
func (r *repository) GetPagotipoChannelByPagotipoId(id int64) (*[]entities.Pagotipochannel, error) {
	var channels []entities.Pagotipochannel
	res := r.SQLClient.Preload("Channel").Where("pagotipos_id = ?", id).Find(&channels)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró tipo de pago con el id: %d", id)
	}
	return &channels, nil
}

func (r *repository) GetPagotipoIntallmentByPagotipoId(id int64) (*[]entities.Pagotipointallment, error) {
	var installmentdetails []entities.Pagotipointallment
	res := r.SQLClient.Where("pagotipos_id = ?", id).Find(&installmentdetails)
	if res.Error != nil {
		return nil, fmt.Errorf("no se encontró cuotas para el tipo de pago con el id: %d", id)
	}

	return &installmentdetails, nil
}

func (r *repository) GetChannelByName(nombre string) (*entities.Channel, error) {
	var channel entities.Channel

	res := r.SQLClient.Where("channel = ?", nombre).Find(&channel)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró metodo de pago con la descripción %s", nombre)
	}

	return &channel, nil
}

func (r *repository) GetCuentaById(id int64) (*entities.Cuenta, error) {
	var cuenta entities.Cuenta

	res := r.SQLClient.Preload("Cliente").Find(&cuenta, id)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró cuenta con el id: %d", id)
	}
	// Incluir datos del cliente

	return &cuenta, nil
}

func (r *repository) CreateResultado(ctx context.Context, resultado *entities.Pagointento) (bool, error) {

	res := r.SQLClient.WithContext(ctx).Create(&resultado)
	if res.RowsAffected <= 0 {
		return false, fmt.Errorf("error al guardar resultado: %s", res.Error.Error())
	}

	err := r.auditarCheckout(res.Statement.Context, resultado.StateComment)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *repository) GetValidPagointentoByPagoId(pagoId int64) (*entities.Pagointento, error) {
	var intento entities.Pagointento
	res := r.SQLClient.Model(entities.Pagointento{}).Where("external_id != '0' AND pagos_id = ?", pagoId).Last(&intento)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró intento con el id de pago: %d", pagoId)
	}

	return &intento, nil
}

func (r *repository) GetMediosDePagos() (*[]entities.Mediopago, error) {
	var medios []entities.Mediopago
	res := r.SQLClient.Model(entities.Mediopago{})
	res.Preload("Mediopagoinstallment")
	res.Preload("Channel")
	res.Where("mediopagos.regexp != ''").Order("longitud_pan DESC")
	res.Order("codigo_bcra")
	res.Find(&medios)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontraron medios de pago")
	}
	return &medios, nil
}

func (r *repository) GetMediopago(filtro map[string]interface{}) (*entities.Mediopago, error) {
	var medio entities.Mediopago
	res := r.SQLClient.Model(entities.Mediopago{})
	res.Where(filtro)
	res.First(&medio)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontraron medios de pago")
	}
	return &medio, nil
}

func (r *repository) GetInstallmentDetailsID(installmentID, numeroCuota int64) int64 {
	var response int64
	res := r.SQLClient.Model(entities.Installmentdetail{})
	res.Select("id")
	res.Where("installments_id = ? AND cuota = ?", installmentID, numeroCuota)
	res.First(&response)
	if res.RowsAffected <= 0 {
		return 1
	}
	return response
}

func (r *repository) GetInstallmentDetails(installmentID, numeroCuota int64) (installmentDetails *dtos.InstallmentDetailsResponse, erro error) {
	var result *entities.Installmentdetail
	res := r.SQLClient.Model(entities.Installmentdetail{})
	res.Where("installments_id = ? AND cuota = ?", installmentID, numeroCuota)
	res.First(&result)
	if res.RowsAffected <= 0 {
		erro = errors.New("no se encontraron detalles de cuotas ")
		return
	}
	installmentDetails = &dtos.InstallmentDetailsResponse{
		Id:             result.Model.ID,
		InstallmentsID: result.InstallmentsID,
		NroCuota:       result.Cuota,
		Coeficiente:    result.Coeficiente,
	}
	return
}

func (r *repository) GetInstallmentsByMedioPagoInstallmentsId(id int64) (installments []entities.Installment, erro error) {
	// and vigencia_hasta is null
	res := r.SQLClient.Model(entities.Installment{})
	res.Where("mediopagoinstallments_id = ? ", id).Order("created_at asc")
	res.Find(&installments)
	if res.RowsAffected <= 0 {
		erro = errors.New("no se encontro plan de cuotas ")
		return
	}
	return
}

func (r *repository) CreatePagoEstadoLog(ctx context.Context, pel *entities.Pagoestadologs) error {
	res := r.SQLClient.WithContext(ctx).Create(&pel)
	if res.RowsAffected <= 0 {
		return fmt.Errorf("error al guardar estado log: %s", res.Error.Error())
	}

	err := r.auditarCheckout(res.Statement.Context, res.RowsAffected)
	if err != nil {
		return err
	}

	return nil
}

// func (r *repository) GetCuentaTelco() (*entities.Configuracione, error) {
// 	var cuentaTelco entities.Configuracione
// 	res := r.SQLClient.Model(&cuentaTelco).Where("nombre = ?", "CBU_CUENTA_TELCO").First(cuentaTelco)
// 	if res.RowsAffected <= 0 {
// 		err := errors.New("error al obtener la cuenta cbu telco")
// 		return nil, err
// 	}
// 	return &cuentaTelco, nil
// }

// obtener estado actual de un pago
func (r *repository) GetPagoEstado(id int64) (*entities.Pagoestado, error) {
	var estado entities.Pagoestado

	res := r.SQLClient.Find(&estado, id)
	if res.RowsAffected <= 0 {
		return nil, fmt.Errorf("no se encontró cuenta con el id: %d", id)
	}

	return &estado, nil
}

func (r *repository) GetPreferencesByIdClienteRepository(id uint) (preference entities.Preference, erro error) {
	res := r.SQLClient.Table("preferences").Where("clientes_id = ?", id).Preload("Cliente")

	res.Last(&preference)
	if res.RowsAffected == 0 {
		return
	}
	if res.Error != nil {
		erro = errors.New("error en la consulta a preferencias de cliente")
		return
	}

	return
}

func (r *repository) GetChannelById(id uint) (channel entities.Channel, erro error) {
	res := r.SQLClient.Model(entities.Channel{}).Where("id = ?", id).Find(&channel)
	if res.RowsAffected <= 0 {
		erro = errors.New("no se pudo encontrar channel relacionados")
		return
	}
	return
}

func (r *repository) SaveHasheado(hasheado *entities.Uuid, pagointento_id uint) (erro error) {
	res := r.SQLClient.Where("uuid = ?", hasheado.Uuid).FirstOrCreate(&hasheado)

	uuid_pagointento := entities.UuidsPagointento{
		UuidsId:        hasheado.ID,
		PagointentosId: pagointento_id,
	}

	res2 := r.SQLClient.Create(&uuid_pagointento)
	if res2.RowsAffected <= 0 {
		return fmt.Errorf("error al guardar uuid_pi: %s", res.Error.Error())
	}

	return
}

func (r *repository) GetHasheado(hash string) (control bool, erro error) {
	var coindidencias []entities.Uuid
	res := r.SQLClient.Model(entities.Uuid{}).Where("uuid = ?", hash)

	res.Where("fecha_bloqueo > '0000-00-00 00:00:00'")

	res.Find(&coindidencias).Limit(1)

	if res.Error != nil {
		return false, fmt.Errorf("error al buscar hash: %s", res.Error.Error())
	}
	if res.RowsAffected > 0 {
		control = true
	}

	return
}

func (cr *repository) UpdateEstadoNotificadoInicial(id uint) error {

	result := cr.SQLClient.Table("pagos").Where("id = ?", id).Update("estado_inicial_notificado", 1)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repository) CreateQrcierrelotesRepository(ctx context.Context, requestQrcierrelote *qrcierrelotesdtos.RequestCreateQrCierrelotes) error {
	return r.SQLClient.Transaction(func(tx *gorm.DB) error {
		res := r.SQLClient.Create(&requestQrcierrelote.QrCierrelote)
		if res.RowsAffected <= 0 {
			return fmt.Errorf("error al guardar qrcierrelote: %s", res.Error.Error())
		}

		res = r.SQLClient.Table("pagos").WithContext(ctx).Updates(&requestQrcierrelote.Pago)
		if res.Error != nil {
			logs.Info(res.Error)
			return errors.New(commons.ERROR_UPDATE_PAGO)
		}

		res = r.SQLClient.Table("pagointentos").WithContext(ctx).Updates(&requestQrcierrelote.Pagointento)
		if res.Error != nil {
			logs.Info(res.Error)
			return errors.New(commons.ERROR_UPDATE_PAGOINTENTO)
		}

		err := r.auditarCheckout(res.Statement.Context, requestQrcierrelote.QrCierrelote.ID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) GetPagointentosRepository(filtro filtros.PagoIntentoFiltro) (pagointentos []entities.Pagointento, err error) {
	result := r.SQLClient.Table("pagointentos").Preload("Pago").Preload("Mediopagos.Channel").Where("external_id in ?", filtro.ExternalIds)

	if filtro.CargarPagoTipo {
		result = result.Preload("Pago.PagosTipo")
	}
	if filtro.CargarPagoEstado {
		result = result.Preload("Pago.PagoEstados")
	}
	if filtro.CargarPagoEstado {
		result = result.Preload("Pago.Pagoitems")
	}

	result.Find(&pagointentos)
	if result.RowsAffected <= 0 {
		err = errors.New("no se pudo encontrar pagointentos")
		return
	}
	return pagointentos, nil
}

func (r *repository) GetPagosEstadosExternos(filtro filtros.PagoEstadoExternoFiltro) (estados []entities.Pagoestadoexterno, erro error) {

	resp := r.SQLClient.Model(entities.Pagoestadoexterno{})

	if len(filtro.Vendor) > 0 {

		resp.Where("vendor", filtro.Vendor)
	}

	if len(filtro.Nombre) > 0 {
		resp.Where("estado", filtro.Nombre)
	}
	if filtro.CargarEstadosInt {
		resp.Preload("PagoEstados")
	}

	resp.Find(&estados)

	if resp.Error != nil {

		logs.Error(resp.Error)

		erro = fmt.Errorf(commons.ERROR_PAGO_ESTADO_EXTERNO)

		err := r.auditarCheckout(resp.Statement.Context, filtro.Nombre)
		if err != nil {
			return nil, err
		}
	}

	return

}

func (cr *repository) UpdateEstadoNotificadoOnline(id uint) error {

	result := cr.SQLClient.Table("pagointentos").Where("id = ?", id).Update("notificado_online", 1)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

/*!SECTION
var intento entities.Pagointento
res := r.SQLClient.Model(entities.Pagointento{}).Where("external_id != '0' AND pagos_id = ?", pagoId).Last(&intento)
if res.RowsAffected <= 0 {
	return nil, fmt.Errorf("no se encontró intento con el id de pago: %d", pagoId)
}

return &intento, nil

*/

func (r *repository) GetPagosRapipago(filtro filtros.RapipagosFiltro) (pagos []entities.Pago, err error) {

	res := r.SQLClient.Model(pagos).Joins("join pagointentos on pagos.id = pagointentos.pagos_id")
	res.Where("pagointentos.barcode != '' ").Where("pagos.pagoestados_id = 2 ")

	if len(filtro.Barcodes) > 0 {
		res.Where("pagointentos.barcode in ?", filtro.Barcodes)
		res.Preload("PagoIntentos", "pagointentos.barcode in ?", filtro.Barcodes)
	}

	if len(filtro.DNI) > 0 {
		res.Where("pagointentos.holder_number = ?", filtro.DNI)
		res.Preload("PagoIntentos", "pagointentos.holder_number = ?", filtro.DNI)
	}

	if filtro.ConEstado {
		res.Preload("PagoEstados")
	}
	if filtro.ConCliente {
		res.Preload("PagosTipo.Cuenta.Cliente")
	}

	res.Find(&pagos)
	if res.RowsAffected <= 0 {
		err = fmt.Errorf("no se encontró pago con la clave enviada")
		return
	}

	return
}

func (r *repository) GetApikeyAdquiriente(apikey string, adquiriente string) (control bool, err error) {

	var adquirienteEntity entities.Adquiriente

	res := r.SQLClient.Model(entities.Adquiriente{})
	res.Where("adquiriente = ? ", adquiriente).Where("apikey = ? ", apikey)

	res.Find(&adquirienteEntity)

	if res.RowsAffected <= 0 {
		err = fmt.Errorf("La apikey no es válida.")
		return
	}

	control = true

	return
}

func (r *repository) GetPagosRapipagoAprobados(filtro filtros.RapipagosFiltro) (pagos []entities.Pago, err error) {

	res := r.SQLClient.Model(pagos).
		Joins("join pagointentos on pagos.id = pagointentos.pagos_id").
		Preload("PagoIntentos", "pagointentos.barcode in ?", filtro.Barcodes).
		Preload("PagoIntentos.Mediopagos").
		Where("pagointentos.barcode in ? and pagointentos.mediopagos_id = 6 ", filtro.Barcodes)

	if filtro.ConEstado {
		res.Preload("PagoEstados")
	}

	res.Find(&pagos)

	if res.Error != nil {
		err = fmt.Errorf("no se encontró pago con los datos enviados")
		return
	}

	return
}

func (r *repository) GetEstadoAplicacionRepository(filtro filtros.ConfiguracionFiltro) (configuracion entities.Configuracione, err error) {
	resp := r.SQLClient.Model(entities.Configuracione{})

	resp.Where("nombre = ?", "ESTADO_APLICACION")
	resp.Find(&configuracion)

	if resp.Error != nil {
		err = fmt.Errorf("error al cargar las configuraciones")

		log := entities.Log{
			Tipo:          entities.Error,
			Mensaje:       resp.Error.Error(),
			Funcionalidad: "GetConfiguraciones",
		}

		r.utilService.CreateLogService(log)
	}

	return
}

func (r *repository) GetCuentaByApiKey(apikey string) (cuenta *entities.Cuenta, erro error) {
	resp := r.SQLClient.Model(entities.Cuenta{}).Where("apikey = ?", apikey)
	// resp.Preload("Pagotipos")
	resp.Find(&cuenta)

	if resp.Error != nil {
		logs.Error("error al consultar cuenta: " + resp.Error.Error())
		erro = errors.New("no se pudo realizar la consulta")
		return
	}
	if resp.RowsAffected <= 0 {
		logs.Error("no existe cuenta")
		erro = errors.New("no se pudo realizar la consulta")
		return
	}
	return
}
