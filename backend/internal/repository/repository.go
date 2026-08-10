package repository

import (
	"context"

	"recap-personalization/internal/model"
	"recap-personalization/pkg/database"
)

type ProfileRepository interface {
	GetProfiles(ctx context.Context) ([]model.ProfileSummary, error)
	GetProfileByID(ctx context.Context, id string) (*model.Profile, error)
}

type RecapRepository interface {
	GetRecapByProfileAndYear(ctx context.Context, profileID string, year int) (*model.Recap, error)
	GetRecapByID(ctx context.Context, id string) (*model.Recap, error)
	CreateRecap(ctx context.Context, value *model.Recap) error
}

type Repository struct {
	DB *database.PostgresDB
}

func NewRepository(db *database.PostgresDB) *Repository {
	return &Repository{DB: db}
}

var (
	_ ProfileRepository = (*Repository)(nil)
	_ RecapRepository   = (*Repository)(nil)
)
