package repository

import (
	"github.com/inxrius/avito-hackathon/internal/model"
)

// GetProfiles — возвращает список профилей с доступными годами
func (r *Repository) GetProfiles() ([]model.ProfileSummary, error) {
	query := `SELECT id, name, description, avatar_url, scenario FROM profiles ORDER BY name`
	rows, err := r.DB.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []model.ProfileSummary
	for rows.Next() {
		var p model.ProfileSummary
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.AvatarURL, &p.Scenario); err != nil {
			return nil, err
		}
		yearsQuery := `SELECT year FROM profile_available_years WHERE profile_id = $1 ORDER BY year`
		yearRows, err := r.DB.DB.Query(yearsQuery, p.ID)
		if err != nil {
			return nil, err
		}
		var years []int
		for yearRows.Next() {
			var y int
			if err := yearRows.Scan(&y); err != nil {
				yearRows.Close()
				return nil, err
			}
			years = append(years, y)
		}
		yearRows.Close()
		p.AvailableYears = years
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// GetProfileByID — возвращает полный профиль
func (r *Repository) GetProfileByID(id string) (*model.Profile, error) {
	query := `
		SELECT id, name, description, avatar_url, scenario, created_at, updated_at
		FROM profiles WHERE id = $1
	`
	var p model.Profile
	err := r.DB.DB.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.AvatarURL, &p.Scenario, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.DB.Query(`SELECT year FROM profile_available_years WHERE profile_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var years []int
	for rows.Next() {
		var y int
		if err := rows.Scan(&y); err != nil {
			return nil, err
		}
		years = append(years, y)
	}
	p.AvailableYears = years
	return &p, nil
}
