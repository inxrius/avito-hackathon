package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	recap "recap-personalization/internal/recap"
	"recap-personalization/internal/recap/pipeline"
)

func main() {
	generator := pipeline.NewGenerator(nil)
	generator.Registry.PublicAvatarHosts["cdn.example"] = struct{}{}
	const profileID = "19369d67-1db3-4e87-b8a8-9c82991ba173"
	events := make([]recap.ActivityEvent, 0, 12)
	for day := 1; day <= 12; day++ {
		events = append(events, recap.ActivityEvent{
			EventID: fmt.Sprintf("event-%02d", day), ProfileID: profileID,
			EventType: recap.EventListingViewed, CategoryCode: "electronics",
			OccurredAt: time.Date(2026, 1, day, 12, 0, 0, 0, time.UTC),
		})
	}
	output, err := generator.Generate(context.Background(), recap.GenerateInput{
		RecapID: "550e8400-e29b-41d4-a716-446655440000",
		Profile: recap.ProfileSnapshot{
			ID: profileID, Name: "Алексей", AvatarURL: "https://cdn.example/avatar.png",
		},
		Year: 2026, Activities: events, GeneratedAt: time.Date(2026, 8, 6, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		log.Fatal(err)
	}
	encoded, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(encoded))
}
