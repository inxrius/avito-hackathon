package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	DecisionKindArchetypeRole  DecisionKind = "archetype_role"
	DecisionKindArchetypeStyle DecisionKind = "archetype_style"
	DecisionKindAchievement    DecisionKind = "achievement"
)

const (
	OpEq  ComparisonOperator = "eq"
	OpNeq ComparisonOperator = "neq"
	OpGt  ComparisonOperator = "gt"
	OpGte ComparisonOperator = "gte"
	OpLt  ComparisonOperator = "lt"
	OpLte ComparisonOperator = "lte"
)

// RecapExplanation — полное объяснение для recap (GET /recaps/{id}/explanation)
type RecapExplanation struct {
	RecapID          uuid.UUID             `json:"recap_id"`
	AlgorithmVersion string                `json:"algorithm_version"`
	ActivityHash     string                `json:"activity_hash"`
	Decisions        []DecisionExplanation `json:"decisions"`
}

// DecisionExplanation — объяснение одного решения (карточки)
type DecisionExplanation struct {
	CardID      string     `json:"card_id"`      // ID карточки, к которой относится объяснение
	Kind        string     `json:"kind"`         // "archetype_role", "archetype_style", "achievement"
	Code        string     `json:"code"`         // код архетипа/роли/достижения (один из enum)
	Reason      string     `json:"reason"`       // человекочитаемое объяснение
	RuleVersion string     `json:"rule_version"` // версия правил
	Facts       []RuleFact `json:"facts"`        // факты, на основе которых принято решение
}

// RuleFact — факт правила (одна метрика с порогом)
type RuleFact struct {
	MetricCode string  `json:"metric_code"`
	Actual     float64 `json:"actual"`    // фактическое значение
	Operator   string  `json:"operator"`  // "eq", "neq", "gt", "gte", "lt", "lte"
	Threshold  float64 `json:"threshold"` // пороговое значение
	Matched    bool    `json:"matched"`   // совпало ли правило
}

// RecapExplanationDB — внутренняя структура для сохранения объяснения в БД
// (не экспортируется в API, используется в репозитории/сервисе)
type RecapExplanationDB struct {
	ID              uuid.UUID    `json:"id"`
	RecapID         uuid.UUID    `json:"recap_id"`
	CardID          string       `json:"card_id"`
	Kind            DecisionKind `json:"kind"`
	RoleCode        *string      `json:"role_code,omitempty"`        // если kind = archetype_role
	StyleCode       *string      `json:"style_code,omitempty"`       // если kind = archetype_style
	AchievementCode *string      `json:"achievement_code,omitempty"` // если kind = achievement
	Reason          string       `json:"reason"`
	RuleVersion     string       `json:"rule_version"`
	CreatedAt       time.Time    `json:"created_at"`
}

// RuleFactDB — внутренняя структура для сохранения факта в БД
type RuleFactDB struct {
	ExplanationID uuid.UUID          `json:"explanation_id"`
	MetricCode    string             `json:"metric_code"`
	Actual        float64            `json:"actual"`
	Operator      ComparisonOperator `json:"operator"`
	Threshold     float64            `json:"threshold"`
	Matched       bool               `json:"matched"`
}

// DecisionKind — тип решения (соответствует enum в БД)
type DecisionKind string

// ComparisonOperator — оператор сравнения для правил
type ComparisonOperator string
