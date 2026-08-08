package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"recap-personalization/internal/model"
	"recap-personalization/internal/repository"
)

func (h *Handler) GetProfiles(c *gin.Context) {
	profiles, err := h.service.GetProfiles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:      model.ErrCodeInternalError,
			Message:   "Failed to get profiles",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

func (h *Handler) GetProfileByID(c *gin.Context) {
	if _, err := uuid.Parse(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid profile id",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "id", Reason: "must_be_uuid"}},
		})
		return
	}
	profile, err := h.service.GetProfile(c.Request.Context(), c.Param("id"))
	if errors.Is(err, repository.ErrProfileNotFound) {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeProfileNotFound,
			Message:   "Profile not found",
			RequestID: generateRequestID(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:      model.ErrCodeInternalError,
			Message:   "Failed to get profile",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}
	c.JSON(http.StatusOK, profile)
}
