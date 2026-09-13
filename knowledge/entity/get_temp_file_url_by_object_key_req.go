package entity

type GetTempFileUrlByObjectKeyRequest struct {
	ObjectKey string `form:"object_key" binding:"required"`
}
