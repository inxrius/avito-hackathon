package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	InteractionRecapOpened          = "recap_opened"
	InteractionCardViewed           = "card_viewed"
	InteractionExplanationOpened    = "explanation_opened"
	InteractionRecapCompleted       = "recap_completed"
	InteractionSharePreviewOpened   = "share_preview_opened"
	InteractionShareClicked         = "share_clicked"
	InteractionAchievementsExpanded = "achievements_expanded"
	InteractionFeedbackSubmitted    = "feedback_submitted"
)

type InteractionRequest struct {
	EventID    uuid.UUID              `json:"event_id" binding:"required"`
	SessionID  uuid.UUID              `json:"session_id" binding:"required"`
	EventName  string                 `json:"event_name" binding:"required"`
	OccurredAt time.Time              `json:"occurred_at" binding:"required"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type InteractionResponse struct {
	Accepted bool      `json:"accepted"`
	RecapID  uuid.UUID `json:"recap_id"`
	EventID  uuid.UUID `json:"event_id"`
}

func IsValidInteractionName(value string) bool {
	switch value {
	case InteractionRecapOpened,
		InteractionCardViewed,
		InteractionExplanationOpened,
		InteractionRecapCompleted,
		InteractionSharePreviewOpened,
		InteractionShareClicked,
		InteractionAchievementsExpanded,
		InteractionFeedbackSubmitted:
		return true
	default:
		return false
	}
}
