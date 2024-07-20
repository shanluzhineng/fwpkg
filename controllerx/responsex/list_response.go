package responsex

type ListResponse struct {
	BaseResponse

	Total int64 `json:"total"`
}

// 构建一个成功的回应
func NewSuccessListResponse(data interface{}, total int64, opts ...func(*ListResponse)) *ListResponse {
	r := &ListResponse{}
	r.SetCode(Code_Successful)
	r.Status = HttpResponseStatusOk
	r.SetMessage(HttpResponseMessageSuccess)
	r.Total = total
	r.SetData(data)
	if len(opts) > 0 {
		for _, eachOpt := range opts {
			eachOpt(r)
		}
	}
	return r
}

// create error list response
func NewErrorListResponse(opts ...func(*ListResponse)) *ListResponse {
	r := &ListResponse{}
	r.SetCode(Code_GeneralError)
	r.Status = HttpResponseStatusOk
	r.SetMessage(HttpResponseMessageSuccess)
	r.Total = 0
	if len(opts) > 0 {
		for _, eachOpt := range opts {
			eachOpt(r)
		}
	}
	return r
}
