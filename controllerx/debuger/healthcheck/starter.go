package healthcheck

import "github.com/shanluzhineng/fwpkg/app"

func init() {
	app.RegisterStartupAction(healthcheckStartup)
}
