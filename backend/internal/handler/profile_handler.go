package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"recap-personalization/internal/model"
)

// GetProfiles возвращает список всех профилей
func (h *Handler) GetProfiles(c *gin.Context) {
	profiles, err := h.service.GetProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:      model.ErrCodeInternalError,
			Message:   "Failed to get profiles",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}
	// По OpenAPI возвращаем массив напрямую
	c.JSON(http.StatusOK, model.ProfileList{Profiles: profiles})
}

// GetProfileByID возвращает профиль по ID
func (h *Handler) GetProfileByID(c *gin.Context) {
	id := c.Param("id")
	profile, err := h.service.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeProfileNotFound,
			Message:   "Profile not found",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}
	c.JSON(http.StatusOK, profile)
}
