package app

var (
	Context ApplicationContext
)

const (
	ApplicationContextName = "app.applicationContext"
)

// ApplicationContext is the alias interface of Application
type ApplicationContext interface {
	IServiceProvider
	GetProperty(name string) (value interface{}, ok bool)
}

type IServiceProvider interface {
	SetInstance(params ...interface{}) (err error)
	GetInstance(params ...interface{}) (instance interface{})
	GetListByBaseInterface(interfaceInstance interface{}) (instanceList []interface{})

	//注册一个服务
	//instance 服务实例
	//typ 要注册的类型实例
	RegistInstanceAs(instance interface{}, typ interface{}) error
	RegistInstance(instance interface{}) error
}
