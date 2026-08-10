package analytics

import (
	"sort"
	"time"

	recap "recap-personalization/internal/recap"
)

func CalculateMetrics(events []recap.NormalizedEvent) recap.Metrics {
	var m recap.Metrics
	m.MeaningfulEvents = len(events)
	days := map[string]time.Time{}
	months := map[string]struct{}{}
	categories := map[string]struct{}{}
	verticals := map[string]struct{}{}

	for _, event := range events {
		day := event.OccurredAt.UTC().Truncate(24 * time.Hour)
		days[day.Format("2006-01-02")] = day
		months[day.Format("2006-01")] = struct{}{}
		if event.ResolvedCategoryCode != "" {
			categories[event.ResolvedCategoryCode] = struct{}{}
		}
		if event.ResolvedVerticalCode != "" {
			verticals[event.ResolvedVerticalCode] = struct{}{}
		}
		switch event.EventType {
		case recap.EventListingViewed:
			m.ViewsCount++
		case recap.EventFavoriteAdded:
			m.FavoritesCount++
		case recap.EventSearchSaved:
			m.SavedSearchesCount++
		case recap.EventChatStarted:
			m.ChatsStartedCount++
		case recap.EventListingPublished:
			m.PublishedListingsCount++
		case recap.EventSaleCompleted:
			m.SalesCount++
		case recap.EventPurchaseCompleted:
			m.PurchasesCount++
		case recap.EventDeliveryUsed:
			m.DeliveryCount++
		}
	}
	m.ActiveDays = len(days)
	m.ActiveMonths = len(months)
	m.UniqueCategories = len(categories)
	m.UniqueVerticals = len(verticals)
	m.MaxActivityStreak = calculateStreak(days)
	m.BuyerActionsCount = m.ViewsCount + m.FavoritesCount + m.SavedSearchesCount + m.ChatsStartedCount + m.PurchasesCount
	m.SellerActionsCount = m.PublishedListingsCount + m.SalesCount
	m.CompletedDealsCount = m.SalesCount + m.PurchasesCount
	m.FavoriteToPurchaseRate = boundedRatio(m.PurchasesCount, m.FavoritesCount)
	m.PublishToSaleRate = boundedRatio(m.SalesCount, m.PublishedListingsCount)
	m.DeliveryUsageRate = boundedRatio(m.DeliveryCount, m.CompletedDealsCount)
	return m
}

func calculateStreak(days map[string]time.Time) int {
	ordered := make([]time.Time, 0, len(days))
	for _, day := range days {
		ordered = append(ordered, day)
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Before(ordered[j]) })
	best, current := 0, 0
	var previous time.Time
	for i, day := range ordered {
		if i == 0 || day.Sub(previous) == 24*time.Hour {
			current++
		} else {
			current = 1
		}
		if current > best {
			best = current
		}
		previous = day
	}
	return best
}

func boundedRatio(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	ratio := float64(numerator) / float64(denominator)
	if ratio > 1 {
		return 1
	}
	if ratio < 0 {
		return 0
	}
	return ratio
}

func PublicMetrics(m recap.Metrics, registry recap.Registry) []recap.MetricValue {
	values := []struct {
		code  recap.MetricCode
		value float64
	}{
		{recap.MetricMeaningfulEvents, float64(m.MeaningfulEvents)},
		{recap.MetricActiveDays, float64(m.ActiveDays)},
		{recap.MetricActiveMonths, float64(m.ActiveMonths)},
		{recap.MetricMaxActivityStreak, float64(m.MaxActivityStreak)},
		{recap.MetricViewsCount, float64(m.ViewsCount)},
		{recap.MetricFavoritesCount, float64(m.FavoritesCount)},
		{recap.MetricSavedSearchesCount, float64(m.SavedSearchesCount)},
		{recap.MetricChatsStartedCount, float64(m.ChatsStartedCount)},
		{recap.MetricPublishedListingsCount, float64(m.PublishedListingsCount)},
		{recap.MetricSalesCount, float64(m.SalesCount)},
		{recap.MetricPurchasesCount, float64(m.PurchasesCount)},
		{recap.MetricDeliveryCount, float64(m.DeliveryCount)},
		{recap.MetricUniqueCategories, float64(m.UniqueCategories)},
		{recap.MetricUniqueVerticals, float64(m.UniqueVerticals)},
		{recap.MetricTopVerticalShare, m.TopVerticalShare},
		{recap.MetricTopCategoryShare, m.TopCategoryShare},
		{recap.MetricBuyerActionsCount, float64(m.BuyerActionsCount)},
		{recap.MetricSellerActionsCount, float64(m.SellerActionsCount)},
		{recap.MetricCompletedDealsCount, float64(m.CompletedDealsCount)},
		{recap.MetricFavoriteToPurchaseRate, m.FavoriteToPurchaseRate},
		{recap.MetricPublishToSaleRate, m.PublishToSaleRate},
		{recap.MetricDeliveryUsageRate, m.DeliveryUsageRate},
	}
	result := make([]recap.MetricValue, 0, len(values))
	for _, item := range values {
		definition := registry.Metrics[item.code]
		result = append(result, recap.MetricValue{MetricCode: item.code, Value: item.value, Unit: definition.Unit})
	}
	return result
}
