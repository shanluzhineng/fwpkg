package consul

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"sync"

	"github.com/shanluzhineng/configurationx"
	"github.com/shanluzhineng/fwpkg/app"
	"github.com/shanluzhineng/fwpkg/app/host"
	"github.com/shanluzhineng/fwpkg/app/web"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/fwpkg/utils/str"
)

var (
	_normalizeOnce sync.Once
)

func init() {
	web.ConfigureService(serviceConfigurator)
	// app.RegisterOneStartupAction(registRegistryAction)
}

func serviceConfigurator(wa web.WebApplication) {
	if app.HostApplication.SystemConfig().App.IsRunInCli {
		return
	}
	normalizeConsulOption()
	consul := configurationx.GetInstance().Consul
	registry := NewRegistry(WithAddress(fmt.Sprintf("%s:%d", consul.Host, consul.Port)),
		WithEnableHealthCheck(*consul.Registration.EnableHealthCheck),
		WithHealthCheckInterval(*consul.Registration.HealthCheckInterval),
		WithHealthCheckTimeout(*consul.Registration.HealthCheckTimeout),
		WithEnableHeartbeatCheck(*consul.Registration.EnableHeartbeatCheck),
		WithHeartbeatCheckInterval(*consul.Registration.HeartbeatCheckInterval),
		WithDeregisterCriticalServiceAfter(*consul.Registration.DeregisterCriticalServiceAfter),
	)
	if registry == nil {
		err := errors.New("无法创建registry对象")
		log.Logger.Error(err.Error())
		return
	}
	//注册registry对象
	app.Context.SetInstance(registry)
}

func normalizeConsulOption() {
	_normalizeOnce.Do(func() {
		consul := configurationx.GetInstance().Consul
		if len(consul.Registration.HealthCheckHTTP) <= 0 {
			consul.Registration.HealthCheckHTTP = host.GetHostEnvironment().GetEnvString(host.ENV_Healthcheck)
		}
		if len(consul.Registration.ServiceName) <= 0 {
			consul.Registration.ServiceName = host.GetHostEnvironment().GetEnvString(host.ENV_AppName)
		}
		if len(consul.Registration.Product) <= 0 {
			consul.Registration.Product = host.GetHostEnvironment().GetEnvString(host.ENV_Product)
			if len(consul.Registration.Product) <= 0 {
				consul.Registration.Product = consul.Registration.ServiceName
			}
		}
		serviceAddressUrl, err := consul.Registration.ParseServiceAddressForScheme("http")
		if err != nil {
			log.Logger.Error("consul.registration.endpoint配置错误")
			panic(err)
		}
		if serviceAddressUrl == nil {
			http := host.GetHostEnvironment().GetEnvString(host.ENV_HTTP)
			httpUrl := str.EnsureStartWith(http, "http://")
			advertiseHost := host.GetHostEnvironment().GetEnvString(host.ENV_AdvertiseHost)
			if len(advertiseHost) > 0 {
				url, err := url.Parse(httpUrl)
				if err == nil {
					_, p, err := net.SplitHostPort(url.Host)
					if err == nil {
						httpUrl = "http://" + advertiseHost + ":" + p
					}
				}
			}
			consul.Registration.Endpoint = append(consul.Registration.Endpoint, httpUrl)
		}
	})
}
