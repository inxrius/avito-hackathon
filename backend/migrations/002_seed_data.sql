-- Seed-данные

-- 1. Заполнение справочников

INSERT INTO verticals (code, title) VALUES
    ('goods', 'Товары'),
    ('transport', 'Транспорт'),
    ('real_estate', 'Недвижимость'),
    ('jobs', 'Вакансии'),
    ('services', 'Услуги');

INSERT INTO categories (code, vertical_code, title) VALUES
    ('electronics', 'goods', 'Электроника'),
    ('home_and_garden', 'goods', 'Для дома'),
    ('clothing_and_accessories', 'goods', 'Одежда и аксессуары'),
    ('hobbies_and_leisure', 'goods', 'Хобби и досуг'),
    ('cars', 'transport', 'Авто'),
    ('apartments', 'real_estate', 'Квартиры'),
    ('vacancies', 'jobs', 'Вакансии'),
    ('personal_services', 'services', 'Личные услуги');

INSERT INTO metric_definitions (code, title, default_unit) VALUES
    ('meaningful_events', 'Значимые действия', 'events'),
    ('active_days', 'Дней активности', 'days'),
    ('active_months', 'Месяцев активности', 'months'),
    ('max_activity_streak', 'Самая длинная серия', 'days'),
    ('views_count', 'Количество просмотров', 'items'),
    ('favorites_count', 'Избранное', 'items'),
    ('saved_searches_count', 'Сохранённые поиски', 'items'),
    ('chats_started_count', 'Начатые чаты', 'items'),
    ('published_listings_count', 'Опубликованные объявления', 'items'),
    ('sales_count', 'Продажи', 'items'),
    ('purchases_count', 'Покупки', 'items'),
    ('delivery_count', 'Доставки', 'items'),
    ('unique_categories', 'Уникальные категории', 'categories'),
    ('unique_verticals', 'Уникальные вертикали', 'verticals'),
    ('top_vertical_share', 'Доля главной вертикали', 'ratio'),
    ('top_category_share', 'Доля главной категории', 'ratio'),
    ('buyer_actions_count', 'Действия покупателя', 'actions'),
    ('seller_actions_count', 'Действия продавца', 'actions'),
    ('completed_deals_count', 'Завершённые сделки', 'items'),
    ('favorite_to_purchase_rate', 'Конверсия из избранного в покупку', 'ratio'),
    ('publish_to_sale_rate', 'Конверсия из публикации в продажу', 'ratio'),
    ('delivery_usage_rate', 'Частота использования доставки', 'ratio');

INSERT INTO archetype_roles (code, title) VALUES
    ('findings_seeker', 'Искатель находок'),
    ('showcase_owner', 'Владелец витрины'),
    ('universal_citizen', 'Гражданин города'),
    ('city_observer', 'Наблюдатель');

INSERT INTO archetype_styles (code, title) VALUES
    ('thoughtful', 'Вдумчивый'),
    ('explorer', 'Исследователь'),
    ('district_expert', 'Эксперт района'),
    ('regular', 'Регулярный гость'),
    ('result_oriented', 'Результативный'),
    ('city_local', 'Местный житель');

INSERT INTO achievement_definitions (code, title, description, icon_code) VALUES
    ('deal_master', 'Мастер сделок', 'Заключил более 5 сделок за год', 'deal-master'),
    ('findings_collector', 'Коллекционер находок', 'Сохранил более 20 объявлений в избранное', 'findings-collector'),
    ('city_navigator', 'Городской навигатор', 'Посетил 3 и более вертикали', 'city-navigator'),
    ('frequent_guest', 'Частый гость', 'Был активен более 50 дней', 'frequent-guest'),
    ('old_timer', 'Старожил', 'Пользовался площадкой более 2 лет', 'old-timer'),
    ('doorstep_delivery', 'Доставка на порог', 'Сделал более 3 доставок', 'doorstep-delivery'),
    ('own_showcase', 'Своя витрина', 'Опубликовал более 3 объявлений', 'own-showcase'),
    ('findings_hunter', 'Охотник за находками', 'Сделал более 5 покупок', 'findings-hunter'),
    ('city_rhythm', 'Ритм города', 'Активен в разное время суток', 'city-rhythm');

INSERT INTO metric_definitions (code, title, default_unit) VALUES
    ('category_Недвижимость', 'Активность в Недвижимости', 'actions'),
    ('category_Для дома', 'Активность в Для дома', 'actions'),
    ('category_Электроника', 'Активность в Электронике', 'actions'),
    ('category_Авто', 'Активность в Авто', 'actions') ON CONFLICT (code) DO NOTHING;

-- 2. Тестовые профили

INSERT INTO profiles (id, name, description, avatar_url, scenario) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Марина', 'Смотрела долго, выбирала придирчиво, взяла своё — и сразу начала обустраивать.', 'https://i.pravatar.cc/150?img=1', 'buyer'),
    ('22222222-2222-2222-2222-222222222222', 'Алексей', 'Превращает лишнее в нужное и делает площадку живее.', 'https://i.pravatar.cc/150?img=2', 'seller'),
    ('33333333-3333-3333-3333-333333333333', 'Ольга', 'Собирает варианты и готовится к большому выбору.', 'https://i.pravatar.cc/150?img=3', 'unfinished'),
    ('44444444-4444-4444-4444-444444444444', 'Дмитрий', 'Город только начинает строиться.', 'https://i.pravatar.cc/150?img=4', 'insufficient');

-- Доступные годы (для available_years)
INSERT INTO profile_available_years (profile_id, year) VALUES
    ('11111111-1111-1111-1111-111111111111', 2025),
    ('11111111-1111-1111-1111-111111111111', 2026),
    ('22222222-2222-2222-2222-222222222222', 2026),
    ('33333333-3333-3333-3333-333333333333', 2026),
    ('44444444-4444-4444-4444-444444444444', 2026);

-- 3. Активности для профилей (примеры)

-- Марина (активный покупатель)
INSERT INTO activities (id, profile_id, type, category, title, description, value, timestamp) VALUES
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'view', 'Недвижимость', 'Квартира в центре', 'Просмотр объявления', 1.0, '2025-01-15 10:00:00+00'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'favorite', 'Недвижимость', 'Квартира в центре', 'Добавлено в избранное', 1.0, '2025-01-15 10:05:00+00'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'view', 'Для дома', 'Диван', 'Просмотр объявления', 1.0, '2025-02-10 14:30:00+00'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'favorite', 'Для дома', 'Диван', 'Добавлено в избранное', 1.0, '2025-02-10 14:35:00+00'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'purchase', 'Электроника', 'iPhone 15', 'Покупка', 1.0, '2025-03-05 16:20:00+00'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'view', 'Авто', 'Toyota Camry', 'Просмотр объявления', 1.0, '2025-03-20 09:15:00+00'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'view', 'Недвижимость', 'Студия', 'Просмотр объявления', 1.0, '2025-04-12 11:45:00+00'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'favorite', 'Авто', 'Toyota Camry', 'Добавлено в избранное', 1.0, '2025-04-15 13:00:00+00'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'purchase', 'Для дома', 'Холодильник', 'Покупка', 1.0, '2025-05-08 10:30:00+00'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'view', 'Электроника', 'MacBook Pro', 'Просмотр объявления', 1.0, '2025-06-01 15:45:00+00');

-- Алексей (продавец)
INSERT INTO activities (id, profile_id, type, category, title, description, value, timestamp) VALUES
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'sale', 'Электроника', 'iPhone 14', 'Продажа', 1.0, '2025-01-20 12:00:00+00'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'sale', 'Для дома', 'Стул', 'Продажа', 1.0, '2025-02-14 10:30:00+00'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'sale', 'Авто', 'Lada Granta', 'Продажа', 1.0, '2025-03-10 14:00:00+00'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'sale', 'Электроника', 'PlayStation 5', 'Продажа', 1.0, '2025-04-05 16:30:00+00'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'sale', 'Недвижимость', 'Гараж', 'Продажа', 1.0, '2025-05-20 11:00:00+00'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'view', 'Электроника', 'Телевизор', 'Просмотр объявления', 1.0, '2025-06-15 09:45:00+00');

-- Ольга (незавершённые действия)
INSERT INTO activities (id, profile_id, type, category, title, description, value, timestamp) VALUES
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'view', 'Недвижимость', 'Квартира', 'Просмотр объявления', 1.0, '2025-01-10 10:00:00+00'),
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'favorite', 'Недвижимость', 'Квартира', 'Добавлено в избранное', 1.0, '2025-01-10 10:10:00+00'),
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'view', 'Авто', 'BMW X5', 'Просмотр объявления', 1.0, '2025-02-05 14:20:00+00'),
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'favorite', 'Авто', 'BMW X5', 'Добавлено в избранное', 1.0, '2025-02-05 14:25:00+00'),
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'view', 'Для дома', 'Диван', 'Просмотр объявления', 1.0, '2025-03-01 11:00:00+00'),
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'favorite', 'Для дома', 'Диван', 'Добавлено в избранное', 1.0, '2025-03-01 11:05:00+00');

-- Дмитрий (недостаточно активности)
INSERT INTO activities (id, profile_id, type, category, title, description, value, timestamp) VALUES
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'view', 'Электроника', 'Наушники', 'Просмотр объявления', 1.0, '2025-06-01 10:00:00+00');