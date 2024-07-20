package controllerx

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	requestLogger "github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/shanluzhineng/fwpkg/app"
	"github.com/shanluzhineng/fwpkg/app/host"
	"github.com/shanluzhineng/fwpkg/app/web"
	"github.com/shanluzhineng/fwpkg/system/log"

	"github.com/shanluzhineng/configurationx"
	"github.com/shanluzhineng/fwpkg/controllerx/middleware/cors"
	errHandler "github.com/shanluzhineng/fwpkg/controllerx/middleware/err"

	"net/http/pprof"

	requestPprof "github.com/kataras/iris/v12/middleware/pprof"
)

func init() {
	app.Register(NewIrisApplication)
}

func requestLogConfig() requestLogger.Config {
	c := requestLogger.DefaultConfig()
	c.AddSkipper(func(ctx *context.Context) bool {
		p := ctx.Path()
		return strings.HasPrefix(p, "/api/health/check")
	})
	return c
}

type IrisApplication struct {
	*iris.Application
	Address string

	isBuilded        bool
	irisConfigurator []iris.Configurator
	Err              error
}

type Configurator func(*IrisApplication)

func NewIrisApplication() *IrisApplication {
	irisNew := iris.New()
	//错误封装
	irisNew.Use(errHandler.New())
	irisNew.Use(recover.New())
	irisNew.Use(requestLogger.New(requestLogConfig()))
	if configurationx.GetInstance().Web != nil {
		cors.UseCors(irisNew.APIBuilder, configurationx.GetInstance().Web.Cors)
	}
	//设置validator
	irisNew.Validator = validator.New()
	irisApp := &IrisApplication{
		Application:      irisNew,
		irisConfigurator: make([]iris.Configurator, 0),
		isBuilded:        false,
	}
	return irisApp
}

func (a *IrisApplication) Configure(configurators ...Configurator) *IrisApplication {
	return a
}

// build IrisApplication environments
func (a *IrisApplication) Build(configurators ...Configurator) *IrisApplication {
	if a.isBuilded {
		return a
	}
	if a.Err != nil {
		return a
	}
	defer func() {
		a.isBuilded = true
	}()
	envHttp := host.GetHostEnvironment().GetEnvString(host.ENV_HTTP)
	if len(envHttp) > 0 {
		a.Address = envHttp
	} else {
		host.GetHostEnvironment().SetHttp(a.Address)
	}
	if len(a.Address) <= 0 {
		msg := "没有配置好app.http参数"
		log.Error(msg)
		panic(msg)
	}

	//配置web应用中间件
	web.SetWebApplication(web.NewWebApplication())
	web.Application.ConfigureService()

	a.pprofStartupAction()
	//运行启动项
	app.HostApplication.RunStartup()

	//构建配置
	appConfigurators := make([]iris.Configurator, 0)
	for _, eachConfigurator := range configurators {
		if eachConfigurator == nil {
			continue
		}
		newAppConfigurator := func(irisApp *iris.Application) {
			eachConfigurator(a)
		}
		appConfigurators = append(appConfigurators, newAppConfigurator)
	}
	a.irisConfigurator = appConfigurators

	//设置启动消耗的时间
	startTime := host.GetHostEnvironment().GetEnv(host.ENV_StartTime).(time.Time)
	interval := time.Since(startTime)
	host.GetHostEnvironment().SetEnv(host.ENV_StartInterval, interval)

	return a
}

func (a *IrisApplication) Run(configurators ...Configurator) *IrisApplication {
	a.Build(configurators...)

	err := a.Application.Run(iris.Addr(a.Address), a.irisConfigurator...)
	a.Err = err
	return a
}

func (a *IrisApplication) pprofStartupAction() {
	if app.HostApplication.SystemConfig().App.IsRunInCli {
		return
	}

	log.Logger.Debug("正在构建pprof路径组件,/debug/pprof...")
	a.Any("/debug/pprof/cmdline", iris.FromStd(pprof.Cmdline))
	a.Any("/debug/pprof/profile", iris.FromStd(pprof.Profile))
	a.Any("/debug/pprof/symbol", iris.FromStd(pprof.Symbol))
	a.Any("/debug/pprof/trace", iris.FromStd(pprof.Trace))
	a.Any("/debug/pprof/debug/pprof/{action:string}", requestPprof.New())

	httpValue := os.Getenv("app.http")
	advertiseHostValue := os.Getenv("app.advertisehost")
	if len(httpValue) > 0 {
		pprofPath := httpValue
		if len(advertiseHostValue) > 0 {
			pprofPath = strings.Replace(httpValue, "0.0.0.0", advertiseHostValue, 1)
		}
		log.Logger.Debug(fmt.Sprintf("已经构建好pprof路径组件,你可以通过 %s/debug/pprof 来访问pprof", pprofPath))
	}
}
