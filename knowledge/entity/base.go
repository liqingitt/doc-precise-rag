package entity

type BasePageReq struct {
	Page     *int64 `json:"page" binding:"required"`
	PageSize *int64 `json:"page_size" binding:"required"`
}

type BaseResp[T any] struct {
	Success bool    `json:"success"`
	Code    *int64  `json:"code"`
	Message *string `json:"message"`
	Data    T       `json:"data"`
}

type BasePageResp[T any] struct {
	Total *int64 `json:"total"`
	List  []T    `json:"list"`
}

func Ptr[T any](v T) *T {
	return &v
}

func GetPageOffset(page *int64, pageSize *int64) int64 {
	if page == nil || pageSize == nil {
		return 0
	}
	return (*page - 1) * (*pageSize)
}
