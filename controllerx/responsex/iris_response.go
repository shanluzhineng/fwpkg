package responsex

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/fwpkg/utils/common"
	"github.com/shanluzhineng/fwpkg/utils/json"
)

var NoLogUriList []string

func SetNoLogUriList(noLogUri []string) {
	NoLogUriList = noLogUri
}

func HandleError(statusCode int, ctx iris.Context, err error) {
	ctx.StopWithError(statusCode, err)
	// ctx.StopWithJSON(statusCode, model.NewErrorResponse(func(br *model.BaseResponse) {
	// 	br.SetMessage(err.Error())
	// }))
}

// handle StatusBadRequest
func HandleErrorBadRequest(ctx iris.Context, err error) {
	HandleError(http.StatusBadRequest, ctx, err)
}

func HandleErrorUnauthorized(ctx iris.Context, err error) {
	HandleError(http.StatusUnauthorized, ctx, err)
}

func HandleErrorNotFound(ctx iris.Context, err error) {
	HandleError(http.StatusNotFound, ctx, err)
}

func HandleErrorInternalServerError(ctx iris.Context, err error) {
	HandleError(http.StatusInternalServerError, ctx, err)
}

func HandleFailWithMsg(ctx iris.Context, errMsg string) {
	debugLogMsg(ctx, errMsg)
	ctx.StopWithJSON(http.StatusOK, NewErrorResponse(func(br *BaseResponse) {
		br.SetMessage(errMsg)
	}))
}

func HandleSuccess(ctx iris.Context) {
	ctx.StopWithJSON(http.StatusOK, NewSuccessResponse())
}

func HandleSuccessWithData(ctx iris.Context, data interface{}) {
	debugLogMsg(ctx, data)
	ctx.StopWithJSON(http.StatusOK, NewSuccessResponse(func(br *BaseResponse) {
		br.SetData(data)
	}))
}

func HandleSuccessWithListData(ctx iris.Context, data interface{}, total int64) {
	ctx.StopWithJSON(http.StatusOK, NewSuccessListResponse(data, total))
}

func HandlerBinary(ctx iris.Context, data []byte) (int, error) {
	return ctx.Binary(data)
}

func debugLogMsg(ctx iris.Context, data interface{}) {
	path := ctx.Path() // path like: /shopping/captcha
	shouldLog := true
	for _, uri := range NoLogUriList {
		if strings.HasPrefix(ctx.Path(), uri) {
			shouldLog = false
			break
		}
	}
	if shouldLog {
		log.Logger.Info(fmt.Sprintf("[%s] resp data: %v, path: %s", common.GetCallerName(3), json.ObjectToJson(data), path))
	}
	return
}
