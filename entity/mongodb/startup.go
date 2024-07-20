package mongodb

import (
	"github.com/shanluzhineng/fwpkg/app/web"
)

func init() {
	web.ConfigureService(initMongodbConfigurator)
}
