CREATE DATABASE IF NOT EXISTS recap;

CREATE TABLE IF NOT EXISTS recap.activity_events
(
    event_id UUID,
    profile_id UUID,
    event_type LowCardinality(String),
    vertical_code LowCardinality(String),
    category_code LowCardinality(String),
    occurred_at DateTime64(3, 'UTC'),
    received_at DateTime64(3, 'UTC') DEFAULT now64(3)
)
ENGINE = ReplacingMergeTree(received_at)
PARTITION BY toYear(occurred_at)
ORDER BY (profile_id, event_id);

CREATE TABLE IF NOT EXISTS recap.interactions
(
    event_id UUID,
    recap_id UUID,
    session_id UUID,
    event_name LowCardinality(String),
    occurred_at DateTime64(3, 'UTC'),
    properties String DEFAULT '{}',
    received_at DateTime64(3, 'UTC') DEFAULT now64(3)
)
ENGINE = ReplacingMergeTree(received_at)
PARTITION BY toYYYYMM(occurred_at)
ORDER BY (recap_id, event_id);
