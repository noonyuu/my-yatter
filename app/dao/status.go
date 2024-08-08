package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type status struct {
	db *sqlx.DB
}

type Status struct {
	ID        int64
	AccountID int64
	Content   string
	CreateAt  time.Time
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
		return fmt.Errorf("failed to delete status from db: %w", err)
	}

	return nil
}

func (s *status) FindByStatus(ctx context.Context, id int64) (*object.Status, error) {
	// sta := &Status{}
	entity := new(object.Status)
	// 投稿情報を取得
	query := "SELECT * FROM status WHERE id = ?"
	err := s.db.QueryRowxContext(ctx, query, id).StructScan(entity)
	if err != nil {
		return nil, fmt.Errorf("error getting status: %v", err)
	}

	// 投稿したユーザー情報を取得
	acc := NewAccount(s.db)
	account, err := acc.FindAccountByID(ctx, entity.AccountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %v", err)
	}
	entity.AccountID = int(account.ID)

	return entity, nil
}
