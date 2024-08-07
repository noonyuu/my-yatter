package usecase

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository/mock"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
)

func TestCreate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockAccount := repository.NewMockAccount(ctrl)

    // 期待値を設定
    mockAccount.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

    // モックを使用したテスト対象のコードを呼び出す
    ctx := context.Background()
    tx := &sqlx.Tx{}
    acc := &object.Account{}
    err := mockAccount.Create(ctx, tx, acc)

    if err != nil {
        t.Errorf("Create() returned an error: %v", err)
    }
}

func TestFindAccountByID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockAccount := repository.NewMockAccount(ctrl)

    // 期待値を設定
    mockAccount.EXPECT().FindAccountByID(gomock.Any(), gomock.Any()).Return(&object.Account{}, nil)

    // モックを使用したテスト対象のコードを呼び出す
    ctx := context.Background()
    id := 1
    acc, err := mockAccount.FindAccountByID(ctx, id)

    if err != nil {
        t.Errorf("FindAccountByID() returned an error: %v", err)
    }
    if acc == nil {
        t.Error("FindAccountByID() returned nil account")
    }
}

func TestFindByUsername(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockAccount := repository.NewMockAccount(ctrl)

    // 期待値を設定
    mockAccount.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(&object.Account{}, nil)

    // モックを使用したテスト対象のコードを呼び出す
    ctx := context.Background()
    username := "testuser"
    acc, err := mockAccount.FindByUsername(ctx, username)

    if err != nil {
        t.Errorf("FindByUsername() returned an error: %v", err)
    }
    if acc == nil {
        t.Error("FindByUsername() returned nil account")
    }
}
