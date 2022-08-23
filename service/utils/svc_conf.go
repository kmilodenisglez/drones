package utils

import (
	"github.com/kmilodenisglez/drones.restapi/lib"
	"github.com/kmilodenisglez/drones.restapi/schema"
	"github.com/tkanos/gonfig"
)

// region ======== TYPES =================================================================

// conf unexported configuration schema holder struct
type conf struct {
	// Environment
	Debug    bool
	ApiDocIp string
	DappPort string

	// Cryptographic conf
	JWTSignKey string
	TkMaxAge   uint8

	// Store
	DbPath     string
}

// SvcConfig exported configuration service struct
type SvcConfig struct {
	Path string `string:"Path to the config YAML file"`
	conf `conf:"Configuration object"`
}

// endregion =============================================================================

// NewSvcConfig create a new configuration service.
func NewSvcConfig() *SvcConfig {
	c := conf{}

	var configPath = lib.GetEnvOrError(schema.EnvConfigPath)
	var jwtSignKey = lib.GetEnvOrError(schema.EnvJWTSignKey)

	err := gonfig.GetConf(configPath, &c) // getting the conf
	if err != nil {
		panic(err)
	} // error check

	c.JWTSignKey = jwtSignKey // saving the sign key into the configuration object

	return &SvcConfig{configPath, c} // We are using struct composition here. Hence, the anonymous field (https://golangbot.com/inheritance/)
}
