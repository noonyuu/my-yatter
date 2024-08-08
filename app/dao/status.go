package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type status struct {
	db *sqlx.DB
}

func NewStatus(db *sqlx.DB) *status {
	return &status{
		db: db,
	}
}

var _ repository.Status = (*status)(nil)

func (s *status) Create(ctx context.Context, tx *sqlx.Tx, status *object.Status) (*object.Status, error) {
	_, err := s.db.Exec("insert into status (content, account_id, create_at) values (?, ?, ?)", status.Content, status.AccountID, status.CreatedAt)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (s *status) FindStatusByID(ctx context.Context, id int) (*object.Status, error) {
	entity := new(object.Status)
	err := s.db.QueryRowxContext(ctx, "select * from status where id = ?", id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find status from db: %w", err)
	}
	return entity, nil
}

func (s *status) DeleteStatus(ctx context.Context, tx *sqlx.Tx, id int) error {
	entity, err := s.FindStatusByID(ctx, id)
	if err != nil {
		return err
	}
	if entity == nil {
		return nil
	}

	_, err = tx.Exec("delete from status where id = ?", id)
	if err != nil {
		return  fmt.Errorf("failed to delete status from db: %w", err)
	}

	return nil
}