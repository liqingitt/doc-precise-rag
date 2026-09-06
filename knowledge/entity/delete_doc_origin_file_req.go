package entity

type DeleteDocOriginFileReq struct {
	Id *int64 `json:"id,string" binding:"required"`
}
