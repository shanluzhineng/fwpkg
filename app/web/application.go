package web

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/shanluzhineng/fwpkg/app"
	"github.com/shanluzhineng/fwpkg/system/log"
)

// 配置服务
type ServiceConfigurator func(WebApplication)

var (
	_registedConfiguratorList []ServiceConfigurator
	_syncOnce                 sync.Once
	Application               WebApplication
)

// web应用
type WebApplication interface {
	GetServiceProvider() app.IServiceProvider
	ConfigureService()
}

type defaultWebApplication struct {
	serviceProvider app.IServiceProvider
}

// new一个web应用
func NewWebApplication() WebApplication {
	newApp := &defaultWebApplication{}
	if app.HostApplication != nil {
		newApp.serviceProvider = app.HostApplication.GetServiceProvider()
	}
	return newApp
}

// 设置Application属性值
func SetWebApplication(webApp WebApplication) {
	Application = webApp
}

func (a *defaultWebApplication) GetServiceProvider() app.IServiceProvider {
	return a.serviceProvider
}

func (a *defaultWebApplication) ConfigureService() {
	_syncOnce.Do(func() {
		for _, eachOption := range _registedConfiguratorList {
			configuratorName := getServiceConfiguratorTypeName(eachOption)
			if !app.HostApplication.SystemConfig().App.IsRunInCli {
				log.Logger.Info(fmt.Sprintf("begin run ServiceConfigurator,%s", configuratorName))
			}
			eachOption(a)
			if !app.HostApplication.SystemConfig().App.IsRunInCli {
				log.Logger.Info(fmt.Sprintf("finish run ServiceConfigurator,%s", configuratorName))
			}
		}
	})
}

func getServiceConfiguratorTypeName(configuratorFunc ServiceConfigurator) string {
	if configuratorFunc == nil {
		return ""
	}
	return runtime.FuncForPC(reflect.ValueOf(configuratorFunc).Pointer()).Name()
}

// 配置服务
func ConfigureService(opts ...ServiceConfigurator) {
	_registedConfiguratorList = append(_registedConfiguratorList, opts...)
}
