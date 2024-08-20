package usecase

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"
)

type MockStatusRepo struct {
	mockCreateFunc func(ctx context.Context, status *object.Status) error
}

func (m *MockStatusRepo) Create(ctx context.Context, status *object.Status) error {
	if m.mockCreateFunc != nil {
		return m.mockCreateFunc(ctx, status)
	}
	return nil
}

func TestStatusUsecase_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("正常値 : Status情報を正常に返すこと", func(t *testing.T) {
		content := "test"
		accountID := &object.Account{ID: 1}
		mockRepo := MockStatusRepo{
			mockCreateFunc: func(ctx context.Context, status *object.Status) error {
				return nil
			},
		}
		sut := NewStatus(&mockRepo)

		got, err := sut.Create(ctx, content, accountID)
		// assert := 
		// want = &object.NewStatus(content, accountID)

		// assert.Equal(t, want, got.Status)
	})

	t.Run("異常値 : StatusRepository.Create()でエラーが発生した場合、エラーを返すこと", func(t *testing.T) {

	})
}
