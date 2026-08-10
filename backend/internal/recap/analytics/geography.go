package analytics

import (
	"sort"

	recap "recap-personalization/internal/recap"
)

type geoAggregate struct {
	Code   string
	Weight int
	Deals  int
}

var geoWeights = map[recap.EventType]int{
	recap.EventListingViewed: 1, recap.EventFavoriteAdded: 2, recap.EventSearchSaved: 2,
	recap.EventChatStarted: 3, recap.EventListingPublished: 3, recap.EventDeliveryUsed: 3,
	recap.EventSaleCompleted: 5, recap.EventPurchaseCompleted: 5,
}

func CalculateGeography(events []recap.NormalizedEvent, metrics *recap.Metrics) recap.Geography {
	verticals := map[string]*geoAggregate{}
	totalVerticalWeight := 0
	for _, event := range events {
		if event.ResolvedVerticalCode == "" {
			continue
		}
		weight := geoWeights[event.EventType]
		totalVerticalWeight += weight
		agg := verticals[event.ResolvedVerticalCode]
		if agg == nil {
			agg = &geoAggregate{Code: event.ResolvedVerticalCode}
			verticals[event.ResolvedVerticalCode] = agg
		}
		agg.Weight += weight
		if event.EventType == recap.EventSaleCompleted || event.EventType == recap.EventPurchaseCompleted {
			agg.Deals++
		}
	}
	mainVertical := bestAggregate(verticals)
	if mainVertical == nil {
		return recap.Geography{}
	}
	if totalVerticalWeight > 0 {
		metrics.TopVerticalShare = float64(mainVertical.Weight) / float64(totalVerticalWeight)
	}

	categories := map[string]*geoAggregate{}
	totalCategoryWeight := 0
	for _, event := range events {
		if event.ResolvedCategoryCode == "" {
			continue
		}
		weight := geoWeights[event.EventType]
		totalCategoryWeight += weight
		if event.ResolvedVerticalCode != mainVertical.Code {
			continue
		}
		agg := categories[event.ResolvedCategoryCode]
		if agg == nil {
			agg = &geoAggregate{Code: event.ResolvedCategoryCode}
			categories[event.ResolvedCategoryCode] = agg
		}
		agg.Weight += weight
		if event.EventType == recap.EventSaleCompleted || event.EventType == recap.EventPurchaseCompleted {
			agg.Deals++
		}
	}
	topCategory := bestAggregate(categories)
	if topCategory != nil && totalCategoryWeight > 0 {
		metrics.TopCategoryShare = float64(topCategory.Weight) / float64(totalCategoryWeight)
	}
	result := recap.Geography{MainVerticalCode: mainVertical.Code}
	if topCategory != nil {
		result.TopCategoryCode = topCategory.Code
	}
	return result
}

func bestAggregate(values map[string]*geoAggregate) *geoAggregate {
	ordered := make([]*geoAggregate, 0, len(values))
	for _, value := range values {
		ordered = append(ordered, value)
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Weight != ordered[j].Weight {
			return ordered[i].Weight > ordered[j].Weight
		}
		if ordered[i].Deals != ordered[j].Deals {
			return ordered[i].Deals > ordered[j].Deals
		}
		return ordered[i].Code < ordered[j].Code
	})
	if len(ordered) == 0 {
		return nil
	}
	return ordered[0]
}
