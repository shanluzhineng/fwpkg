package esconnector

import (
	"github.com/shanluzhineng/fwpkg/app"

	"github.com/shanluzhineng/fwpkg/opevents/esconnector"
	opevent "github.com/shanluzhineng/fwpkg/opevents/pkg"
)

func init() {
	//初始化es
	app.RegisterOneStartupAction(initElasticsearchStartup)
}

func initElasticsearchStartup() app.IStartupAction {
	return app.NewStartupAction(func() {
		if app.HostApplication.SystemConfig().App.IsRunInCli {
			return
		}
		esconnector.InitElasticsearch()
		app.Context.RegistInstanceAs(esconnector.NewOpEventLogESConnectService(), new(opevent.IOpEventLogStoreService))
	})
}
