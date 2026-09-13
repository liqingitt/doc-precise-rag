package repository

import (
	"context"
	"database/sql"
	"doc-precise-rag/knowledge/entity"
	"doc-precise-rag/knowledge/module"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/bwmarrin/snowflake"
	"github.com/jmoiron/sqlx"
)

type DocOriginFileRepository struct {
	db              *sqlx.DB
	snowflakeClient *snowflake.Node
}

func (r *DocOriginFileRepository) FindDocOriginFileById(ctx context.Context, id int64) (*module.DocOriginFileModule, error) {
	sqlStr, args, err := squirrel.Select("id", "object_key", "doc_title", "create_time", "update_time").From("doc_origin_file").Where("id = ?", id).ToSql()
	if err != nil {
		return nil, err
	}
	var module module.DocOriginFileModule
	err = r.db.GetContext(ctx, &module, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	return &module, nil
}

func NewDocOriginFileRepository(db *sqlx.DB, snowflakeClient *snowflake.Node) *DocOriginFileRepository {
	return &DocOriginFileRepository{db: db, snowflakeClient: snowflakeClient}
}

func (r *DocOriginFileRepository) QueryDocOriginFileListCount(ctx context.Context) (int64, error) {
	totalSql, args, err := squirrel.
		Select("COUNT(*) as total").
		From("doc_origin_file").
		ToSql()
	if err != nil {
		return 0, err
	}

	var total int64
	err = r.db.GetContext(ctx, &total, totalSql, args...)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *DocOriginFileRepository) QueryDocOriginFileList(ctx context.Context, page int64, pageSize int64) ([]*module.DocOriginFileModule, error) {
	offset := entity.GetPageOffset(&page, &pageSize)
	limit := pageSize

	sql, args, err := squirrel.
		Select("id", "object_key", "analyze_doc_object_key", "doc_title", "create_time", "update_time").
		From("doc_origin_file").OrderBy("create_time DESC").
		Limit(uint64(limit)).Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, err
	}
	var rows []*module.DocOriginFileModule

	err = r.db.SelectContext(ctx, &rows, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil

}

func (r *DocOriginFileRepository) AddDocOriginFile(ctx context.Context, module *module.DocOriginFileModule) (int64, error) {
	id := r.snowflakeClient.Generate().Int64()
	sql, args, err := squirrel.
		Insert("doc_origin_file").
		Columns("id", "object_key", "doc_title").
		Values(id, module.ObjectKey, module.DocTitle).
		ToSql()
	if err != nil {
		return 0, err
	}
	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *DocOriginFileRepository) FindDocOriginFileByDocTitle(ctx context.Context, DocTitle string) (*module.DocOriginFileModule, error) {
	sqlStr, args, err := squirrel.
		Select("id", "object_key", "doc_title", "create_time", "update_time").
		From("doc_origin_file").
		Where("doc_title = ?", DocTitle).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	var module module.DocOriginFileModule
	err = r.db.GetContext(ctx, &module, sqlStr, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *DocOriginFileRepository) UpdateDocOriginFileById(ctx context.Context, id int64, module *module.DocOriginFileModule) error {
	sql, args, err := squirrel.
		Update("doc_origin_file").
		Where("id = ?", id).
		SetMap(map[string]any{
			"object_key":             module.ObjectKey,
			"doc_title":              module.DocTitle,
			"analyze_doc_object_key": module.AnalyzeDocObjectKey,
		}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *DocOriginFileRepository) DeleteDocOriginFileById(ctx context.Context, id int64) error {
	sql, args, err := squirrel.
		Delete("doc_origin_file").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
