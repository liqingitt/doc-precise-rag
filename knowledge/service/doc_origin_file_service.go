package service

import (
	"context"
	"doc-precise-rag/knowledge/entity"
	"doc-precise-rag/knowledge/module"
	"doc-precise-rag/knowledge/repository"
)

type DocOriginFileService struct {
	docOriginFileRepository *repository.DocOriginFileRepository
}

func NewDocOriginFileService(docOriginFileRepository *repository.DocOriginFileRepository) *DocOriginFileService {
	return &DocOriginFileService{docOriginFileRepository: docOriginFileRepository}
}

func (s *DocOriginFileService) QueryDocOriginFileList(ctx context.Context, queryReq *entity.QueryDocOriginFileReq) ([]*entity.QueryDocOriginFileResp, int64, error) {
	total, err := s.docOriginFileRepository.QueryDocOriginFileListCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	rows, err := s.docOriginFileRepository.QueryDocOriginFileList(ctx, *queryReq.Page, *queryReq.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var list []*entity.QueryDocOriginFileResp
	for _, row := range rows {
		list = append(list, entity.NewQueryDocOriginFileRespByModule(row))
	}
	return list, total, nil
}

func (s *DocOriginFileService) AddDocOriginFile(ctx context.Context, addReq *entity.AddDocOriginFileReq) (int64, error) {
	module := &module.DocOriginFileModule{
		ObjectKey: addReq.ObjectKey,
		DocTitle:  addReq.DocTitle,
	}

	docOriginFile, err := s.docOriginFileRepository.FindDocOriginFileByDocTitle(ctx, *module.DocTitle)
	if err != nil {
		return 0, err
	}
	if docOriginFile == nil {
		return s.docOriginFileRepository.AddDocOriginFile(ctx, module)
	}

	err = s.docOriginFileRepository.UpdateDocOriginFileById(ctx, *docOriginFile.Id, module)
	if err != nil {
		return 0, err
	}
	return *docOriginFile.Id, nil
}

func (s *DocOriginFileService) DeleteDocOriginFile(ctx context.Context, id int64) error {
	return s.docOriginFileRepository.DeleteDocOriginFileById(ctx, id)
}
