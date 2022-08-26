package endpoints

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kmilodenisglez/drones.restapi/service/cron"
	"github.com/kmilodenisglez/drones.restapi/service/utils"
)

// EventLogHandler  endpoint handler struct for EventLog
type EventLogHandler struct {
	response *utils.SvcResponse
	service  *cron.ISvcEventLog
}

// NewEventLogHandler create and register the handler for EventLog
//
// - app [*iris.Application] ~ Iris App instance
//
// - MdwAuthChecker [*context.Handler] ~ Authentication checker middleware
//
// - svcR [*utils.SvcResponse] ~ GrantIntentResponse service instance
//
// - svcC [utils.SvcConfig] ~ Configuration service instance
func NewEventLogHandler(app *iris.Application, mdwAuthChecker *context.Handler, svcR *utils.SvcResponse, svcC *utils.SvcConfig) EventLogHandler { // --- VARS SETUP ---
	svc := cron.NewSvcRepoEventLog(svcC)
	// registering protected / guarded router
	h := EventLogHandler{svcR, &svc}

	// Simple group: v1
	v1 := app.Party("/api/v1")
	{
		// registering protected / guarded router
		guardTxsRouter := v1.Party("/logs")
		{
			// --- GROUP / PARTY MIDDLEWARES ---
			guardTxsRouter.Use(*mdwAuthChecker)
			guardTxsRouter.Get("/", h.GetEventLog)
		}
	}
	return h
}

// GetEventLog get event logs
// @Summary Get event logs
// @description.markdown GetEventLogDescription
// @Tags logs
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param	Authorization	header	string	true 	"Insert access token" default(Bearer <Add access token here>)
// @Success 200 {object} []dto.LogEvent "OK"
// @Failure 400 {object} dto.Problem "err.processing_param"
// @Failure 500 {object} dto.Problem "err.database_related"
// @Failure 504 {object} dto.Problem "err.network"
// @Router /logs [get]
func (h EventLogHandler) GetEventLog(ctx iris.Context) {
	logs, problem := (*h.service).GetEventLogs()
	if problem != nil {
		h.response.ResErr(problem, &ctx)
		return
	}
	h.response.ResOKWithData(logs, &ctx)
}

// region ======== LOCAL DEPENDENCIES ====================================================
