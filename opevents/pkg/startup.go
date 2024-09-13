package pkg

import "github.com/shanluzhineng/fwpkg/app"

func init() {
	app.RegisterOneStartupAction(registServiceStartup)
}

func registServiceStartup() app.IStartupAction {
	return app.NewStartupAction(func() {
		if app.HostApplication.SystemConfig().App.IsRunInCli {
			return
		}

		registedEventLogServiceList := app.Context.GetListByBaseInterface(new(IOpEventLogService))
		//注册默认的写log文件的
		composedOpEventLogService := newComposedOpEventLogService(newDefaultOpEventLogService())
		for _, eachService := range registedEventLogServiceList {
			currentLogService, ok := eachService.(IOpEventLogService)
			if !ok {
				continue
			}
			composedOpEventLogService.registEventLogService(currentLogService)
		}
		app.Context.RegistInstanceAs(composedOpEventLogService, new(IOpEventLogService))
	})
}
