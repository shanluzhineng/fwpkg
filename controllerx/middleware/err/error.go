package err

import (
	"encoding/json"
	"fmt"

	"github.com/kataras/iris/v12/context"

	"github.com/shanluzhineng/fwpkg/controllerx/responsex"
	"github.com/shanluzhineng/fwpkg/system/log"
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
	log.Logger.Warn(fmt.Sprintf("ServeHTTP >> resp: %s", string(responseData)))
}

func (v *errWrapperMiddleware) responseIsIgnore(responseData []byte) bool {
	if len(responseData) <= 0 {
		return false
	}
	baseResponse := &responsex.BaseResponse{}
	marshalErr := json.Unmarshal(responseData, baseResponse)
	return marshalErr == nil
}
