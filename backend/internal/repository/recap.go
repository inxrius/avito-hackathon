package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/inxrius/avito-hackathon/internal/model"
)

// GetRecapByProfileAndYear — проверяет существование recap для пары (profile, year)
func (r *Repository) GetRecapByProfileAndYear(profileID string, year int) (*model.Recap, error) {
	query := `
		SELECT id, profile_id, status, year, algorithm_version, feature_schema_version,
			activity_hash, generated_at, schema_version, summary_title, summary_text,
			narrative_source, prompt_version, narrative_model, main_vertical_code, accent_token
		FROM recaps WHERE profile_id = $1 AND year = $2 LIMIT 1
	`
	var recap model.Recap
	var narrativeSource string
	var mainVerticalCode, accentToken sql.NullString
	err := r.DB.DB.QueryRow(query, profileID, year).Scan(
		&recap.ID, &recap.ProfileID, &recap.Status, &recap.Year,
		&recap.AlgorithmVersion, &recap.FeatureSchemaVersion,
		&recap.ActivityHash, &recap.GeneratedAt, &recap.SchemaVersion,
		&recap.SummaryTitle, &recap.SummaryText,
		&narrativeSource, &recap.PromptVersion, &recap.NarrativeModel,
		&mainVerticalCode, &accentToken,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	recap.NarrativeSource = narrativeSource
	if mainVerticalCode.Valid {
		recap.MainVerticalCode = mainVerticalCode.String
	}
	if accentToken.Valid {
		recap.AccentToken = &accentToken.String
	}
	// Загружаем профиль
	profile, err := r.GetProfileByID(profileID)
	if err != nil {
		return nil, err
	}
	recap.Profile = model.RecapProfile{
		ID:        profile.ID,
		Name:      profile.Name,
		AvatarURL: profile.AvatarURL,
	}
	// Загружаем связанные данные (карточки, архетипы, достижения)
	if err := r.loadRecapDetails(&recap); err != nil {
		return nil, err
	}
	return &recap, nil
}

// GetRecapByID — возвращает recap по ID
func (r *Repository) GetRecapByID(id string) (*model.Recap, error) {
	query := `
		SELECT id, profile_id, status, year, algorithm_version, feature_schema_version,
			activity_hash, generated_at, schema_version, summary_title, summary_text,
			narrative_source, prompt_version, narrative_model, main_vertical_code, accent_token
		FROM recaps WHERE id = $1
	`
	var recap model.Recap
	var narrativeSource string
	var mainVerticalCode, accentToken sql.NullString
	err := r.DB.DB.QueryRow(query, id).Scan(
		&recap.ID, &recap.ProfileID, &recap.Status, &recap.Year,
		&recap.AlgorithmVersion, &recap.FeatureSchemaVersion,
		&recap.ActivityHash, &recap.GeneratedAt, &recap.SchemaVersion,
		&recap.SummaryTitle, &recap.SummaryText,
		&narrativeSource, &recap.PromptVersion, &recap.NarrativeModel,
		&mainVerticalCode, &accentToken,
	)
	if err != nil {
		return nil, err
	}
	recap.NarrativeSource = narrativeSource
	if mainVerticalCode.Valid {
		recap.MainVerticalCode = mainVerticalCode.String
	}
	if accentToken.Valid {
		recap.AccentToken = &accentToken.String
	}
	profile, err := r.GetProfileByID(recap.ProfileID.String())
	if err != nil {
		return nil, err
	}
	recap.Profile = model.RecapProfile{
		ID:        profile.ID,
		Name:      profile.Name,
		AvatarURL: profile.AvatarURL,
	}
	if err := r.loadRecapDetails(&recap); err != nil {
		return nil, err
	}
	return &recap, nil
}

// CreateRecap — создаёт recap со всеми связями в одной транзакции
func (r *Repository) CreateRecap(recap *model.Recap) error {
	tx, err := r.DB.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Основная запись в recaps
	queryRecap := `
		INSERT INTO recaps (
			id, profile_id, status, year, algorithm_version, feature_schema_version,
			activity_hash, generated_at, schema_version, summary_title, summary_text,
			narrative_source, prompt_version, narrative_model, main_vertical_code, accent_token
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err = tx.Exec(queryRecap,
		recap.ID, recap.ProfileID, recap.Status, recap.Year,
		recap.AlgorithmVersion, recap.FeatureSchemaVersion,
		recap.ActivityHash, recap.GeneratedAt, recap.SchemaVersion,
		recap.SummaryTitle, recap.SummaryText,
		recap.NarrativeSource, recap.PromptVersion, recap.NarrativeModel,
		recap.MainVerticalCode, recap.AccentToken,
	)
	if err != nil {
		return err
	}

	// 2. Карточки
	for _, card := range recap.Cards {
		dataJSON, err := json.Marshal(card.Data)
		if err != nil {
			return err
		}
		queryCard := `
			INSERT INTO recap_cards (
				recap_id, card_id, type, position, visibility, eyebrow, title, description,
				visual_kind, visual_asset_code, explainable, data
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`
		_, err = tx.Exec(queryCard,
			recap.ID, card.ID, card.Type, card.Position, card.Visibility,
			card.Eyebrow, card.Title, card.Description,
			card.Visual.Kind, card.Visual.AssetCode, card.Explainable, dataJSON,
		)
		if err != nil {
			return err
		}
	}

	// 3. Метрики (из карточек типа metric)
	// for _, card := range recap.Cards {
	// 	if card.Type == "metric" {
	// 		metricCode, _ := card.Data["metric_code"].(string)
	// 		value, _ := card.Data["value"].(float64)
	// 		unit, _ := card.Data["unit"].(string)
	// 		secondaryLabel, _ := card.Data["secondary_label"].(string)
	// 		if metricCode != "" {
	// 			queryMetric := `
	// 				INSERT INTO recap_metrics (recap_id, metric_code, value, unit, secondary_label)
	// 				VALUES ($1, $2, $3, $4, $5)
	// 			`
	// 			_, err = tx.Exec(queryMetric, recap.ID, metricCode, value, unit, secondaryLabel)
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// }

	// 4. Архетип (из карточки типа archetype)
	for _, card := range recap.Cards {
		if card.Type == "archetype" {
			roleData, ok := card.Data["role"].(map[string]interface{})
			if !ok {
				continue
			}
			styleData, ok := card.Data["style"].(map[string]interface{})
			if !ok {
				continue
			}
			roleCode, _ := roleData["code"].(string)
			styleCode, _ := styleData["code"].(string)
			if roleCode != "" && styleCode != "" {
				queryArchetype := `
					INSERT INTO recap_archetypes (recap_id, role_code, style_code)
					VALUES ($1, $2, $3)
				`
				_, err = tx.Exec(queryArchetype, recap.ID, roleCode, styleCode)
				if err != nil {
					return err
				}
			}
			break
		}
	}

	// 5. Достижения (из карточки типа achievements)
	for _, card := range recap.Cards {
		if card.Type == "achievements" {
			items, ok := card.Data["items"].([]interface{})
			if !ok {
				continue
			}
			for idx, item := range items {
				achMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				ach := model.Achievement{
					Code:        achMap["code"].(string),
					Title:       achMap["title"].(string),
					Description: achMap["description"].(string),
					Level:       model.AchievementLevel(achMap["level"].(string)),
					Icon:        achMap["icon"].(string),
				}
				if mc, ok := achMap["metric_code"].(string); ok && mc != "" {
					ach.MetricCode = &mc
				}
				if cv, ok := achMap["current_value"].(float64); ok {
					ach.CurrentValue = &cv
				}
				if nt, ok := achMap["next_level_threshold"].(float64); ok {
					ach.NextLevelThreshold = &nt
				}
				pos := idx + 1
				queryAchievement := `
					INSERT INTO recap_achievements (
						recap_id, achievement_code, level, position, title, description, icon_code,
						metric_code, current_value, next_level_threshold
					) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				`
				_, err = tx.Exec(queryAchievement,
					recap.ID, ach.Code, ach.Level, pos, ach.Title, ach.Description, ach.Icon,
					ach.MetricCode, ach.CurrentValue, ach.NextLevelThreshold,
				)
				if err != nil {
					return err
				}
			}
			break
		}
	}

	// TODO: сохранение share_card, если нужно (пока пропустим)

	return tx.Commit()
}

// loadRecapDetails — загружает связанные данные (карточки, архетипы, достижения) для заполнения Recap
func (r *Repository) loadRecapDetails(recap *model.Recap) error {
	// Карточки
	rows, err := r.DB.DB.Query(`
		SELECT card_id, type, position, visibility, eyebrow, title, description,
			visual_kind, visual_asset_code, explainable, data
		FROM recap_cards WHERE recap_id = $1 ORDER BY position
	`, recap.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var card model.RecapCard
		var visualKind, visualAssetCode, description sql.NullString
		var visibility string
		var explainable bool
		var data []byte
		err := rows.Scan(
			&card.ID, &card.Type, &card.Position, &visibility, &card.Eyebrow,
			&card.Title, &description, &visualKind, &visualAssetCode,
			&explainable, &data,
		)
		if err != nil {
			return err
		}
		card.Visibility = visibility
		card.Explainable = explainable
		if visualKind.Valid {
			card.Visual.Kind = visualKind.String
		}
		if visualAssetCode.Valid {
			card.Visual.AssetCode = &visualAssetCode.String
		}
		if description.Valid {
			card.Description = &description.String
		}
		if len(data) > 0 {
			var dataMap map[string]interface{}
			if err := json.Unmarshal(data, &dataMap); err != nil {
				return err
			}
			card.Data = dataMap
		} else {
			card.Data = make(map[string]interface{})
		}
		recap.Cards = append(recap.Cards, card)
	}

	// Архетип
	var roleCode, styleCode string
	err = r.DB.DB.QueryRow(`
		SELECT role_code, style_code FROM recap_archetypes WHERE recap_id = $1
	`, recap.ID).Scan(&roleCode, &styleCode)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		// Загружаем названия
		var roleTitle, styleTitle string
		_ = r.DB.DB.QueryRow(`SELECT title FROM archetype_roles WHERE code = $1`, roleCode).Scan(&roleTitle)
		_ = r.DB.DB.QueryRow(`SELECT title FROM archetype_styles WHERE code = $1`, styleCode).Scan(&styleTitle)
		// Ищем карточку archetype и вставляем данные
		for i, card := range recap.Cards {
			if card.Type == "archetype" {
				recap.Cards[i].Data = map[string]interface{}{
					"role":  model.ArchetypeRole{Code: roleCode, Title: roleTitle},
					"style": model.ArchetypeStyle{Code: styleCode, Title: styleTitle},
				}
				break
			}
		}
	}

	// Достижения
	rows2, err := r.DB.DB.Query(`
		SELECT achievement_code, level, position, title, description, icon_code,
			metric_code, current_value, next_level_threshold
		FROM recap_achievements WHERE recap_id = $1 ORDER BY position
	`, recap.ID)
	if err != nil {
		return err
	}
	defer rows2.Close()
	var achievements []model.Achievement
	for rows2.Next() {
		var ach model.Achievement
		var metricCode, currentValue, nextThreshold sql.NullString
		var level string
		err := rows2.Scan(
			&ach.Code, &level, &ach.Position, &ach.Title, &ach.Description,
			&ach.Icon, &metricCode, &currentValue, &nextThreshold,
		)
		if err != nil {
			return err
		}
		ach.Level = model.AchievementLevel(level)
		if metricCode.Valid {
			ach.MetricCode = &metricCode.String
		}
		if currentValue.Valid {
			var v float64
			fmt.Sscan(currentValue.String, &v)
			ach.CurrentValue = &v
		}
		if nextThreshold.Valid {
			var v float64
			fmt.Sscan(nextThreshold.String, &v)
			ach.NextLevelThreshold = &v
		}
		achievements = append(achievements, ach)
	}
	// Ищем карточку achievements
	for i, card := range recap.Cards {
		if card.Type == "achievements" {
			recap.Cards[i].Data = map[string]interface{}{
				"items":       achievements,
				"total_count": len(achievements),
			}
			break
		}
	}

	return nil
}
