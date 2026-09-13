package module

import "time"

type DocOriginFileModule struct {
	Id                  *int64     `db:"id"`
	ObjectKey           *string    `db:"object_key"`
	DocTitle            *string    `db:"doc_title"`
	AnalyzeDocObjectKey *string    `db:"analyze_doc_object_key"`
	CreateTime          *time.Time `db:"create_time"`
	UpdateTime          *time.Time `db:"update_time"`
}
