package cron

import (
	"github.com/go-co-op/gocron"
	"github.com/kataras/iris/v12"
	"github.com/kmilodenisglez/drones.restapi/repo/db"
	"github.com/kmilodenisglez/drones.restapi/schema"
	"github.com/kmilodenisglez/drones.restapi/schema/dto"
	"github.com/kmilodenisglez/drones.restapi/service/utils"
	"log"
	"time"
)

// ISvcEventLog EventLog request service interface
type ISvcEventLog interface {
	GetEventLogs() (*[]dto.LogEvent, *dto.Problem)
	MeinerCronJob() error
}

type svcEventLogReqs struct {
	svcConf       *utils.SvcConfig
	reposEventLog *db.RepoEventLog
	reposDrones   *db.RepoDrones
}

// endregion =============================================================================

// NewSvcRepoEventLog instantiate the Drones request services
func NewSvcRepoEventLog(svcConf *utils.SvcConfig) ISvcEventLog {
	reposEventLog := db.NewRepoEventLog(svcConf)
	reposDrones := db.NewRepoDrones(svcConf)
	return &svcEventLogReqs{svcConf, &reposEventLog, &reposDrones}
}

// GetEventLogs get event log
func (e svcEventLogReqs) GetEventLogs() (*[]dto.LogEvent, *dto.Problem) {
	logs, err := (*e.reposEventLog).GetEventLogs()
	if err != nil {
		return nil, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return logs, nil
}

// MeinerCronJob periodic task to check drones battery levels and create history/audit event log for this
func (e svcEventLogReqs) MeinerCronJob() error {
	// cron job is started only if it is active in configuration
	if e.svcConf.CronEnabled {
		log.Printf("schedules a new periodic Job with an interval: %d seconds", e.svcConf.EveryTime)
		cron := gocron.NewScheduler(time.UTC)

		_, err := cron.Every(e.svcConf.EveryTime).Seconds().WaitForSchedule().Do(e.doFunc)
		if err != nil {
			return err
		}
		// starts the scheduler asynchronously
		cron.StartAsync()
	}
	return nil
}

func (e svcEventLogReqs) doFunc() {
	log.Println("cron job executing")
	// If the drone database has not been populated then the cron is skipped
	isPopulated := (*e.reposDrones).IsPopulated()
	if !isPopulated {
		return
	}
	// drones are requested to populate the event log database
	drones, err := (*e.reposDrones).GetDrones("")
	if err != nil || drones == nil {
		return
	}
	err = (*e.reposEventLog).CheckBatteryLevelsDrones(drones)
	if err != nil {
		return
	}
	log.Println("cron job ending")
}
