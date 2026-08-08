package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"recap-personalization/internal/model"
	"recap-personalization/internal/repository"
	"recap-personalization/internal/service"
)

func (h *Handler) SaveInteraction(c *gin.Context) {
	recapID := c.Param("id")
	parsedRecapID, err := uuid.Parse(recapID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid recap id",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "id", Reason: "must_be_uuid"}},
		})
		return
	}

	var request model.InteractionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid request body",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "body", Reason: err.Error()}},
		})
		return
	}
	if request.EventID == uuid.Nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid event id",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "event_id", Reason: "must_be_uuid"}},
		})
		return
	}
	if request.SessionID == uuid.Nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid session id",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "session_id", Reason: "must_be_uuid"}},
		})
		return
	}
	if !model.IsValidInteractionName(request.EventName) {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid interaction event name",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "event_name", Reason: "unsupported_value"}},
		})
		return
	}
	if request.OccurredAt.IsZero() {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid occurred_at",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "occurred_at", Reason: "required"}},
		})
		return
	}

	if _, err := h.service.GetRecap(c.Request.Context(), recapID); errors.Is(err, repository.ErrRecapNotFound) {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeRecapNotFound,
			Message:   "Recap not found",
			RequestID: generateRequestID(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:      model.ErrCodeInternalError,
			Message:   "Failed to verify recap",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}

	if err := h.service.SaveInteraction(c.Request.Context(), recapID, request); err != nil {
		code := model.ErrCodeInternalError
		status := http.StatusInternalServerError
		message := "Failed to save interaction"
		if errors.Is(err, service.ErrInteractionSinkMissing) || errors.Is(err, service.ErrInteractionSinkUnavailable) {
			code = model.ErrCodeDependencyUnavailable
			status = http.StatusServiceUnavailable
			message = "Interaction storage is unavailable"
		}
		c.JSON(status, model.APIError{
			Code:      code,
			Message:   message,
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}

	c.JSON(http.StatusAccepted, model.InteractionResponse{
		Accepted: true,
		RecapID:  parsedRecapID,
		EventID:  request.EventID,
	})
}
