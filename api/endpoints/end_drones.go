package endpoints

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/hero"
	"github.com/kmilodenisglez/drones.restapi/repo/db"
	"github.com/kmilodenisglez/drones.restapi/schema"
	"github.com/kmilodenisglez/drones.restapi/schema/dto"
	"github.com/kmilodenisglez/drones.restapi/service"
	"github.com/kmilodenisglez/drones.restapi/service/utils"
)

// DronesHandler  endpoint handler struct for Drones
type DronesHandler struct {
	response *utils.SvcResponse
	service  *service.ISvcDrones
}

// NewDronesHandler create and register the handler for Drones
//
// - app [*iris.Application] ~ Iris App instance
//
// - MdwAuthChecker [*context.Handler] ~ Authentication checker middleware
//
// - svcR [*utils.SvcResponse] ~ GrantIntentResponse service instance
//
// - svcC [utils.SvcConfig] ~ Configuration service instance
func NewDronesHandler(app *iris.Application, mdwAuthChecker *context.Handler, svcR *utils.SvcResponse, svcC *utils.SvcConfig) DronesHandler { // --- VARS SETUP ---
	repoDrones := db.NewRepoDrones(svcC)
	svc := service.NewSvcDronesReqs(&repoDrones)
	// registering protected / guarded router
	h := DronesHandler{svcR, &svc}

	// registering unprotected router
	authRouter := app.Party("/database") // unauthorized
	{
		authRouter.Post("/populate", h.PopulateDB)
	}

	// registering protected / guarded router
	guardTxsDatabase := app.Party("/database")
	{
		// --- GROUP / PARTY MIDDLEWARES ---
		guardTxsDatabase.Use(*mdwAuthChecker)

		// --- DEPENDENCIES ---
		hero.Register(DepObtainUserDid)
	}

	// registering protected / guarded router
	guardTxsRouter := app.Party("/drones")
	{
		// --- GROUP / PARTY MIDDLEWARES ---
		guardTxsRouter.Use(*mdwAuthChecker)

		guardTxsRouter.Get("/", h.GetDrones)
		guardTxsRouter.Get("/{serialNumber:string}", h.GetADrone)
		guardTxsRouter.Post("/", h.RegisterADrone)

		// --- DEPENDENCIES ---
		hero.Register(DepObtainUserDid)
	}

	// registering protected / guarded router
	guardMedicationsRouter := app.Party("/medications")
	{
		// --- GROUP / PARTY MIDDLEWARES ---
		guardMedicationsRouter.Use(*mdwAuthChecker)

		guardMedicationsRouter.Get("/", h.GetMedications)

		// --- DEPENDENCIES ---
		hero.Register(DepObtainUserDid)
	}

	return h
}

// PopulateDB
// @Summary Populate the database with fake data
// @description.markdown PopulateDbDescription
// @Tags database
// @Accept  json
// @Produce json
// @Success 204 "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 500 {object} dto.Problem "err.database_related"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /database/populate [post]
func (h DronesHandler) PopulateDB(ctx iris.Context) {
	problem := (*h.service).PopulateDBSvc()
	if problem != nil {
		h.response.ResErr(problem, &ctx)
		return
	}
	h.response.ResOK(&ctx)
}

// GetDrones get drones
// @Summary Get drones
// @description.markdown GetDronesDescription
// @Tags drones
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	true 	"Insert access token" default(Bearer <Add access token here>)
// @Param   state           query   int     false   "drone state"         Enums(0, 1, 2, 3, 4, 5)
// @Success 200 {object} []dto.Drone "OK"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 500 {object} dto.Problem "err.database_related"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /drones [get]
func (h DronesHandler) GetDrones(ctx iris.Context) {
	qState, err := ctx.URLParamInt("state")
	if err != nil && err != iris.ErrNotFound {
		h.response.ResErr(dto.NewProblem(iris.StatusInternalServerError, schema.ErrParamURL, err.Error()), &ctx)
		return
	}

	var state = ""
	// if no query parameter is passed then we show all drones
	if qState != -1 {
		state = fmt.Sprintf("\"state\":%d", qState)
	}
	drones, problem := (*h.service).GetDronesSvc(state)
	if problem != nil {
		h.response.ResErr(problem, &ctx)
		return
	}
	h.response.ResOKWithData(drones, &ctx)
}

// GetADrone get a drone
// @Summary Get a drone by serialNumber
// @description.markdown GetADroneDescription
// @Tags drones
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	true 	"Insert access token"          default(Bearer <Add access token here>)
// @Param   serialNumber    path    string  true    "Serial number of a drone"     Format(string)
// @Success 200 {object} dto.Drone "OK"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 500 {object} dto.Problem "err.database_related"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /drones/{serialNumber} [get]
func (h DronesHandler) GetADrone(ctx iris.Context) {
	// checking the serialNumber param
	serialNumber := ctx.Params().GetString("serialNumber")
	if serialNumber == "" {
		h.response.ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: schema.ErrDetInvalidField}, &ctx)
		return
	}
	drone, problem := (*h.service).GetADroneSvc(serialNumber)
	if problem != nil {
		h.response.ResErr(problem, &ctx)
		return
	}
	h.response.ResOKWithData(drone, &ctx)
}

// RegisterADrone registers a new drone
// @Summary Registers a new drone, also updates a previously inserted drone
// @description.markdown RegisterADroneDescription
// @Tags drones
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string 			true 	"Insert access token" default(Bearer <Add access token here>)
// @Param	tx				body	dto.Drone		true	"Drone data"
// @Success 204 "OK"
// @Failure 401 {object} dto.Problem "err.unauthorized"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 500 {object} dto.Problem "err.database_related"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /drones [post]
func (h DronesHandler) RegisterADrone(ctx iris.Context) {
	drone := new(dto.Drone)

	// unmarshalling the json and check
	if err := ctx.ReadJSON(drone); err != nil {
		h.response.ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrProcParam, Detail: err.Error()}, &ctx)
		return
	}

	// validate drone fields
	_, err := govalidator.ValidateStruct(drone)
	if err != nil {
		h.response.ResErr(&dto.Problem{Status: iris.StatusBadRequest, Title: schema.ErrValidationField, Detail: err.Error()}, &ctx)
		return
	}

	problem := (*h.service).RegisterDroneSvc(drone)
	if problem != nil {
		h.response.ResErr(problem, &ctx)
		return
	}
	h.response.ResOK(&ctx)
}

// endregion =============================================================================

// region ======== Medications ======================================================

// GetMedications get medications
// @Summary Get medications
// @description.markdown GetMedicationsDescription
// @Tags medications
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	true 	"Insert access token" default(Bearer <Add access token here>)
// @Success 200 {object} []dto.Medication "OK"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 500 {object} dto.Problem "err.database_related"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /medications [get]
func (h DronesHandler) GetMedications(ctx iris.Context) {
	medications, problem := (*h.service).GetMedications()
	if problem != nil {
		h.response.ResErr(problem, &ctx)
		return
	}
	h.response.ResOKWithData(medications, &ctx)
}

// endregion ======== Medications ======================================================

// region ======== LOCAL DEPENDENCIES ====================================================

// DepObtainUserDid this tries to get the user DID store in the previously generated auth Bearer token.
func DepObtainUserDid(ctx iris.Context) dto.InjectedParam {
	tkData := ctx.Values().Get("iris.jwt.claims").(*dto.AccessTokenData)

	// returning the DID and Identifier (Username)
	return tkData.Claims
}

// endregion =============================================================================
