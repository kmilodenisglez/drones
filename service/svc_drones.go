package service

import (
	"github.com/kataras/iris/v12"
	"github.com/kmilodenisglez/drones.restapi/repo/db"
	"github.com/kmilodenisglez/drones.restapi/schema"
	"github.com/kmilodenisglez/drones.restapi/schema/dto"
	"github.com/tidwall/buntdb"
	"log"
)

// region ======== SETUP =================================================================

// ISvcDrones Drones request service interface
type ISvcDrones interface {
	IsPopulateDBSvc() bool
	PopulateDBSvc() *dto.Problem

	// user functions

	ExistUserSvc(id string) (bool, *dto.Problem)
	GetUserSvc(id string, filter bool) (*dto.User, *dto.Problem)
	GetUsersSvc() (*[]dto.User, *dto.Problem)

	// drone functions

	GetADroneSvc(serialNumber string) (*dto.Drone, *dto.Problem)
	GetDronesSvc(filters ...string) (*[]dto.Drone, *dto.Problem)
	RegisterDroneSvc(drone *dto.Drone) *dto.Problem

	// medication functions

	GetMedications() (*[]dto.Medication, *dto.Problem)
}

type svcDronesReqs struct {
	reposDrones *db.RepoDrones
}

// endregion =============================================================================

// NewSvcDronesReqs instantiate the Drones request services
func NewSvcDronesReqs(reposDrones *db.RepoDrones) ISvcDrones {
	return &svcDronesReqs{reposDrones}
}

// region ======== METHODS ======================================================

func (s *svcDronesReqs) IsPopulateDBSvc() bool {
	return (*s.reposDrones).IsPopulated()
}

func (s *svcDronesReqs) PopulateDBSvc() *dto.Problem {
	log.Println("1")
	err := (*s.reposDrones).PopulateDB()

	switch {
	case err == buntdb.ErrNotFound:
		return dto.NewProblem(iris.StatusPreconditionFailed, schema.ErrBuntdbItemNotFound, err.Error())
	case err != nil:
		if err.Error() == schema.ErrBuntdbPopulated {
			return dto.NewProblem(iris.StatusInternalServerError, schema.ErrBuntdbPopulated, "the database has already been populated")
		}
		return dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return nil
}

func (s *svcDronesReqs) ExistUserSvc(id string) (bool, *dto.Problem) {
	err := (*s.reposDrones).Exist(id)
	// Getting non-existent values will cause an ErrNotFound error.
	if err == buntdb.ErrNotFound {
		return false, dto.NewProblem(iris.StatusPreconditionFailed, schema.ErrBuntdbItemNotFound, err.Error())
	} else if err != nil {
		return false, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}

	return true, nil
}

func (s *svcDronesReqs) GetUserSvc(id string, filter bool) (*dto.User, *dto.Problem) {
	res, err := (*s.reposDrones).GetUser(id, filter)
	if err != nil {
		return nil, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return res, nil
}

func (s *svcDronesReqs) GetUsersSvc() (*[]dto.User, *dto.Problem) {
	res, err := (*s.reposDrones).GetUsers()
	if err != nil {
		return nil, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return res, nil
}

// GetADroneSvc get a specific drone
func (s *svcDronesReqs) GetADroneSvc(serialNumber string) (*dto.Drone, *dto.Problem) {
	res, err := (*s.reposDrones).GetDrone(serialNumber)
	// Getting non-existent values will cause an ErrNotFound error.
	if err == buntdb.ErrNotFound {
		return nil, dto.NewProblem(iris.StatusPreconditionFailed, schema.ErrBuntdbItemNotFound, err.Error())
	} else if err != nil {
		return nil, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}

	return res, nil
}

func (s *svcDronesReqs) GetDronesSvc(filters ...string) (*[]dto.Drone, *dto.Problem) {
	var filter = ""
	if len(filters) > 0 {
		filter = filters[0]
	}

	res, err := (*s.reposDrones).GetDrones(filter)
	if err != nil {
		return nil, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}

	return res, nil
}

func (s *svcDronesReqs) RegisterDroneSvc(drone *dto.Drone) *dto.Problem {
	err := (*s.reposDrones).RegisterDrone(drone)
	if err != nil {
		return dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return nil
}

func (s *svcDronesReqs) GetMedications() (*[]dto.Medication, *dto.Problem) {
	res, err := (*s.reposDrones).GetMedications()
	if err != nil {
		return nil, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return res, nil
}
