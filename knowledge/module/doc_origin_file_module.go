package module

import "time"

type DocOriginFileModule struct {
	Id         *int64     `db:"id"`
	ObjectKey  *string    `db:"object_key"`
	DocTitle   *string    `db:"doc_title"`
	CreateTime *time.Time `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
}
