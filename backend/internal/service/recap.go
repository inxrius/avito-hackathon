package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"recap-personalization/internal/model"
	recap "recap-personalization/internal/recap"
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

	// 2. Загружаем активности за указанный год из ClickHouse
	activities, err := s.clickHouseActivities.GetActivitiesByProfileIDAndYear(profileID, year)
	if err != nil {
		return nil, err
	}
	profile, err := s.repo.GetProfileByID(profileID)
	if err != nil {
		return nil, err
	}

	// 3. Конвертируем активности в формат recap
	recapEvents := make([]recap.ActivityEvent, 0, len(activities))
	for _, a := range activities {
		recapEvents = append(recapEvents, recap.ActivityEvent{
			EventID:      a.EventID,
			ProfileID:    a.ProfileID,
			EventType:    convertEventType(a.EventType),
			CategoryCode: a.CategoryCode,
			OccurredAt:   time.Unix(a.OccurredAt, 0),
		})
	}

	// 4. Генерируем recap через recap pipeline
	output, err := s.recapGenerator.Generate(context.Background(), recap.GenerateInput{
		RecapID:     generateRecapID(),
		Profile:     convertProfile(profile),
		Year:        year,
		Activities:  recapEvents,
		GeneratedAt: time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	// 5. Конвертируем результат в нашу модель (включая explanation и share)
	recapModel := convertRecapOutput(&output, profile)
	
	// 6. Сохраняем recap
	if err := s.repo.CreateRecap(recapModel); err != nil {
		return nil, err
	}

	return recapModel, nil
}

// convertEventType конвертирует тип события из ClickHouse в модель recap
func convertEventType(eventType string) recap.EventType {
	switch eventType {
	case "favorite_added":
		return recap.EventFavoriteAdded
	case "listing_viewed":
		return recap.EventListingViewed
	case "purchase_completed", "sale_completed":
		return recap.EventPurchaseCompleted
	case "chat_started":
		return recap.EventChatStarted
	case "listing_published":
		return recap.EventListingPublished
	default:
		return recap.EventListingViewed
	}
}

// convertCategory конвертирует категорию в код категории recap
func convertCategory(category string) string {
	categoryMap := map[string]string{
		"Недвижимость": "apartments",
		"Авто":         "cars",
		"Электроника":  "electronics",
		"Для дома":     "home_and_garden",
	}
	if code, ok := categoryMap[category]; ok {
		return code
	}
	return "electronics"
}

// convertProfile конвертирует профиль в формат recap
func convertProfile(profile *model.Profile) recap.ProfileSnapshot {
	avatarURL := ""
	if profile.AvatarURL != nil {
		avatarURL = *profile.AvatarURL
	}
	return recap.ProfileSnapshot{
		ID:        profile.ID.String(),
		Name:      profile.Name,
		AvatarURL: avatarURL,
	}
}

// generateRecapID генерирует UUID для recap
func generateRecapID() string {
	return uuid.New().String()
}

// convertRecapOutput конвертирует вывод pipeline в нашу модель Recap
func convertRecapOutput(output *recap.GenerateOutput, profile *model.Profile) *model.Recap {
	accentToken := ""
	if output.Recap.Theme.AccentToken != nil {
		accentToken = string(*output.Recap.Theme.AccentToken)
	}
	
	recap := &model.Recap{
		ID:                   uuid.MustParse(output.Recap.ID),
		ProfileID:            uuid.MustParse(output.Recap.ProfileID),
		Status:               "completed",
		Year:                 output.Recap.Year,
		SchemaVersion:        output.Recap.SchemaVersion,
		AlgorithmVersion:     output.Recap.Generation.AlgorithmVersion,
		FeatureSchemaVersion: output.Recap.Generation.FeatureSchemaVersion,
		ActivityHash:         output.Recap.Generation.ActivityHash,
		GeneratedAt:          output.Recap.Generation.GeneratedAt,
		NarrativeSource:      output.Recap.Generation.Narrative.Source,
		PromptVersion:        output.Recap.Generation.Narrative.PromptVersion,
		NarrativeModel:       output.Recap.Generation.Narrative.Model,
		MainVerticalCode:     string(output.Recap.Theme.MainDistrict.Code),
		AccentToken:          &accentToken,
		SummaryTitle:         output.Narrative.SummaryTitle,
		SummaryText:          output.Narrative.SummaryText,
		Profile: model.RecapProfile{
			ID:        uuid.MustParse(output.Recap.Profile.ID),
			Name:      output.Recap.Profile.Name,
			AvatarURL: output.Recap.Profile.AvatarURL,
		},
		Generation: model.RecapGeneration{
			AlgorithmVersion:     output.Recap.Generation.AlgorithmVersion,
			FeatureSchemaVersion: output.Recap.Generation.FeatureSchemaVersion,
			ActivityHash:         output.Recap.Generation.ActivityHash,
			GeneratedAt:          output.Recap.Generation.GeneratedAt,
			Narrative: model.NarrativeGeneration{
				Source:        output.Recap.Generation.Narrative.Source,
				PromptVersion: output.Recap.Generation.Narrative.PromptVersion,
				Model:         output.Recap.Generation.Narrative.Model,
			},
		},
		Theme: model.RecapTheme{
			Code:         output.Recap.Theme.Code,
			MainDistrict: convertVertical(output.Recap.Theme.MainDistrict),
			AccentToken:  &accentToken,
		},
		Cards:        convertCards(output.Recap.Cards),
		Capabilities: convertCapabilities(output.Recap.Capabilities),
	}
	
	// Добавляем explanation и share в recap
	recap.Explanation = convertExplanation(output.Explanation)
	recap.Share = convertShare(output.Share)

	return recap
}

// convertVertical конвертирует Vertical из recap в model
func convertVertical(v recap.Vertical) model.Vertical {
	return model.Vertical{
		Code:  string(v.Code),
		Title: v.Title,
	}
}

// convertCards конвертирует карточки из модели recap в model
func convertCards(cards []recap.RecapCard) []model.RecapCard {
	result := make([]model.RecapCard, 0, len(cards))
	for _, card := range cards {
		base := card.Base()
		dataMap := make(map[string]interface{})
		
		// Извлекаем данные из конкретного типа карточки
		switch c := card.(type) {
		case recap.IntroCard:
			dataMap["year"] = c.Data.Year
		case recap.MetricCard:
			dataMap["metric_code"] = c.Data.MetricCode
			dataMap["value"] = c.Data.Value
			dataMap["unit"] = c.Data.Unit
		case recap.DistrictCard:
			dataMap["vertical"] = c.Data.Vertical
			dataMap["activity_share"] = c.Data.ActivityShare
			if c.Data.TopCategory != nil {
				dataMap["top_category"] = c.Data.TopCategory
			}
		case recap.ArchetypeCard:
			dataMap["role"] = c.Data.Role
			dataMap["style"] = c.Data.Style
		case recap.AchievementsCard:
			dataMap["items"] = c.Data.Items
			dataMap["total_count"] = c.Data.TotalCount
		case recap.FinalCard:
			dataMap["show_share_button"] = c.Data.ShowShareButton
			dataMap["show_feedback"] = c.Data.ShowFeedback
		}
		
		result = append(result, model.RecapCard{
			ID:          base.ID,
			Type:        string(base.Type),
			Position:    base.Position,
			Visibility:  string(base.Visibility),
			Eyebrow:     base.Eyebrow,
			Title:       base.Title,
			Description: base.Description,
			Visual:      convertVisual(base.Visual),
			Explainable: base.Explainable,
			Data:        dataMap,
		})
	}
	return result
}

// convertVisual конвертирует визуал карточки
func convertVisual(v *recap.CardVisual) model.CardVisual {
	if v == nil {
		return model.CardVisual{}
	}
	return model.CardVisual{
		Kind:      string(v.Kind),
		AssetCode: v.AssetCode,
	}
}

// convertCapabilities конвертирует возможности
func convertCapabilities(c recap.Capabilities) model.RecapCapabilities {
	return model.RecapCapabilities{
		ShareAvailable:       c.ShareAvailable,
		ExplanationAvailable: c.ExplanationAvailable,
		FeedbackAvailable:    c.FeedbackAvailable,
	}
}

// convertExplanation конвертирует объяснение из модели recap
func convertExplanation(explanation *recap.RecapExplanation) *model.RecapExplanation {
	if explanation == nil {
		return nil
	}
	
	decisions := make([]model.DecisionExplanation, 0, len(explanation.Decisions))
	for _, d := range explanation.Decisions {
		facts := make([]model.RuleFact, 0, len(d.Facts))
		for _, f := range d.Facts {
			facts = append(facts, model.RuleFact{
				MetricCode: f.MetricCode,
				Actual:     f.Actual,
				Operator:   f.Operator,
				Threshold:  f.Threshold,
				Matched:    f.Matched,
			})
		}
		decisions = append(decisions, model.DecisionExplanation{
			CardID:      d.CardID,
			Kind:        d.Kind,
			Code:        d.Code,
			Reason:      d.Reason,
			RuleVersion: d.RuleVersion,
			Facts:       facts,
		})
	}
	
	return &model.RecapExplanation{
		RecapID:          explanation.RecapID,
		AlgorithmVersion: explanation.AlgorithmVersion,
		ActivityHash:     explanation.ActivityHash,
		Decisions:        decisions,
	}
}

// convertShare конвертирует share карточку из модели recap
func convertShare(share *recap.ShareCard) *recap.ShareCard {
	if share == nil {
		return nil
	}
	
	avatarURL := ""
	if share.AvatarURL != nil {
		avatarURL = *share.AvatarURL
	}
	
	return &recap.ShareCard{
		SchemaVersion: share.SchemaVersion,
		RecapID:       share.RecapID,
		ProfileName:   share.ProfileName,
		AvatarURL:     &avatarURL,
		Year:          share.Year,
		Title:         share.Title,
		Subtitle:      share.Subtitle,
		MainDistrict:  share.MainDistrict,
		Facts:         share.Facts,
		Achievements:  share.Achievements,
		Visual:        share.Visual,
	}
}

// strPtr возвращает указатель на строку
func strPtr(s string) *string {
	return &s
}
