package inject

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/shanluzhineng/fwpkg/system/factory"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/fwpkg/system/reflector"
	"github.com/shanluzhineng/fwpkg/utils/str"
)

var (
	// ErrNotImplemented the interface is not implemented
	ErrNotImplemented = errors.New("[inject] interface is not implemented")

	// ErrInvalidObject the object is invalid
	ErrInvalidObject = errors.New("[inject] invalid object")

	// ErrInvalidTagName the tag name is invalid
	ErrInvalidTagName = errors.New("[inject] invalid tag name, e.g. exampleTag")

	// ErrSystemConfiguration system is not configured
	ErrSystemConfiguration = errors.New("[inject] system is not configured")

	// ErrInvalidFunc the function is invalid
	ErrInvalidFunc = errors.New("[inject] invalid func")

	// ErrInvalidMethod the function is invalid
	ErrInvalidMethod = errors.New("[inject] invalid method")

	// ErrFactoryIsNil factory is invalid
	ErrFactoryIsNil = errors.New("[inject] factory is nil")

	ErrAnnotationsIsNil = fmt.Errorf("err: annotations is nil")

	tagsContainer []Tag
)

// Inject is the interface for inject tag
type Inject interface {
	IntoObject(instance factory.Instance, object interface{}) error
	IntoObjectValue(instance factory.Instance, object reflect.Value, property string, tags ...Tag) error
	IntoMethod(instance factory.Instance, object interface{}, m interface{}) (retVal interface{}, err error)
	IntoFunc(instance factory.Instance, object interface{}) (retVal interface{}, err error)
}

type inject struct {
	factory factory.InstantiateFactory
}

// NewInject is the constructor of inject
func NewInject(factory factory.InstantiateFactory) Inject {
	return &inject{factory: factory}
}

// InitTag init tag implements
func InitTag(tag Tag) (t Tag) {
	return
}

// AddTag add new tag
func AddTag(tag Tag) {
	if tag != nil {
		t := InitTag(tag)
		if t != nil {
			tagsContainer = append(tagsContainer, t)
		}
	}
}

func (i *inject) getInstance(instance factory.Instance, typ reflect.Type) (inst interface{}) {
	n := reflector.GetLowerCamelFullNameByType(typ)
	inst = i.factory.GetInstance(instance, n)
	return
}

// IntoObject injects instance into the tagged field with `inject:"instanceName"`
func (i *inject) IntoObject(instance factory.Instance, object interface{}) (err error) {
	// inject into value
	err = i.IntoObjectValue(instance, reflect.ValueOf(object), "")
	return
}

func (i *inject) convert(f reflect.StructField, src interface{}) (fov reflect.Value) {
	fov = reflect.ValueOf(src)

	// convert slice
	switch src.(type) {
	//case []string:
	case []interface{}:
		switch f.Type.Elem().Kind() {
		case reflect.String:
			var sv []string
			src := src.([]interface{})
			for _, elm := range src {
				sv = append(sv, elm.(string))
			}
			fov = reflect.ValueOf(sv)
		}
	}
	//log.Debugf("Injected slice %v.(%v) into %v.%v", src, fov.Type(), fov.Type(), f.Name)
	return
}

// IntoObjectValue injects instance into the tagged field with `inject:"instanceName"`
func (i *inject) IntoObjectValue(instance factory.Instance, object reflect.Value, property string, tags ...Tag) error {
	var err error

	//// TODO refactor IntoObject
	//if appFactory == nil {
	//	return ErrSystemConfiguration
	//}

	obj := reflector.Indirect(object)
	if obj.Kind() != reflect.Struct {
		log.Debugf("[inject] ignore object: %v, kind: %v", object, obj.Kind())
		return ErrInvalidObject
	}

	var targetTags []Tag
	if len(tags) != 0 {
		targetTags = tags
	} else {
		targetTags = tagsContainer
	}

	// field injection
	for _, f := range reflector.DeepFields(object.Type()) {
		var injectedObject interface{}
		var prop string
		pn, ok := f.Tag.Lookup("json")
		if ok {
			pns := strings.Split(pn, ",")
			if len(pns) > 1 {
				prop = pns[0]
			}
		} else {
			prop = str.ToLowerCamel(f.Name)
		}

		// debug prints
		//n := reflector.GetLowerCamelFullNameByType(f.Type)
		//log.Debugf("%v : %v",property + "." + prop, n)

		if property != "" {
			prop = property + "." + prop
		}

		//log.Debugf("parent: %v, name: %v, type: %v, tag: %v", obj.Type(), f.Name, f.Type, f.Tag)
		// check if object has value field to be injected

		ft := f.Type
		if f.Type.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		// set field object
		var fieldObjValue reflect.Value
		if obj.IsValid() && obj.Kind() == reflect.Struct {
			// TODO: consider embedded property
			//log.Debugf("inject: %v.%v", obj.Type(), f.Name)
			fieldObjValue = obj.FieldByName(f.Name)
		}

		// TODO: assume that the f.Name of value and inject tag is not the same
		injectedObject = i.getInstance(instance, f.Type)
		if injectedObject == nil {
			for _, tagImpl := range targetTags {
				tagImpl.Init(i.factory)
				injectedObject = tagImpl.Decode(object, f, prop)
				if injectedObject != nil {
					break
				}
			}
		}

		// assign value to struct field
		// if ft.Kind() != reflect.Struct || annotation.Contains(injectedObject, at.AutoWired{}) {
		if ft.Kind() != reflect.Struct {
			if injectedObject != nil && fieldObjValue.CanSet() {
				fov := i.convert(f, injectedObject)
				if fov.Type().AssignableTo(fieldObjValue.Type()) {
					fieldObjValue.Set(fov)
					//} else {
					//	log.Errorf("unmatched type %v against %v", fov.Type(), fieldObj.Type())
				}
			}
		}

		//log.Debugf("- kind: %v, %v, %v, %v", obj.Kind(), object.Type(), fieldObj.Type(), f.Name)
		//log.Debugf("isValid: %v, canSet: %v", fieldObj.IsValid(), fieldObj.CanSet())
		filedObject := reflect.Indirect(fieldObjValue)
		filedKind := filedObject.Kind()
		canNested := filedKind == reflect.Struct
		if canNested && fieldObjValue.IsValid() && fieldObjValue.CanSet() && filedObject.Type() != obj.Type() {
			err = i.IntoObjectValue(instance, fieldObjValue, prop, tags...)
		}
	}

	return err
}

func (i *inject) parseFuncOrMethodInput(instance factory.Instance, inType reflect.Type) (paramValue reflect.Value, ok bool) {
	inType = reflector.IndirectType(inType)
	inst := i.getInstance(instance, inType)
	ok = true
	if inst == nil {
		//log.Debug(inType.Kind())
		switch inType.Kind() {
		// interface and slice creation is not supported
		case reflect.Interface, reflect.Slice:
			ok = false
			break
		default:
			// should find instance in the component container first

			// if it is not found, then create new instance
			paramValue = reflect.New(inType)
			inst = paramValue.Interface()
			// TODO: inTypeName
			i.factory.SetInstance(instance, inst)
		}
	}

	if inst != nil {
		paramValue = reflect.ValueOf(inst)
	}
	return
}

// IntoFunc inject object into func and return instance
func (i *inject) IntoFunc(instance factory.Instance, object interface{}) (retVal interface{}, err error) {
	fn := reflect.ValueOf(object)
	if fn.Kind() == reflect.Func {
		numIn := fn.Type().NumIn()
		inputs := make([]reflect.Value, numIn)
		// TODO: should load function inputs when resolving dependencies to improve performance
		for n := 0; n < numIn; n++ {
			fnInType := fn.Type().In(n)
			//expectedTypName := reflector.GetLowerCamelFullNameByType(fnInType)
			//log.Debugf("expected: %v", expectedTypName)
			val, ok := i.parseFuncOrMethodInput(instance, fnInType)
			if ok {

				inputs[n] = val
				//log.Debugf("Injected %v into func parameter %v", val, fnInType)
			} else {
				return nil, fmt.Errorf("%v is not injected", fnInType.Name())
			}

			paramValue := reflect.Indirect(val)
			if val.IsValid() && paramValue.IsValid() && paramValue.Kind() == reflect.Struct {
				err = i.IntoObject(instance, val.Interface())
			}
		}
		results := fn.Call(inputs)
		if len(results) != 0 {
			retVal = results[0].Interface()
			return
		}
		return
	}
	err = ErrInvalidFunc
	return
}

// IntoMethod inject object into func and return instance
// TODO: IntoMethod or IntoFunc should accept metaData, because it contains dependencies
func (i *inject) IntoMethod(instance factory.Instance, object interface{}, m interface{}) (retVal interface{}, err error) {
	if object != nil && m != nil {
		switch m.(type) {
		case reflect.Method:
			method := m.(reflect.Method)
			numIn := method.Type.NumIn()
			inputs := make([]reflect.Value, numIn)
			inputs[0] = reflect.ValueOf(object)
			// var ann interface{}
			for n := 1; n < numIn; n++ {
				fnInType := method.Type.In(n)
				val, ok := i.parseFuncOrMethodInput(instance, fnInType)
				if ok {
					inputs[n] = val
				} else {
					return nil, fmt.Errorf("%v is not injected", fnInType.Name())
				}

				paramObject := reflect.Indirect(val)
				if val.IsValid() && paramObject.IsValid() && paramObject.Kind() == reflect.Struct {
					err = i.IntoObject(instance, val.Interface())
				}
			}
			results := method.Func.Call(inputs)
			if len(results) != 0 {
				retVal = results[0].Interface()
				if len(results) > 1 {
					errObj := results[1].Interface()
					switch errObj.(type) {
					case error:
						err = errObj.(error)
					}
				}
				return
			}
		}
	}
	err = ErrInvalidMethod
	return
}
