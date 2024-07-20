package app

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/shanluzhineng/fwpkg/system"
	"github.com/shanluzhineng/fwpkg/system/cmap"
	"github.com/shanluzhineng/fwpkg/system/factory"
	"github.com/shanluzhineng/fwpkg/system/factory/autoconfigure"
	"github.com/shanluzhineng/fwpkg/system/factory/instantiate"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/fwpkg/utils/io"
)

var (
	HostApplication Application
)

type Application interface {
	Initialize() error
	SetProperty(name string, value ...interface{}) Application
	GetProperty(name string) (value interface{}, ok bool)
	SetAddCommandLineProperties(enabled bool) Application
	Run()

	ConfigurableFactory() factory.ConfigurableFactory
	SystemConfig() *system.Configuration

	//获取serviceProvider
	GetServiceProvider() IServiceProvider
	// 运行启动时行为
	Build() Application
	RunStartup()
	//shutdown
	Shutdown()
}

// BaseApplication is the base application
type BaseApplication struct {
	WorkDir                  string
	configurations           cmap.ConcurrentMap
	instances                cmap.ConcurrentMap
	configurableFactory      factory.ConfigurableFactory
	systemConfig             *system.Configuration
	postProcessor            *postProcessor
	defaultProperties        cmap.ConcurrentMap
	mu                       sync.Mutex
	addCommandLineProperties bool

	//启动行为,shutdown行为
	startupAction  *startupAction
	shutdownAction *shutdownAction
}

var (
	configContainer    []*factory.MetaData
	componentContainer []*factory.MetaData
	// Profiles include profiles initially
	Profiles []string

	// ErrInvalidObjectType indicates that configuration type is invalid
	ErrInvalidObjectType = errors.New("[app] invalid Configuration type, one of app.Configuration need to be embedded")
)

// SetProperty set application property
// should be able to set property from source code by SetProperty, it can be override by program argument, e.g. myapp --app.profiles.active=dev
func (a *BaseApplication) SetProperty(name string, value ...interface{}) Application {
	var val interface{}
	if len(value) == 1 {
		val = value[0]
	} else {
		val = value
	}

	kind := reflect.TypeOf(val).Kind()
	if kind == reflect.String && strings.Contains(val.(string), ",") {
		val = strings.SplitN(val.(string), ",", -1)
	}
	a.defaultProperties.Set(name, val)

	return a
}

// #region ApplicationContext Members

// 获取应用属性
func (a *BaseApplication) GetProperty(name string) (value interface{}, ok bool) {
	value, ok = a.defaultProperties.Get(name)
	return
}

// #endregion

// 初始化应用
func (a *BaseApplication) Initialize() (err error) {
	log.SetLevel(log.InfoLevel)
	a.defaultProperties = cmap.New()
	a.configurations = cmap.New()
	a.instances = cmap.New()

	a.SetAddCommandLineProperties(true)
	return nil
}

func (a *BaseApplication) GetServiceProvider() IServiceProvider {
	return a
}

// 构建应用运行所需的环境
func (a *BaseApplication) Build() Application {
	a.mu.Lock()
	defer a.mu.Unlock()

	//构建默认的日志
	log.BuildDefaultLogger()

	a.WorkDir = io.GetWorkDir()

	instantiateFactory := instantiate.NewInstantiateFactory(a.instances, componentContainer, a.defaultProperties)
	configurableFactory := autoconfigure.NewConfigurableFactory(instantiateFactory, a.configurations)
	a.configurableFactory = configurableFactory

	a.postProcessor = newPostProcessor(instantiateFactory)
	a.systemConfig, _ = configurableFactory.BuildProperties()
	a.startupAction = newStartupAction(instantiateFactory)
	a.shutdownAction = newShutdown(instantiateFactory)

	// set logging level
	log.SetLevel(a.systemConfig.Logging.Level)

	//先设置其是运行在cli模式
	a.systemConfig.App.IsRunInCli = true

	//全局保存到app.HostApplication属性中
	HostApplication = a
	Context = a

	return a
}

// 返回应用系统配置
func (a *BaseApplication) SystemConfig() *system.Configuration {
	return a.systemConfig
}

// 获取BuildConfigurations对象
func (a *BaseApplication) BuildConfigurations() (err error) {
	// build configurations
	a.configurableFactory.Build(configContainer)
	// build components
	err = a.configurableFactory.BuildComponents()

	return
}

func (a *BaseApplication) ConfigurableFactory() factory.ConfigurableFactory {
	return a.configurableFactory
}

// 初始化完后执行的post行为
func (a *BaseApplication) AfterInitialization(configs ...cmap.ConcurrentMap) {
	// pass user's instances
	a.postProcessor.Init()
	a.postProcessor.AfterInitialization()
	if a.systemConfig != nil && a.systemConfig.App != nil && !a.systemConfig.App.IsRunInCli {
		log.Logger.Debug(fmt.Sprintf("command line properties is enabled: %t", a.addCommandLineProperties))
	}
}

// 运行启动时行为
func (a *BaseApplication) RunStartup() {
	a.startupAction.Init()
	a.startupAction.Run()
	if a.systemConfig != nil && a.systemConfig.App != nil && !a.systemConfig.App.IsRunInCli {
		log.Logger.Info("run startup action finished")
	}
}

// shutdown 应用
func (a *BaseApplication) Shutdown() {
	if a.systemConfig != nil && a.systemConfig.App != nil && !a.systemConfig.App.IsRunInCli {
		log.Logger.Info("准备shutdown application...")
	}
	a.shutdownAction.Init()
	a.shutdownAction.Shutdown()
	if a.systemConfig != nil && a.systemConfig.App != nil && !a.systemConfig.App.IsRunInCli {
		log.Logger.Info("shutdown application完成")
	}
}

// SetAddCommandLineProperties set add command line properties to be enabled or disabled
func (a *BaseApplication) SetAddCommandLineProperties(enabled bool) Application {
	a.addCommandLineProperties = enabled
	return a
}

// 启动应用，继承类实现
func (a *BaseApplication) Run() {
	log.Logger.Warn("application is not implemented!")
}

// #region IServiceProvider Members

// 从ioc容器中获取对象实例
func (a *BaseApplication) GetInstance(params ...interface{}) (instance interface{}) {
	if a.configurableFactory != nil {
		instance = a.configurableFactory.GetInstance(params...)
	}
	return
}

func (a *BaseApplication) GetListByBaseInterface(interfaceInstance interface{}) (instanceList []interface{}) {
	if a.configurableFactory != nil {
		instanceList = a.configurableFactory.GetListByBaseInterface(interfaceInstance)
	}
	return
}

func (a *BaseApplication) SetInstance(params ...interface{}) (err error) {
	if a.configurableFactory == nil {
		return errors.New("在执行此方法前必须先调用Build方法")
	}
	return a.configurableFactory.SetInstance(params...)
}

// 注册一个服务
// instance 服务实例
// typ 要注册的类型实例
func (a *BaseApplication) RegistInstanceAs(instance interface{}, typ interface{}) error {
	if a.configurableFactory == nil {
		return errors.New("在执行此方法前必须先调用Build方法")
	}
	name, _ := factory.ParseParams(typ)
	return a.SetInstance(name, instance)
}

func (a *BaseApplication) RegistInstance(instance interface{}) error {
	if a.configurableFactory == nil {
		return errors.New("在执行此方法前必须先调用Build方法")
	}
	return a.SetInstance(instance)
}

// #endregion
