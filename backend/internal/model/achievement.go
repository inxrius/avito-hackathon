package model

// Achievement — достижение пользователя (используется внутри RecapCard типа "achievements")
type Achievement struct {
	Code               string           `json:"code"`                           // уникальный код достижения (из achievement_definitions)
	Title              string           `json:"title"`                          // название
	Description        string           `json:"description"`                    // описание
	Level              AchievementLevel `json:"level"`                          // "newcomer", "local", "expert", "guru"
	Icon               string           `json:"icon"`                           // иконка (code)
	MetricCode         *string          `json:"metric_code,omitempty"`          // метрика, по которой считается достижение
	CurrentValue       *float64         `json:"current_value,omitempty"`        // текущее значение метрики
	NextLevelThreshold *float64         `json:"next_level_threshold,omitempty"` // порог для следующего уровня
	Position           int              `json:"-"`                              // позиция в списке (не сериализуется, используется для БД)
}

// AchievementDefinition — справочная информация о достижении (хранится в БД)
type AchievementDefinition struct {
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IconCode    string `json:"icon_code"`
}

// AchievementLevel — тип для уровня (можно использовать как enum)
type AchievementLevel string

const (
	LevelNewcomer AchievementLevel = "newcomer"
	LevelLocal    AchievementLevel = "local"
	LevelExpert   AchievementLevel = "expert"
	LevelGuru     AchievementLevel = "guru"
)

// AchievementCode — тип для кода достижения (справочник)
type AchievementCode string

// Примеры кодов (для удобства)
const (
	AchievementDealMaster        AchievementCode = "deal_master"
	AchievementFindingsCollector AchievementCode = "findings_collector"
	AchievementCityNavigator     AchievementCode = "city_navigator"
	AchievementFrequentGuest     AchievementCode = "frequent_guest"
	AchievementOldTimer          AchievementCode = "old_timer"
	AchievementDoorstepDelivery  AchievementCode = "doorstep_delivery"
	AchievementOwnShowcase       AchievementCode = "own_showcase"
	AchievementFindingsHunter    AchievementCode = "findings_hunter"
	AchievementCityRhythm        AchievementCode = "city_rhythm"
)
