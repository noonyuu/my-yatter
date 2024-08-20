package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Timeline interface {
	GetPublicTimeline(ctx context.Context, sinceID, limit int64) ([]*object.Status, error)
	GetHomeTimeline(ctx context.Context, id, maxID, sinceID, limit int64) ([]*object.Status, error)
}