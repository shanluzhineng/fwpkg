package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/shanluzhineng/fwpkg/app"
)

type Application interface {
	app.Application
}

type application struct {
	app.BaseApplication
	root Command
}

type CommandNameValue struct {
	Name    string
	Command interface{}
}

const (
	// RootCommandName the instance name of cli.rootCommand
	RootCommandName = "cli.rootCommand"
)

func NewApplication(cmd ...interface{}) Application {
	a := new(application)
	a.initialize(cmd...)
	return a
}

func (a *application) initialize(cmd ...interface{}) (err error) {
	if len(cmd) > 0 {
		app.Register(RootCommandName, cmd[0])
	}
	err = a.Initialize()
	return
}

// 构建应用运行所需的环境
func (a *application) Build() app.Application {
	//先调用基类的构建函数
	a.BaseApplication.Build()

	basename := filepath.Base(os.Args[0])
	basename = strings.ToLower(basename)
	basename = strings.TrimSuffix(basename, ".exe")

	f := a.ConfigurableFactory()
	f.SetInstance(app.ApplicationContextName, a)

	// 处理自动注入配置
	a.BuildConfigurations()

	// cli root command
	r := f.GetInstance(RootCommandName)
	var root Command
	if r != nil {
		root = r.(Command)
		Register(root)
		a.root = root
		root.EmbeddedCommand().Use = basename
	}

	a.AfterInitialization()
	return a
}

// 设置应用属性名
func (a *application) SetProperty(name string, value ...interface{}) app.Application {
	a.BaseApplication.SetProperty(name, value...)
	return a
}

func (a *application) SetAddCommandLineProperties(enabled bool) app.Application {
	a.BaseApplication.SetAddCommandLineProperties(enabled)
	return a
}

// 初始化应用
func (a *application) Initialize() error {
	return a.BaseApplication.Initialize()
}

// 运行应用
func (a *application) Run() {
	if a.root != nil {
		a.root.Exec()
	}
	a.Shutdown()
}
