package app

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/shanluzhineng/fwpkg/system/factory"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/fwpkg/system/reflector"
	"github.com/shanluzhineng/fwpkg/utils/slicex"
)

const (
	//执行在最后的优先维
	lastPriority = 9999
)

// 用于post处理
type IStartupAction interface {
	Run()
}

type startupAction struct {
	factory    factory.InstantiateFactory
	subscribes []IStartupAction
}

func newStartupAction(factory factory.InstantiateFactory) *startupAction {
	return &startupAction{
		factory: factory,
	}
}

// 封装一个IStartupAction对象信息
type IStartupActionInfo interface {
	//设置优先级，越小越高，优先级越高的执行在最前面
	// SetPriority(priority int32)
	SetName(name string) IStartupActionInfo
	//设置最后执行，最后调用的执行在最后
	SetLast() IStartupActionInfo
}

type startupActionInfo struct {
	//名称
	name string
	//优先级，默认值为0
	priority int32
	//用来构建IStartupAction的函数
	actionFunc interface{}
}

func newStartupActionInfo(p interface{}) *startupActionInfo {
	return &startupActionInfo{
		priority:   0,
		actionFunc: p,
		name:       reflector.GetFullName(p), //默认使用类名来做name
	}
}

func (s *startupActionInfo) SetPriority(priority int32) {
	s.priority = priority
	//设置完成后，根据优先级进行排序
	comparer := newStartupActionInfoComparer(_startupActions, true)
	comparer.Sort()

}

func (s *startupActionInfo) SetName(name string) IStartupActionInfo {
	s.name = name
	return s
}

func (s *startupActionInfo) SetLast() IStartupActionInfo {
	if len(s.name) <= 0 {
		return s
	}
	//检测是否已经加入到了列表中，没有加入不进行任何处理
	exist := slicex.InWithFunc(_startupActions, func(info *startupActionInfo) bool {
		return info.name == s.name
	})
	if !exist {
		return s
	}
	defer func() {
		//设置完成后，根据优先级进行排序
		comparer := newStartupActionInfoComparer(_startupActions, true)
		comparer.Sort()
	}()
	if len(_startupActions) <= 1 {
		//第一个人设的，则直接使用最后的优先级值
		s.priority = lastPriority
		return s
	}
	//将最后的优先级加1,因为已经排序过，所以最后一个一定是最小优先级的
	minPriority := _startupActions[len(_startupActions)-1].priority
	if minPriority < lastPriority {
		//最小一个的优先级值比预定义的最后的优先级值要小，则使用预定义的最后优先级值
		s.priority = lastPriority
		return s
	}
	s.priority = lastPriority + 1
	return s
}

var (
	_startupActions []*startupActionInfo
)

// 注册一个startupAction
func RegisterOneStartupAction(p interface{}) IStartupActionInfo {
	startupActionInfo := newStartupActionInfo(p)
	//注册后，根据优先级进行排序
	startupActionInfo.SetPriority(int32(len(_startupActions)))
	_startupActions = append(_startupActions, startupActionInfo)

	return startupActionInfo
}

// 注册一组startupAction
func RegisterStartupAction(p ...interface{}) {
	for _, eachP := range p {
		startupActionInfo := newStartupActionInfo(eachP)
		startupActionInfo.priority = int32(len(_startupActions))
		_startupActions = append(_startupActions, startupActionInfo)
	}
	//注册后，根据优先维进行排序
	comparer := newStartupActionInfoComparer(_startupActions, true)
	comparer.Sort()
}

func (p *startupAction) Init() {
	for _, eachStartupAction := range _startupActions {
		ss, err := p.factory.InjectIntoFunc(nil, eachStartupAction.actionFunc)
		if err == nil {
			p.subscribes = append(p.subscribes, ss.(IStartupAction))
		}
	}
}

// 运行启动时的行为
func (p *startupAction) Run() {
	for _, eachStartupAction := range p.subscribes {
		startupName := getStartupActionTypeName(eachStartupAction)
		if !HostApplication.SystemConfig().App.IsRunInCli {
			log.Logger.Info(fmt.Sprintf("begin run startupaction,%s", startupName))
		}
		p.factory.InjectIntoFunc(nil, eachStartupAction)
		eachStartupAction.Run()
		if !HostApplication.SystemConfig().App.IsRunInCli {
			log.Logger.Info(fmt.Sprintf("finish run startupaction,%s", startupName))
		}
	}
}

func getStartupActionTypeName(startupAction IStartupAction) string {
	actionFunc, ok := startupAction.(*startupActionFunc)
	if ok && actionFunc != nil && actionFunc.runFunc != nil {
		return runtime.FuncForPC(reflect.ValueOf(actionFunc.runFunc).Pointer()).Name()
	}
	return ""
}

// 使用函数来实现IStartupAction
type startupActionFunc struct {
	runFunc func()
}

func (p *startupActionFunc) Run() {
	if p.runFunc == nil {
		return
	}
	p.runFunc()
}

// 使用函数创建一个PostProcessor对象
func NewStartupAction(runFunc func()) IStartupAction {
	return &startupActionFunc{
		runFunc: runFunc,
	}
}
