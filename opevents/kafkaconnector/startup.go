package kafkaconnector

import (
	"github.com/shanluzhineng/fwpkg/app"
	"github.com/shanluzhineng/fwpkg/app/web"
)

func init() {
	web.ConfigureService(func(wa web.WebApplication) {
		if app.HostApplication.SystemConfig().App.IsRunInCli {
			return
		}

		app.Context.RegistInstance(newOpEventLogService())
	})
}
