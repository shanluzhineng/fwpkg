// Package instantiate implement InstantiateFactory
package instantiate

import (
	"errors"
	"path/filepath"
	"sync"

	"github.com/shanluzhineng/fwpkg/system"
	"github.com/shanluzhineng/fwpkg/system/cmap"
	"github.com/shanluzhineng/fwpkg/system/factory"
	"github.com/shanluzhineng/fwpkg/system/factory/depends"
	"github.com/shanluzhineng/fwpkg/system/inject"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/fwpkg/system/reflector"
	"github.com/shanluzhineng/fwpkg/utils/io"
)

var (
	// ErrNotInitialized InstantiateFactory is not initialized
	ErrNotInitialized = errors.New("[factory] InstantiateFactory is not initialized")

	// ErrInvalidObjectType invalid object type
	ErrInvalidObjectType = errors.New("[factory] invalid object type")
)

const (
	application = "application"
	config      = "config"
	yaml        = "yaml"
)

// InstantiateFactory is the factory that responsible for object instantiation
type instantiateFactory struct {
	instance          factory.Instance
	components        []*factory.MetaData
	resolved          []*factory.MetaData
	defaultProperties cmap.ConcurrentMap
	inject            inject.Inject
	builder           system.Builder
	mutex             sync.Mutex
}

// NewInstantiateFactory the constructor of instantiateFactory
func NewInstantiateFactory(instanceMap cmap.ConcurrentMap, components []*factory.MetaData, defaultProperties cmap.ConcurrentMap) factory.InstantiateFactory {
	if defaultProperties == nil {
		defaultProperties = cmap.New()
	}

	f := &instantiateFactory{
		instance:          newInstance(instanceMap),
		components:        components,
		defaultProperties: defaultProperties,
	}
	f.inject = inject.NewInject(f)

	// create new builder
	workDir := io.GetWorkDir()

	sa := new(system.App)
	sl := new(system.Logging)
	syscfg := system.NewConfiguration()

	customProps := defaultProperties.Items()
	f.builder = system.NewPropertyBuilder(
		filepath.Join(workDir, config),
		customProps,
	)

	f.Append(syscfg, sa, sl, f, f.builder)

	return f
}

// Initialized check if factory is initialized
func (f *instantiateFactory) Initialized() bool {
	return f.instance != nil
}

// Builder get builder
func (f *instantiateFactory) Builder() (builder system.Builder) {
	return f.builder
}

// GetProperty get property
func (f *instantiateFactory) GetProperty(name string) (retVal interface{}) {
	retVal = f.builder.GetProperty(name)
	return
}

// SetProperty get property
func (f *instantiateFactory) SetProperty(name string, value interface{}) factory.InstantiateFactory {
	f.builder.SetProperty(name, value)
	return f
}

// SetDefaultProperty set default property
func (f *instantiateFactory) SetDefaultProperty(name string, value interface{}) factory.InstantiateFactory {
	f.builder.SetDefaultProperty(name, value)
	return f
}

// Append append to component and instance container
func (f *instantiateFactory) Append(i ...interface{}) {
	for _, inst := range i {
		f.AppendComponent(inst)
		_ = f.SetInstance(inst)
	}
}

// AppendComponent append component
func (f *instantiateFactory) AppendComponent(c ...interface{}) {
	metaData := factory.NewMetaData(c...)
	f.components = append(f.components, metaData)
}

// injectDependency inject dependency
func (f *instantiateFactory) injectDependency(instance factory.Instance, item *factory.MetaData) (err error) {
	var name string
	var inst interface{}
	switch item.Kind {
	case factory.Func:
		inst, err = f.inject.IntoFunc(instance, item.MetaObject)
		name = item.Name
		// TODO: should report error when err is not nil
		if err == nil {
			log.Debugf("inject into func: %v %v", item.ShortName, item.Type)
		}
	case factory.Method:
		inst, err = f.inject.IntoMethod(instance, item.ObjectOwner, item.MetaObject)
		name = item.Name
		if err != nil {
			return
		}
		log.Debugf("inject into method: %v %v", item.ShortName, item.Type)
	default:
		name, inst = item.Name, item.MetaObject
	}
	if inst != nil {
		// inject into object
		err = f.inject.IntoObject(instance, inst)
		if name != "" {
			// save object
			item.Instance = inst
			// set item
			err = f.SetInstance(instance, name, item)
		}
	}
	return
}

// InjectDependency inject dependency
func (f *instantiateFactory) InjectDependency(instance factory.Instance, object interface{}) (err error) {
	return f.injectDependency(instance, factory.CastMetaData(object))
}

// BuildComponents build all registered components
func (f *instantiateFactory) BuildComponents() (err error) {
	// first resolve the dependency graph
	var resolved []*factory.MetaData
	log.Debugf("Resolving dependencies")
	resolved, err = depends.Resolve(f.components)
	f.resolved = resolved
	log.Debugf("Injecting dependencies")
	// then build components
	for _, item := range resolved {
		// log.Debugf("build component: %v %v", idx, item.Type)
		// inject dependencies into function
		// components, controllers
		// TODO: should save the upstream dependencies that contains item.ContextAware annotation for runtime injection
		err = f.injectDependency(f.instance, item)
	}
	if err == nil {
		log.Debugf("Injected dependencies")
	}
	return
}

// SetInstance save instance
func (f *instantiateFactory) SetInstance(params ...interface{}) (err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	var instance factory.Instance
	switch params[0].(type) {
	case factory.Instance:
		instance = params[0].(factory.Instance)
		params = params[1:]
	default:
		instance = f.instance
		if len(params) > 1 && params[0] == nil {
			params = params[1:]
		}
	}

	name, inst := factory.ParseParams(params...)

	if inst == nil {
		return ErrNotInitialized
	}

	metaData := factory.CastMetaData(inst)
	if metaData == nil {
		metaData = factory.NewMetaData(inst)
	}

	if metaData != nil {
		err = instance.Set(name, inst)
	}

	return
}

// GetInstance get instance by name
func (f *instantiateFactory) GetInstance(params ...interface{}) (retVal interface{}) {
	switch params[0].(type) {
	case factory.Instance:
		instance := params[0].(factory.Instance)
		params = params[1:]
		retVal = instance.Get(params...)
	default:
		if len(params) > 1 && params[0] == nil {
			params = params[1:]
		}

	}
	// if it does not found from instance, try to find it from f.instance
	if retVal == nil {
		retVal = f.instance.Get(params...)
	}
	return
}

func (f *instantiateFactory) GetListByBaseInterface(interfaceInstance interface{}) []interface{} {
	return f.instance.GetListByBaseInterface(interfaceInstance)
}

// Items return instance map
func (f *instantiateFactory) Items() map[string]interface{} {
	return f.instance.Items()
}

// DefaultProperties return default properties
func (f *instantiateFactory) DefaultProperties() map[string]interface{} {
	dp := f.defaultProperties.Items()
	return dp
}

// InjectIntoObject inject into object
func (f *instantiateFactory) InjectIntoObject(instance factory.Instance, object interface{}) error {
	return f.inject.IntoObject(instance, object)
}

// InjectIntoFunc inject into func
func (f *instantiateFactory) InjectIntoFunc(instance factory.Instance, object interface{}) (retVal interface{}, err error) {
	return f.inject.IntoFunc(instance, object)
}

// InjectIntoMethod inject into method
func (f *instantiateFactory) InjectIntoMethod(instance factory.Instance, owner, object interface{}) (retVal interface{}, err error) {
	return f.inject.IntoMethod(instance, owner, object)
}

func (f *instantiateFactory) Replace(source string) (retVal interface{}) {
	retVal = f.builder.Replace(source)
	return
}

// InjectContextAwareObject inject context aware objects
func (f *instantiateFactory) injectContextAwareDependencies(instance factory.Instance, dps []*factory.MetaData) (err error) {
	for _, d := range dps {
		if len(d.DepMetaData) > 0 {
			err = f.injectContextAwareDependencies(instance, d.DepMetaData)
			if err != nil {
				return
			}
		}
	}
	return
}

// InjectContextAwareObjects inject context aware objects
func (f *instantiateFactory) InjectContextAwareObjects(instanceFunc func() interface{}, dps []*factory.MetaData) (instance factory.Instance, err error) {
	injectInstance := instanceFunc()
	if injectInstance == nil {
		err = errors.New("can't inject nil object")
		log.Error()
		return
	}
	log.Debugf(">>> InjectContextAwareObjects(%x) ...", &injectInstance)

	// create new runtime instance
	instance = newInstance(nil)

	// update context
	err = instance.Set(reflector.GetLowerCamelFullName(instance), injectInstance)
	if err != nil {
		log.Error(err)
		return
	}

	err = f.injectContextAwareDependencies(instance, dps)
	if err != nil {
		log.Error(err)
	}

	return
}
