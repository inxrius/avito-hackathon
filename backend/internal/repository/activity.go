package repository

import (
	"github.com/inxrius/avito-hackathon/internal/model"
)

// GetActivitiesByProfileID — возвращает активности профиля
func (r *Repository) GetActivitiesByProfileID(profileID string) ([]model.Activity, error) {
	query := `
		SELECT id, profile_id, type, category, title, description, value, timestamp
		FROM activities WHERE profile_id = $1 ORDER BY timestamp DESC
	`
	rows, err := r.DB.DB.Query(query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []model.Activity
	for rows.Next() {
		var a model.Activity
		if err := rows.Scan(&a.ID, &a.ProfileID, &a.Type, &a.Category, &a.Title, &a.Description, &a.Value, &a.Timestamp); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}
