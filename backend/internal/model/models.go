package model

import (
	"time"
	"github.com/google/uuid"
)

type Profile struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatar_url"`
	Scenario    string    `json:"scenario"`
	CreatedAt   time.Time `json:"created_at"`
}

type Activity struct {
	ID          uuid.UUID `json:"id"`
	ProfileID   uuid.UUID `json:"profile_id"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Value       float64   `json:"value"`
	Timestamp   time.Time `json:"timestamp"`
}

type Recap struct {
	ID            uuid.UUID   `json:"id"`
	ProfileID     uuid.UUID   `json:"profile_id"`
	Status        string      `json:"status"`
	Year          int         `json:"year"`
	Algorithm     string      `json:"algorithm_version"`
	ActivityHash  string      `json:"activity_hash"`
	GeneratedAt   time.Time   `json:"generated_at"`
	Profile       Profile     `json:"profile"`
	Cards         []RecapCard `json:"cards"`
}

type RecapCard struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Position int                    `json:"position"`
	Title    string                 `json:"title"`
	Data     map[string]interface{} `json:"data"`
}

type ProfileSummary struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatar_url"`
	Scenario    string    `json:"scenario"`
}

type ProfileList struct {
	Profiles []ProfileSummary `json:"profiles"`
}

type CreateRecapRequest struct {
	ProfileID uuid.UUID `json:"profile_id"`
	Year      int       `json:"year"`
}

type CreateRecapResponse struct {
	RecapID uuid.UUID `json:"recap_id"`
	Status  string    `json:"status"`
	Recap   *Recap    `json:"recap,omitempty"`
}

type InteractionRequest struct {
	Type     string                 `json:"type" binding:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Explanation models
type RecapExplanation struct {
	RecapID      uuid.UUID             `json:"recap_id"`
	Algorithm    string                `json:"algorithm_version"`
	ActivityHash string                `json:"activity_hash"`
	Decisions    []DecisionExplanation `json:"decisions"`
}

type DecisionExplanation struct {
	CardID string     `json:"card_id"`
	Reason string     `json:"reason"`
	Facts  []RuleFact `json:"facts"`
}

type RuleFact struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Share card models
type ShareCard struct {
	RecapID      uuid.UUID          `json:"recap_id"`
	ProfileName  string             `json:"profile_name"`
	Year         int                `json:"year"`
	Facts        []ShareFact        `json:"facts"`
	Achievements []ShareAchievement `json:"achievements"`
	Visual       ShareVisual        `json:"visual"`
}

type ShareFact struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ShareAchievement struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

type ShareVisual struct {
	Theme  string   `json:"theme"`
	Colors []string `json:"colors"`
}

// Interaction models
type InteractionResponse struct {
	Success bool      `json:"success"`
	RecapID uuid.UUID `json:"recap_id"`
}
