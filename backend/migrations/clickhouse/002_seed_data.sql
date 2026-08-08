INSERT INTO recap.activity_events
SELECT
    generateUUIDv4(),
    toUUID('11111111-1111-1111-1111-111111111111'),
    multiIf(number % 6 = 0, 'purchase_completed', number % 4 = 0, 'favorite_added', 'listing_viewed'),
    'goods',
    multiIf(number % 3 = 0, 'home_and_garden', 'electronics'),
    toDateTime64('2026-01-01 12:00:00', 3, 'UTC') + toIntervalDay(number),
    now64(3)
FROM numbers(18);

INSERT INTO recap.activity_events
SELECT
    generateUUIDv4(),
    toUUID('22222222-2222-2222-2222-222222222222'),
    multiIf(number % 3 = 0, 'sale_completed', 'listing_published'),
    'goods',
    multiIf(number % 2 = 0, 'electronics', 'home_and_garden'),
    toDateTime64('2026-02-01 10:00:00', 3, 'UTC') + toIntervalDay(number),
    now64(3)
FROM numbers(18);

INSERT INTO recap.activity_events
SELECT
    generateUUIDv4(),
    toUUID('33333333-3333-3333-3333-333333333333'),
    multiIf(number % 5 = 0, 'favorite_added', 'listing_viewed'),
    multiIf(number % 4 = 0, 'transport', number % 4 = 1, 'real_estate', 'goods'),
    multiIf(number % 4 = 0, 'cars', number % 4 = 1, 'apartments', number % 4 = 2, 'electronics', 'hobbies_and_leisure'),
    toDateTime64('2026-03-01 09:00:00', 3, 'UTC') + toIntervalDay(number),
    now64(3)
FROM numbers(20);

INSERT INTO recap.activity_events
SELECT
    generateUUIDv4(),
    toUUID('44444444-4444-4444-4444-444444444444'),
    'listing_viewed',
    'goods',
    'electronics',
    toDateTime64('2026-04-01 12:00:00', 3, 'UTC'),
    now64(3);
