CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE recap_card_type AS ENUM (
    'intro',
    'metric',
    'district',
    'archetype',
    'achievements',
    'summary',
    'final'
);

CREATE TYPE card_visibility AS ENUM ('personal', 'shareable');

CREATE TYPE card_visual_kind AS ENUM (
    'illustration',
    'district',
    'street',
    'calendar',
    'badge',
    'chart',
    'character',
    'skyline'
);

CREATE TYPE achievement_level AS ENUM ('newcomer', 'local', 'expert', 'guru');
CREATE TYPE narrative_source AS ENUM ('mistral', 'template');
CREATE TYPE decision_kind AS ENUM ('archetype_role', 'archetype_style', 'achievement');
CREATE TYPE comparison_operator AS ENUM ('eq', 'neq', 'gt', 'gte', 'lt', 'lte');
CREATE TYPE share_fact_kind AS ENUM ('main_district', 'active_days', 'top_achievement');

CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL CHECK (char_length(name) BETWEEN 1 AND 100),
    description VARCHAR(500) NOT NULL CHECK (char_length(description) BETWEEN 1 AND 500),
    avatar_url TEXT,
    scenario VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE profile_available_years (
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    year SMALLINT NOT NULL CHECK (year BETWEEN 2000 AND 2100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (profile_id, year)
);

CREATE INDEX idx_profile_available_years_year
    ON profile_available_years (year, profile_id);

CREATE TABLE verticals (
    code VARCHAR(50) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE categories (
    code VARCHAR(50) PRIMARY KEY,
    vertical_code VARCHAR(50) NOT NULL REFERENCES verticals(code),
    title VARCHAR(100) NOT NULL,
    UNIQUE (vertical_code, title)
);

CREATE INDEX idx_categories_vertical
    ON categories (vertical_code, code);

CREATE TABLE metric_definitions (
    code VARCHAR(50) PRIMARY KEY,
    title VARCHAR(120) NOT NULL,
    default_unit VARCHAR(30)
);

CREATE TABLE archetype_roles (
    code VARCHAR(50) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE archetype_styles (
    code VARCHAR(50) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE achievement_definitions (
    code VARCHAR(50) PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(300) NOT NULL,
    icon_code VARCHAR(100) NOT NULL
);

CREATE TABLE recaps (
    id UUID PRIMARY KEY,
    profile_id UUID NOT NULL,
    year SMALLINT NOT NULL CHECK (year BETWEEN 2000 AND 2100),
    profile_name VARCHAR(100) NOT NULL,
    profile_avatar_url TEXT,
    schema_version VARCHAR(20) NOT NULL,
    algorithm_version VARCHAR(50) NOT NULL,
    feature_schema_version VARCHAR(50) NOT NULL,
    activity_hash VARCHAR(100) NOT NULL CHECK (activity_hash ~ '^sha256:[a-f0-9]{64}$'),
    generated_at TIMESTAMPTZ NOT NULL,
    narrative_source narrative_source NOT NULL,
    prompt_version VARCHAR(50) NOT NULL,
    narrative_model VARCHAR(100),
    main_vertical_code VARCHAR(50) NOT NULL REFERENCES verticals(code),
    main_vertical_title VARCHAR(100) NOT NULL,
    accent_token VARCHAR(30),
    share_available BOOLEAN NOT NULL,
    explanation_available BOOLEAN NOT NULL,
    feedback_available BOOLEAN NOT NULL,
    summary_title VARCHAR(120) NOT NULL,
    summary_text VARCHAR(500) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (profile_id, year),
    FOREIGN KEY (profile_id, year)
        REFERENCES profile_available_years(profile_id, year)
);

CREATE INDEX idx_recaps_generated_at ON recaps (generated_at);

CREATE TABLE recap_cards (
    recap_id UUID NOT NULL REFERENCES recaps(id) ON DELETE CASCADE,
    card_id VARCHAR(100) NOT NULL,
    type recap_card_type NOT NULL,
    position SMALLINT NOT NULL CHECK (position >= 0),
    visibility card_visibility NOT NULL,
    eyebrow VARCHAR(100),
    title VARCHAR(120) NOT NULL,
    description VARCHAR(500),
    visual_kind card_visual_kind NOT NULL,
    visual_asset_code VARCHAR(100),
    explainable BOOLEAN NOT NULL,
    data JSONB NOT NULL,
    PRIMARY KEY (recap_id, card_id),
    UNIQUE (recap_id, position)
);

CREATE TABLE recap_metrics (
    recap_id UUID NOT NULL REFERENCES recaps(id) ON DELETE CASCADE,
    metric_code VARCHAR(50) NOT NULL REFERENCES metric_definitions(code),
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(30),
    secondary_label VARCHAR(200),
    PRIMARY KEY (recap_id, metric_code)
);

CREATE TABLE recap_archetypes (
    recap_id UUID PRIMARY KEY REFERENCES recaps(id) ON DELETE CASCADE,
    role_code VARCHAR(50) NOT NULL REFERENCES archetype_roles(code),
    style_code VARCHAR(50) NOT NULL REFERENCES archetype_styles(code)
);

CREATE TABLE recap_achievements (
    recap_id UUID NOT NULL REFERENCES recaps(id) ON DELETE CASCADE,
    achievement_code VARCHAR(50) NOT NULL REFERENCES achievement_definitions(code),
    level achievement_level NOT NULL,
    position SMALLINT NOT NULL CHECK (position >= 0),
    title VARCHAR(100) NOT NULL,
    description VARCHAR(300) NOT NULL,
    icon_code VARCHAR(100) NOT NULL,
    metric_code VARCHAR(50) REFERENCES metric_definitions(code),
    current_value DOUBLE PRECISION,
    next_level_threshold DOUBLE PRECISION,
    PRIMARY KEY (recap_id, achievement_code),
    UNIQUE (recap_id, position)
);

CREATE TABLE recap_explanations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recap_id UUID NOT NULL,
    card_id VARCHAR(100) NOT NULL,
    position SMALLINT NOT NULL CHECK (position >= 0),
    kind decision_kind NOT NULL,
    role_code VARCHAR(50) REFERENCES archetype_roles(code),
    style_code VARCHAR(50) REFERENCES archetype_styles(code),
    achievement_code VARCHAR(50) REFERENCES achievement_definitions(code),
    reason VARCHAR(500) NOT NULL,
    rule_version VARCHAR(50) NOT NULL,
    FOREIGN KEY (recap_id, card_id)
        REFERENCES recap_cards(recap_id, card_id)
        ON DELETE CASCADE,
    UNIQUE (recap_id, position),
    CHECK (
        (kind = 'archetype_role' AND role_code IS NOT NULL AND style_code IS NULL AND achievement_code IS NULL)
        OR
        (kind = 'archetype_style' AND role_code IS NULL AND style_code IS NOT NULL AND achievement_code IS NULL)
        OR
        (kind = 'achievement' AND role_code IS NULL AND style_code IS NULL AND achievement_code IS NOT NULL)
    )
);

CREATE TABLE recap_rule_facts (
    explanation_id UUID NOT NULL REFERENCES recap_explanations(id) ON DELETE CASCADE,
    position SMALLINT NOT NULL CHECK (position >= 0),
    metric_code VARCHAR(50) NOT NULL REFERENCES metric_definitions(code),
    actual DOUBLE PRECISION NOT NULL,
    operator comparison_operator NOT NULL,
    threshold DOUBLE PRECISION NOT NULL,
    matched BOOLEAN NOT NULL,
    PRIMARY KEY (explanation_id, position)
);

CREATE TABLE share_cards (
    recap_id UUID PRIMARY KEY REFERENCES recaps(id) ON DELETE CASCADE,
    schema_version VARCHAR(20) NOT NULL,
    profile_name VARCHAR(100) NOT NULL,
    avatar_url TEXT,
    year SMALLINT NOT NULL CHECK (year BETWEEN 2000 AND 2100),
    title VARCHAR(120) NOT NULL,
    subtitle VARCHAR(200) NOT NULL,
    main_vertical_code VARCHAR(50) NOT NULL REFERENCES verticals(code),
    main_vertical_title VARCHAR(100) NOT NULL,
    visual_theme VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE share_facts (
    recap_id UUID NOT NULL REFERENCES share_cards(recap_id) ON DELETE CASCADE,
    position SMALLINT NOT NULL CHECK (position >= 0),
    kind share_fact_kind NOT NULL,
    label VARCHAR(100) NOT NULL,
    value VARCHAR(200) NOT NULL,
    PRIMARY KEY (recap_id, position),
    UNIQUE (recap_id, kind)
);

CREATE TABLE share_achievements (
    recap_id UUID NOT NULL REFERENCES share_cards(recap_id) ON DELETE CASCADE,
    achievement_code VARCHAR(50) NOT NULL REFERENCES achievement_definitions(code),
    position SMALLINT NOT NULL CHECK (position >= 0),
    title VARCHAR(100) NOT NULL,
    level achievement_level NOT NULL,
    icon_code VARCHAR(100) NOT NULL,
    PRIMARY KEY (recap_id, achievement_code),
    UNIQUE (recap_id, position)
);

CREATE INDEX idx_recap_cards_type ON recap_cards (type);
CREATE INDEX idx_recap_metrics_recap ON recap_metrics (recap_id);
CREATE INDEX idx_recap_achievements_recap ON recap_achievements (recap_id);
CREATE INDEX idx_recap_explanations_recap ON recap_explanations (recap_id);
