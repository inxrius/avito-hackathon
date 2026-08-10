package recap

type MetricDefinition struct {
	Code MetricCode
	Unit MetricUnit
}

type MetricCardTemplate struct {
	TitleOne    string
	TitleFew    string
	TitleMany   string
	Description string
}

type AchievementDefinition struct {
	Code        AchievementCode
	Title       string
	Description string
	MetricCode  MetricCode
	Thresholds  [4]int
	Group       string
	Priority    int
}

type Registry struct {
	Verticals           map[VerticalCode]Vertical
	Categories          map[CategoryCode]Category
	Roles               map[ArchetypeRoleCode]ArchetypeRole
	Styles              map[ArchetypeStyleCode]ArchetypeStyle
	Metrics             map[MetricCode]MetricDefinition
	MetricCardTemplates map[MetricCode]MetricCardTemplate
	Achievements        []AchievementDefinition
	ThemeAccents        map[VerticalCode]AccentToken
	PublicAvatarHosts   map[string]struct{}
}

func DefaultRegistry() Registry {
	return Registry{
		Verticals: map[VerticalCode]Vertical{
			VerticalGoods:      {Code: VerticalGoods, Title: "Товары"},
			VerticalTransport:  {Code: VerticalTransport, Title: "Транспорт"},
			VerticalRealEstate: {Code: VerticalRealEstate, Title: "Недвижимость"},
			VerticalJobs:       {Code: VerticalJobs, Title: "Работа"},
			VerticalServices:   {Code: VerticalServices, Title: "Услуги"},
		},
		Categories: map[CategoryCode]Category{
			CategoryElectronics:            {Code: CategoryElectronics, Title: "Электроника", VerticalCode: VerticalGoods},
			CategoryHomeAndGarden:          {Code: CategoryHomeAndGarden, Title: "Для дома и дачи", VerticalCode: VerticalGoods},
			CategoryClothingAndAccessories: {Code: CategoryClothingAndAccessories, Title: "Одежда и аксессуары", VerticalCode: VerticalGoods},
			CategoryHobbiesAndLeisure:      {Code: CategoryHobbiesAndLeisure, Title: "Хобби и отдых", VerticalCode: VerticalGoods},
			CategoryCars:                   {Code: CategoryCars, Title: "Автомобили", VerticalCode: VerticalTransport},
			CategoryApartments:             {Code: CategoryApartments, Title: "Квартиры", VerticalCode: VerticalRealEstate},
			CategoryVacancies:              {Code: CategoryVacancies, Title: "Вакансии", VerticalCode: VerticalJobs},
			CategoryPersonalServices:       {Code: CategoryPersonalServices, Title: "Услуги", VerticalCode: VerticalServices},
		},
		Roles: map[ArchetypeRoleCode]ArchetypeRole{
			RoleCityObserver:     {Code: RoleCityObserver, Title: "Городской наблюдатель"},
			RoleUniversalCitizen: {Code: RoleUniversalCitizen, Title: "Универсальный горожанин"},
			RoleShowcaseOwner:    {Code: RoleShowcaseOwner, Title: "Хозяин витрины"},
			RoleFindingsSeeker:   {Code: RoleFindingsSeeker, Title: "Искатель находок"},
		},
		Styles: map[ArchetypeStyleCode]ArchetypeStyle{
			StyleDistrictExpert: {Code: StyleDistrictExpert, Title: "Знаток района"},
			StyleExplorer:       {Code: StyleExplorer, Title: "Исследователь"},
			StyleResultOriented: {Code: StyleResultOriented, Title: "Результативный"},
			StyleThoughtful:     {Code: StyleThoughtful, Title: "Вдумчивый"},
			StyleCityLocal:      {Code: StyleCityLocal, Title: "Свой в городе"},
			StyleRegular:        {Code: StyleRegular, Title: "Завсегдатай"},
		},
		Metrics: defaultMetricDefinitions(),
		MetricCardTemplates: map[MetricCode]MetricCardTemplate{
			MetricFavoritesCount: {
				TitleOne: "%d находка в избранном", TitleFew: "%d находки в избранном", TitleMany: "%d находок в избранном",
				Description: "Ты сохранял объявления, к которым хотелось вернуться",
			},
			MetricViewsCount: {
				TitleOne: "%d объявление изучено", TitleFew: "%d объявления изучено", TitleMany: "%d объявлений изучено",
				Description: "Ты внимательно исследовал предложения города",
			},
			MetricSalesCount: {
				TitleOne: "%d успешная продажа", TitleFew: "%d успешные продажи", TitleMany: "%d успешных продаж",
				Description: "Твоя витрина приносила завершённые сделки",
			},
			MetricPublishedListingsCount: {
				TitleOne: "%d объявление в твоей витрине", TitleFew: "%d объявления в твоей витрине", TitleMany: "%d объявлений в твоей витрине",
				Description: "Ты регулярно пополнял город своими предложениями",
			},
			MetricCompletedDealsCount: {
				TitleOne: "%d завершённая сделка", TitleFew: "%d завершённые сделки", TitleMany: "%d завершённых сделок",
				Description: "Ты был активен в покупательских и продавцовских сценариях",
			},
			MetricActiveDays: {
				TitleOne: "%d активный день", TitleFew: "%d активных дня", TitleMany: "%d активных дней",
				Description: "Твоя активность складывалась из разных городских сценариев",
			},
		},
		Achievements: []AchievementDefinition{
			{Code: AchievementDealMaster, Title: "Мастер сделок", Description: "Ты успешно завершал продажи и уверенно управлял своей витриной.", MetricCode: MetricSalesCount, Thresholds: [4]int{1, 3, 7, 15}, Group: "transactions", Priority: 90},
			{Code: AchievementFindingsCollector, Title: "Коллекционер находок", Description: "Ты собрал коллекцию объявлений, к которым хотелось возвращаться.", MetricCode: MetricFavoritesCount, Thresholds: [4]int{5, 15, 40, 80}, Group: "catalog", Priority: 70},
			{Code: AchievementCityNavigator, Title: "Городской навигатор", Description: "Ты исследовал разные улицы и категории города Авито.", MetricCode: MetricUniqueCategories, Thresholds: [4]int{3, 5, 8, 12}, Group: "catalog", Priority: 80},
			{Code: AchievementFrequentGuest, Title: "Частый гость", Description: "Ты регулярно возвращался в город в течение года.", MetricCode: MetricActiveDays, Thresholds: [4]int{7, 30, 90, 180}, Group: "loyalty", Priority: 50},
			{Code: AchievementOldTimer, Title: "Старожил", Description: "Твоя активность охватила значительную часть года.", MetricCode: MetricActiveMonths, Thresholds: [4]int{3, 6, 9, 12}, Group: "loyalty", Priority: 55},
			{Code: AchievementDoorstepDelivery, Title: "Доставка до двери", Description: "Ты использовал доставку, чтобы сделки проходили удобнее.", MetricCode: MetricDeliveryCount, Thresholds: [4]int{1, 3, 7, 15}, Group: "logistics", Priority: 65},
			{Code: AchievementOwnShowcase, Title: "Своя витрина", Description: "Ты активно пополнял собственную витрину объявлениями.", MetricCode: MetricPublishedListingsCount, Thresholds: [4]int{2, 5, 10, 25}, Group: "seller_presence", Priority: 75},
			{Code: AchievementFindingsHunter, Title: "Охотник за находками", Description: "Ты находил подходящие предложения и завершал покупки.", MetricCode: MetricPurchasesCount, Thresholds: [4]int{1, 3, 7, 15}, Group: "transactions", Priority: 85},
			{Code: AchievementCityRhythm, Title: "В ритме города", Description: "Ты сохранял серию активности несколько дней подряд.", MetricCode: MetricMaxActivityStreak, Thresholds: [4]int{3, 7, 14, 30}, Group: "loyalty", Priority: 60},
		},
		ThemeAccents: map[VerticalCode]AccentToken{
			VerticalGoods: AccentViolet, VerticalTransport: AccentBlue, VerticalRealEstate: AccentGreen,
			VerticalJobs: AccentOrange, VerticalServices: AccentViolet,
		},
		PublicAvatarHosts: map[string]struct{}{},
	}
}

func defaultMetricDefinitions() map[MetricCode]MetricDefinition {
	units := map[MetricCode]MetricUnit{
		MetricMeaningfulEvents: UnitEvents, MetricActiveDays: UnitDays, MetricActiveMonths: UnitMonths,
		MetricMaxActivityStreak: UnitDays, MetricViewsCount: UnitItems, MetricFavoritesCount: UnitItems,
		MetricSavedSearchesCount: UnitItems, MetricChatsStartedCount: UnitItems, MetricPublishedListingsCount: UnitItems,
		MetricSalesCount: UnitItems, MetricPurchasesCount: UnitItems, MetricDeliveryCount: UnitItems,
		MetricUniqueCategories: UnitCategories, MetricUniqueVerticals: UnitVerticals, MetricTopVerticalShare: UnitRatio,
		MetricTopCategoryShare: UnitRatio, MetricBuyerActionsCount: UnitActions, MetricSellerActionsCount: UnitActions,
		MetricCompletedDealsCount: UnitItems, MetricFavoriteToPurchaseRate: UnitRatio, MetricPublishToSaleRate: UnitRatio,
		MetricDeliveryUsageRate: UnitRatio,
	}
	result := make(map[MetricCode]MetricDefinition, len(units))
	for code, unit := range units {
		result[code] = MetricDefinition{Code: code, Unit: unit}
	}
	return result
}
