package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"

	"github.com/jmoiron/sqlx"
)

type Relationship interface {
	// フォーローする
	FollowUser(ctx context.Context, tx *sqlx.Tx, followerID, followeeID *object.Account) error

	// リレーションを取得
	GetRelationship(ctx context.Context, myAcc *object.Account, otherAcc []*object.Account) ([]*object.Relationship, error)
}
