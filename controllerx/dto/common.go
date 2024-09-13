package dto

// import "go.mongodb.org/mongo-driver/bson/primitive"

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   //关键字

	OrderBy string `json:"orderBy,omitempty" form:"orderBy"`
}

type Order struct {
	Field     string `json:"field"`
	Ascending bool   `json:"ascending,omitempty"`
}

func (pageInfo *PageInfo) Normalize() {
	if pageInfo.PageSize <= 0 {
		//使用默认值
		pageInfo.PageSize = 10
	}
	if pageInfo.Page <= 0 {
		//默认设置为第一页
		pageInfo.Page = 1
	}
}

// GetById Find by id structure
type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

// GetAuthorityId Get role by id structure
type GetAuthorityId struct {
	AuthorityId string `json:"authorityId" form:"authorityId"` // 角色ID
}

type Empty struct{}

type UUIDInput struct {
	Id string `json:"id" form:"id"`
}

type ObjectIdInput struct {
	ObjectId string `json:"objectId" form:"objectId"`
}

//分页结果页
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

//列表结果页
type ListResult struct {
	List interface{} `json:"list"`
}
