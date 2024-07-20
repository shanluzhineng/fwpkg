package system

// App is the properties of the application, it hold the base info of the application
type App struct {
	//是否运行在命令行模式
	IsRunInCli bool
	// project name
	Title string `json:"title,omitempty"`
	// project name
	Project string `json:"project,omitempty" `
	// app name
	Name string `json:"name,omitempty" `
	// app description
	Description string `json:"description,omitempty"`
	// Version
	Version string `json:"version,omitempty" default:"${APP_VERSION:v1}"`
	// TermsOfService
	TermsOfService string `json:"termsOfService,omitempty"`
}

func newApp() *App {
	return &App{
		IsRunInCli: false,
		Version:    "v1",
	}
}

func (app *App) WithName(name string) *App {
	app.Name = name
	return app
}

func (app *App) WithVersion(version string) *App {
	app.Version = version
	return app
}

func (app *App) WithTitle(title string) *App {
	app.Title = title
	return app
}

func (app *App) WithDescription(description string) *App {
	app.Description = description
	return app
}

// Logging is the properties of logging
type Logging struct {
	Level string `json:"level,omitempty" default:"info"`
}

func newLogging() *Logging {
	return &Logging{
		Level: "info",
	}
}
