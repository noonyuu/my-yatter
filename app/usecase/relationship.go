package usecase

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type Relationship interface {
	FollowUser(ctx context.Context, followerID, followeeID *object.Account) error
	GetRelationships(ctx context.Context, myAcc *object.Account, otherAcc []*object.Account) ([]*object.Relationship, error)
}

type relationship struct {
	db *sqlx.DB
	rr repository.Relationship
	ar repository.Account
}

var _ Status = (*status)(nil)

func NewRelationship(db *sqlx.DB, rr repository.Relationship, ar repository.Account) *relationship {
	return &relationship{
		db: db,
		rr: rr,
		ar: ar,
	}
}

func (r *relationship) FollowUser(ctx context.Context, follower, followee *object.Account) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	err = r.rr.FollowUser(ctx, tx, follower, followee)
	if err != nil {
		return err
	}

	return nil
}

func (r *relationship) GetRelationships(ctx context.Context, myAcc *object.Account, otherAcc []*object.Account) ([]*object.Relationship, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil ,err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	relationships, err := r.rr.GetRelationship(ctx, myAcc, otherAcc)
	if err != nil {
		return nil, err
	}

	return relationships, nil
}
