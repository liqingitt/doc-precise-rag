package entity

import (
	"doc-precise-rag/knowledge/module"
	"time"
)

type QueryDocOriginFileResp struct {
	Id                  *int64     `json:"id,string"`
	ObjectKey           *string    `json:"object_key"`
	DocTitle            *string    `json:"doc_title"`
	CreateTime          *time.Time `json:"create_time"`
	UpdateTime          *time.Time `json:"update_time"`
	AnalyzeDocObjectKey *string    `json:"analyze_doc_object_key"`
}

func NewQueryDocOriginFileRespByModule(module *module.DocOriginFileModule) *QueryDocOriginFileResp {
	return &QueryDocOriginFileResp{
		Id:                  module.Id,
		ObjectKey:           module.ObjectKey,
		DocTitle:            module.DocTitle,
		CreateTime:          module.CreateTime,
		UpdateTime:          module.UpdateTime,
		AnalyzeDocObjectKey: module.AnalyzeDocObjectKey,
	}
}
