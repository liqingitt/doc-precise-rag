package router

import (
	"doc-precise-rag/knowledge/handler"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	CommonHandler *handler.CommonHandler
	DocHandler    *handler.DocHandler
}

func NewRouter(handler *Handler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	commonApi := api.Group("/common")
	registerCommonRouter(commonApi, handler.CommonHandler)

	docApi := api.Group("/doc")
	registerDocRouter(docApi, handler.DocHandler)
	return r
}

func registerCommonRouter(api *gin.RouterGroup, handler *handler.CommonHandler) {
	api.POST("/upload-url", handler.GetUploadURL)
}

func registerDocRouter(api *gin.RouterGroup, handler *handler.DocHandler) {
	api.POST("/query-doc-origin-file-list", handler.QueryDocOriginFileList)
	api.POST("/add-doc-origin-file", handler.AddDocOriginFile)
	api.POST("/analysis-doc-origin-file", handler.AnalysisDocOriginFile)
	api.POST("/delete-doc-origin-file", handler.DeleteDocOriginFile)
}
