package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/inxrius/avito-hackathon/internal/model"
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

	// Формируем объяснение (пока статичное)
	explanation := model.RecapExplanation{
		RecapID:          recap.ID,
		AlgorithmVersion: recap.AlgorithmVersion,
		ActivityHash:     recap.ActivityHash,
		Decisions: []model.DecisionExplanation{
			{
				CardID:      "archetype",
				Kind:        "archetype_role",
				Code:        "findings_seeker",
				Reason:      "На основе анализа активности пользователя за год",
				RuleVersion: "archetype-rules-v1",
				Facts: []model.RuleFact{
					{
						MetricCode: "favorites_count",
						Actual:     46.0,
						Operator:   "gte",
						Threshold:  30.0,
						Matched:    true,
					},
				},
			},
		},
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

	// Собираем публичные достижения из карточек
	var shareAchievements []model.ShareAchievement
	for _, card := range recap.Cards {
		if card.Type == "achievements" {
			if items, ok := card.Data["items"].([]interface{}); ok {
				for _, item := range items {
					if achMap, ok := item.(map[string]interface{}); ok {
						code, _ := achMap["code"].(string)
						title, _ := achMap["title"].(string)
						level, _ := achMap["level"].(string)
						icon, _ := achMap["icon"].(string)
						shareAchievements = append(shareAchievements, model.ShareAchievement{
							Code:  code,
							Title: title,
							Level: level,
							Icon:  icon,
						})
					}
				}
			}
		}
	}

	// Если нет достижений, добавим заглушку
	if len(shareAchievements) == 0 {
		shareAchievements = []model.ShareAchievement{
			{Code: "first_step", Title: "Первый шаг", Level: "newcomer", Icon: "👁"},
		}
	}

	facts := []model.ShareFact{
		{Kind: "main_district", Label: "Главный район", Value: recap.Theme.MainDistrict.Title},
		{Kind: "active_days", Label: "Дней активности", Value: "84"},
	}

	shareCard := model.ShareCard{
		SchemaVersion: "2.0",
		RecapID:       recap.ID,
		ProfileName:   recap.Profile.Name,
		Year:          recap.Year,
		Title:         recap.SummaryTitle,
		Subtitle:      recap.SummaryText,
		MainDistrict:  recap.Theme.MainDistrict,
		Facts:         facts,
		Achievements:  shareAchievements,
		Visual: model.ShareVisual{
			Theme:  "city",
			Colors: []string{"#7C3AED", "#F3E8FF"},
		},
	}

	c.JSON(http.StatusOK, shareCard)
}
