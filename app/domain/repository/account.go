package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"

	"github.com/jmoiron/sqlx"
)

type Account interface {
	// Fetch account which has specified username
	FindByUsername(ctx context.Context, username string) (*object.Account, error)
	FindAccountByID(ctx context.Context, id int) (*object.Account, error)
	// TODO: Add Other APIs
	Create(ctx context.Context, tx *sqlx.Tx, acc *object.Account) error
	UpdateAccountCredential(ctx context.Context, tx *sqlx.Tx, account *object.Account) error
	// フォローしているアカウントを取得する
	FolloweeAccount(ctx context.Context, followee *object.Account, limit int) ([]*object.Account, error)
}
