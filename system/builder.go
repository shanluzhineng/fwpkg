// TODO: app config should be generic kit

package system

import (
	"github.com/mitchellh/mapstructure"
)

// Builder is the config file (yaml, json) builder
type Builder interface {
	Init() error
	Build(profiles ...string) (p interface{}, err error)
	BuildWithProfile(profile string) (interface{}, error)
	Load(properties interface{}, opts ...func(*mapstructure.DecoderConfig)) (err error)
	Save(p interface{}) (err error)
	Replace(source string) (retVal interface{})
	GetProperty(name string) (retVal interface{})
	SetProperty(name string, val interface{}) Builder
	SetDefaultProperty(name string, val interface{}) Builder
	SetConfiguration(in interface{})
}
