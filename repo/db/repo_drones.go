package db

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/kmilodenisglez/drones.restapi/schema/dto"
	"github.com/kmilodenisglez/drones.restapi/service/utils"
	"github.com/tidwall/buntdb"
	"strings"
)

// region ======== SETUP =================================================================

type RepoDrones interface {
	GetUser(field string, filterOptional ...bool) (*dto.User, error)
	GetUsers() (*[]dto.User, error)
	Exist(id string) error
}

type repoDrones struct {
	DBUserLocation string
}

// endregion =============================================================================

func NewRepoDrones(svcConf *utils.SvcConfig) RepoDrones {
	return &repoDrones{DBUserLocation: svcConf.DbPath}
}

// region ======== METHODS ===============================================================

// GetUser get the user from the DB file that should be compliant with the dto.UserList struct
// return a list of dto.User
func (r *repoDrones) GetUser(field string, filterOptional ...bool) (*dto.User, error) {
	filter := false
	if len(filterOptional) > 0 {
		filter = filterOptional[0]
	}
	user := dto.User{}

	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open(r.DBUserLocation)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.CreateIndex("username", "*", buntdb.IndexString)
	if err != nil {
		return nil, err
	}
	err = db.View(func(tx *buntdb.Tx) error {
		if filter {
			err := tx.Ascend("username", func(key, value string) bool {
				if strings.Contains(value, field) {
					err := jsoniter.UnmarshalFromString(value, &user)
					if err != nil {
						return false
					}
					return false
				}

				return true
			})
			return err
		}
		// filter = false
		value, err := tx.Get(field)
		if err != nil {
			return err
		}
		err = jsoniter.UnmarshalFromString(value, &user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repoDrones) GetUsers() (*[]dto.User, error) {
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open(r.DBUserLocation)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	user := dto.User{}
	var list []dto.User

	err = db.CreateIndex("username", "*", buntdb.IndexString)
	if err != nil {
		return nil, err
	}
	err = db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("username", func(key, value string) bool {
			err = jsoniter.UnmarshalFromString(value, &user)
			if err == nil {
				list = append(list, user)
			}
			return err == nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (r *repoDrones) Exist(id string) error {
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open(r.DBUserLocation)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(id)
		if err != nil {
			return err
		}
		return nil
	})
	// Getting non-existent values will cause an ErrNotFound error.
	if err != nil {
		return err
	}

	return nil
}

// endregion =============================================================================

// region ======== PRIVATE AUX ===========================================================
// endregion =============================================================================
