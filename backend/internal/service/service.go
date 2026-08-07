package service

import (
	"github.com/inxrius/avito-hackathon/internal/model"
	"github.com/inxrius/avito-hackathon/internal/repository"
)

// Service — объединяет все бизнес-сервисы
type Service struct {
	repo *repository.Repository
}

// NewService создаёт новый экземпляр Service
func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// GetProfiles возвращает список профилей
func (s *Service) GetProfiles() ([]model.ProfileSummary, error) {
	return s.repo.GetProfiles()
}

// GetProfile возвращает профиль по ID
func (s *Service) GetProfile(id string) (*model.Profile, error) {
	return s.repo.GetProfileByID(id)
}

// GetRecap возвращает Recap по ID
func (s *Service) GetRecap(id string) (*model.Recap, error) {
	return s.repo.GetRecapByID(id)
}

// SaveInteraction сохраняет событие взаимодействия
func (s *Service) SaveInteraction(recapID, eventType string, metadata map[string]interface{}) error {
	return s.repo.SaveInteraction(recapID, eventType, metadata)
}
