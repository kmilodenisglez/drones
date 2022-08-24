package dto

type DroneState uint

const (
	IDLE DroneState = iota
	LOADING
	LOADED
	DELIVERING
	DELIVERED
	RETURNING
)

type DroneModel uint

const (
	Lightweight DroneModel = iota
	Middleweight
	Cruiserweight
	Heavyweight
)

func (droneState DroneState) String() string {
	names := []string{"IDLE", "LOADING", "LOADED", "DELIVERING", "DELIVERED", "RETURNING"}
	if droneState < IDLE || droneState > RETURNING {
		return "unknown"
	}
	return names[droneState]
}
func (droneModel DroneModel) String() string {
	names := []string{"Lightweight", "Middleweight", "Cruiserweight", "Heavyweight"}
	if droneModel < Lightweight || droneModel > Heavyweight {
		return "unknown"
	}
	return names[droneModel]
}

var DroneStateID = map[string]int{"IDLE": 0, "LOADING": 1, "LOADED": 2, "DELIVERING": 3, "DELIVERED": 4, "RETURNING": 5}

type ConfigDB struct {
	IsPopulated bool `json:"isPopulated"`
}

type Drone struct {
	SerialNumber    string     `json:"serialNumber" valid:"maxstringlength(100)"`
	Model           DroneModel `json:"model"`
	WeightLimit     int        `json:"weightLimit" valid:"range(0|500)"`
	BatteryCapacity float64    `json:"batteryCapacity" valid:"range(0|100)"`
	State           DroneState `json:"state"`
}
type Medication struct {
	Name   string  `json:"name" valid:"customnamevalidation"`
	Weight float64 `json:"weight"`
	Code   string  `json:"code" valid:"uppercase,customcodevalidation"` // we assume that the code is unique
	Image  string  `json:"image" valid:"base64"`
}
