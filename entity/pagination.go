package entity

var PaginationDefaultPage = 1
var PaginationDefaultSize = 10

type Pagination struct {
	Page int `form:"page" url:"page"`
	Size int `form:"size" url:"size"`
}

func (p *Pagination) IsZero() (ok bool) {
	return p.Page == 0 && p.Size == 0
}

func (p *Pagination) IsDefault() (ok bool) {
	return p.Page == PaginationDefaultPage &&
		p.Size == PaginationDefaultSize
}
