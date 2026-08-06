package model

import (
	"time"

	"github.com/google/uuid"
)

// Recap — полная структура итогов года (ответ POST /recaps и GET /recaps/{id})
type Recap struct {
	ID                   uuid.UUID         `json:"id"`
	ProfileID            uuid.UUID         `json:"profile_id"`
	Status               string            `json:"status,omitempty"` // "completed" | "pending"
	Year                 int               `json:"year"`
	SchemaVersion        string            `json:"schema_version"` // "2.0"
	AlgorithmVersion     string            `json:"algorithm_version"`
	FeatureSchemaVersion string            `json:"feature_schema_version"`
	ActivityHash         string            `json:"activity_hash"`
	GeneratedAt          time.Time         `json:"generated_at"`
	NarrativeSource      string            `json:"narrative_source"` // "mistral" | "template"
	PromptVersion        string            `json:"prompt_version"`
	NarrativeModel       *string           `json:"narrative_model,omitempty"`
	MainVerticalCode     string            `json:"main_vertical_code,omitempty"`
	AccentToken          *string           `json:"accent_token,omitempty"`
	SummaryTitle         string            `json:"summary_title"`
	SummaryText          string            `json:"summary_text"`
	Profile              RecapProfile      `json:"profile"`
	Generation           RecapGeneration   `json:"generation"`
	Theme                RecapTheme        `json:"theme"`
	Cards                []RecapCard       `json:"cards"`
	Capabilities         RecapCapabilities `json:"capabilities"`
}

// RecapProfile — информация о пользователе внутри Recap (без лишних полей)
type RecapProfile struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
}

// RecapGeneration — метаданные генерации
type RecapGeneration struct {
	AlgorithmVersion     string              `json:"algorithm_version"`
	FeatureSchemaVersion string              `json:"feature_schema_version"`
	ActivityHash         string              `json:"activity_hash"`
	GeneratedAt          time.Time           `json:"generated_at"`
	Narrative            NarrativeGeneration `json:"narrative"`
}

// NarrativeGeneration — данные о генерации текстов
type NarrativeGeneration struct {
	Source        string  `json:"source"` // "mistral" | "template"
	PromptVersion string  `json:"prompt_version"`
	Model         *string `json:"model,omitempty"`
}

// RecapTheme — тема итогов (город)
type RecapTheme struct {
	Code         string   `json:"code"` // "city"
	MainDistrict Vertical `json:"main_district"`
	AccentToken  *string  `json:"accent_token,omitempty"` // "violet", "blue", "green", "orange"
}

// Vertical — вертикаль (район города)
type Vertical struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

// RecapCapabilities — возможности взаимодействия
type RecapCapabilities struct {
	ShareAvailable       bool `json:"share_available"`
	ExplanationAvailable bool `json:"explanation_available"`
	FeedbackAvailable    bool `json:"feedback_available"`
}

// RecapCard — универсальная карточка (discriminator по полю Type)
type RecapCard struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`       // "intro", "metric", "district", "archetype", "achievements", "summary", "final"
	Position    int                    `json:"position"`   // порядок сортировки
	Visibility  string                 `json:"visibility"` // "personal" | "shareable"
	Eyebrow     *string                `json:"eyebrow,omitempty"`
	Title       string                 `json:"title"`
	Description *string                `json:"description,omitempty"`
	Visual      CardVisual             `json:"visual"`
	Explainable bool                   `json:"explainable"`
	Data        map[string]interface{} `json:"data"`
}

// CardVisual — визуальные характеристики карточки
type CardVisual struct {
	Kind      string  `json:"kind"` // "illustration", "district", "street", "calendar", "badge", "chart", "character", "skyline"
	AssetCode *string `json:"asset_code,omitempty"`
}

// CreateRecapRequest — запрос на создание итогов
type CreateRecapRequest struct {
	ProfileID uuid.UUID `json:"profile_id"`
	Year      int       `json:"year"`
}
