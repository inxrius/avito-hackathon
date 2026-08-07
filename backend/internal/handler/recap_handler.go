package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"recap-personalization/internal/model"
)

// CreateRecap создаёт или возвращает существующий recap (идемпотентный)
func (h *Handler) CreateRecap(c *gin.Context) {
	var req model.CreateRecapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Invalid request body",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "body", Reason: err.Error()}},
		})
		return
	}

	// Валидация года
	if req.Year < 2000 || req.Year > 2100 {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:      model.ErrCodeInvalidArgument,
			Message:   "Year must be between 2000 and 2100",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Field: "year", Reason: "out_of_range"}},
		})
		return
	}

	// Проверка существования профиля
	_, err := h.service.GetProfile(req.ProfileID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeProfileNotFound,
			Message:   "Profile not found",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}

	// Генерация (или получение существующего) recap
	recap, err := h.service.GenerateRecap(req.ProfileID.String(), req.Year)
	if err != nil {
		if strings.Contains(err.Error(), "insufficient_activity") || strings.Contains(err.Error(), "no activities") {
			c.JSON(http.StatusUnprocessableEntity, model.APIError{
				Code:      model.ErrCodeInsufficientActivity,
				Message:   "Not enough activity data for the selected year",
				RequestID: generateRequestID(),
				Details:   []model.ErrorDetail{{Reason: err.Error()}},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:      model.ErrCodeInternalError,
			Message:   "Failed to generate recap",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}

	c.JSON(http.StatusCreated, recap)
}

// GetRecap возвращает recap по ID
func (h *Handler) GetRecap(c *gin.Context) {
	id := c.Param("id")
	recap, err := h.service.GetRecap(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeRecapNotFound,
			Message:   "Recap not found",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}
	c.JSON(http.StatusOK, recap)
}

// GetRecapExplanation возвращает объяснение решений
func (h *Handler) GetRecapExplanation(c *gin.Context) {
	id := c.Param("id")
	recap, err := h.service.GetRecap(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeRecapNotFound,
			Message:   "Recap not found",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}

	// Проверяем, доступно ли объяснение
	if !recap.Capabilities.ExplanationAvailable {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeExplanationNotAvailable,
			Message:   "Explanation not available for this recap",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: "explanation_disabled"}},
		})
		return
	}

	// Используем динамическое объяснение из генератора, если доступно
	if recap.Explanation != nil {
		c.JSON(http.StatusOK, recap.Explanation)
		return
	}

	// Fallback: если объяснение не сгенерировано, возвращаем базовую информацию
	explanation := model.RecapExplanation{
		RecapID:          recap.ID.String(),
		AlgorithmVersion: recap.AlgorithmVersion,
		ActivityHash:     recap.ActivityHash,
		Decisions:        []model.DecisionExplanation{},
	}
	
	c.JSON(http.StatusOK, explanation)
}

// GetShareCard возвращает публичную карточку
func (h *Handler) GetShareCard(c *gin.Context) {
	id := c.Param("id")
	recap, err := h.service.GetRecap(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeRecapNotFound,
			Message:   "Recap not found",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: err.Error()}},
		})
		return
	}

	// Проверяем, доступен ли share
	if !recap.Capabilities.ShareAvailable {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:      model.ErrCodeShareNotAvailable,
			Message:   "Share not available for this recap",
			RequestID: generateRequestID(),
			Details:   []model.ErrorDetail{{Reason: "share_disabled"}},
		})
		return
	}

	// Используем динамический share из генератора, если доступен
	if recap.Share != nil {
		c.JSON(http.StatusOK, recap.Share)
		return
	}

	// Fallback: если share не сгенерирован, возвращаем ошибку
	c.JSON(http.StatusNotFound, model.APIError{
		Code:      model.ErrCodeShareNotAvailable,
		Message:   "Share not available for this recap",
		RequestID: generateRequestID(),
		Details:   []model.ErrorDetail{{Reason: "share_not_generated"}},
	})
}
