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

	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM relationship WHERE follower_id = ? AND followee_id = ?", follower.ID, followee.ID).Scan(&count)
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

func (r *relationship) UnFollowUser(ctx context.Context, tx *sqlx.Tx, follower *object.Account, followee *object.Account) error {
	var count int

	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM relationship WHERE follower_id = ? AND followee_id = ?", follower.ID, followee.ID).Scan(&count)
	if err != nil {
		return err
	}

	// まだフォローしていない場合
	if count == 1 {
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

		_, err = r.db.ExecContext(ctx, "DELETE FROM relationship WHERE follower_id = ? AND followee_id = ?", follower.ID, followee.ID)
		if err != nil {
			return fmt.Errorf("failed to delete relationship: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		return nil
	}
	return fmt.Errorf("not following")
}

func (r *relationship) GetRelationship(ctx context.Context, myAcc *object.Account, otherAcc []*object.Account) ([]*object.Relationship, error) {
	results := make([]*object.Relationship, 0, len(otherAcc))

	for _, other := range otherAcc {
		rel := &object.Relationship{ID: other.ID}

		var err error
		rel.FollowedBy, rel.Following, err = r.checkRelationship(ctx, myAcc.ID, other.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}

		results = append(results, rel)
	}

	return results, nil
}

// 関係性をチェックする
func (r *relationship) checkRelationship(ctx context.Context, myID, otherID int64) (followedBy, following bool, err error) {
	var count int

	// FollowedByをチェック
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM relationship WHERE followee_id = ? AND follower_id = ?", myID, otherID).Scan(&count)
	if err != nil {
		return false, false, err
	}
	followedBy = (count == 1)

	// Followingをチェック
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM relationship WHERE followee_id = ? AND follower_id = ?", otherID, myID).Scan(&count)
	if err != nil {
		return false, false, err
	}
	following = (count == 1)

	return followedBy, following, nil
}
