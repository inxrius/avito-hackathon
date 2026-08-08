package model

import (
	"time"

	"github.com/google/uuid"
	recap "recap-personalization/internal/recap"
)

type RecapExplanation = recap.RecapExplanation
type DecisionExplanation = recap.DecisionExplanation
type RuleFact = recap.RuleFact
type ShareCard = recap.ShareCard
type ShareFact = recap.ShareFact
type ShareAchievement = recap.ShareAchievement
type ShareVisual = recap.ShareVisual
type Achievement = recap.Achievement
type AchievementLevel = recap.AchievementLevel
type AchievementCode = recap.AchievementCode
type ArchetypeRole = recap.ArchetypeRole
type ArchetypeStyle = recap.ArchetypeStyle
type MetricCode = recap.MetricCode
type ArchetypeRoleCode = recap.ArchetypeRoleCode
type ArchetypeStyleCode = recap.ArchetypeStyleCode

type Recap struct {
	SchemaVersion string                      `json:"schema_version"`
	ID            uuid.UUID                   `json:"id"`
	ProfileID     uuid.UUID                   `json:"profile_id"`
	Year          int                         `json:"year"`
	Profile       RecapProfile                `json:"profile"`
	Generation    RecapGeneration             `json:"generation"`
	Theme         RecapTheme                  `json:"theme"`
	Cards         []RecapCard                 `json:"cards"`
	Capabilities  RecapCapabilities           `json:"capabilities"`
	Metrics       []recap.MetricValue         `json:"-"`
	Archetype     recap.ArchetypeDecision     `json:"-"`
	Achievements  []recap.AchievementDecision `json:"-"`
	Narrative     recap.Narrative             `json:"-"`
	Explanation   *RecapExplanation           `json:"-"`
	Share         *ShareCard                  `json:"-"`
}

type RecapProfile struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
}

type RecapGeneration struct {
	AlgorithmVersion     string              `json:"algorithm_version"`
	FeatureSchemaVersion string              `json:"feature_schema_version"`
	ActivityHash         string              `json:"activity_hash"`
	GeneratedAt          time.Time           `json:"generated_at"`
	Narrative            NarrativeGeneration `json:"narrative"`
}

type NarrativeGeneration struct {
	Source        string  `json:"source"`
	PromptVersion string  `json:"prompt_version"`
	Model         *string `json:"model,omitempty"`
}

type RecapTheme struct {
	Code         string   `json:"code"`
	MainDistrict Vertical `json:"main_district"`
	AccentToken  *string  `json:"accent_token,omitempty"`
}

type Vertical struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

type RecapCapabilities struct {
	ShareAvailable       bool `json:"share_available"`
	ExplanationAvailable bool `json:"explanation_available"`
	FeedbackAvailable    bool `json:"feedback_available"`
}

type RecapCard struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Position    int                    `json:"position"`
	Visibility  string                 `json:"visibility"`
	Eyebrow     *string                `json:"eyebrow,omitempty"`
	Title       string                 `json:"title"`
	Description *string                `json:"description,omitempty"`
	Visual      CardVisual             `json:"visual"`
	Explainable bool                   `json:"explainable"`
	Data        map[string]interface{} `json:"data"`
}

type CardVisual struct {
	Kind      string  `json:"kind"`
	AssetCode *string `json:"asset_code,omitempty"`
}

type CreateRecapRequest struct {
	ProfileID uuid.UUID `json:"profile_id" binding:"required"`
	Year      int       `json:"year" binding:"required"`
}
