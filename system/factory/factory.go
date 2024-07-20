// Package factory provides InstantiateFactory and ConfigurableFactory interface
package factory

import (
	"reflect"

	"github.com/shanluzhineng/fwpkg/system"
	"github.com/shanluzhineng/fwpkg/system/reflector"
)

const (
	// InstantiateFactoryName is the instance name of factory.instantiateFactory
	InstantiateFactoryName = "factory.instantiateFactory"
	// ConfigurableFactoryName is the instance name of factory.configurableFactory
	ConfigurableFactoryName = "factory.configurableFactory"
)

// Factory interface
type Factory interface{}

type Instance interface {
	Get(params ...interface{}) (retVal interface{})
	GetListByBaseInterface(interfaceInstance interface{}) []interface{}
	Set(params ...interface{}) (err error)
	Items() map[string]interface{}
}

// InstantiateFactory instantiate factory interface
type InstantiateFactory interface {
	Initialized() bool
	SetInstance(params ...interface{}) (err error)
	GetInstance(params ...interface{}) (retVal interface{})
	GetListByBaseInterface(interfaceInstance interface{}) []interface{}
	Items() map[string]interface{}
	Append(i ...interface{})
	AppendComponent(c ...interface{})
	BuildComponents() (err error)
	Builder() (builder system.Builder)
	GetProperty(name string) interface{}
	SetProperty(name string, value interface{}) InstantiateFactory
	SetDefaultProperty(name string, value interface{}) InstantiateFactory
	DefaultProperties() map[string]interface{}
	InjectIntoFunc(instance Instance, object interface{}) (retVal interface{}, err error)
	InjectIntoMethod(instance Instance, owner, object interface{}) (retVal interface{}, err error)
	InjectIntoObject(instance Instance, object interface{}) error
	InjectDependency(instance Instance, object interface{}) (err error)
	Replace(name string) interface{}
	InjectContextAwareObjects(instanceFunc func() interface{}, dps []*MetaData) (runtimeInstance Instance, err error)
}

// ConfigurableFactory configurable factory interface
type ConfigurableFactory interface {
	InstantiateFactory
	SystemConfiguration() *system.Configuration
	Configuration(name string) interface{}
	BuildProperties() (systemConfig *system.Configuration, err error)
	Build(configs []*MetaData)
}

// Configuration configuration interface
type Configuration interface {
}

type depsMap map[string][]string

// Deps the dependency mapping of configuration
type Deps struct {
	deps depsMap
}

func (c *Deps) ensure() {
	if c.deps == nil {
		c.deps = make(depsMap)
	}
}

// Get get the dependencies mapping
func (c *Deps) Get(name string) (deps []string) {
	c.ensure()

	deps = c.deps[name]

	return
}

// Set set dependencies
func (c *Deps) Set(dep interface{}, value []string) {
	c.ensure()
	var name string
	val := reflect.ValueOf(dep)
	kind := val.Kind()
	switch kind {
	case reflect.Func:
		name = reflector.GetFuncName(dep)
	case reflect.String:
		name = dep.(string)
	default:
		return
	}
	c.deps[name] = value
}

// CastMetaData cast object to *factory.MetaData
func CastMetaData(object interface{}) (metaData *MetaData) {
	switch object.(type) {
	case *MetaData:
		metaData = object.(*MetaData)
	}
	return
}
