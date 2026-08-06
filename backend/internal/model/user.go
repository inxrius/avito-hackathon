package model

import (
	"time"

	"github.com/google/uuid"
)

// Profile — полная информация о пользователе (возвращается по /profiles/{id})
type Profile struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	AvatarURL      string    `json:"avatar_url"`
	Scenario       string    `json:"scenario,omitempty"` // для внутреннего использования (тип пользователя)
	AvailableYears []int     `json:"available_years"`    // годы, за которые доступны итоги
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ProfileSummary — сокращённая версия для списка (/profiles)
type ProfileSummary struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	AvatarURL      string    `json:"avatar_url"`
	Scenario       string    `json:"scenario,omitempty"`
	AvailableYears []int     `json:"available_years"`
}

// ProfileList — используется для ответа /profiles (массив ProfileSummary)
// В OpenAPI мы возвращаем массив напрямую, но для удобства можно использовать обёртку
// или возвращать []ProfileSummary. Оставлю здесь на случай, если понадобится.
type ProfileList struct {
	Profiles []ProfileSummary `json:"profiles"`
}
