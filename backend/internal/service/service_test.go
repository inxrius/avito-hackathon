package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/inxrius/avito-hackathon/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCountDistricts(t *testing.T) {
	tests := []struct {
		name       string
		activities []model.Activity
		want       int
	}{
		{
			name: "empty activities",
			activities: []model.Activity{},
			want: 0,
		},
		{
			name: "single category",
			activities: []model.Activity{
				{Category: "Электроника"},
				{Category: "Электроника"},
			},
			want: 1,
		},
		{
			name: "multiple categories",
			activities: []model.Activity{
				{Category: "Электроника"},
				{Category: "Авто"},
				{Category: "Недвижимость"},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countDistricts(tt.activities)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCountFavorites(t *testing.T) {
	tests := []struct {
		name       string
		activities []model.Activity
		want       int
	}{
		{
			name: "no favorites",
			activities: []model.Activity{
				{Type: "view"},
				{Type: "purchase"},
			},
			want: 0,
		},
		{
			name: "with favorites",
			activities: []model.Activity{
				{Type: "view"},
				{Type: "favorite"},
				{Type: "favorite"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countFavorites(tt.activities)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCountPurchases(t *testing.T) {
	tests := []struct {
		name       string
		activities []model.Activity
		want       int
	}{
		{
			name: "no purchases",
			activities: []model.Activity{
				{Type: "view"},
				{Type: "favorite"},
			},
			want: 0,
		},
		{
			name: "with purchases",
			activities: []model.Activity{
				{Type: "purchase"},
				{Type: "sale"},
				{Type: "purchase"},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countPurchases(tt.activities)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResolveArchetype(t *testing.T) {
	tests := []struct {
		name      string
		scenario  string
		activities []model.Activity
		wantName  string
	}{
		{
			name:     "seller scenario",
			scenario: "seller",
			activities: []model.Activity{},
			wantName: "Мастер обмена",
		},
		{
			name:     "unfinished scenario",
			scenario: "unfinished",
			activities: []model.Activity{},
			wantName: "Исследователь возможностей",
		},
		{
			name:     "insufficient scenario",
			scenario: "insufficient",
			activities: []model.Activity{},
			wantName: "Новый житель",
		},
		{
			name:     "default buyer scenario",
			scenario: "buyer",
			activities: []model.Activity{},
			wantName: "Охотник за домом",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveArchetype(tt.scenario, tt.activities)
			assert.Equal(t, tt.wantName, got.Name)
		})
	}
}

func TestBuildAchievements(t *testing.T) {
	tests := []struct {
		name       string
		activities []model.Activity
		wantCount  int
	}{
		{
			name:       "no activities",
			activities: []model.Activity{},
			wantCount:  0,
		},
		{
			name: "first view only",
			activities: []model.Activity{
				{Type: "view", Category: "Электроника"},
			},
			wantCount: 1,
		},
		{
			name: "collector achievement",
			activities: func() []model.Activity {
				acts := []model.Activity{}
				for i := 0; i < 5; i++ {
					acts = append(acts, model.Activity{Type: "favorite"})
				}
				return acts
			}(),
			wantCount: 2,
		},
		{
			name: "active buyer achievement",
			activities: func() []model.Activity {
				acts := []model.Activity{}
				for i := 0; i < 3; i++ {
					acts = append(acts, model.Activity{Type: "purchase"})
				}
				return acts
			}(),
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildAchievements(tt.activities)
			assert.Len(t, got, tt.wantCount)
		})
	}
}

func TestBuildOpportunities(t *testing.T) {
	tests := []struct {
		name       string
		activities []model.Activity
		wantCount  int
	}{
		{
			name: "no unfinished",
			activities: []model.Activity{
				{Type: "purchase", Category: "Электроника"},
			},
			wantCount: 1,
		},
		{
			name: "with unfinished",
			activities: []model.Activity{
				{Type: "favorite", Category: "Электроника"},
				{Type: "view", Category: "Авто"},
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildOpportunities(tt.activities)
			assert.Len(t, got, tt.wantCount)
		})
	}
}

func TestComputeActivityHash(t *testing.T) {
	activities1 := []model.Activity{
		{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), Type: "view", Category: "test", Value: 1.0},
	}
	activities2 := []model.Activity{
		{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), Type: "view", Category: "test", Value: 1.0},
	}
	activities3 := []model.Activity{
		{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), Type: "view", Category: "test", Value: 1.0},
	}

	hash1 := computeActivityHash(activities1)
	hash2 := computeActivityHash(activities2)
	hash3 := computeActivityHash(activities3)

	assert.Equal(t, hash1, hash2, "same activities should produce same hash")
	assert.NotEqual(t, hash1, hash3, "different activities should produce different hash")
	assert.NotEmpty(t, hash1, "hash should not be empty")
	assert.Len(t, hash1, 64, "SHA256 hash should be 64 characters")
}