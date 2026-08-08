package recap

import "time"

const (
	SchemaVersion          = "2.0"
	OpenAPIVersion         = "2.0.1"
	AlgorithmVersion       = "recap-rules-2026.08.2-spike"
	FeatureSchemaVersion   = "features-v1"
	PromptVersion          = "city-summary-v2-spike"
	ArchetypeRuleVersion   = "archetype-rules-v2-spike"
	AchievementRuleVersion = "achievement-rules-v1-spike"
)

type EventType string

const (
	EventListingViewed     EventType = "listing_viewed"
	EventFavoriteAdded     EventType = "favorite_added"
	EventSearchSaved       EventType = "search_saved"
	EventChatStarted       EventType = "chat_started"
	EventListingPublished  EventType = "listing_published"
	EventSaleCompleted     EventType = "sale_completed"
	EventPurchaseCompleted EventType = "purchase_completed"
	EventDeliveryUsed      EventType = "delivery_used"
)

type VerticalCode string

const (
	VerticalGoods      VerticalCode = "goods"
	VerticalTransport  VerticalCode = "transport"
	VerticalRealEstate VerticalCode = "real_estate"
	VerticalJobs       VerticalCode = "jobs"
	VerticalServices   VerticalCode = "services"
)

type CategoryCode string

const (
	CategoryElectronics            CategoryCode = "electronics"
	CategoryHomeAndGarden          CategoryCode = "home_and_garden"
	CategoryClothingAndAccessories CategoryCode = "clothing_and_accessories"
	CategoryHobbiesAndLeisure      CategoryCode = "hobbies_and_leisure"
	CategoryCars                   CategoryCode = "cars"
	CategoryApartments             CategoryCode = "apartments"
	CategoryVacancies              CategoryCode = "vacancies"
	CategoryPersonalServices       CategoryCode = "personal_services"
)

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

type ArchetypeRoleCode string

const (
	RoleFindingsSeeker   ArchetypeRoleCode = "findings_seeker"
	RoleShowcaseOwner    ArchetypeRoleCode = "showcase_owner"
	RoleUniversalCitizen ArchetypeRoleCode = "universal_citizen"
	RoleCityObserver     ArchetypeRoleCode = "city_observer"
)

type ArchetypeStyleCode string

const (
	StyleThoughtful     ArchetypeStyleCode = "thoughtful"
	StyleExplorer       ArchetypeStyleCode = "explorer"
	StyleDistrictExpert ArchetypeStyleCode = "district_expert"
	StyleRegular        ArchetypeStyleCode = "regular"
	StyleResultOriented ArchetypeStyleCode = "result_oriented"
	StyleCityLocal      ArchetypeStyleCode = "city_local"
)

type AchievementCode string

const (
	AchievementDealMaster        AchievementCode = "deal_master"
	AchievementFindingsCollector AchievementCode = "findings_collector"
	AchievementCityNavigator     AchievementCode = "city_navigator"
	AchievementFrequentGuest     AchievementCode = "frequent_guest"
	AchievementOldTimer          AchievementCode = "old_timer"
	AchievementDoorstepDelivery  AchievementCode = "doorstep_delivery"
	AchievementOwnShowcase       AchievementCode = "own_showcase"
	AchievementFindingsHunter    AchievementCode = "findings_hunter"
	AchievementCityRhythm        AchievementCode = "city_rhythm"
)

type AchievementLevel string

const (
	LevelNewcomer AchievementLevel = "newcomer"
	LevelLocal    AchievementLevel = "local"
	LevelExpert   AchievementLevel = "expert"
	LevelGuru     AchievementLevel = "guru"
)

type CardType string

const (
	CardTypeIntro        CardType = "intro"
	CardTypeMetric       CardType = "metric"
	CardTypeDistrict     CardType = "district"
	CardTypeArchetype    CardType = "archetype"
	CardTypeAchievements CardType = "achievements"
	CardTypeSummary      CardType = "summary"
	CardTypeFinal        CardType = "final"
)

type CardVisibility string

const (
	VisibilityPersonal  CardVisibility = "personal"
	VisibilityShareable CardVisibility = "shareable"
)

type CardVisualKind string

const (
	VisualIllustration CardVisualKind = "illustration"
	VisualDistrict     CardVisualKind = "district"
	VisualStreet       CardVisualKind = "street"
	VisualCalendar     CardVisualKind = "calendar"
	VisualBadge        CardVisualKind = "badge"
	VisualChart        CardVisualKind = "chart"
	VisualCharacter    CardVisualKind = "character"
	VisualSkyline      CardVisualKind = "skyline"
)

type AccentToken string

const (
	AccentViolet AccentToken = "violet"
	AccentBlue   AccentToken = "blue"
	AccentGreen  AccentToken = "green"
	AccentOrange AccentToken = "orange"
)

type ExplanationKind string

const (
	ExplanationArchetypeRole  ExplanationKind = "archetype_role"
	ExplanationArchetypeStyle ExplanationKind = "archetype_style"
	ExplanationAchievement    ExplanationKind = "achievement"
)

type RuleOperator string

const (
	OperatorEQ  RuleOperator = "eq"
	OperatorNEQ RuleOperator = "neq"
	OperatorGT  RuleOperator = "gt"
	OperatorGTE RuleOperator = "gte"
	OperatorLT  RuleOperator = "lt"
	OperatorLTE RuleOperator = "lte"
)

type ShareFactKind string

const (
	ShareFactMainDistrict   ShareFactKind = "main_district"
	ShareFactActiveDays     ShareFactKind = "active_days"
	ShareFactTopAchievement ShareFactKind = "top_achievement"
)

type ActivityEvent struct {
	EventID      string
	ProfileID    string
	EventType    EventType
	VerticalCode string
	CategoryCode string
	OccurredAt   time.Time
}

type NormalizedEvent struct {
	EventID              string
	EventType            EventType
	ResolvedVerticalCode string
	ResolvedCategoryCode string
	OccurredAt           time.Time
}

type ProfileSnapshot struct {
	ID        string
	Name      string
	AvatarURL string
}

type GenerateInput struct {
	RecapID     string
	Profile     ProfileSnapshot
	Year        int
	Activities  []ActivityEvent
	GeneratedAt time.Time
}

type Metrics struct {
	MeaningfulEvents       int
	ActiveDays             int
	ActiveMonths           int
	MaxActivityStreak      int
	ViewsCount             int
	FavoritesCount         int
	SavedSearchesCount     int
	ChatsStartedCount      int
	PublishedListingsCount int
	SalesCount             int
	PurchasesCount         int
	DeliveryCount          int
	UniqueCategories       int
	UniqueVerticals        int
	TopVerticalShare       float64
	TopCategoryShare       float64
	BuyerActionsCount      int
	SellerActionsCount     int
	CompletedDealsCount    int
	FavoriteToPurchaseRate float64
	PublishToSaleRate      float64
	DeliveryUsageRate      float64
}

type Geography struct {
	MainVerticalCode string
	TopCategoryCode  string
}

type Vertical struct {
	Code  VerticalCode `json:"code"`
	Title string       `json:"title"`
}

type Category struct {
	Code         CategoryCode `json:"code"`
	Title        string       `json:"title"`
	VerticalCode VerticalCode `json:"vertical_code"`
}

type ArchetypeRole struct {
	Code  ArchetypeRoleCode `json:"code"`
	Title string            `json:"title"`
}

type ArchetypeStyle struct {
	Code  ArchetypeStyleCode `json:"code"`
	Title string             `json:"title"`
}

type ArchetypeDecision struct {
	Role  ArchetypeRole  `json:"role"`
	Style ArchetypeStyle `json:"style"`
}

type AchievementDecision struct {
	Code               AchievementCode  `json:"code"`
	Title              string           `json:"title"`
	Description        string           `json:"description"`
	MetricCode         MetricCode       `json:"metric_code"`
	Value              int              `json:"value"`
	Level              AchievementLevel `json:"level"`
	CurrentThreshold   int              `json:"current_level_threshold"`
	NextLevelThreshold *int             `json:"next_level_threshold"`
	Group              string           `json:"-"`
	Priority           int              `json:"-"`
	Progress           float64          `json:"-"`
	Icon               string           `json:"icon"`
}

type Achievement struct {
	Code               AchievementCode  `json:"code"`
	Title              string           `json:"title"`
	Description        string           `json:"description"`
	Level              AchievementLevel `json:"level"`
	Icon               string           `json:"icon"`
	MetricCode         MetricCode       `json:"metric_code,omitempty"`
	CurrentValue       *float64         `json:"current_value,omitempty"`
	NextLevelThreshold *float64         `json:"next_level_threshold"`
}

type Narrative struct {
	SummaryTitle  string  `json:"summary_title"`
	SummaryText   string  `json:"summary_text"`
	Source        string  `json:"source"`
	Model         *string `json:"model"`
	PromptVersion string  `json:"prompt_version"`
}

type NarrativeGeneration struct {
	Source        string  `json:"source"`
	PromptVersion string  `json:"prompt_version"`
	Model         *string `json:"model"`
}

type MetricValue struct {
	MetricCode     MetricCode `json:"metric_code"`
	Value          float64    `json:"value"`
	Unit           MetricUnit `json:"unit,omitempty"`
	SecondaryLabel *string    `json:"secondary_label,omitempty"`
}

type CardVisual struct {
	Kind      CardVisualKind `json:"kind"`
	AssetCode *string        `json:"asset_code,omitempty"`
}

type BaseRecapCard struct {
	ID          string         `json:"id"`
	Type        CardType       `json:"type"`
	Position    int            `json:"position"`
	Visibility  CardVisibility `json:"visibility"`
	Eyebrow     *string        `json:"eyebrow,omitempty"`
	Title       string         `json:"title"`
	Description *string        `json:"description,omitempty"`
	Visual      *CardVisual    `json:"visual,omitempty"`
	Explainable bool           `json:"explainable"`
}

type RecapCard interface {
	Base() BaseRecapCard
}

type IntroCardData struct {
	Year int `json:"year"`
}

type IntroCard struct {
	BaseRecapCard
	Data IntroCardData `json:"data"`
}

func (c IntroCard) Base() BaseRecapCard { return c.BaseRecapCard }

type MetricCard struct {
	BaseRecapCard
	Data MetricValue `json:"data"`
}

func (c MetricCard) Base() BaseRecapCard { return c.BaseRecapCard }

type DistrictCardData struct {
	Vertical      Vertical  `json:"vertical"`
	ActivityShare float64   `json:"activity_share"`
	TopCategory   *Category `json:"top_category,omitempty"`
}

type DistrictCard struct {
	BaseRecapCard
	Data DistrictCardData `json:"data"`
}

func (c DistrictCard) Base() BaseRecapCard { return c.BaseRecapCard }

type ArchetypeCardData struct {
	Role  ArchetypeRole  `json:"role"`
	Style ArchetypeStyle `json:"style"`
}

type ArchetypeCard struct {
	BaseRecapCard
	Data ArchetypeCardData `json:"data"`
}

func (c ArchetypeCard) Base() BaseRecapCard { return c.BaseRecapCard }

type AchievementsCardData struct {
	Items      []Achievement `json:"items"`
	TotalCount int           `json:"total_count"`
}

type AchievementsCard struct {
	BaseRecapCard
	Data AchievementsCardData `json:"data"`
}

func (c AchievementsCard) Base() BaseRecapCard { return c.BaseRecapCard }

type SummaryCardData struct {
	RoleCode         ArchetypeRoleCode  `json:"role_code"`
	StyleCode        ArchetypeStyleCode `json:"style_code"`
	AchievementCodes []AchievementCode  `json:"achievement_codes"`
}

type SummaryCard struct {
	BaseRecapCard
	Data SummaryCardData `json:"data"`
}

func (c SummaryCard) Base() BaseRecapCard { return c.BaseRecapCard }

type FinalCardData struct {
	ShowShareButton bool `json:"show_share_button"`
	ShowFeedback    bool `json:"show_feedback"`
}

type FinalCard struct {
	BaseRecapCard
	Data FinalCardData `json:"data"`
}

func (c FinalCard) Base() BaseRecapCard { return c.BaseRecapCard }

type Capabilities struct {
	ShareAvailable       bool `json:"share_available"`
	ExplanationAvailable bool `json:"explanation_available"`
	FeedbackAvailable    bool `json:"feedback_available"`
}

type RecapTheme struct {
	Code         string       `json:"code"`
	MainDistrict Vertical     `json:"main_district"`
	AccentToken  *AccentToken `json:"accent_token,omitempty"`
}

type RecapProfile struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type RecapGeneration struct {
	AlgorithmVersion     string              `json:"algorithm_version"`
	FeatureSchemaVersion string              `json:"feature_schema_version"`
	ActivityHash         string              `json:"activity_hash"`
	GeneratedAt          time.Time           `json:"generated_at"`
	Narrative            NarrativeGeneration `json:"narrative"`
}

type RuleFact struct {
	MetricCode MetricCode   `json:"metric_code"`
	Actual     float64      `json:"actual"`
	Operator   RuleOperator `json:"operator"`
	Threshold  float64      `json:"threshold"`
	Matched    bool         `json:"matched"`
}

type DecisionExplanation struct {
	CardID      string          `json:"card_id"`
	Kind        ExplanationKind `json:"kind"`
	Code        string          `json:"code"`
	Reason      string          `json:"reason"`
	RuleVersion string          `json:"rule_version"`
	Facts       []RuleFact      `json:"facts"`
}

type RecapExplanation struct {
	RecapID          string                `json:"recap_id"`
	AlgorithmVersion string                `json:"algorithm_version"`
	ActivityHash     string                `json:"activity_hash"`
	Decisions        []DecisionExplanation `json:"decisions"`
}

type ShareFact struct {
	Kind  ShareFactKind `json:"kind"`
	Label string        `json:"label"`
	Value string        `json:"value"`
}

type ShareAchievement struct {
	Code  AchievementCode  `json:"code"`
	Title string           `json:"title"`
	Level AchievementLevel `json:"level"`
	Icon  string           `json:"icon"`
}

type ShareVisual struct {
	Theme string `json:"theme"`
}

type ShareCard struct {
	SchemaVersion string             `json:"schema_version"`
	RecapID       string             `json:"recap_id"`
	ProfileName   string             `json:"profile_name"`
	AvatarURL     *string            `json:"avatar_url,omitempty"`
	Year          int                `json:"year"`
	Title         string             `json:"title"`
	Subtitle      string             `json:"subtitle"`
	MainDistrict  Vertical           `json:"main_district"`
	Facts         []ShareFact        `json:"facts"`
	Achievements  []ShareAchievement `json:"achievements"`
	Visual        ShareVisual        `json:"visual"`
}

type RecapSnapshot struct {
	SchemaVersion string          `json:"schema_version"`
	ID            string          `json:"id"`
	ProfileID     string          `json:"profile_id"`
	Year          int             `json:"year"`
	Profile       RecapProfile    `json:"profile"`
	Generation    RecapGeneration `json:"generation"`
	Theme         RecapTheme      `json:"theme"`
	Cards         []RecapCard     `json:"cards"`
	Capabilities  Capabilities    `json:"capabilities"`
}

type GenerateOutput struct {
	Recap        RecapSnapshot         `json:"recap"`
	Metrics      []MetricValue         `json:"metrics"`
	Archetype    ArchetypeDecision     `json:"archetype"`
	Achievements []AchievementDecision `json:"achievements"`
	Narrative    Narrative             `json:"narrative"`
	Explanation  *RecapExplanation     `json:"explanation,omitempty"`
	Share        *ShareCard            `json:"share,omitempty"`
}
