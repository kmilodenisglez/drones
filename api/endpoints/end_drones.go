package endpoints

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/hero"
	"github.com/kmilodenisglez/drones.restapi/repo/db"
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
	svc := service.NewSvcDronesTxs(&repoDrones)
	// registering protected / guarded router
	h := DronesHandler{svcR, &svc}

	// registering unprotected router
	authRouter := app.Party("/drones") // unauthorized
	{
		authRouter.Get("/get", h.Get)
	}

	// registering protected / guarded router
	guardTxsRouter := app.Party("/drones")
	{
		// --- GROUP / PARTY MIDDLEWARES ---
		guardTxsRouter.Use(*mdwAuthChecker)

		// --- DEPENDENCIES ---
		hero.Register(DepObtainUserDid)
	}

	return h
}

// Get test get
// @Tags Txs.drones
// @Accept  json
// @Produce json
// @Success 200 {object} []dto.User "OK"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 502 {object} dto.Problem "err.bad_gateway"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /drones/get [get]
func (h DronesHandler) Get(ctx iris.Context) {
	//users, problem := (*h.service).GetUsersRegisterSvc()
	//if problem != nil {
	//	(*h.response).ResErr(problem, &ctx)
	//	return
	//}

	h.response.ResOKWithData("users", &ctx)
}


// endregion =============================================================================


// region ======== LOCAL DEPENDENCIES ====================================================

// DepObtainUserDid this tries to get the user DID store in the previously generated auth Bearer token.
func DepObtainUserDid(ctx iris.Context) dto.InjectedParam {
	tkData := ctx.Values().Get("iris.jwt.claims").(*dto.AccessTokenData)

	// returning the DID and Identifier (Username)
	return tkData.Claims
}

// endregion =============================================================================
