# Модели данных PostgreSQL и ClickHouse

SQL-миграции в `backend/migrations` — канонический источник схемы.

## Граница хранилищ

- PostgreSQL хранит пользователей, справочники и однажды сформированный неизменяемый recap.
- ClickHouse хранит только исходные действия и продуктовые события просмотра recap.
- `UNIQUE (user_id, year)` в PostgreSQL гарантирует один recap на пользователя и год.
- Composite FK `(user_id, year)` разрешает генерацию только для доступного пользователю года.
- Ссылки из ClickHouse на PostgreSQL логические; их проверяет приложение.

Актуальная ER-диаграмма находится в `docs/database.dbml`.

## PostgreSQL

### Пользователи

- `users` — имя, описание и аватар пользователя.
- `user_available_years` — множество годов, за которые пользователю доступны итоги.

Происхождение данных не влияет на доменную модель пользователя.

### Справочники

Коды, которые определяют тип или категорию, не хранятся произвольным текстом:

- `verticals` — вертикали Авито;
- `categories` — категории, каждая принадлежит одной вертикали;
- `metric_definitions` — поддерживаемые расчётные метрики и их единицы;
- `archetype_roles` и `archetype_styles` — допустимые части архетипа;
- `achievement_definitions` — допустимые достижения.

Для небольших стабильных множеств используются PostgreSQL enum: тип и видимость
карточки, вид визуала, уровень достижения, источник текста, тип решения,
оператор сравнения и тип публичного факта.

### Сформированные итоги

`recaps` — корень неизменяемого snapshot. Помимо идентификаторов и версий в нём
явно сохраняются сгенерированные `summary_title` и `summary_text`, метаданные
генерации и выбранная главная вертикаль.

Backend должен брать транзакционную advisory-блокировку по паре пользователь–год
перед расчётом. Вместе с unique constraint это исключает параллельный двойной расчёт.

Дочерние данные:

- `recap_cards` — упорядоченные presentation-карточки;
- `recap_metrics` — рассчитанные значения из `metric_definitions`;
- `recap_archetypes` — выбранные роль и стиль;
- `recap_achievements` — до трёх главных достижений;
- `recap_explanations` и `recap_rule_facts` — объяснения решений;
- `share_cards`, `share_facts`, `share_achievements` — безопасная публичная проекция.

`recap_cards.data` остаётся JSONB только для специфичных данных разных типов
карточек. Поля, по которым выполняются связи, фильтрация или проверка множества,
вынесены в типизированные таблицы и колонки.

## ClickHouse

### `activity_events`

- Engine: `MergeTree`.
- Partition key: `toYYYYMM(occurred_at)`.
- Sorting key: `(profile_id, occurred_at, event_type, event_id)`.
- `event_type`, `vertical_code` и `category_code` ограничены ClickHouse enum.
- `profile_id` логически ссылается на `PostgreSQL.users.id`.

### `interactions`

- Engine: `MergeTree`.
- Partition key: `toYYYYMM(occurred_at)`.
- Sorting key: `(recap_id, occurred_at, event_name, event_id)`.
- `event_name` ограничен ClickHouse enum.
- `recap_id` логически ссылается на `PostgreSQL.recaps.id`.

## Изменение схемы

1. Для применённых окружений добавлять новую нумерованную миграцию.
2. Вместе с миграцией обновлять `docs/database.dbml`, этот документ и OpenAPI,
   если меняется публичный контракт.
3. Новые типы и категории сначала добавлять в соответствующий справочник или enum.
4. Не переносить готовый recap в ClickHouse: это транзакционная сущность PostgreSQL.
