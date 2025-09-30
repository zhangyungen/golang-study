package param

type PageParam struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
	OrderField string `json:"order_field" form:"order_field"`
	OrderAsc   bool   `json:"order_asc" form:"order_asc"`
}

type Id[T any] struct {
	Id T `json:"id" form:"id" binding:"required"`
}
