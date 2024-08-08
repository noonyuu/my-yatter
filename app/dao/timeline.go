package dao

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type timeline struct {
	db *sqlx.DB
}

func NewTimeline(db *sqlx.DB) repository.Timeline {
	return &timeline{db: db}
}

// GetPublicTimeline implements repository.Status.
func (t *timeline) GetPublicTimeline(ctx context.Context, maxId int, sinceId int, limit int) ([]*object.Status, error) {
	mainQuery := "select * from status "
	query := mainQuery + "ORDER BY id DESC LIMIT ?"
	// クエリパラメータとして使用されるlimitを追加する
	args := []interface{}{limit}
	if sinceId != 0 && maxId != 0 {
		query = mainQuery + "where id >= ? AND id <= ? ORDER BY id DESC LIMIT ?"
		args = []interface{}{sinceId, maxId, limit}
	}
	if maxId == 0 {
		query = mainQuery + "where id <= ? ORDER BY id DESC LIMIT ?"
		args = []interface{}{maxId, limit}
	}
	if sinceId == 0 {
		query = mainQuery + "where id >= ? ORDER BY id DESC LIMIT ?"
		args = []interface{}{sinceId, limit}
	}

	// func (c *Conn) QueryxContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	rows, err := t.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	// rowsをクローズする
	defer rows.Close()

	var status []*object.Status
	for rows.Next() {
		entity := new(object.Status)
		if err := rows.StructScan(entity); err != nil {
			return nil, err
		}
		status = append(status, entity)
	}
	return status, nil
}