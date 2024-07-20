package inject

import (
	"reflect"
	"strings"

	"github.com/shanluzhineng/fwpkg/system/cmap"
	"github.com/shanluzhineng/fwpkg/system/factory"
)

// Tag the interface of Tag
type Tag interface {
	// Init init tag
	Init(configurableFactory factory.InstantiateFactory)
	// Decode parse tag and do dependency injection
	Decode(object reflect.Value, field reflect.StructField, property string) (retVal interface{})
	// Properties get properties
	Properties() cmap.ConcurrentMap
	// IsSingleton check if it is Singleton
	IsSingleton() bool
}

// BaseTag is the base struct of tag
type BaseTag struct {
	instantiateFactory factory.InstantiateFactory
	properties         cmap.ConcurrentMap
}

// IsSingleton check if it is Singleton
func (t *BaseTag) IsSingleton() bool {
	return false
}

// Init init the tag
func (t *BaseTag) Init(configurableFactory factory.InstantiateFactory) {
	t.instantiateFactory = configurableFactory
}

// ParseProperties parse properties
func (t *BaseTag) ParseProperties(tag string) cmap.ConcurrentMap {
	t.properties = cmap.New()

	args := strings.Split(tag, ",")
	for _, v := range args {
		//log.Debug(v)
		n := strings.Index(v, "=")
		if n > 0 {
			key := v[:n]
			val := v[n+1:]
			if key != "" && val != "" {
				// check if val contains reference or env
				// TODO: should lookup certain config instead of for loop
				replacedVal := t.instantiateFactory.Replace(val)
				t.properties.Set(key, replacedVal)
			}
		}
	}
	return t.properties
}

// Properties get properties
func (t *BaseTag) Properties() cmap.ConcurrentMap {
	return t.properties
}

// Decode no implementation for base tag
func (t *BaseTag) Decode(object reflect.Value, field reflect.StructField, property string) (retVal interface{}) {
	return nil
}
