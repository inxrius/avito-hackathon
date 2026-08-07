-- PostgreSQL миграция для Avito Recap
-- Версия: 001
-- Описание: Полная схема БД (таблицы, enum, индексы, связи)

-- 1. ENUM-типы
CREATE TYPE recap_card_type AS ENUM (
    'intro', 'metric', 'district', 'archetype', 
    'achievements', 'summary', 'final'
);

CREATE TYPE card_visibility AS ENUM (
    'personal', 'shareable'
);

CREATE TYPE card_visual_kind AS ENUM (
    'illustration', 'district', 'street', 'calendar',
    'badge', 'chart', 'character', 'skyline'
);

CREATE TYPE achievement_level AS ENUM (
    'newcomer', 'local', 'expert', 'guru'
);

CREATE TYPE narrative_source AS ENUM (
    'mistral', 'template'
);

CREATE TYPE decision_kind AS ENUM (
    'archetype_role', 'archetype_style', 'achievement'
);

CREATE TYPE comparison_operator AS ENUM (
    'eq', 'neq', 'gt', 'gte', 'lt', 'lte'
);

CREATE TYPE share_fact_kind AS ENUM (
    'main_district', 'active_days', 'top_achievement'
);

-- 2. Основные таблицы

-- Пользователи (соответствует модели Profile)
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    avatar_url TEXT,
    scenario VARCHAR(100) NOT NULL,          -- для MVP (тип пользователя)
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Доступные годы для пользователя (для available_years)
CREATE TABLE profile_available_years (
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    year SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (profile_id, year)
);
CREATE INDEX idx_profile_available_years_year ON profile_available_years (year, profile_id);

-- Справочник вертикалей
CREATE TABLE verticals (
    code VARCHAR(50) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE
);

-- Справочник категорий
CREATE TABLE categories (
    code VARCHAR(50) PRIMARY KEY,
    vertical_code VARCHAR(50) NOT NULL REFERENCES verticals(code),
    title VARCHAR(100) NOT NULL,
    UNIQUE (vertical_code, title)
);
CREATE INDEX idx_categories_vertical ON categories (vertical_code, code);

-- Справочник метрик
CREATE TABLE metric_definitions (
    code VARCHAR(50) PRIMARY KEY,
    title VARCHAR(120) NOT NULL,
    default_unit VARCHAR(30)
);

-- Архетипы: роли
CREATE TABLE archetype_roles (
    code VARCHAR(50) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE
);

-- Архетипы: стили
CREATE TABLE archetype_styles (
    code VARCHAR(50) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE
);

-- Определения достижений
CREATE TABLE achievement_definitions (
    code VARCHAR(50) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(300) NOT NULL,
    icon_code VARCHAR(100) NOT NULL
);

-- 3. Таблица recaps (неизменяемый снепшот)
CREATE TABLE recaps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,                    -- "completed", "pending" и т.д.
    year SMALLINT NOT NULL,
    algorithm_version VARCHAR(50) NOT NULL,
    feature_schema_version VARCHAR(50) NOT NULL,
    activity_hash VARCHAR(100) NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    schema_version VARCHAR(20) NOT NULL,
    summary_title VARCHAR(120) NOT NULL,
    summary_text VARCHAR(1000) NOT NULL,
    narrative_source narrative_source NOT NULL,
    prompt_version VARCHAR(50) NOT NULL,
    narrative_model VARCHAR(100),
    main_vertical_code VARCHAR(50) REFERENCES verticals(code),
    accent_token VARCHAR(30),
    UNIQUE (profile_id, year)
);
CREATE INDEX idx_recaps_generated_at ON recaps (generated_at);

-- 4. Карточки recap (связь один-ко-многим)
CREATE TABLE recap_cards (
    recap_id UUID NOT NULL REFERENCES recaps(id) ON DELETE CASCADE,
    card_id VARCHAR(100) NOT NULL,
    type recap_card_type NOT NULL,
    position SMALLINT NOT NULL,
    visibility card_visibility NOT NULL DEFAULT 'personal',
    eyebrow VARCHAR(100),
    title VARCHAR(120) NOT NULL,
    description VARCHAR(500),
    visual_kind card_visual_kind,
    visual_asset_code VARCHAR(100),
    explainable BOOLEAN NOT NULL DEFAULT FALSE,
    data JSONB NOT NULL DEFAULT '{}',
    PRIMARY KEY (recap_id, card_id),
    UNIQUE (recap_id, position)
);

-- 5. Метрики recap
CREATE TABLE recap_metrics (
    recap_id UUID NOT NULL REFERENCES recaps(id) ON DELETE CASCADE,
    metric_code VARCHAR(50) NOT NULL REFERENCES metric_definitions(code),
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(30),
    secondary_label VARCHAR(200),
    PRIMARY KEY (recap_id, metric_code)
);

-- 6. Архетипы recap
CREATE TABLE recap_archetypes (
    recap_id UUID PRIMARY KEY REFERENCES recaps(id) ON DELETE CASCADE,
    role_code VARCHAR(50) NOT NULL REFERENCES archetype_roles(code),
    style_code VARCHAR(50) NOT NULL REFERENCES archetype_styles(code)
);

-- 7. Достижения recap
CREATE TABLE recap_achievements (
    recap_id UUID NOT NULL REFERENCES recaps(id) ON DELETE CASCADE,
    achievement_code VARCHAR(50) NOT NULL REFERENCES achievement_definitions(code),
    level achievement_level NOT NULL,
    position SMALLINT NOT NULL,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(300) NOT NULL,
    icon_code VARCHAR(100) NOT NULL,
    metric_code VARCHAR(50) REFERENCES metric_definitions(code),
    current_value DOUBLE PRECISION,
    next_level_threshold DOUBLE PRECISION,
    PRIMARY KEY (recap_id, achievement_code),
    UNIQUE (recap_id, position)
);

-- 8. Объяснения для карточек
CREATE TABLE recap_explanations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recap_id UUID NOT NULL,
    card_id VARCHAR(100) NOT NULL,
    kind decision_kind NOT NULL,
    role_code VARCHAR(50) REFERENCES archetype_roles(code),
    style_code VARCHAR(50) REFERENCES archetype_styles(code),
    achievement_code VARCHAR(50) REFERENCES achievement_definitions(code),
    reason VARCHAR(500) NOT NULL,
    rule_version VARCHAR(50) NOT NULL,
    FOREIGN KEY (recap_id, card_id) REFERENCES recap_cards(recap_id, card_id) ON DELETE CASCADE
);

-- 9. Факты правил
CREATE TABLE recap_rule_facts (
    explanation_id UUID NOT NULL REFERENCES recap_explanations(id) ON DELETE CASCADE,
    metric_code VARCHAR(50) NOT NULL REFERENCES metric_definitions(code),
    actual DOUBLE PRECISION NOT NULL,
    operator comparison_operator NOT NULL,
    threshold DOUBLE PRECISION NOT NULL,
    matched BOOLEAN NOT NULL,
    PRIMARY KEY (explanation_id, metric_code)
);

-- 10. Публичные share-карточки
CREATE TABLE share_cards (
    recap_id UUID PRIMARY KEY REFERENCES recaps(id) ON DELETE CASCADE,
    schema_version VARCHAR(20) NOT NULL,
    title VARCHAR(120) NOT NULL,
    subtitle VARCHAR(200) NOT NULL,
    main_vertical_code VARCHAR(50) NOT NULL REFERENCES verticals(code),
    colors VARCHAR(7)[] NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 11. Факты для share-карточки
CREATE TABLE share_facts (
    recap_id UUID NOT NULL REFERENCES share_cards(recap_id) ON DELETE CASCADE,
    position SMALLINT NOT NULL,
    kind share_fact_kind NOT NULL,
    label VARCHAR(100) NOT NULL,
    value VARCHAR(200) NOT NULL,
    PRIMARY KEY (recap_id, position),
    UNIQUE (recap_id, kind)
);

-- 12. Достижения на share-карточке
CREATE TABLE share_achievements (
    recap_id UUID NOT NULL REFERENCES share_cards(recap_id) ON DELETE CASCADE,
    achievement_code VARCHAR(50) NOT NULL,
    position SMALLINT NOT NULL,
    PRIMARY KEY (recap_id, achievement_code),
    UNIQUE (recap_id, position)
);

-- 13. Таблица активностей (события пользователя)
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    value FLOAT NOT NULL DEFAULT 1.0,
    timestamp TIMESTAMPTZ NOT NULL
);
CREATE INDEX idx_activities_profile_id ON activities (profile_id);
CREATE INDEX idx_activities_timestamp ON activities (timestamp);

-- 14. Таблица взаимодействий с recap (продуктовые события)
CREATE TABLE interactions (
    id SERIAL PRIMARY KEY,
    recap_id UUID NOT NULL REFERENCES recaps(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_interactions_recap_id ON interactions (recap_id);
CREATE INDEX idx_interactions_event_type ON interactions (event_type);

-- 15. Дополнительные индексы для производительности
CREATE INDEX idx_recap_cards_type ON recap_cards (type);
CREATE INDEX idx_recap_metrics_recap ON recap_metrics (recap_id);
CREATE INDEX idx_recap_achievements_recap ON recap_achievements (recap_id);
CREATE INDEX idx_recap_explanations_recap ON recap_explanations (recap_id);