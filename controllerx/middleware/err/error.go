package err

import (
	"encoding/json"
	"strings"

	"github.com/kataras/iris/v12/context"

	"github.com/shanluzhineng/fwpkg/controllerx/responsex"
)

type errWrapperMiddleware struct {
}

func New() context.Handler {
	v := &errWrapperMiddleware{}
	return v.ServeHTTP
}

func (v *errWrapperMiddleware) ServeHTTP(ctx *context.Context) {
	ctx.Record()
	ctx.Next()

	responseData := ctx.Recorder().Body()
	statusCode := ctx.GetStatusCode()
	if context.StatusCodeNotSuccessful(statusCode) && !v.responseIsIgnore(responseData) {
		ctx.Recorder().ResetBody()
		err := ctx.GetErr()
		ctx.StopWithJSON(statusCode, responsex.NewErrorResponse(func(br *responsex.BaseResponse) {
			if err != nil {
				br.SetMessage(err.Error())
			} else {
				br.SetMessage(string(responseData))
			}
		}))
	}
	respHeader := ctx.ResponseWriter().Header()
	ctype := respHeader.Get("Content-Type")
	if strings.Contains(ctype, "text/plain") || strings.Contains(ctype, "application/json") && !strings.Contains(ctx.Path(), "captcha") {
		// log.Logger.Warn(fmt.Sprintf("ServeHTTP >> resp: %s, ctype: %s", strings.Trim(string(responseData), "\n"), ctype))
		return
	}
	// log.Logger.Debug(fmt.Sprintf("ServeHTTP >> resp:: %s url: %s, ctype: %s", ctx.Request().Method, ctx.RequestPath(false), ctype))
}

func (v *errWrapperMiddleware) responseIsIgnore(responseData []byte) bool {
	if len(responseData) <= 0 {
		return false
	}
	baseResponse := &responsex.BaseResponse{}
	marshalErr := json.Unmarshal(responseData, baseResponse)
	return marshalErr == nil
}
