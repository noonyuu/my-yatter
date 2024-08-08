package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"

	"github.com/jmoiron/sqlx"
)

type Relationship interface {
	FollowUser(ctx context.Context, tx *sqlx.Tx, followerID, followeeID *object.Account) error
}
