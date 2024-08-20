package dao

import (
	"context"
	"fmt"
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
func (t *timeline) GetPublicTimeline(ctx context.Context, sinceId, limit int64) ([]*object.Status, error) {
	query := "SELECT s.id FROM status AS s WHERE s.id > ? ORDER BY s.id LIMIT ?"

	rows, err := t.db.QueryxContext(ctx, query, sinceId, limit)
	if err != nil {
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()
	var statusesId []int64
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		statusesId = append(statusesId, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	sta := NewStatus(t.db)
	tls := make([]*object.Status, 0)
	for _, id := range statusesId {
		status, err := sta.FindByStatus(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("error getting status: %v", err)
		}
		tls = append(tls, status)
	}
	return tls, nil
}

func (t *timeline) GetHomeTimeline(ctx context.Context, id, maxId, sinceId, limit int64) ([]*object.Status, error) {
	query := "SELECT s.id FROM status AS s JOIN (SELECT followee_id FROM relationship WHERE follower_id = ?) AS f ON f.followee_id = s.account_id WHERE s.account_id <= ? AND s.account_id >= ? LIMIT ?"

	rows, err := t.db.QueryxContext(ctx, query, id, maxId, sinceId, limit)
	if err != nil {
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()
	var statusesId []int64
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		statusesId = append(statusesId, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	sta := NewStatus(t.db)
	tls := make([]*object.Status, 0)
	for _, id := range statusesId {
		status, err := sta.FindByStatus(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("error getting status: %v", err)
		}
		tls = append(tls, status)
	}
	return tls, nil
}
// SELECT s.id FROM status AS s JOIN (SELECT followee_id FROM relationship WHERE follower_id = 2) AS f ON f.followee_id = s.account_id WHERE s.account_id >= 1 AND s.account_id <= 3 LIMIT 5;
