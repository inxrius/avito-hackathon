package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"recap-personalization/internal/service"
)

// Handler — основной обработчик HTTP-запросов
type Handler struct {
	service *service.Service
}

// NewHandler создаёт новый экземпляр Handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes регистрирует все маршруты API
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Profiles
		api.GET("/profiles", h.GetProfiles)
		api.GET("/profiles/:id", h.GetProfileByID)

		// Recaps
		api.POST("/recaps", h.CreateRecap)
		api.GET("/recaps/:id", h.GetRecap)
		api.GET("/recaps/:id/explanation", h.GetRecapExplanation)
		api.GET("/recaps/:id/share", h.GetShareCard)

		// Interactions
		api.POST("/recaps/:id/interactions", h.SaveInteraction)
	}
}

// generateRequestID генерирует UUID для трассировки запросов
func generateRequestID() string {
	return uuid.New().String()
}
