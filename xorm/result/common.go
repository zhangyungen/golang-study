package result

import "zyj.com/golang-study/xorm/param"

type PageVO[T any] struct {
	Page      int   `json:"page" form:"page"`
	PageSize  int   `json:"pageSize" form:"pageSize"`
	DataSlice []T   `json:"dataSlice" form:"dataSlice"`
	Total     int64 `json:"total" form:"total"`
}

func Convert2PageVO[T any](param *param.PageParam, total int64, dataSlice []T) *PageVO[T] {
	return &PageVO[T]{
		Page:      param.Page,
		PageSize:  param.PageSize,
		DataSlice: dataSlice,
		Total:     total,
	}
}
