package handler

import (
	"doc-precise-rag/knowledge/entity"
	"doc-precise-rag/knowledge/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DocHandler struct {
	DocOriginFileService *service.DocOriginFileService
}

func NewDocHandler(docOriginFileService *service.DocOriginFileService) *DocHandler {
	return &DocHandler{
		DocOriginFileService: docOriginFileService,
	}
}

func (h *DocHandler) QueryDocOriginFileList(ctx *gin.Context) {
	var queryReq entity.QueryDocOriginFileReq
	if err := ctx.ShouldBindJSON(&queryReq); err != nil {
		ctx.JSON(http.StatusOK, entity.BaseResp[any]{
			Success: false,
			Code:    entity.Ptr(int64(http.StatusBadRequest)),
			Message: entity.Ptr(err.Error()),
			Data:    nil,
		})
		return
	}

	list, total, err := h.DocOriginFileService.QueryDocOriginFileList(ctx, &queryReq)
	if err != nil {
		ctx.JSON(http.StatusOK, entity.BaseResp[any]{
			Success: false,
			Code:    entity.Ptr(int64(http.StatusBadRequest)),
			Message: entity.Ptr(err.Error()),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, entity.BaseResp[*entity.BasePageResp[*entity.QueryDocOriginFileResp]]{
		Success: true,
		Code:    entity.Ptr(int64(http.StatusOK)),
		Message: entity.Ptr("查询成功"),
		Data: &entity.BasePageResp[*entity.QueryDocOriginFileResp]{
			Total: entity.Ptr(total),
			List:  list,
		},
	})
}

func (h *DocHandler) AddDocOriginFile(ctx *gin.Context) {
	var addReq entity.AddDocOriginFileReq
	if err := ctx.ShouldBindJSON(&addReq); err != nil {
		ctx.JSON(http.StatusOK, entity.BaseResp[any]{
			Success: false,
			Code:    entity.Ptr(int64(http.StatusBadRequest)),
			Message: entity.Ptr(err.Error()),
			Data:    nil,
		})
		return
	}
	id, err := h.DocOriginFileService.AddDocOriginFile(ctx, &addReq)
	if err != nil {
		ctx.JSON(http.StatusOK, entity.BaseResp[any]{
			Success: false,
			Code:    entity.Ptr(int64(http.StatusBadRequest)),
			Message: entity.Ptr(err.Error()),
			Data:    nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, entity.BaseResp[string]{
		Success: true,
		Code:    entity.Ptr(int64(http.StatusOK)),
		Message: entity.Ptr("添加成功"),
		Data:    strconv.FormatInt(id, 10),
	})
}

func (h *DocHandler) AnalysisDocOriginFile(ctx *gin.Context) {
	var analysisReq entity.AnalysisDocOriginFileReq
	if err := ctx.ShouldBindJSON(&analysisReq); err != nil {
		ctx.JSON(http.StatusOK, entity.BaseResp[any]{
			Success: false,
			Code:    entity.Ptr(int64(http.StatusBadRequest)),
			Message: entity.Ptr(err.Error()),
			Data:    nil,
		})
		return
	}

	err := h.DocOriginFileService.AnalysisDocOriginFile(ctx.Request.Context(), *analysisReq.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, entity.BaseResp[any]{
			Success: false,
			Code:    entity.Ptr(int64(http.StatusBadRequest)),
			Message: entity.Ptr(err.Error()),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, entity.BaseResp[any]{
		Success: true,
		Code:    entity.Ptr(int64(http.StatusOK)),
		Message: entity.Ptr("解析成功"),
		Data:    nil,
	})
}

func (h *DocHandler) DeleteDocOriginFile(ctx *gin.Context) {
	var deleteReq entity.DeleteDocOriginFileReq
	if err := ctx.ShouldBindJSON(&deleteReq); err != nil {
		ctx.JSON(http.StatusOK, entity.BaseResp[any]{
			Success: false,
			Code:    entity.Ptr(int64(http.StatusBadRequest)),
			Message: entity.Ptr(err.Error()),
			Data:    nil,
		})
		return
	}
	err := h.DocOriginFileService.DeleteDocOriginFile(ctx, *deleteReq.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, entity.BaseResp[any]{
			Success: false,
			Code:    entity.Ptr(int64(http.StatusBadRequest)),
			Message: entity.Ptr(err.Error()),
			Data:    nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, entity.BaseResp[any]{
		Success: true,
		Code:    entity.Ptr(int64(http.StatusOK)),
		Message: entity.Ptr("删除成功"),
		Data:    nil,
	})
}
