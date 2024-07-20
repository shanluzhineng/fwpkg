package app

import (
	"github.com/kataras/iris/v12/core/host"
)

var (
	_registedServe []func(host.TaskHost)
)

func ConfigureHostServe(cb ...func(host.TaskHost)) {
	_registedServe = append(_registedServe, cb...)
}

func hostServe(taskHost host.TaskHost) {
	for _, eachFunc := range _registedServe {
		eachFunc(taskHost)
	}
}
