package usecase

import (
	"context"
	"strconv"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type Timeline interface {
	GetPublicTimeline(ctx context.Context, sinceId, limit string) (*GetPublicStatusDTO, error)
	GetHomeTimeline(ctx context.Context, id int64, maxId, sinceId, limit string) (*GetPublicStatusDTO, error)
}

type timeline struct {
	db *sqlx.DB
	ar repository.Account
	tr repository.Timeline
}

type GetPublicStatusDTO struct {
	Account []*object.Account
	Status  []*object.Status
}

var _ Timeline = (*timeline)(nil)

func NewTimeline(db *sqlx.DB, ar repository.Account, tr repository.Timeline) *timeline {
	return &timeline{
		db: db,
		ar: ar,
		tr: tr,
	}
}

func (t *timeline) GetPublicTimeline(ctx context.Context, sinceId, limit string) (*GetPublicStatusDTO, error) {
	// sinceId, limitをint64に変換
	sinceID, err := strconv.ParseInt(sinceId, 10, 64)
	if err != nil {
		return nil, err
	}
	lmt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, err
	}

	// statusの取得
	sta, err := t.tr.GetPublicTimeline(ctx, sinceID, lmt)
	if err != nil {
		return nil, err
	}

	// statusに紐づくaccountの取得
	// statusの数のスライスを作成
	acc := make([]*object.Account, len(sta))
	for i := range sta {
		acc[i], err = t.ar.FindAccountByID(ctx, sta[i].AccountID)
		if err != nil {
			return nil, err
		}
	}

	return &GetPublicStatusDTO{
		Account: acc,
		Status:  sta,
	}, nil
}

func (t *timeline) GetHomeTimeline(ctx context.Context, id int64, maxId, sinceId, limit string) (*GetPublicStatusDTO, error) {
	// maxId, sinceId, limitをintに変換
	maxID, err := strconv.ParseInt(maxId, 10, 64)
	if err != nil {
		return nil, err
	}
	sinceID, err := strconv.ParseInt(sinceId, 10, 64)
	if err != nil {
		return nil, err
	}
	lmt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, err
	}

	// statusの取得
	sta, err := t.tr.GetHomeTimeline(ctx, id, maxID, sinceID, lmt)
	if err != nil {
		return nil, err
	}

	// statusに紐づくaccountの取得
	// statusの数のスライスを作成
	acc := make([]*object.Account, len(sta))
	for i := range sta {
		acc[i], err = t.ar.FindAccountByID(ctx, sta[i].AccountID)
		if err != nil {
			return nil, err
		}
	}

	return &GetPublicStatusDTO{
		Account: acc,
		Status:  sta,
	}, nil
}
