package model

// ArchetypeRole — роль (архетип) пользователя
type ArchetypeRole struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

// ArchetypeStyle — стиль (архетип) пользователя
type ArchetypeStyle struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

// ArchetypeData — структура данных для карточки типа "archetype"
// Используется в RecapCard.Data для типа "archetype"
type ArchetypeData struct {
	Role  ArchetypeRole  `json:"role"`
	Style ArchetypeStyle `json:"style"`
}

// Предопределённые коды ролей (из БД)
const (
	ArchetypeRoleFindingsSeeker   = "findings_seeker"
	ArchetypeRoleShowcaseOwner    = "showcase_owner"
	ArchetypeRoleUniversalCitizen = "universal_citizen"
	ArchetypeRoleCityObserver     = "city_observer"
)

// Предопределённые коды стилей (из БД)
const (
	ArchetypeStyleThoughtful     = "thoughtful"
	ArchetypeStyleExplorer       = "explorer"
	ArchetypeStyleDistrictExpert = "district_expert"
	ArchetypeStyleRegular        = "regular"
	ArchetypeStyleResultOriented = "result_oriented"
	ArchetypeStyleCityLocal      = "city_local"
)
