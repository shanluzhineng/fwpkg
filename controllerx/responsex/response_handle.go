package responsex

import (
	"net/http"

	"github.com/kataras/iris/v12"
)

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

func HandleSuccess(ctx iris.Context) {
	ctx.StopWithJSON(http.StatusOK, NewSuccessResponse())
}

func HandleSuccessWithData(ctx iris.Context, data interface{}) {
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
