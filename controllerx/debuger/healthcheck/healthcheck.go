package healthcheck

import (
	"strings"

	"github.com/shanluzhineng/fwpkg/app"
	"github.com/shanluzhineng/fwpkg/app/host"
	"github.com/shanluzhineng/fwpkg/controllerx"
	"github.com/shanluzhineng/fwpkg/controllerx/responsex"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/fwpkg/utils/str"

	"github.com/kataras/iris/v12"
)

func healthcheckStartup(webApp *controllerx.IrisApplication) app.IStartupAction {
	return app.NewStartupAction(func() {
		if app.HostApplication.SystemConfig().App.IsRunInCli {
			return
		}
		log.Logger.Debug("正在构建healthcheck路径组件,api/health/check...")
		healthRouterParty := webApp.Party("/api/health")
		{
			healthRouterParty.Get("/check", healthcheck)
		}

		healthcheck := host.GetHostEnvironment().GetEnvString(host.ENV_Healthcheck)
		if len(healthcheck) <= 0 {
			http := host.GetHostEnvironment().GetEnvString(host.ENV_HTTP)
			if len(http) > 0 {
				//设置健康检查地址
				url := strings.Join([]string{"http://", str.EnsureEndWith(http, "/"), "api/health/check"}, "")
				host.GetHostEnvironment().SetEnv(host.ENV_Healthcheck, url)
			}
		}
	})
}

func healthcheck(ctx iris.Context) {
	response := responsex.NewSuccessResponse(func(br *responsex.BaseResponse) {
		br.SetMessage("Hi,I am a OK ,and I am running")

		envValue := make(map[string]interface{})
		envKeyList := host.GetHostEnvironment().AllKey()
		for _, eachKey := range envKeyList {
			if !strings.HasPrefix(eachKey, "app.") {
				continue
			}
			val := host.GetHostEnvironment().GetEnv(eachKey)
			if val == nil {
				continue
			}
			envValue[eachKey] = val
		}
		br.SetData(envValue)
	})
	ctx.JSON(response)
}
