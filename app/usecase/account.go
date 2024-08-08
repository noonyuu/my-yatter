package usecase

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type Account interface {
	Create(ctx context.Context, username, password string) (*CreateAccountDTO, error)
	FindByUsername(ctx context.Context, username string) (*GetAccountDTO, error)
	UpdateCredentials(ctx context.Context, account *object.Account) (*CreateAccountDTO, error)
}

type account struct {
	db *sqlx.DB
	ar repository.Account
}

type CreateAccountDTO struct {
	Account *object.Account
}

type GetAccountDTO struct {
	Account *object.Account
}

var _ Account = (*account)(nil)

func NewAccount(db *sqlx.DB, ar repository.Account) *account {
	return &account{
		db: db,
		ar: ar,
	}
}

func (a *account) Create(ctx context.Context, username, password string) (*CreateAccountDTO, error) {
	acc, err := object.NewAccount(username, password)
	if err != nil {
		return nil, err
	}

	tx, err := a.db.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	if err := a.ar.Create(ctx, tx, acc); err != nil {
		return nil, err
	}

	return &CreateAccountDTO{
		Account: acc,
	}, nil
}

func (a *account) FindByUsername(ctx context.Context, username string) (*GetAccountDTO, error) {

	acc, err := a.ar.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &GetAccountDTO{
		Account: acc,
	}, nil
}

func (a *account) UpdateCredentials(ctx context.Context, acc *object.Account) (*CreateAccountDTO, error) {
	tx, err := a.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	upCre, err := object.UpdateCredential(*acc.DisplayName, *acc.Note, *acc.Avatar, *acc.Header)
	upCre.ID = acc.ID
	if err != nil {
		return nil, err
	}

	if err := a.ar.UpdateAccountCredential(ctx, tx, upCre); err != nil {
		return nil, err
	}

	return &CreateAccountDTO{
		Account: acc,
	}, nil
}
