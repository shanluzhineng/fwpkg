// Package autoconfigure implement ConfigurableFactory
package autoconfigure

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/shanluzhineng/fwpkg/system"
	"github.com/shanluzhineng/fwpkg/system/cmap"
	"github.com/shanluzhineng/fwpkg/system/factory"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/fwpkg/system/reflector"
	"github.com/shanluzhineng/fwpkg/utils/io"
	"github.com/shanluzhineng/fwpkg/utils/str"
)

const (
	// System configuration name
	System = "system"

	// PropAppProfilesActive is the property name "app.profiles.active"
	PropAppProfilesActive = "app.profiles.active"

	// EnvAppProfilesActive is the environment variable name APP_PROFILES_ACTIVE
	EnvAppProfilesActive = "APP_PROFILES_ACTIVE"

	// PostfixConfiguration is the Configuration postfix
	PostfixConfiguration = "Configuration"

	defaultProfileName = "default"
)

var (
	// ErrInvalidMethod method is invalid
	ErrInvalidMethod = errors.New("[factory] method is invalid")

	// ErrFactoryCannotBeNil means that the InstantiateFactory can not be nil
	ErrFactoryCannotBeNil = errors.New("[factory] InstantiateFactory can not be nil")

	// ErrFactoryIsNotInitialized means that the InstantiateFactory is not initialized
	ErrFactoryIsNotInitialized = errors.New("[factory] InstantiateFactory is not initialized")

	// ErrInvalidObjectType means that the Configuration type is invalid, it should embeds app.Configuration
	ErrInvalidObjectType = errors.New("[factory] invalid Configuration type, one of app.Configuration need to be embedded")

	// ErrConfigurationNameIsTaken means that the configuration name is already taken
	ErrConfigurationNameIsTaken = errors.New("[factory] configuration name is already taken")

	// ErrComponentNameIsTaken means that the component name is already taken
	ErrComponentNameIsTaken = errors.New("[factory] component name is already taken")
)

type configurableFactory struct {
	// at.Qualifier `value:"factory.configurableFactory"`

	factory.InstantiateFactory
	configurations cmap.ConcurrentMap
	systemConfig   *system.Configuration

	// preConfigureContainer  []*factory.MetaData
	configureContainer []*factory.MetaData
	// postConfigureContainer []*factory.MetaData
	builder system.Builder
}

// NewConfigurableFactory is the constructor of configurableFactory
func NewConfigurableFactory(instantiateFactory factory.InstantiateFactory, configurations cmap.ConcurrentMap) factory.ConfigurableFactory {
	f := &configurableFactory{
		InstantiateFactory: instantiateFactory,
		configurations:     configurations,
	}

	f.configurations = configurations
	_ = f.SetInstance("configurations", configurations)

	f.builder = f.Builder()

	f.Append(f)
	return f
}

// SystemConfiguration getter
func (f *configurableFactory) SystemConfiguration() *system.Configuration {
	return f.systemConfig
}

// Configuration getter
func (f *configurableFactory) Configuration(name string) interface{} {
	cfg, ok := f.configurations.Get(name)
	if ok {
		return cfg
	}
	return nil
}

// BuildProperties build all properties
func (f *configurableFactory) BuildProperties() (systemConfig *system.Configuration, err error) {
	// manually inject systemConfiguration
	systemConfig = f.GetInstance(system.Configuration{}).(*system.Configuration)

	profile := os.Getenv(EnvAppProfilesActive)
	if profile == "" {
		profile = defaultProfileName
	}
	f.builder.SetDefaultProperty(PropAppProfilesActive, profile)

	for prop, val := range f.DefaultProperties() {
		f.builder.SetDefaultProperty(prop, val)
	}

	_, err = f.builder.Build(profile)
	if err == nil {
		_ = f.InjectIntoObject(nil, systemConfig)

		f.configurations.Set(System, systemConfig)

		f.systemConfig = systemConfig
	}

	return
}

// Build build all auto configurations
func (f *configurableFactory) Build(configs []*factory.MetaData) {
	// categorize configurations first, then inject object if necessary
	for _, item := range configs {
		err := ErrInvalidObjectType
		log.Errorf("item: %v err: %v", item, err)
	}

	f.build(f.configureContainer)

}

// Instantiate run instantiation by method
func (f *configurableFactory) Instantiate(configuration interface{}) (err error) {
	cv := reflect.ValueOf(configuration)
	icv := reflector.Indirect(cv)

	configType := cv.Type()
	//log.Debug("type: ", configType)
	//name := configType.Elem().Name()
	//log.Debug("fieldName: ", name)
	pkgName := io.DirName(icv.Type().PkgPath())
	var runtimeDeps factory.Deps
	rd := icv.FieldByName("RuntimeDeps")
	if rd.IsValid() {
		runtimeDeps = rd.Interface().(factory.Deps)
	}
	// call Init
	numOfMethod := cv.NumMethod()
	//log.Debug("methods: ", numOfMethod)
	for mi := 0; mi < numOfMethod; mi++ {
		// get method
		// find the dependencies of the method
		method := configType.Method(mi)
		methodName := str.LowerFirst(method.Name)
		if rd.IsValid() {
			// append inst to f.components
			deps := runtimeDeps.Get(method.Name)

			metaData := &factory.MetaData{
				Name:       pkgName + "." + methodName,
				MetaObject: method,
				DepNames:   deps,
			}
			f.AppendComponent(configuration, metaData)
		} else {
			f.AppendComponent(configuration, method)
		}
	}
	return
}

func (f *configurableFactory) parseName(item *factory.MetaData) string {

	//return item.PkgName
	name := strings.Replace(item.TypeName, PostfixConfiguration, "", -1)
	name = str.ToLowerCamel(name)

	if name == "" || name == strings.ToLower(PostfixConfiguration) {
		name = item.PkgName
	}
	return name
}

func (f *configurableFactory) build(cfgContainer []*factory.MetaData) {
	var err error
	for _, item := range cfgContainer {
		name := f.parseName(item)
		config := item.MetaObject

		// inject into func
		var cf interface{}
		if item.Kind == factory.Func {
			cf, err = f.InjectIntoFunc(nil, config)
		}
		if err == nil && cf != nil {

			// inject other fields
			_ = f.InjectIntoObject(nil, cf)

			// instantiation
			_ = f.Instantiate(cf)

			// save configuration
			configName := name
			// TODO: should set full name instead
			f.configurations.Set(configName, cf)
		} else {
			log.Warn(err)
		}
	}
}
