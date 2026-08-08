package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"recap-personalization/internal/model"
	recap "recap-personalization/internal/recap"
	"recap-personalization/internal/repository"
	"recap-personalization/internal/service"
)

func (h *Handler) CreateRecap(c *gin.Context) {
	var request model.CreateRecapRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid request body",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "body", Reason: err.Error()}},
		})
		return
	}
	if request.Year < 2000 || request.Year > 2100 {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Year must be between 2000 and 2100",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "year", Reason: "out_of_range"}},
		})
		return
	}

	value, created, err := h.service.GenerateRecap(
		c.Request.Context(),
		request.ProfileID.String(),
		request.Year,
	)
	if err != nil {
		h.writeRecapGenerationError(c, err)
		return
	}
	if created {
		c.Header("Location", "/api/v1/recaps/"+value.ID.String())
		c.JSON(http.StatusCreated, value)
		return
	}
	c.JSON(http.StatusOK, value)
}

func (h *Handler) GetRecap(c *gin.Context) {
	if _, err := uuid.Parse(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid recap id",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "id", Reason: "must_be_uuid"}},
		})
		return
	}
	value, err := h.service.GetRecap(c.Request.Context(), c.Param("id"))
	if errors.Is(err, repository.ErrRecapNotFound) {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeRecapNotFound,
			Message:   "Recap not found",
			RequestID: generateRequestID(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:      model.ErrCodeInternalError,
			Message:   "Failed to get recap",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}
	c.JSON(http.StatusOK, value)
}

func (h *Handler) GetRecapExplanation(c *gin.Context) {
	if _, err := uuid.Parse(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid recap id",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "id", Reason: "must_be_uuid"}},
		})
		return
	}
	value, err := h.service.GetRecap(c.Request.Context(), c.Param("id"))
	if errors.Is(err, repository.ErrRecapNotFound) {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeRecapNotFound,
			Message:   "Recap not found",
			RequestID: generateRequestID(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:      model.ErrCodeInternalError,
			Message:   "Failed to get recap explanation",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}
	if !value.Capabilities.ExplanationAvailable || value.Explanation == nil {
		c.JSON(http.StatusConflict, model.APIError{
			Code:      model.ErrCodeExplanationNotAvailable,
			Message:   "Explanation not available for this recap",
			RequestID: generateRequestID(),
		})
		return
	}
	c.JSON(http.StatusOK, value.Explanation)
}

func (h *Handler) GetShareCard(c *gin.Context) {
	if _, err := uuid.Parse(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid recap id",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "id", Reason: "must_be_uuid"}},
		})
		return
	}
	value, err := h.service.GetRecap(c.Request.Context(), c.Param("id"))
	if errors.Is(err, repository.ErrRecapNotFound) {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeRecapNotFound,
			Message:   "Recap not found",
			RequestID: generateRequestID(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:      model.ErrCodeInternalError,
			Message:   "Failed to get share card",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}
	if !value.Capabilities.ShareAvailable || value.Share == nil {
		c.JSON(http.StatusConflict, model.APIError{
			Code:      model.ErrCodeShareNotAvailable,
			Message:   "Share card not available for this recap",
			RequestID: generateRequestID(),
		})
		return
	}
	c.JSON(http.StatusOK, value.Share)
}

func (h *Handler) writeRecapGenerationError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrProfileNotFound):
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeProfileNotFound,
			Message:   "Profile not found",
			RequestID: generateRequestID(),
		})
	case errors.Is(err, service.ErrYearNotAvailable):
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "The selected year is not available for this profile",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "year", Reason: "not_available"}},
		})
	case errors.Is(err, recap.ErrInsufficientActivity):
		c.JSON(http.StatusUnprocessableEntity, model.APIError{
			Code:      model.ErrCodeInsufficientActivity,
			Message:   "Not enough activity data for the selected year",
			RequestID: generateRequestID(),
		})
	case errors.Is(err, service.ErrActivitySourceMissing), errors.Is(err, service.ErrActivitySourceUnavailable):
		c.JSON(http.StatusServiceUnavailable, model.APIError{
			Code:      model.ErrCodeDependencyUnavailable,
			Message:   "Activity source is unavailable",
			RequestID: generateRequestID(),
		})
	default:
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:      model.ErrCodeInternalError,
			Message:   "Failed to generate recap",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
	}
}
