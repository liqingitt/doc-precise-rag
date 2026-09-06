package handler

import (
	"errors"
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "filename 不能为空", "success": false})
		return
	}

	result, err := h.CommonService.GetUploadURL(ctx.Request.Context(), req.Filename)
	if err != nil {
		status := http.StatusInternalServerError
		code := 500
		if errors.Is(err, service.ErrInvalidUpload) {
			status = http.StatusBadRequest
			code = 400
		}
		ctx.JSON(status, gin.H{"code": code, "message": err.Error(), "success": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": result, "success": true})
}
