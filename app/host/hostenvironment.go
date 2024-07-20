package host

import (
	"encoding/json"
	"fmt"
	golog "log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shanluzhineng/fwpkg/system/log"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapio"
)

const (
	FrameworkVersion = "1.0.0"

	HostEnvironment_Windows    = "windows"
	HostEnvironment_Linux      = "linux"
	HostEnvironment_Supervisor = "supervisor"
	HostEnvironment_Docker     = "docker"
	HostEnvironment_Systemd    = "systemd"
	HostEnvironment_Other      = "other"
)

var (
	_hostEnvironment *hostEnvironment = &hostEnvironment{
		os:         GOOS_OS(runtime.GOOS),
		properties: make(map[string]interface{}),
	}
	_zapWriter *zapio.Writer
)

func GetHostEnvironment() IHostEnvironment {
	return _hostEnvironment
}

type IHostEnvironment interface {
	GetOS() GOOS_OS
	GetCurrentPath() string
	GetExecFileName() string
	SetProduct(product string)
	SetAppName(appName string)
	SetAppVersion(appVersion string)
	GetHttp() string
	SetHttp(http string)

	SetEnv(key string, value interface{})
	GetEnv(key string) interface{}
	GetEnvString(key string) string

	//设置os环境 变量
	SetOSEnv(key string, value string)
	AllKey() []string
}

type Option func(IHostEnvironment)

type hostEnvironment struct {
	os         GOOS_OS
	properties map[string]interface{}
}

func SetupHostEnvironment(opts ...Option) IHostEnvironment {
	//使用默认日志
	log.BuildDefaultLogger()

	//将go标准库中的log与zap集成
	level, _ := zapcore.ParseLevel(log.DefaultLogConfiguration.Level)
	_zapWriter = &zapio.Writer{
		Log:   log.Logger,
		Level: level,
	}
	golog.SetOutput(_zapWriter)

	if _hostEnvironment.properties == nil {
		_hostEnvironment.properties = make(map[string]interface{})
	}

	//设置操作系统
	if len(_hostEnvironment.os) <= 0 {
		_hostEnvironment.os = GOOS_OS(runtime.GOOS)
	}
	filePath, _ := os.Executable()
	if len(filePath) > 0 {
		dir, file := filepath.Split(filePath)
		//应用路径与文件名
		_hostEnvironment.SetEnv(ENV_Path, dir)
		//linux应用名直接使用文件名
		_hostEnvironment.SetEnv(ENV_FileName, file)
		if _hostEnvironment.os.IsLinux() {
			_hostEnvironment.SetEnv(ENV_AppName, file)
			_hostEnvironment.SetEnv(ENV_HostEnvironment, HostEnvironment_Linux)
		} else if _hostEnvironment.os.IsWindows() {
			filesuffix := path.Ext(file)
			_hostEnvironment.SetEnv(ENV_AppName, file[0:len(file)-len(filesuffix)])
			_hostEnvironment.SetEnv(ENV_HostEnvironment, HostEnvironment_Windows)
		}
	}
	_hostEnvironment.SetEnv(ENV_Product, _hostEnvironment.GetEnvString(ENV_AppName))
	_hostEnvironment.SetEnv(ENV_IsHostInABMP, true)
	_hostEnvironment.SetEnv(ENV_Description, _hostEnvironment.GetEnvString(ENV_AppName))
	_hostEnvironment.SetEnv(ENV_StartTime, time.Now())
	//版本号
	_hostEnvironment.SetEnv(ENV_FrameworkVersion, FrameworkVersion)
	//读取环境变量
	hostEnv := os.Getenv(strings.Replace(ENV_HostEnvironment, ".", "_", 1))
	if len(hostEnv) > 0 {
		_hostEnvironment.SetEnv(ENV_HostEnvironment, hostEnv)
	}

	for _, eachOpt := range opts {
		eachOpt(_hostEnvironment)
	}
	return _hostEnvironment
}

func (e *hostEnvironment) GetOS() GOOS_OS {
	return GOOS_OS(e.os)
}

func (e *hostEnvironment) GetCurrentPath() string {
	return e.GetEnvString(ENV_Path)
}

func (e *hostEnvironment) GetExecFileName() string {
	return e.GetEnvString(ENV_FileName)
}
func (e *hostEnvironment) SetProduct(product string) {
	e.SetEnv(ENV_Product, product)
}

func (e *hostEnvironment) SetAppName(appName string) {
	e.SetEnv(ENV_AppName, appName)
}

func (e *hostEnvironment) SetAppVersion(appVersion string) {
	e.SetEnv(ENV_AppVersion, appVersion)
}

func (e *hostEnvironment) GetHttp() string {
	return e.GetEnvString(ENV_HTTP)
}

func (e *hostEnvironment) SetHttp(http string) {
	e.SetEnv(ENV_HTTP, http)
}

func (e *hostEnvironment) SetEnv(key string, value interface{}) {
	if len(key) <= 0 {
		return
	}
	lowerKey := strings.ToLower(key)
	if value == nil {
		delete(e.properties, lowerKey)
	} else {
		e.properties[lowerKey] = value
	}
	//处理环境变量
	s, err := castToString(value)
	if err == nil {
		e.SetOSEnv(key, s)
	}
}

func (e *hostEnvironment) GetEnv(key string) interface{} {
	if len(key) <= 0 {
		return nil
	}
	propValue, ok := _hostEnvironment.properties[strings.ToLower(key)]
	if ok {
		return propValue
	}
	//从环境变量中获取
	envValue := os.Getenv(key)
	if len(envValue) <= 0 {
		return nil
	}
	return propValue
}

func (e *hostEnvironment) GetEnvString(key string) string {
	key = strings.ToLower(key)
	sValue, ok := e.GetEnv(key).(string)
	if ok {
		return sValue
	}
	return ""
}

func (e *hostEnvironment) SetOSEnv(key string, value string) {
	if len(value) <= 0 {
		os.Unsetenv(key)
		return
	}
	os.Setenv(key, value)
}

func (e *hostEnvironment) AllKey() []string {
	keyList := make([]string, 0)
	for eachKey := range e.properties {
		keyList = append(keyList, eachKey)
	}
	sort.Strings(keyList)
	return keyList
}

// 常用的环境变量
const (
	ENV_ConsulPath = "app.abmpconsul.path"
	ENV_Product    = "app.product"
	ENV_AppName    = "app.name"
	//应用运行目录
	ENV_Path = "app.path"
	//应用名称
	ENV_FileName     = "app.filename"
	ENV_IsHostInABMP = "app.isHostInABMP"
	ENV_Description  = "app.description"
	ENV_StartTime    = "app.startTime"
	//启动app消耗的时间
	ENV_StartInterval = "app.startInterval"
	ENV_AppVersion    = "app.appVersion"
	//abmp框架版本
	ENV_FrameworkVersion = "app.frameworkVersion"
	//应用的运行环境meta key,值主要有windows,linux,supervisor,docker,systemd,other
	ENV_HostEnvironment = "app.hostEnvironment"
	ENV_HTTP            = "app.http"
	//公告主机地址
	ENV_AdvertiseHost = "app.advertiseHost"
	//健康检查的地址
	ENV_Healthcheck = "app.healthcheck"

	ENV_Plugininstaller_sourceUrl = "plugininstaller_sourceUrl"
	ENV_Plugininstaller_feedName  = "plugininstaller_feedName"
	ENV_Plugininstaller_apiKey    = "plugininstaller_apiKey"
)

var (
	builtinKeyList map[string]bool
)

func init() {
	builtinKeyList = make(map[string]bool)
	builtinKeyList[strings.ToLower(ENV_Product)] = true
	builtinKeyList[strings.ToLower(ENV_AppName)] = true
	builtinKeyList[strings.ToLower(ENV_Path)] = true
	builtinKeyList[strings.ToLower(ENV_FileName)] = true
	builtinKeyList[strings.ToLower(ENV_IsHostInABMP)] = true
	builtinKeyList[strings.ToLower(ENV_Description)] = true
	builtinKeyList[strings.ToLower(ENV_StartTime)] = true
	builtinKeyList[strings.ToLower(ENV_StartInterval)] = true
	builtinKeyList[strings.ToLower(ENV_AppVersion)] = true
	builtinKeyList[strings.ToLower(ENV_FrameworkVersion)] = true
	builtinKeyList[strings.ToLower(ENV_HostEnvironment)] = true
	builtinKeyList[strings.ToLower(ENV_HTTP)] = true
	builtinKeyList[strings.ToLower(ENV_AdvertiseHost)] = true
	builtinKeyList[strings.ToLower(ENV_Healthcheck)] = true

	builtinKeyList[strings.ToLower(ENV_Plugininstaller_sourceUrl)] = true
	builtinKeyList[strings.ToLower(ENV_Plugininstaller_feedName)] = true
	builtinKeyList[strings.ToLower(ENV_Plugininstaller_apiKey)] = true
}

// 获取指定的key是否是内置的key
func IsEnvKey(key string) bool {
	_, ok := builtinKeyList[strings.ToLower(key)]
	return ok
}

func castToString(i interface{}) (string, error) {
	i = indirectToStringerOrError(i)
	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(2), 10), nil
	case json.Number:
		return s.String(), nil
	case []byte:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirectToStringerOrError returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil) or an implementation of fmt.Stringer
// or error,
func indirectToStringerOrError(a interface{}) interface{} {
	if a == nil {
		return nil
	}

	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
