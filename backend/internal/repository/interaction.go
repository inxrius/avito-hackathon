package repository

import (
	"encoding/json"
)

// SaveInteraction — сохраняет событие взаимодействия
func (r *Repository) SaveInteraction(recapID, eventType string, metadata map[string]interface{}) error {
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	query := `INSERT INTO interactions (recap_id, event_type, metadata) VALUES ($1, $2, $3)`
	_, err = r.DB.DB.Exec(query, recapID, eventType, metadataJSON)
	return err
}
