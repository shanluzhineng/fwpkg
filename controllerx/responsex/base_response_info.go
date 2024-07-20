package responsex

const (
	//成功
	Code_Successful = 1
	//通用错误
	Code_GeneralError = 0
)

const (
	HttpResponseStatusOk       = "ok"
	HttpResponseMessageSuccess = "success"
	HttpResponseMessageError   = "error"
)

type ResponseInfo interface {
	SetCode(code int)
	GetCode() int
	SetMessage(message string)
	GetMessage() string
	GetStatus() string
}

type BaseResponseInfo struct {
	Code    int    `json:"code" schema:"HTTP response code"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty" schema:"HTTP response message"`
}

// #region ResponseInfo Members

// 设置回应code
func (r *BaseResponseInfo) SetCode(code int) *BaseResponseInfo {
	r.Code = code
	return r
}

func (r *BaseResponseInfo) GetCode() int {
	return r.Code
}

func (r *BaseResponseInfo) SetMessage(message string) *BaseResponseInfo {
	r.Message = message
	return r
}

// 获取消息
func (r *BaseResponseInfo) GetMessage() string {
	return r.Message
}

func (r *BaseResponseInfo) GetStatus() string {
	return r.Status
}

// #endregion
