package cors

import (
	"github.com/shanluzhineng/configurationx/options/web"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

func UseCors(apiBuilder *router.APIBuilder, opts web.CORS) {
	options := allowedAllOptions()
	if opts.Mode == web.CorsMode_Whitelist {
		options.AllowedOrigins = opts.GetAllowedOrigins()
	}
	cors := cors.New(options)
	apiBuilder.UseRouter(cors)
}

func allowedAllOptions() cors.Options {
	return cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{iris.MethodPost,
			iris.MethodGet,
			iris.MethodOptions,
			iris.MethodDelete,
			iris.MethodOptions,
			iris.MethodPut},
		AllowedHeaders: []string{"Content-Type",
			"AccessToken",
			"X-CSRF-Token",
			"Authorization",
			"Token",
			"X-Token",
			"X-User-Id"},
		ExposedHeaders: []string{"Content-Length",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"New-Token",
			"New-Expires-At"},
		AllowCredentials: true,
	}
}
