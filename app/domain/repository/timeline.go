package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Timeline interface {
	GetPublicTimeline(ctx context.Context, maxID, sinceID, limit int) ([]*object.Status, error)
}