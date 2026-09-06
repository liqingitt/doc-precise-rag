package entity

type AddDocOriginFileReq struct {
	ObjectKey *string `json:"object_key" binding:"required"`
	DocTitle  *string `json:"doc_title" binding:"required"`
}
