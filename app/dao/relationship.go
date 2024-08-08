package dao

import (
	"context"
	"fmt"
	"yatter-backend-go/app/domain/object"

	"github.com/jmoiron/sqlx"
)

type relationship struct {
	db *sqlx.DB
}

func NewRelationship(db *sqlx.DB) *relationship {
	return &relationship{db: db}
}

func (r *relationship) FollowUser(ctx context.Context, tx *sqlx.Tx, follower *object.Account, followee *object.Account) error {
	var count int

	err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM relationship WHERE follower_id = ? AND followee_id = ?", follower.ID, followee.ID).Scan(&count)
	if err != nil {
		return err
	}

	// まだフォローしていない場合
	if count == 0 {
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

		_, err = r.db.ExecContext(ctx, "INSERT INTO relationship (follower_id, followee_id) VALUES (?, ?)", follower.ID, followee.ID)
		if err != nil {
			return fmt.Errorf("failed to insert relationship: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		return nil
	}

	return fmt.Errorf("already following")
}
