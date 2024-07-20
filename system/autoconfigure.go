package system

// Configuration is the system configuration
type Configuration struct {
	App     *App
	Logging *Logging
}

func NewConfiguration() *Configuration {
	return &Configuration{
		App:     newApp(),
		Logging: newLogging(),
	}
}
