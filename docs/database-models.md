# Модели данных PostgreSQL и ClickHouse

## Назначение

Документ описывает физические модели хранения MVP:

- PostgreSQL — пользователи и доступные годы;
- ClickHouse, база `avito_recap` — активность, результаты recap и продуктовая аналитика.

Все миграции находятся в `backend/migrations` и разделены по СУБД:

- PostgreSQL — `backend/migrations/postgres`;
- ClickHouse — `backend/migrations/clickhouse`.

SQL-миграции являются каноническим источником схемы. При добавлении или изменении таблицы соответствующее описание необходимо обновлять в этом документе в том же pull request.

## PostgreSQL

### `users`

Пользователи и тестовые профили, которые возвращаются через `ProfileSummary` и `Profile` в OpenAPI.

- Primary key: `id`.
- Business key: уникальный `code`.
- `id` передаётся в ClickHouse как внешний `profile_id`.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `id` | `UUID` | `gen_random_uuid()` | Идентификатор пользователя и `Profile.id` в API. |
| `code` | `VARCHAR(100)` | обязательное, unique | Стабильный машинный код тестового сценария. |
| `display_name` | `VARCHAR(100)` | обязательное | Отображаемое имя. |
| `avatar_url` | `TEXT` | необязательное | URL или относительный путь к аватару. |
| `teaser` | `VARCHAR(240)` | обязательное | Краткое описание без раскрытия recap. |
| `description` | `VARCHAR(500)` | обязательное | Расширенное описание профиля. |
| `is_test` | `BOOLEAN` | `true` | Признак тестового профиля MVP. |
| `created_at` | `TIMESTAMPTZ` | текущее время | Время создания. |
| `updated_at` | `TIMESTAMPTZ` | текущее время | Время последнего обновления. |

Для `code`, `display_name`, `teaser` и `description` установлены проверки на непустое значение после `BTRIM`.

### `user_available_years`

Нормализованный список годов, для которых пользователь может сформировать recap.

- Primary key: `(user_id, year)`.
- Foreign key: `user_id -> users.id ON DELETE CASCADE`.
- `year` ограничен диапазоном 2000–2100, совпадающим с OpenAPI.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `user_id` | `UUID` | обязательное | Ссылка на `users.id`. |
| `year` | `SMALLINT` | обязательное | Доступный год итогов. |
| `created_at` | `TIMESTAMPTZ` | текущее время | Время добавления года. |

Поле `available_years` объекта `Profile` собирается из всех строк пользователя. Поле `year` объекта `ProfileSummary` выбирается прикладным слоем из доступных годов, например как максимальный доступный год.

## ClickHouse

### Общие соглашения

- `UUID` используется для идентификаторов сущностей и событий.
- `DateTime64(3, 'UTC')` хранит время в UTC с точностью до миллисекунд.
- `LowCardinality(String)` используется для строк с небольшим количеством повторяющихся значений: кодов, типов и версий.
- `Nullable(T)` обозначает необязательное значение.
- Полиморфные структуры сохраняются как JSON в колонках типа `String`.
- `MergeTree` используется для append-only событий и фактов.
- `ReplacingMergeTree` используется для сущностей, у которых могут появляться новые версии.
- Связи между таблицами являются логическими: ClickHouse не создаёт внешние ключи.
- `ReplacingMergeTree` выполняет дедупликацию во время фоновых merge и не обеспечивает мгновенную транзакционную уникальность.

### Логические связи

```text
PostgreSQL users/profiles
├── activity_events.profile_id
└── recaps.profile_id
    ├── recap_cards
    ├── recommendations
    ├── decision_explanations
    │   └── rule_facts
    ├── share_cards
    ├── interactions
    └── idempotency_keys
```

### Граница с PostgreSQL

Пользователи и тестовые профили не хранятся в ClickHouse. Их основное хранилище — PostgreSQL.

Колонка `profile_id` в аналитических таблицах ClickHouse содержит внешний UUID пользователя или тестового профиля из PostgreSQL. Это логическая ссылка без `FOREIGN KEY`: проверка существования пользователя выполняется вне ClickHouse.

### `activity_events`

Исходные действия тестовых пользователей, из которых рассчитываются метрики, архетипы, ачивки и рекомендации.

- Engine: `MergeTree`.
- Partition key: `toYYYYMM(occurred_at)`.
- Sorting key: `(profile_id, occurred_at, event_type, event_id)`.
- Модель записи: append-only.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `event_id` | `UUID` | обязательное | Идентификатор события. |
| `profile_id` | `UUID` | обязательное | Внешний UUID пользователя или тестового профиля из PostgreSQL. |
| `event_type` | `LowCardinality(String)` | обязательное | Тип пользовательского действия. |
| `category_code` | `LowCardinality(Nullable(String))` | необязательное | Код связанной категории. |
| `listing_id` | `Nullable(UUID)` | необязательное | Связанное объявление; приватное поле. |
| `occurred_at` | `DateTime64(3, 'UTC')` | обязательное | Фактическое время действия. |
| `properties_json` | `String` | `{}` | Дополнительные свойства события в JSON. |
| `ingested_at` | `DateTime64(3, 'UTC')` | текущее время | Время загрузки события в ClickHouse. |

### `recaps`

Состояние генерации и метаданные неизменяемого snapshot «Итогов года».

- Engine: `ReplacingMergeTree(version)`.
- Partition key: `year`.
- Sorting key: `(profile_id, year, id)`.
- Один recap может иметь несколько физических строк, отражающих смену статуса.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `id` | `UUID` | обязательное | Идентификатор recap. |
| `schema_version` | `LowCardinality(String)` | обязательное | Версия публичной схемы результата. |
| `profile_id` | `UUID` | обязательное | Внешний UUID пользователя или тестового профиля из PostgreSQL. |
| `year` | `UInt16` | обязательное | Год итогов. |
| `status` | `Enum8` | обязательное | `queued`, `processing`, `ready` или `failed`. |
| `algorithm_version` | `LowCardinality(String)` | обязательное | Версия правил генерации. |
| `feature_schema_version` | `LowCardinality(String)` | обязательное | Версия схемы рассчитанных признаков. |
| `data_snapshot_id` | `Nullable(UUID)` | необязательное | Идентификатор зафиксированного набора входных данных. |
| `activity_hash` | `String` | обязательное | Hash исходной активности для воспроизводимости. |
| `generated_at` | `Nullable(DateTime64(3, 'UTC'))` | необязательное | Время успешного завершения генерации. |
| `theme_code` | `LowCardinality(String)` | обязательное | Код визуальной темы. |
| `theme_variant` | `LowCardinality(String)` | обязательное | Вариант визуальной темы. |
| `accent_token` | `Nullable(LowCardinality(String))` | необязательное | Токен акцентного цвета дизайн-системы. |
| `share_available` | `Bool` | `true` | Доступна ли публичная share-card. |
| `explanation_available` | `Bool` | `true` | Доступно ли объяснение решений. |
| `feedback_available` | `Bool` | `true` | Доступна ли обратная связь. |
| `version` | `UInt64` | Unix time в миллисекундах | Версия состояния для `ReplacingMergeTree`. |
| `updated_at` | `DateTime64(3, 'UTC')` | текущее время | Время записи версии состояния. |

### `recap_cards`

Упорядоченные presentation-карточки готового recap. Поле `type` является discriminator и определяет структуру `data_json`.

- Engine: `ReplacingMergeTree(created_at)`.
- Sorting key: `(recap_id, position, card_id)`.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `recap_id` | `UUID` | обязательное | Логическая ссылка на recap. |
| `card_id` | `String` | обязательное | Идентификатор карточки внутри recap. |
| `type` | `Enum8` | обязательное | `intro`, `metric`, `category`, `archetype`, `achievement`, `opportunity` или `final`. |
| `position` | `UInt16` | обязательное | Порядок отображения. |
| `visibility` | `Enum8` | обязательное | `personal`, `shareable` или `private`. |
| `eyebrow` | `Nullable(String)` | необязательное | Короткая надпись над заголовком. |
| `title` | `String` | обязательное | Заголовок карточки. |
| `description` | `Nullable(String)` | необязательное | Поясняющий текст. |
| `visual_kind` | `Nullable(LowCardinality(String))` | необязательное | Тип визуального представления. |
| `visual_asset_code` | `Nullable(LowCardinality(String))` | необязательное | Код локального визуального ассета. |
| `explainable` | `Bool` | `false` | Есть ли объяснение бизнес-решения. |
| `data_json` | `String` | `{}` | Типизированные данные карточки в JSON. |
| `created_at` | `DateTime64(3, 'UTC')` | текущее время | Время создания версии карточки. |

### `recommendations`

Следующее полезное действие после просмотра recap.

- Engine: `ReplacingMergeTree(created_at)`.
- Sorting key: `recap_id`.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `recap_id` | `UUID` | обязательное | Recap, для которого выбрана рекомендация. |
| `code` | `LowCardinality(String)` | обязательное | Машинный код рекомендации. |
| `title` | `String` | обязательное | Пользовательский заголовок CTA. |
| `description` | `String` | обязательное | Объяснение следующего действия. |
| `target_kind` | `Enum8` | обязательное | `mock_route`, `deep_link` или `web_url`. |
| `target_value` | `String` | обязательное | Маршрут или URL назначения. |
| `score` | `Nullable(Float32)` | необязательное | Внутренняя оценка релевантности. |
| `reasons_json` | `String` | `[]` | Причины выбора рекомендации в JSON. |
| `created_at` | `DateTime64(3, 'UTC')` | текущее время | Время создания версии рекомендации. |

### `decision_explanations`

Объяснения отдельных решений: архетипа, ачивки, opportunity или рекомендации.

- Engine: `MergeTree`.
- Sorting key: `(recap_id, kind, decision_id)`.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `recap_id` | `UUID` | обязательное | Логическая ссылка на recap. |
| `decision_id` | `UUID` | обязательное | Идентификатор решения. |
| `kind` | `Enum8` | обязательное | `archetype`, `achievement`, `opportunity` или `recommendation`. |
| `code` | `LowCardinality(String)` | обязательное | Код результата. |
| `selected` | `Bool` | обязательное | Было ли решение выбрано правилами. |
| `score` | `Nullable(Float32)` | необязательное | Итоговая оценка решения. |
| `rule_version` | `LowCardinality(String)` | обязательное | Версия применённого правила. |
| `created_at` | `DateTime64(3, 'UTC')` | текущее время | Время фиксации объяснения. |

### `rule_facts`

Конкретные факты и сравнения с порогами, на которых основано решение.

- Engine: `MergeTree`.
- Sorting key: `(recap_id, decision_id, metric)`.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `recap_id` | `UUID` | обязательное | Логическая ссылка на recap. |
| `decision_id` | `UUID` | обязательное | Логическая ссылка на объяснение решения. |
| `metric` | `LowCardinality(String)` | обязательное | Код проверяемого признака. |
| `actual` | `Float64` | обязательное | Фактическое значение признака. |
| `operator` | `Enum8` | обязательное | `eq`, `neq`, `gt`, `gte`, `lt` или `lte`. |
| `threshold` | `Float64` | обязательное | Порог правила. |
| `matched` | `Bool` | обязательное | Выполнилось ли условие. |
| `contribution` | `Nullable(Float32)` | необязательное | Вклад факта в итоговую оценку. |
| `created_at` | `DateTime64(3, 'UTC')` | текущее время | Время фиксации факта. |

### `share_cards`

Отдельная безопасная публичная проекция recap.

- Engine: `ReplacingMergeTree(version)`.
- Sorting key: `recap_id`.
- Запрещены profile ID, listing ID, точные покупки и цены, переписки, внутренние score и персональные рекомендации.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `recap_id` | `UUID` | обязательное | Идентификатор исходного recap. |
| `schema_version` | `LowCardinality(String)` | обязательное | Версия публичной схемы. |
| `year` | `UInt16` | обязательное | Год итогов. |
| `title` | `String` | обязательное | Публичный заголовок. |
| `subtitle` | `String` | обязательное | Публичный подзаголовок или архетип. |
| `facts_json` | `String` | `[]` | Безопасные публичные факты в JSON. |
| `achievement_json` | `Nullable(String)` | необязательное | Публичная ачивка в JSON. |
| `visual_theme` | `LowCardinality(String)` | обязательное | Код публичной темы. |
| `visual_variant` | `LowCardinality(String)` | обязательное | Вариант публичного визуала. |
| `version` | `UInt64` | Unix time в миллисекундах | Версия публичной проекции. |
| `created_at` | `DateTime64(3, 'UTC')` | текущее время | Время создания версии. |

### `interactions`

Продуктовые события просмотра и использования recap.

- Engine: `MergeTree`.
- Partition key: `toYYYYMM(occurred_at)`.
- Sorting key: `(recap_id, occurred_at, event_name, event_id)`.
- Модель записи: append-only.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `event_id` | `UUID` | обязательное | Клиентский идентификатор события. |
| `recap_id` | `UUID` | обязательное | Связанный recap. |
| `session_id` | `UUID` | обязательное | Идентификатор пользовательской сессии. |
| `event_name` | `LowCardinality(String)` | обязательное | Название продуктового события. |
| `occurred_at` | `DateTime64(3, 'UTC')` | обязательное | Клиентское время события. |
| `properties_json` | `String` | `{}` | Дополнительные свойства события в JSON. |
| `received_at` | `DateTime64(3, 'UTC')` | текущее время | Серверное время приёма события. |

### `idempotency_keys`

Соответствие `Idempotency-Key` запросу создания recap и закреплённому результату.

- Engine: `ReplacingMergeTree(created_at)`.
- Sorting key: `key`.
- TTL: `expires_at DELETE`.
- Физическая дедупликация асинхронна; атомарный контроль конфликтов выполняется прикладным слоем.

| Поле | Тип | Обязательность / default | Назначение |
|---|---|---|---|
| `key` | `String` | обязательное | Значение HTTP-заголовка `Idempotency-Key`. |
| `request_hash` | `String` | обязательное | Hash нормализованного тела запроса. |
| `recap_id` | `UUID` | обязательное | Результат, закреплённый за ключом. |
| `created_at` | `DateTime64(3, 'UTC')` | текущее время | Время регистрации ключа. |
| `expires_at` | `DateTime64(3, 'UTC')` | обязательное | Время удаления записи по TTL. |

## Правила изменения схемы

1. Не изменять уже применённую миграцию для обновления существующих окружений.
2. Каждое изменение оформлять новым нумерованным SQL-файлом.
3. Одновременно обновлять этот документ, `docs/openapi.yaml` и `context/API_AND_DATA_MODEL.md`, если меняется публичная модель.
4. Для новой таблицы указывать назначение, engine, sorting key, partition key и TTL.
5. Для новой колонки указывать физический тип, nullable/default и назначение.
6. Не считать `ReplacingMergeTree` заменой `UNIQUE` constraint или транзакционной блокировки.
