package fwauth

import (
	"sync"

	"github.com/shanluzhineng/configurationx"
	optCasdoor "github.com/shanluzhineng/configurationx/options/casdoor"
)

var (
	_casdoorOptions CasdoorOptions
	_cm             *Middleware
	_sync           sync.Once
)

func GetCasdoorMiddleware() *Middleware {
	_sync.Do(func() {
		casdoorOpt := &optCasdoor.CasdoorOptions{}
		configurationx.GetInstance().UnmarshalPropertiesTo(optCasdoor.ConfigurationKey, casdoorOpt)
		_casdoorOptions = CasdoorOptions{
			CasdoorOptions: *casdoorOpt,
			Extractor: FromFirst(FromHeader("Authorization"),
				FromAuthHeader),
		}
		_cm = New(_casdoorOptions)
	})
	return _cm
}
