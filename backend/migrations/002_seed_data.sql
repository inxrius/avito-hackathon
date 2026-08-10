INSERT INTO verticals (code, title) VALUES
    ('goods', 'Товары'),
    ('transport', 'Транспорт'),
    ('real_estate', 'Недвижимость'),
    ('jobs', 'Работа'),
    ('services', 'Услуги');

INSERT INTO categories (code, vertical_code, title) VALUES
    ('electronics', 'goods', 'Электроника'),
    ('home_and_garden', 'goods', 'Для дома и дачи'),
    ('clothing_and_accessories', 'goods', 'Одежда и аксессуары'),
    ('hobbies_and_leisure', 'goods', 'Хобби и отдых'),
    ('cars', 'transport', 'Автомобили'),
    ('apartments', 'real_estate', 'Квартиры'),
    ('vacancies', 'jobs', 'Вакансии'),
    ('personal_services', 'services', 'Услуги');

INSERT INTO metric_definitions (code, title, default_unit) VALUES
    ('meaningful_events', 'Значимые действия', 'events'),
    ('active_days', 'Активные дни', 'days'),
    ('active_months', 'Активные месяцы', 'months'),
    ('max_activity_streak', 'Максимальная серия активности', 'days'),
    ('views_count', 'Просмотры объявлений', 'items'),
    ('favorites_count', 'Добавления в избранное', 'items'),
    ('saved_searches_count', 'Сохранённые поиски', 'items'),
    ('chats_started_count', 'Начатые чаты', 'items'),
    ('published_listings_count', 'Опубликованные объявления', 'items'),
    ('sales_count', 'Продажи', 'items'),
    ('purchases_count', 'Покупки', 'items'),
    ('delivery_count', 'Использования доставки', 'items'),
    ('unique_categories', 'Уникальные категории', 'categories'),
    ('unique_verticals', 'Уникальные вертикали', 'verticals'),
    ('top_vertical_share', 'Доля главной вертикали', 'ratio'),
    ('top_category_share', 'Доля главной категории', 'ratio'),
    ('buyer_actions_count', 'Покупательские действия', 'actions'),
    ('seller_actions_count', 'Продавцовские действия', 'actions'),
    ('completed_deals_count', 'Завершённые сделки', 'items'),
    ('favorite_to_purchase_rate', 'Отношение избранного к покупкам', 'ratio'),
    ('publish_to_sale_rate', 'Отношение публикаций к продажам', 'ratio'),
    ('delivery_usage_rate', 'Доля сделок с доставкой', 'ratio');

INSERT INTO archetype_roles (code, title) VALUES
    ('findings_seeker', 'Искатель находок'),
    ('showcase_owner', 'Хозяин витрины'),
    ('universal_citizen', 'Универсальный горожанин'),
    ('city_observer', 'Городской наблюдатель');

INSERT INTO archetype_styles (code, title) VALUES
    ('thoughtful', 'Вдумчивый'),
    ('explorer', 'Исследователь'),
    ('district_expert', 'Знаток района'),
    ('regular', 'Завсегдатай'),
    ('result_oriented', 'Результативный'),
    ('city_local', 'Свой в городе');

INSERT INTO achievement_definitions (code, title, description, icon_code) VALUES
    ('deal_master', 'Мастер сделок', 'Ты успешно завершал продажи и уверенно управлял своей витриной.', 'deal-master'),
    ('findings_collector', 'Коллекционер находок', 'Ты собрал коллекцию объявлений, к которым хотелось возвращаться.', 'findings-collector'),
    ('city_navigator', 'Городской навигатор', 'Ты исследовал разные улицы и категории города Авито.', 'city-navigator'),
    ('frequent_guest', 'Частый гость', 'Ты регулярно возвращался в город в течение года.', 'frequent-guest'),
    ('old_timer', 'Старожил', 'Твоя активность охватила значительную часть года.', 'old-timer'),
    ('doorstep_delivery', 'Доставка до двери', 'Ты использовал доставку, чтобы сделки проходили удобнее.', 'doorstep-delivery'),
    ('own_showcase', 'Своя витрина', 'Ты активно пополнял собственную витрину объявлениями.', 'own-showcase'),
    ('findings_hunter', 'Охотник за находками', 'Ты находил подходящие предложения и завершал покупки.', 'findings-hunter'),
    ('city_rhythm', 'В ритме города', 'Ты сохранял серию активности несколько дней подряд.', 'city-rhythm');

INSERT INTO profiles (id, name, description, avatar_url, scenario) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Марина', 'Изучает предложения, сохраняет находки и совершает покупки.', 'https://cdn.example/marina.png', 'buyer'),
    ('22222222-2222-2222-2222-222222222222', 'Алексей', 'Публикует объявления и завершает продажи.', 'https://cdn.example/alexey.png', 'seller'),
    ('33333333-3333-3333-3333-333333333333', 'Ольга', 'Активно исследует разные категории и районы.', 'https://cdn.example/olga.png', 'explorer'),
    ('44444444-4444-4444-4444-444444444444', 'Дмитрий', 'Только начинает пользоваться городом Авито.', 'https://cdn.example/dmitry.png', 'insufficient');

INSERT INTO profile_available_years (profile_id, year) VALUES
    ('11111111-1111-1111-1111-111111111111', 2026),
    ('22222222-2222-2222-2222-222222222222', 2026),
    ('33333333-3333-3333-3333-333333333333', 2026),
    ('44444444-4444-4444-4444-444444444444', 2026);
