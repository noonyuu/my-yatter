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

type (
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

var _ repository.Account = (*account)(nil)

// Create accout repository
func NewAccount(db *sqlx.DB) *account {
	return &account{db: db}
}

// FindByUsername : ユーザ名からユーザを取得
func (a *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := a.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find account from db: %w", err)
	}
	err = a.FollowerAndFollowingCount(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (a *account) Create(ctx context.Context, tx *sqlx.Tx, acc *object.Account) error {
	_, err := a.db.Exec("insert into account (username, password_hash, display_name, avatar, header, note, create_at) values (?, ?, ?, ?, ?, ?, ?)",
		acc.Username, acc.PasswordHash, acc.DisplayName, acc.Avatar, acc.Header, acc.Note, acc.CreateAt)
	if err != nil {
		return fmt.Errorf("failed to insert account: %w", err)
	}

	return nil
}

func (a *account) FindAccountByID(ctx context.Context, id int) (*object.Account, error) {
	entity := new(object.Account)
	err := a.db.QueryRowxContext(ctx, "select * from account where id = ?", id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find account from db: %w", err)
	}
	err = a.FollowerAndFollowingCount(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (a *account) UpdateAccountCredential(ctx context.Context, x *sqlx.Tx, account *object.Account) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	_, err = tx.ExecContext(ctx, "update account set display_name = ?, note = ?, avatar = ?, header = ? where id = ?",
		account.DisplayName, account.Note, account.Avatar, account.Header, account.ID)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	return nil
}

func (a *account) FolloweeAccount(ctx context.Context, follower *object.Account, limit int) ([]*object.Account, error) {
	query := "select * from account where id in (select followee_id from Relationship where follower_id = ?) ORDER BY id DESC LIMIT ?"
	rows, err := a.db.QueryxContext(ctx, query, follower.ID, limit)
	if err != nil {
		return nil, err
	}

	var accounts []*object.Account
	for rows.Next() {
		entity := new(object.Account)
		if err := rows.StructScan(entity); err != nil {
			return nil, err
		}
		accounts = append(accounts, entity)
	}
	for _, account := range accounts {
		err := a.FollowerAndFollowingCount(ctx, account)
		if err != nil {
			return nil, err
		}
	}
	return accounts, nil
}

func (a *account) FollowerAndFollowingCount(ctx context.Context, entity *object.Account) error {
	err := a.db.QueryRowxContext(ctx, "select count(*) from relationship where followee_id = ?", entity.ID).Scan(&entity.FolloweeCount)
	if err != nil {
		return err
	}
	err = a.db.QueryRowxContext(ctx, "select count(*) from relationship where follower_id = ?", entity.ID).Scan(&entity.FollowerCount)
	if err != nil {
		return err
	}
	return nil
}
