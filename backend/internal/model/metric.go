package model

// MetricCode — код метрики (соответствует enum в OpenAPI)
type MetricCode string

const (
	MetricMeaningfulEvents       MetricCode = "meaningful_events"
	MetricActiveDays             MetricCode = "active_days"
	MetricActiveMonths           MetricCode = "active_months"
	MetricMaxActivityStreak      MetricCode = "max_activity_streak"
	MetricViewsCount             MetricCode = "views_count"
	MetricFavoritesCount         MetricCode = "favorites_count"
	MetricSavedSearchesCount     MetricCode = "saved_searches_count"
	MetricChatsStartedCount      MetricCode = "chats_started_count"
	MetricPublishedListingsCount MetricCode = "published_listings_count"
	MetricSalesCount             MetricCode = "sales_count"
	MetricPurchasesCount         MetricCode = "purchases_count"
	MetricDeliveryCount          MetricCode = "delivery_count"
	MetricUniqueCategories       MetricCode = "unique_categories"
	MetricUniqueVerticals        MetricCode = "unique_verticals"
	MetricTopVerticalShare       MetricCode = "top_vertical_share"
	MetricTopCategoryShare       MetricCode = "top_category_share"
	MetricBuyerActionsCount      MetricCode = "buyer_actions_count"
	MetricSellerActionsCount     MetricCode = "seller_actions_count"
	MetricCompletedDealsCount    MetricCode = "completed_deals_count"
	MetricFavoriteToPurchaseRate MetricCode = "favorite_to_purchase_rate"
	MetricPublishToSaleRate      MetricCode = "publish_to_sale_rate"
	MetricDeliveryUsageRate      MetricCode = "delivery_usage_rate"
)

// MetricUnit — единица измерения метрики
type MetricUnit string

const (
	UnitEvents     MetricUnit = "events"
	UnitDays       MetricUnit = "days"
	UnitMonths     MetricUnit = "months"
	UnitItems      MetricUnit = "items"
	UnitCategories MetricUnit = "categories"
	UnitVerticals  MetricUnit = "verticals"
	UnitActions    MetricUnit = "actions"
	UnitRatio      MetricUnit = "ratio"
)

// MetricDefinition — справочная запись (таблица metric_definitions)
type MetricDefinition struct {
	Code        MetricCode  `json:"code"`
	Title       string      `json:"title"`
	DefaultUnit *MetricUnit `json:"default_unit,omitempty"`
}

// MetricValue — структура для карточки типа "metric"
// Используется в RecapCard.Data
type MetricValue struct {
	MetricCode     MetricCode  `json:"metric_code"`
	Value          float64     `json:"value"`
	Unit           *MetricUnit `json:"unit,omitempty"`
	SecondaryLabel *string     `json:"secondary_label,omitempty"`
}
