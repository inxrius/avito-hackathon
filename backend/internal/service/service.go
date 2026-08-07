package service

import (
	"recap-personalization/internal/model"
	"recap-personalization/internal/repository"
	"recap-personalization/internal/recap/pipeline"
	"recap-personalization/internal/recap/ports"
)

// Service — объединяет все бизнес-сервисы
type Service struct {
	repo              *repository.Repository
	recapGenerator    *pipeline.Service
	clickHouseActivities ports.ActivityRepository
	clickHouseInteractions ports.InteractionRepository
}

// NewService создаёт новый экземпляр Service
func NewService(repo *repository.Repository, clickHouseActivities ports.ActivityRepository, clickHouseInteractions ports.InteractionRepository) *Service {
	recapGen := pipeline.NewGenerator(nil)
	recapGen.Registry.PublicAvatarHosts = map[string]struct{}{}
	
	return &Service{
		repo:                   repo,
		recapGenerator:         recapGen,
		clickHouseActivities:   clickHouseActivities,
		clickHouseInteractions: clickHouseInteractions,
	}
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
