package responsex

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// type Response struct {
// 	Code int         `json:"code"`
// 	Data interface{} `json:"data"`
// 	Msg  string      `json:"msg"`
// }

const (
	Msg_SUCCESS = "Success"
)

const (
	NotFound = http.StatusNoContent
	ERROR    = 7
	SUCCESS  = 0
	Created  = 201
	//兼容旧系统或其它系统的成功ID
	SUCCESSV0 = 1
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	r := &BaseResponse{}
	r.SetCode(code)
	r.Status = HttpResponseStatusOk
	r.SetMessage(msg)
	r.SetData(data)
	c.JSON(http.StatusOK, r)
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "success", c)
}

func OkWithCodeAndMessage(code int, message string, c *gin.Context) {
	Result(code, map[string]interface{}{}, message, c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "success", c)
}

func OkWithCodeAndDetailed(code int, data interface{}, message string, c *gin.Context) {
	Result(code, data, message, c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "fail", c)
}

func FailWithCodeAndMessage(code int, message string, c *gin.Context) {
	Result(code, map[string]interface{}{}, message, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
