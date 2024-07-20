package app

import "github.com/shanluzhineng/fwpkg/system/factory"

// 用于shutdown处理
type IShutdownAction interface {
	Run()
}

type shutdownAction struct {
	factory    factory.InstantiateFactory
	subscribes []IShutdownAction
}

func newShutdown(factory factory.InstantiateFactory) *shutdownAction {
	return &shutdownAction{
		factory: factory,
	}
}

var (
	_shutdownActions []interface{}
)

// 注册一个
func RegisterOneShutdown(p interface{}) {
	_shutdownActions = append(_shutdownActions, p)
}

// 注册一组
func RegisterShutdown(p ...interface{}) {
	_shutdownActions = append(_shutdownActions, p...)
}

// 初始化
func (p *shutdownAction) Init() {
	for _, eachShutdownAction := range _shutdownActions {
		ss, err := p.factory.InjectIntoFunc(nil, eachShutdownAction)
		if err == nil {
			p.subscribes = append(p.subscribes, ss.(IShutdownAction))
		}
	}
}

// shutdown
func (p *shutdownAction) Shutdown() {
	for _, eachShutdownAction := range p.subscribes {
		p.factory.InjectIntoFunc(nil, eachShutdownAction)
		eachShutdownAction.Run()
	}
}

// 使用函数来实现IShutdownAction
type shutdownActionFunc struct {
	runFunc func()
}

func (p *shutdownActionFunc) Run() {
	if p.runFunc == nil {
		return
	}
	p.runFunc()
}

// 使用函数创建一个IShutdownAction对象
func NewShutdownAction(runFunc func()) IShutdownAction {
	return &shutdownActionFunc{
		runFunc: runFunc,
	}
}
