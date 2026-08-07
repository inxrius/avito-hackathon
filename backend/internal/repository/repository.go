package repository

import (
	"recap-personalization/internal/model"
	"recap-personalization/pkg/database"
)

// Интерфейсы репозиториев
// ProfileRepository — работа с профилями
type ProfileRepository interface {
	GetProfiles() ([]model.ProfileSummary, error)
	GetProfileByID(id string) (*model.Profile, error)
}

// ActivityRepository — работа с активностями
type ActivityRepository interface {
	GetActivitiesByProfileIDAndYear(profileID string, year int) ([]model.Activity, error)
}

// RecapRepository — работа с итогами (recap)
type RecapRepository interface {
	GetRecapByProfileAndYear(profileID string, year int) (*model.Recap, error)
	GetRecapByID(id string) (*model.Recap, error)
	CreateRecap(recap *model.Recap) error
}

// InteractionRepository — работа с взаимодействиями
type InteractionRepository interface {
	SaveInteraction(recapID, eventType string, metadata map[string]interface{}) error
}

// Реализация
// Repository — объединяет все репозитории в одной структуре
type Repository struct {
	DB *database.PostgresDB
}

// NewRepository создаёт новый экземпляр Repository
func NewRepository(db *database.PostgresDB) *Repository {
	return &Repository{DB: db}
}

// Проверяем, что Repository реализует все интерфейсы
var (
	_ ProfileRepository     = (*Repository)(nil)
	_ ActivityRepository    = (*Repository)(nil)
	_ RecapRepository       = (*Repository)(nil)
	_ InteractionRepository = (*Repository)(nil)
)
