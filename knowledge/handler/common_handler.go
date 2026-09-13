package handler

import (
	"errors"
	"net/http"

	"doc-precise-rag/knowledge/entity"
	"doc-precise-rag/knowledge/service"

	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	CommonService *service.CommonService
}

func NewCommonHandler(commonService *service.CommonService) *CommonHandler {
	return &CommonHandler{
		CommonService: commonService,
	}
}

type getUploadURLRequest struct {
	Filename string `json:"filename"`
}

func (h *CommonHandler) GetUploadURL(ctx *gin.Context) {
	var req getUploadURLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "message": "filename 不能为空", "success": false})
		return
	}

	result, err := h.CommonService.GetUploadURL(ctx.Request.Context(), req.Filename)
	if err != nil {
		status := http.StatusOK
		code := 500
		if errors.Is(err, service.ErrInvalidUpload) {
			status = http.StatusOK
			code = 400
		}
		ctx.JSON(status, gin.H{"code": code, "message": err.Error(), "success": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": result, "success": true})
}

func (h *CommonHandler) GetTempFileUrlByObjectKey(ctx *gin.Context) {
	var req entity.GetTempFileUrlByObjectKeyRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 400, "message": "objectKey 不能为空", "success": false})
		return
	}
	url, err := h.CommonService.GetTempFileUrlByObjectKey(ctx.Request.Context(), req.ObjectKey)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error(), "success": false})
		return
	}
	ctx.JSON(http.StatusOK, entity.BaseResp[*entity.GetTempFileUrlByObjectKeyResp]{
		Success: true,
		Code:    entity.Ptr(int64(http.StatusOK)),
		Message: entity.Ptr("获取成功"),
		Data: &entity.GetTempFileUrlByObjectKeyResp{
			URL: url,
		},
	})
}
