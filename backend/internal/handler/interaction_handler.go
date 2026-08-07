package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/inxrius/avito-hackathon/internal/model"
)

// SaveInteraction сохраняет событие взаимодействия
func (h *Handler) SaveInteraction(c *gin.Context) {
	recapID := c.Param("id")
	var req model.InteractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid request body",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "body", Reason: err.Error()}},
		})
		return
	}

	// Валидация event_id
	if _, err := uuid.Parse(req.EventID); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid event_id",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "event_id", Reason: "must be UUID"}},
		})
		return
	}
	if _, err := uuid.Parse(recapID); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid recap_id",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "recap_id", Reason: "must be UUID"}},
		})
		return
	}

	// Проверяем существование recap
	if _, err := h.service.GetRecap(recapID); err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeRecapNotFound,
			Message:   "Recap not found",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}

	// Сохраняем событие
	if err := h.service.SaveInteraction(recapID, req.EventName, req.Properties); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:      model.ErrCodeInternalError,
			Message:   "Failed to save interaction",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}

	eventID, _ := uuid.Parse(req.EventID)
	recapUUID, _ := uuid.Parse(recapID)
	c.JSON(http.StatusAccepted, model.InteractionResponse{
		Accepted: true,
		RecapID:  recapUUID,
		EventID:  eventID,
	})
}
