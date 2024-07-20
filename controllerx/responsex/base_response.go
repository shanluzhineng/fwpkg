package responsex

type Response interface {
	ResponseInfo
	SetData(data interface{})
	GetData() interface{}
}

// web回应基类
type BaseResponse struct {
	BaseResponseInfo
	Data interface{} `json:"data,omitempty" schema:"HTTP response data"`
}

func (r *BaseResponse) SetData(data interface{}) {
	r.Data = data
}

func (r *BaseResponse) GetData() interface{} {
	return r.Data
}

func (r *BaseResponse) IsSuccessful() bool {
	return r.Code == Code_Successful
}

// 构建一个成功的回应
func NewSuccessResponse(opts ...func(*BaseResponse)) *BaseResponse {
	r := &BaseResponse{}
	r.SetCode(Code_Successful)
	r.Status = HttpResponseStatusOk
	r.SetMessage(HttpResponseMessageSuccess)
	if len(opts) > 0 {
		for _, eachOpt := range opts {
			eachOpt(r)
		}
	}
	return r
}

// 构建一个错误的回应
func NewErrorResponse(opts ...func(*BaseResponse)) *BaseResponse {
	r := &BaseResponse{}
	r.SetCode(Code_GeneralError)
	r.Status = HttpResponseStatusOk
	r.SetMessage(HttpResponseMessageError)
	if len(opts) > 0 {
		for _, eachOpt := range opts {
			eachOpt(r)
		}
	}
	return r
}
