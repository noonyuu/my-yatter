package usecase

import (
	"context"
	"strconv"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type Status interface {
	Create(ctx context.Context, status string, acc *object.Account) (*CreateStatusDTO, error)
	FindByStatus(ctx context.Context, id string) (*GetStatusDTO, error)
	GetPublicTimeline(ctx context.Context, maxId, sinceId, limit string) (*GetPublicStatusDTO, error)
	DeleteStatus(ctx context.Context, id string) error
}

type status struct {
	db *sqlx.DB
	ar repository.Account
	sr repository.Status
}

type CreateStatusDTO struct {
	Account *object.Account
	Status  *object.Status
}

type GetStatusDTO struct {
	Account *object.Account
	Status  *object.Status
}

type GetPublicStatusDTO struct {
	Account []*object.Account
	Status  []*object.Status
}

var _ Status = (*status)(nil)

func NewStatus(db *sqlx.DB, ar repository.Account, sr repository.Status) *status {
	return &status{
		db: db,
		ar: ar,
		sr: sr,
	}
}

func (s *status) Create(ctx context.Context, status string, acc *object.Account) (*CreateStatusDTO, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	sta := object.NewStatus(status)
	sta.AccountID = int(acc.ID)

	sta, err = s.sr.Create(ctx, tx, sta)
	if err != nil {
		return nil, err
	}

	return &CreateStatusDTO{
		Account: acc,
		Status:  sta,
	}, nil
}

func (s *status) FindByStatus(ctx context.Context, sid string) (*GetStatusDTO, error) {
	id, err := strconv.Atoi(sid) // stringで受け取ったidをintに変換
	if err != nil {
		return nil, err
	}

	// statusの取得
	sta, err := s.sr.FindStatusByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// statusに紐づくaccountの取得
	acc, err := s.ar.FindAccountByID(ctx, sta.AccountID)
	if err != nil {
		return nil, err
	}

	return &GetStatusDTO{
		Account: acc,
		Status:  sta,
	}, nil
}

func (s *status) GetPublicTimeline(ctx context.Context, maxId, sinceId, limit string) (*GetPublicStatusDTO, error) {
	// maxId, sinceId, limitをintに変換
	maxID, err := strconv.Atoi(maxId)
	if err != nil {
		return nil, err
	}
	sinceID, err := strconv.Atoi(sinceId)
	if err != nil {
		return nil, err
	}
	lmt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	// statusの取得
	sta, err := s.sr.GetPublicTimeline(ctx, maxID, sinceID, lmt)
	if err != nil {
		return nil, err
	}

	// statusに紐づくaccountの取得
	// statusの数のスライスを作成
	acc := make([]*object.Account, len(sta))
	for i := range sta {
		acc[i], err = s.ar.FindAccountByID(ctx, sta[i].AccountID)
		if err != nil {
			return nil, err
		}
	}

	return &GetPublicStatusDTO{
		Account: acc,
		Status:  sta,
	}, nil
}

func (s *status) DeleteStatus(ctx context.Context, id string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	// idをintに変換
	ID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	// statusの削除
	err = s.sr.DeleteStatus(ctx, tx, ID)
	if err != nil {
		return err
	}
	return nil
}
