package service

import (
	"context"
	"errors"
	"fmt"

	"recap-personalization/internal/model"
	recap "recap-personalization/internal/recap"
	"recap-personalization/internal/recap/ports"
	"recap-personalization/internal/repository"
)

var (
	ErrYearNotAvailable           = errors.New("year_not_available")
	ErrActivitySourceMissing      = errors.New("activity_source_missing")
	ErrInteractionSinkMissing     = errors.New("interaction_sink_missing")
	ErrActivitySourceUnavailable  = errors.New("activity_source_unavailable")
	ErrInteractionSinkUnavailable = errors.New("interaction_sink_unavailable")
)

type Service struct {
	profiles               repository.ProfileRepository
	recaps                 repository.RecapRepository
	recapGenerator         recap.Generator
	clickHouseActivities   ports.ActivityRepository
	clickHouseInteractions ports.InteractionRepository
}

func NewService(
	repo *repository.Repository,
	clickHouseActivities ports.ActivityRepository,
	clickHouseInteractions ports.InteractionRepository,
	generator recap.Generator,
) *Service {
	return &Service{
		profiles:               repo,
		recaps:                 repo,
		recapGenerator:         generator,
		clickHouseActivities:   clickHouseActivities,
		clickHouseInteractions: clickHouseInteractions,
	}
}

func (s *Service) GetProfiles(ctx context.Context) ([]model.ProfileSummary, error) {
	return s.profiles.GetProfiles(ctx)
}

func (s *Service) GetProfile(ctx context.Context, id string) (*model.Profile, error) {
	return s.profiles.GetProfileByID(ctx, id)
}

func (s *Service) GetRecap(ctx context.Context, id string) (*model.Recap, error) {
	return s.recaps.GetRecapByID(ctx, id)
}

func (s *Service) SaveInteraction(ctx context.Context, recapID string, request model.InteractionRequest) error {
	if s.clickHouseInteractions == nil {
		return ErrInteractionSinkMissing
	}
	err := s.clickHouseInteractions.SaveInteraction(ctx, ports.InteractionEvent{
		EventID:    request.EventID.String(),
		RecapID:    recapID,
		SessionID:  request.SessionID.String(),
		EventName:  request.EventName,
		OccurredAt: request.OccurredAt.UTC(),
		Properties: request.Properties,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInteractionSinkUnavailable, err)
	}
	return nil
}

func containsYear(values []int, year int) bool {
	for _, value := range values {
		if value == year {
			return true
		}
	}
	return false
}
