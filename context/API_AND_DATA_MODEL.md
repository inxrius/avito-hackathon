# API и модель данных

## Контракт

Канонический API-контракт — `docs/openapi.yaml` (OpenAPI 3.0.3).

Backend принимает пользователя и год, синхронно формирует recap один раз и
возвращает готовый неизменяемый snapshot. Повторный `POST /recaps` для той же
пары пользователь–год возвращает сохранённый результат без пересчёта.

## Endpoint

| Метод и путь | Назначение | Результат |
|---|---|---|
| `GET /profiles` | Пользователи с доступными годами | `Profile[]` |
| `GET /profiles/{id}` | Один пользователь | `Profile` |
| `POST /recaps` | Однократно сформировать итог | `201 Recap` или `200 Recap` |
| `GET /recaps/{id}` | Получить сохранённый итог | `Recap` |
| `GET /recaps/{id}/explanation` | Объяснить выборы | `RecapExplanation` |
| `GET /recaps/{id}/share` | Получить публичную проекцию | `ShareCard` |
| `POST /recaps/{id}/interactions` | Записать продуктовое событие | `InteractionResponse` |

Список пользователей возвращается обычным массивом `Profile[]`. Статусы фоновой
генерации и polling отсутствуют: они не соответствуют синхронной однократной генерации.

## Хранение

- PostgreSQL: пользователи, справочники, recap, суммаризация, карточки, метрики,
  архетип, достижения, объяснения и share-card.
- ClickHouse: исходная активность и продуктовые interaction events.

Один итог на год гарантируют транзакционная advisory-блокировка backend и
`recaps UNIQUE (user_id, year)`. Composite FK разрешает генерацию только за год
из `user_available_years`. Суммаризация сохраняется в `recaps.summary_title` и
`recaps.summary_text`.

Вертикали, категории, метрики, роли, стили и достижения задаются справочниками.
Малые закрытые множества задаются enum в PostgreSQL, ClickHouse и OpenAPI.

## Карточки

`RecapCard` — discriminated union по обязательному полю `type`:

- `intro`;
- `metric`;
- `district`;
- `archetype`;
- `achievements`;
- `summary`;
- `final`.

## Ошибки

Единый формат — `APIError`. Поддерживаются коды `invalid_argument`,
`profile_not_found`, `recap_not_found`, `insufficient_activity`,
`rate_limit_exceeded`, `dependency_unavailable`, `internal_error`.
