package usecase

import (
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type Relationship interface {
	// Follow(ctx context.Context, status string, acc *object.Account) (*CreateStatusDTO, error)
}

type relationship struct {
	db *sqlx.DB
	rr repository.Relationship
	ar repository.Account
}

type CreateRelationshipDTO struct {
	Account *object.Account
	Status  *object.Status
}

type RelationshipDTO struct {
	Account *object.Account
	Status  *object.Status
}

type GetPublicRelationshipDTO struct {
	Account []*object.Account
	Status  []*object.Status
}

var _ Status = (*status)(nil)

func NewRelationship(db *sqlx.DB, rr repository.Relationship, ar repository.Account) *relationship {
	return &relationship{
		db: db,
		rr: rr,
		ar: ar,
	}
}

// func (s *status) Create(ctx context.Context, status string, acc *object.Account) (*CreateStatusDTO, error) {
