package db

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/kmilodenisglez/drones.restapi/lib"
	"github.com/kmilodenisglez/drones.restapi/schema/dto"
	"github.com/kmilodenisglez/drones.restapi/service/utils"
	"github.com/tidwall/buntdb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)


// region ======== SETUP =================================================================

type RepoEventLog interface {
	GetEventLogs() (*[]dto.LogEvent, error)
	CheckBatteryLevelsDrones(drones *[]dto.Drone) error
}

type repoEventLog struct {
	LogDBLocation string
}

// endregion =============================================================================

func NewRepoEventLog(svcConf *utils.SvcConfig) RepoEventLog {
	return &repoEventLog{LogDBLocation: svcConf.LogDBPath}
}

// region ======== METHODS ===============================================================

// GetEventLogs A read-only transaction, return events in db
func (r *repoEventLog) GetEventLogs() (*[]dto.LogEvent, error) {
	db, err := r.loadEventDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	eventLog := dto.LogEvent{}
	eventLogList := make([]dto.LogEvent, 0)
	lastTen := 4
	// custom index: sort drones descending by battery capacity
	db.CreateIndex("log", "event_log:*", buntdb.IndexString)
	err = db.View(func(tx *buntdb.Tx) error {
		err := tx.Descend("log", func(key, value string) bool {
			err = jsoniter.UnmarshalFromString(value, &eventLog)
			if err == nil {
				eventLogList = append(eventLogList, eventLog)
				lastTen--
				// return only the last 4 LogEvent
				if lastTen == 0 {
					return false
				}
			}
			return err == nil
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	return &eventLogList, nil
}

// CheckBatteryLevelsDrones check drones battery levels and create history/audit event log for this
func (r *repoEventLog) CheckBatteryLevelsDrones(drones *[]dto.Drone) error {
	db, err := r.loadEventDB()
	if err != nil {
		return err
	}
	defer db.Close()

	dronesBatteryLevelList := make([]dto.DroneBatteryLevel, 0)
	for _, v := range *drones {
		dronesBatteryLevelList = append(dronesBatteryLevelList, dto.DroneBatteryLevel{
			SerialNumber:    v.SerialNumber,
			BatteryCapacity: v.BatteryCapacity,
		})
	}

	// it is also used as a key for db
	timestamp := timestamppb.Now().AsTime().Format("20060102-150405")
	logEvent := dto.LogEvent{
		Created:             timestamp,
		UUID:                lib.GenerateUUIDStr(),
		DronesBatteryLevels: dronesBatteryLevelList,
	}

	log.Printf("writing event log")
	err = db.Update(func(tx *buntdb.Tx) error {
		res, err := jsoniter.MarshalToString(logEvent)
		if err != nil {
			return err
		}
		_, _, err = tx.Set("event_log:"+timestamp, res, nil)
		return err
	})
	if err != nil {
		return err
	}
	log.Println("successfully added drone")

	return nil
}

// region ======== PRIVATE AUX ===========================================================

func (r *repoEventLog) loadEventDB() (*buntdb.DB, error) {
	log.Println("Load Event Log DB ", r.LogDBLocation)
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open(r.LogDBLocation)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

// endregion =============================================================================
