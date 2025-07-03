package dto

// Pagination 分页参数
type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// Offset 获取偏移量
func (p Pagination) Offset() int {
	return (p._Page() - 1) * p._PageSize()
}

// Limit 获取限制数量
func (p Pagination) Limit() int {
	return p._PageSize()
}

func (p Pagination) _Page() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p Pagination) _PageSize() int {
	if p.PageSize <= 0 {
		p.PageSize = 5
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

func (p Pagination) Value() *Pagination {
	return &Pagination{
		Page:     p._Page(),
		PageSize: p._PageSize(),
	}
}
