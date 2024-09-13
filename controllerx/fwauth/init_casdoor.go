package fwauth

import (
	"fmt"
	"sync"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/shanluzhineng/configurationx"
	optCasdoor "github.com/shanluzhineng/configurationx/options/casdoor"
	"github.com/shanluzhineng/fwpkg/system/log"
)

var (
	_casdoorOptions CasdoorOptions
	_cm             *Middleware
	_sync           sync.Once
)

// 根据configurationx中的配置解析参数，并初始化casdoor，返回iris格式的ctx
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

// 根据fwpkg框架内的配置解析参数，并初始化casdoor
func InitCasdoor() error {
	options := &CasdoorOptions{}
	if ok := configurationx.GetInstance().UnmarshalPropertiesTo("casdoor", options); !ok {
		log.Logger.Warn("casdoor params parser fail")
		return fmt.Errorf("casdoor params parser fail")
	}
	options.Normalize()
	log.Logger.Info(fmt.Sprintf("[initCasdoor]>>> url: %s, clientId: %s", options.Endpoint, options.ClientId))
	if !options.Disabled {
		casdoorsdk.InitConfig(options.Endpoint,
			options.ClientId,
			options.ClientSecret,
			options.Certificate,
			options.OrganizationName,
			options.ApplicationName)
	} else {
		log.Logger.Warn(fmt.Sprintf("oauth配置已被禁用, options: %+v", options))
	}
	return nil
}
