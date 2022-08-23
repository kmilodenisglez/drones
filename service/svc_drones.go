package service

import (
	"github.com/kataras/iris/v12"
	"github.com/kmilodenisglez/drones.restapi/repo/db"
	"github.com/kmilodenisglez/drones.restapi/schema"
	"github.com/kmilodenisglez/drones.restapi/schema/dto"
	"github.com/tidwall/buntdb"
)

// region ======== SETUP =================================================================

// ISvcDrones Drones request service interface
type ISvcDrones interface {
	// drones functions

	ExistUserSvc(id string)  (bool, *dto.Problem)
	GetUserSvc(id string, filter bool) (*dto.User, *dto.Problem)
	GetUsersSvc() (*[]dto.User, *dto.Problem)
}

type svcDronesTxs struct {
	reposUser *db.RepoDrones
}

// endregion =============================================================================

// NewSvcDronesTxs instantiate the Drones request services
func NewSvcDronesTxs(reposUser *db.RepoDrones) ISvcDrones {
	return &svcDronesTxs{reposUser }
}

// region ======== METHODS ======================================================

func (s *svcDronesTxs) ExistUserSvc(id string) (bool, *dto.Problem) {
	err := (*s.reposUser).Exist(id)
	// Getting non-existent values will cause an ErrNotFound error.
	if err == buntdb.ErrNotFound {
		return false, dto.NewProblem(iris.StatusPreconditionFailed, schema.ErrBuntdbItemNotFound, err.Error())
	} else if err != nil {
		return false, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}

	return true, nil
}

func (s *svcDronesTxs) GetUserSvc(id string, filter bool)  (*dto.User, *dto.Problem) {
	res, err := (*s.reposUser).GetUser(id, filter)
	if err != nil {
		return nil, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return res, nil
}

func (s *svcDronesTxs) GetUsersSvc()  (*[]dto.User, *dto.Problem) {
	res, err := (*s.reposUser).GetUsers()
	if err != nil {
		return nil, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return res, nil
}