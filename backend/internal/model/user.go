package model

import "github.com/google/uuid"

type Profile struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	AvatarURL      *string   `json:"avatar_url,omitempty"`
	Scenario       string    `json:"-"`
	AvailableYears []int     `json:"available_years"`
}

type ProfileSummary = Profile
