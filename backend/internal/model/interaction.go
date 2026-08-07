package model

import (
	"time"

	"github.com/google/uuid"
)

// InteractionRequest — запрос на сохранение события взаимодействия
// POST /recaps/{id}/interactions
type InteractionRequest struct {
	EventID    string                 `json:"event_id" binding:"required"`    // UUID события
	SessionID  string                 `json:"session_id" binding:"required"`  // UUID сессии
	EventName  string                 `json:"event_name" binding:"required"`  // одно из enum: recap_opened, card_viewed, ...
	OccurredAt time.Time              `json:"occurred_at" binding:"required"` // время события (UTC)
	Properties map[string]interface{} `json:"properties,omitempty"`           // дополнительные данные (card_id, position и т.д.)
}

// InteractionResponse — ответ на сохранение события
type InteractionResponse struct {
	Accepted bool      `json:"accepted"` // всегда true при успехе
	RecapID  uuid.UUID `json:"recap_id"` // ID recap, к которому относится событие
	EventID  uuid.UUID `json:"event_id"` // ID события (возвращаем как UUID)
}