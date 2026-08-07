package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/inxrius/avito-hackathon/internal/model"
)

// GenerateRecap — основная бизнес-логика генерации итогов
func (s *Service) GenerateRecap(profileID string, year int) (*model.Recap, error) {
	// 1. Проверка идемпотентности: если recap уже есть — возвращаем его
	existing, err := s.repo.GetRecapByProfileAndYear(profileID, year)
	if err == nil && existing != nil {
		return existing, nil
	}
	if err != nil && err.Error() != "recap_not_found" {
		return nil, err
	}

	// 2. Загружаем активности и профиль
	activities, err := s.repo.GetActivitiesByProfileID(profileID)
	if err != nil {
		return nil, err
	}
	profile, err := s.repo.GetProfileByID(profileID)
	if err != nil {
		return nil, err
	}

	// 3. Вычисляем хеш активностей
	activityHash := computeActivityHash(activities)

	// 4. Формируем карточки
	cards := buildCards(profile, activities, year)

	// 5. Определяем главный район
	mainDistrict := determineMainDistrict(activities)

	// 6. Собираем Recap
	recap := &model.Recap{
		ID:                   uuid.New(),
		ProfileID:            profile.ID,
		Status:               "completed",
		Year:                 year,
		SchemaVersion:        "2.0",
		AlgorithmVersion:     "recap-rules-2026.08.2",
		FeatureSchemaVersion: "features-v1",
		ActivityHash:         activityHash,
		GeneratedAt:          time.Now().UTC(),
		NarrativeSource:      "template",
		PromptVersion:        "city-summary-v1",
		NarrativeModel:       nil,
		MainVerticalCode:     mainDistrict.Code,
		AccentToken:          strPtr("violet"),
		SummaryTitle:         fmt.Sprintf("Год в городе %s", profile.Name),
		SummaryText:          buildSummaryText(profile, activities),
		Profile: model.RecapProfile{
			ID:        profile.ID,
			Name:      profile.Name,
			AvatarURL: profile.AvatarURL,
		},
		Generation: model.RecapGeneration{
			AlgorithmVersion:     "recap-rules-2026.08.2",
			FeatureSchemaVersion: "features-v1",
			ActivityHash:         activityHash,
			GeneratedAt:          time.Now().UTC(),
			Narrative: model.NarrativeGeneration{
				Source:        "template",
				PromptVersion: "city-summary-v1",
				Model:         nil,
			},
		},
		Theme: model.RecapTheme{
			Code:         "city",
			MainDistrict: mainDistrict,
			AccentToken:  strPtr("violet"),
		},
		Cards: cards,
		Capabilities: model.RecapCapabilities{
			ShareAvailable:       true,
			ExplanationAvailable: true,
			FeedbackAvailable:    true,
		},
	}

	// 7. Сохраняем recap (включая карточки, метрики, архетипы, достижения)
	if err := s.repo.CreateRecap(recap); err != nil {
		return nil, err
	}

	return recap, nil
}

// computeActivityHash — хеш активностей для детекции изменений
func computeActivityHash(activities []model.Activity) string {
	h := sha256.New()
	for _, a := range activities {
		fmt.Fprintf(h, "%s|%s|%s|%f|%d|", a.ID, a.Type, a.Category, a.Value, a.Timestamp.UnixNano())
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}

// determineMainDistrict — определяет главный район на основе активностей
func determineMainDistrict(activities []model.Activity) model.Vertical {
	districtCounts := make(map[string]int)
	for _, a := range activities {
		districtCounts[a.Category]++
	}

	maxCount := 0
	topDistrict := "goods"
	districtMap := map[string]string{
		"Недвижимость": "real_estate",
		"Авто":         "transport",
		"Электроника":  "goods",
		"Для дома":     "goods",
	}

	for district, count := range districtCounts {
		if count > maxCount {
			maxCount = count
			if code, ok := districtMap[district]; ok {
				topDistrict = code
			} else {
				topDistrict = "goods"
			}
		}
	}

	titleMap := map[string]string{
		"real_estate": "Недвижимость",
		"transport":   "Транспорт",
		"goods":       "Товары",
	}

	return model.Vertical{
		Code:  topDistrict,
		Title: titleMap[topDistrict],
	}
}

// buildSummaryText — формирует краткое описание итогов
func buildSummaryText(profile *model.Profile, activities []model.Activity) string {
	archetype := resolveArchetype(profile.Scenario, activities)
	return fmt.Sprintf("%s — %s", profile.Name, archetype.Description)
}

// buildCards — собирает все карточки для recap
func buildCards(profile *model.Profile, activities []model.Activity, year int) []model.RecapCard {
	cards := []model.RecapCard{
		{
			ID:          "intro",
			Type:        "intro",
			Position:    0,
			Visibility:  "shareable",
			Title:       fmt.Sprintf("Город %s", profile.Name),
			Explainable: false,
			Visual:      model.CardVisual{Kind: "skyline", AssetCode: strPtr("city-intro")},
			Data: map[string]interface{}{
				"year":        year,
				"subtitle":    fmt.Sprintf("%d действий · %d района · %d строек на окраине", len(activities), countDistricts(activities), countUnfinished(activities)),
				"archetype":   profile.Scenario,
				"description": profile.Description,
			},
		},
	}

	// Метрики (районы)
	metrics := aggregateMetrics(activities)
	cardID := 1
	for district, value := range metrics {
		cards = append(cards, model.RecapCard{
			ID:          fmt.Sprintf("metric-%d", cardID),
			Type:        "metric",
			Position:    cardID,
			Visibility:  "shareable",
			Title:       district,
			Explainable: false,
			Visual:      model.CardVisual{Kind: "district", AssetCode: strPtr("district-" + district)},
			Data: map[string]interface{}{
				"metric_code": "category_" + district,
				"value":       value,
				"unit":        "действий",
			},
		})
		cardID++
	}

	// Архетип
	archetype := resolveArchetype(profile.Scenario, activities)
	cards = append(cards, model.RecapCard{
		ID:          "archetype",
		Type:        "archetype",
		Position:    cardID,
		Visibility:  "shareable",
		Title:       "Твой архетип",
		Explainable: true,
		Visual:      model.CardVisual{Kind: "character", AssetCode: strPtr("archetype-" + profile.Scenario)},
		Data: map[string]interface{}{
			"role":  model.ArchetypeRole{Code: profile.Scenario, Title: archetype.Name},
			"style": model.ArchetypeStyle{Code: "thoughtful", Title: "Вдумчивый"},
		},
	})
	cardID++

	// Достижения
	achievements := buildAchievements(activities)
	if len(achievements) > 0 {
		cards = append(cards, model.RecapCard{
			ID:          "achievements",
			Type:        "achievements",
			Position:    cardID,
			Visibility:  "shareable",
			Title:       "Твои городские звания",
			Explainable: true,
			Visual:      model.CardVisual{Kind: "badge", AssetCode: strPtr("top-achievements")},
			Data: map[string]interface{}{
				"items":       achievements,
				"total_count": len(achievements),
			},
		})
		cardID++
	}

	// Возможности (недостроенные элементы)
	opportunities := buildOpportunities(activities)
	for _, opp := range opportunities {
		cards = append(cards, model.RecapCard{
			ID:          fmt.Sprintf("opportunity-%d", cardID),
			Type:        "opportunity",
			Position:    cardID,
			Visibility:  "personal",
			Title:       opp.Title,
			Explainable: false,
			Visual:      model.CardVisual{Kind: "illustration", AssetCode: nil},
			Data: map[string]interface{}{
				"description": opp.Description,
				"action":      opp.Action,
			},
		})
		cardID++
	}

	// Финальная карточка
	cards = append(cards, model.RecapCard{
		ID:          "final",
		Type:        "final",
		Position:    cardID,
		Visibility:  "shareable",
		Title:       "До встречи в городе",
		Explainable: false,
		Visual:      model.CardVisual{Kind: "skyline", AssetCode: strPtr("city-final")},
		Data: map[string]interface{}{
			"show_share_button": true,
			"show_feedback":     true,
			"actions": []map[string]string{
				{"label": "Вернуться в избранное", "url": "/favorites"},
				{"label": "Открыть рекомендации", "url": "/recommendations"},
			},
		},
	})

	// Сортируем по позиции
	sort.Slice(cards, func(i, j int) bool { return cards[i].Position < cards[j].Position })
	return cards
}

// Вспомогательные функции для подсчётов

func countDistricts(activities []model.Activity) int {
	districts := map[string]struct{}{}
	for _, a := range activities {
		districts[a.Category] = struct{}{}
	}
	return len(districts)
}

func countUnfinished(activities []model.Activity) int {
	count := 0
	for _, a := range activities {
		if a.Type == "favorite" || a.Type == "view" {
			count++
		}
	}
	return count
}

func aggregateMetrics(activities []model.Activity) map[string]float64 {
	metrics := map[string]float64{}
	for _, a := range activities {
		metrics[a.Category] += a.Value
	}
	return metrics
}

type archetypeResult struct {
	Name        string
	Description string
	Badges      []string
}

func resolveArchetype(scenario string, activities []model.Activity) archetypeResult {
	switch scenario {
	case "seller":
		return archetypeResult{
			Name:        "Мастер обмена",
			Description: "Ты превращаешь лишнее в нужное и делаешь площадку живее.",
			Badges:      []string{"Продавец года", "Надёжный партнёр"},
		}
	case "unfinished":
		return archetypeResult{
			Name:        "Исследователь возможностей",
			Description: "Ты собираешь варианты и готовишься к большому выбору.",
			Badges:      []string{"Коллекционер идей", "Внимательный планировщик"},
		}
	case "insufficient":
		return archetypeResult{
			Name:        "Новый житель",
			Description: "Город только начинает строиться, и у тебя впереди много интересных сценариев.",
			Badges:      []string{"Первый шаг"},
		}
	default:
		return archetypeResult{
			Name:        "Охотник за домом",
			Description: "Ты долго выбираешь, но когда наступает время — начинаешь обустраивать.",
			Badges:      []string{"Ранняя пташка", "Мастер торга", "Новосёл"},
		}
	}
}

type achievement struct {
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Level       string  `json:"level"`
	Icon        string  `json:"icon"`
	MetricCode  *string `json:"metric_code,omitempty"`
}

func buildAchievements(activities []model.Activity) []achievement {
	result := []achievement{}
	if len(activities) > 0 {
		result = append(result, achievement{
			Code:        "first_step",
			Title:       "Первый взгляд",
			Description: "Ты начал исследовать город.",
			Level:       "newcomer",
			Icon:        "👁",
		})
	}
	if countFavorites(activities) >= 5 {
		result = append(result, achievement{
			Code:        "collector",
			Title:       "Коллекционер",
			Description: "Ты сохранил много понравившегося.",
			Level:       "local",
			Icon:        "⭐",
		})
	}
	if countPurchases(activities) >= 3 {
		result = append(result, achievement{
			Code:        "active_buyer",
			Title:       "Активный покупатель",
			Description: "Ты совершил несколько значимых покупок.",
			Level:       "expert",
			Icon:        "🛒",
		})
	}
	return result
}

type opportunity struct {
	Title       string
	Description string
	Action      string
}

func buildOpportunities(activities []model.Activity) []opportunity {
	result := []opportunity{}
	if countUnfinished(activities) > 0 {
		result = append(result, opportunity{
			Title:       "Незавершённые сценарии",
			Description: "У тебя есть сохранённые варианты, которые ещё не обсуждены.",
			Action:      "Вернуться к диалогам",
		})
	}
	if countCategories(activities) < 3 {
		result = append(result, opportunity{
			Title:       "Новый район",
			Description: "Попробуй explore категорию, которую ты ещё не посещал.",
			Action:      "Открыть рекомендации",
		})
	}
	return result
}

func countFavorites(activities []model.Activity) int {
	count := 0
	for _, a := range activities {
		if a.Type == "favorite" {
			count++
		}
	}
	return count
}

func countPurchases(activities []model.Activity) int {
	count := 0
	for _, a := range activities {
		if a.Type == "purchase" || a.Type == "sale" {
			count++
		}
	}
	return count
}

func countCategories(activities []model.Activity) int {
	cats := map[string]struct{}{}
	for _, a := range activities {
		cats[a.Category] = struct{}{}
	}
	return len(cats)
}

// strPtr возвращает указатель на строку
func strPtr(s string) *string {
	return &s
}
