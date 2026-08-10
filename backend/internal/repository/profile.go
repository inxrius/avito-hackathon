package repository

import (
	"context"
	"database/sql"
	"errors"

	"recap-personalization/internal/model"
)

var ErrProfileNotFound = errors.New("profile_not_found")

func (r *Repository) GetProfiles(ctx context.Context) ([]model.ProfileSummary, error) {
	rows, err := r.DB.DB.QueryContext(ctx, `
		SELECT id, name, description, avatar_url
		FROM profiles
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := make([]model.ProfileSummary, 0)
	for rows.Next() {
		var profile model.ProfileSummary
		if err := rows.Scan(&profile.ID, &profile.Name, &profile.Description, &profile.AvatarURL); err != nil {
			return nil, err
		}
		years, err := r.getAvailableYears(ctx, profile.ID.String())
		if err != nil {
			return nil, err
		}
		profile.AvailableYears = years
		profiles = append(profiles, profile)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (r *Repository) GetProfileByID(ctx context.Context, id string) (*model.Profile, error) {
	var profile model.Profile
	err := r.DB.DB.QueryRowContext(ctx, `
		SELECT id, name, description, avatar_url, scenario
		FROM profiles
		WHERE id = $1
	`, id).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Description,
		&profile.AvatarURL,
		&profile.Scenario,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProfileNotFound
	}
	if err != nil {
		return nil, err
	}

	years, err := r.getAvailableYears(ctx, id)
	if err != nil {
		return nil, err
	}
	profile.AvailableYears = years
	return &profile, nil
}

func (r *Repository) getAvailableYears(ctx context.Context, profileID string) ([]int, error) {
	rows, err := r.DB.DB.QueryContext(ctx, `
		SELECT year
		FROM profile_available_years
		WHERE profile_id = $1
		ORDER BY year
	`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	years := make([]int, 0)
	for rows.Next() {
		var year int
		if err := rows.Scan(&year); err != nil {
			return nil, err
		}
		years = append(years, year)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return years, nil
}
