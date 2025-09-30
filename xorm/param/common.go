package param

type PageParam struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
	OrderField string `json:"orderField" form:"orderField"`
	OrderAsc   bool   `json:"orderAsc" form:"orderAsc"`
}

type Id[T any] struct {
	Id T `json:"id" form:"id" binding:"required"`
}
