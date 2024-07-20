package app

import (
	"reflect"

	"github.com/shanluzhineng/fwpkg/system/factory"
)

// 应用配置基类
type Configuration struct {
	RuntimeDeps factory.Deps
}

// appendParam is the common func to append meta data to meta data slice
func appendParam(container []*factory.MetaData, params ...interface{}) (retVal []*factory.MetaData, err error) {

	retVal = container

	// parse meta data
	metaData := factory.NewMetaData(params...)

	// append meta data
	if metaData.MetaObject != nil {
		retVal = append(retVal, metaData)
	}
	return
}

// appendParams is the common func to append params to component or configuration containers
func appendParams(container []*factory.MetaData, params ...interface{}) (retVal []*factory.MetaData, err error) {
	retVal = container
	if len(params) == 0 || params[0] == nil {
		err = ErrInvalidObjectType
		return
	}

	if len(params) > 1 && reflect.TypeOf(params[0]).Kind() != reflect.String {
		for _, param := range params {
			retVal, err = appendParam(retVal, param)
		}
	} else {
		retVal, err = appendParam(retVal, params...)
	}
	return
}

// IncludeProfiles include specific profiles
func IncludeProfiles(profiles ...string) {
	Profiles = append(Profiles, profiles...)
}

// Register register a struct instance or constructor (func), so that it will be injectable.
func Register(params ...interface{}) {
	// appendParams will append the object that annotated with at.AutoConfiguration
	componentContainer, _ = appendParams(componentContainer, params...)
}
