package main

import (
	"doc-precise-rag/knowledge/clients"
	"doc-precise-rag/knowledge/handler"
	processimport "doc-precise-rag/knowledge/process/import"
	"doc-precise-rag/knowledge/repository"
	"doc-precise-rag/knowledge/router"
	"doc-precise-rag/knowledge/service"
)

func main() {

	defer clients.DB.Close()
	commonSvc := service.NewCommonService(clients.CosClient)
	commonHandler := handler.NewCommonHandler(commonSvc)
	importGraphCompiledInstance := processimport.ImportGraphCompiledInstance
	docOriginFileRepository := repository.NewDocOriginFileRepository(clients.DB, clients.SnowflakeClient)
	docOriginFileService := service.NewDocOriginFileService(docOriginFileRepository, importGraphCompiledInstance)
	docOriginFileHandler := handler.NewDocHandler(docOriginFileService)

	h := &router.Handler{
		CommonHandler: commonHandler,
		DocHandler:    docOriginFileHandler,
	}
	r := router.NewRouter(h)
	r.Run(":8080")

}
