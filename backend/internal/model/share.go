package model

import (
	"github.com/google/uuid"
)

// ShareCard — публичная карточка итогов (без приватных данных)
// Возвращается по GET /recaps/{id}/share
type ShareCard struct {
	SchemaVersion string             `json:"schema_version"` // "2.0"
	RecapID       uuid.UUID          `json:"recap_id"`
	ProfileName   string             `json:"profile_name"` // только имя, без ID
	Year          int                `json:"year"`
	Title         string             `json:"title"`         // например, "Мой год в городе Авито"
	Subtitle      string             `json:"subtitle"`      // например, "Вдумчивый искатель находок"
	MainDistrict  Vertical           `json:"main_district"` // главный район
	Facts         []ShareFact        `json:"facts"`         // факты для отображения
	Achievements  []ShareAchievement `json:"achievements"`  // публичные достижения (без уровней)
	Visual        ShareVisual        `json:"visual"`        // визуальное оформление
}

// ShareFact — факт в публичной карточке
type ShareFact struct {
	Kind  string `json:"kind"`  // "main_district", "active_days", "top_achievement"
	Label string `json:"label"` // например, "Главный район"
	Value string `json:"value"` // например, "Товары"
}

// ShareAchievement — публичное достижение (без полной информации)
type ShareAchievement struct {
	Code  string `json:"code"`  // код достижения (например, "findings_collector")
	Title string `json:"title"` // название
	Level string `json:"level"` // уровень (newcomer, local, expert, guru)
	Icon  string `json:"icon"`  // иконка
}

// ShareVisual — визуальное оформление карточки
type ShareVisual struct {
	Theme  string   `json:"theme"`  // "city"
	Colors []string `json:"colors"` // массив hex-цветов (например, ["#7C3AED", "#F3E8FF"])
}
