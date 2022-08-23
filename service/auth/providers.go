package auth

import (
	"github.com/kmilodenisglez/drones.restapi/lib"
	"github.com/kmilodenisglez/drones.restapi/repo/db"
	"github.com/kmilodenisglez/drones.restapi/schema"
	"github.com/kmilodenisglez/drones.restapi/schema/dto"
	"github.com/kataras/iris/v12"
)

type Provider interface {
	GrantIntent(userCredential *dto.UserCredIn, data interface{}) (*dto.GrantIntentResponse, *dto.Problem)
}

// region ======== EVOTE AUTHENTICATION PROVIDER =========================================

type ProviderDrone struct {
	// walletLocations string
	repo *db.RepoDrones
}

func (p *ProviderDrone) GrantIntent(uCred *dto.UserCredIn, options interface{}) (*dto.GrantIntentResponse, *dto.Problem) {
	// getting the users
	user, err := (*p.repo).GetUser(uCred.Username, true)
	if err != nil {
		return nil, dto.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	checksum, _ := lib.Checksum("SHA256", []byte(uCred.Password))
	if user.Passphrase == checksum {
		return &dto.GrantIntentResponse{Identifier: user.Username, DID: user.Username}, nil
	}

	return nil, dto.NewProblem(iris.StatusNotFound, schema.ErrFile, schema.ErrCredsNotFound)
}

// endregion =============================================================================
