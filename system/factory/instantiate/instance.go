package instantiate

import (
	"reflect"

	"github.com/shanluzhineng/fwpkg/system/cmap"
	"github.com/shanluzhineng/fwpkg/system/factory"
)

type instance struct {
	instMap cmap.ConcurrentMap
}

var _ factory.Instance = (*instance)(nil)

func newInstance(instMap cmap.ConcurrentMap) factory.Instance {
	if instMap == nil {
		instMap = cmap.New()
	}
	return &instance{
		instMap: instMap,
	}
}

// #region factory.Instance Members

// Get get instance
func (i *instance) Get(params ...interface{}) (retVal interface{}) {
	name, obj := factory.ParseParams(params...)

	// get from instance map if external instance map does not have it
	if md, ok := i.instMap.Get(name); ok {
		metaData := factory.CastMetaData(md)
		if metaData != nil {
			switch obj.(type) {
			case factory.MetaData:
				retVal = metaData
			default:
				retVal = metaData.Instance
			}
		}
	}

	return
}

func (i *instance) GetListByBaseInterface(interfaceInstance interface{}) []interface{} {
	result := make([]interface{}, 0)
	items := i.instMap.Items()
	if len(items) <= 0 {
		return result
	}
	interfaceT := reflect.TypeOf(interfaceInstance).Elem()
	for _, eachInstance := range items {
		currentMetaData := factory.CastMetaData(eachInstance)
		if currentMetaData == nil || currentMetaData.Instance == nil {
			continue
		}
		reflectT := reflect.TypeOf(currentMetaData.Instance)
		if reflectT.Implements(interfaceT) {
			result = append(result, currentMetaData.Instance)
		}
	}
	return result
}

// Set save instance
func (i *instance) Set(params ...interface{}) (err error) {
	name, inst := factory.ParseParams(params...)

	metaData := factory.CastMetaData(inst)
	if metaData == nil {
		metaData = factory.NewMetaData(inst)
	}

	// old, ok := i.instMap.Get(name)
	// if ok {
	// 	err = fmt.Errorf("instance %v is already taken by %v", name, old)
	// 	// TODO: should handle such error
	// 	log.Debugf("%+v", err)
	// 	return
	// }

	i.instMap.Set(name, metaData)
	return
}

// Items return map items
func (i *instance) Items() map[string]interface{} {
	return i.instMap.Items()
}

// #endregion
