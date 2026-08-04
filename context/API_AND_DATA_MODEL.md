# Черновик API и модели данных MVP

## Статус

Документ фиксирует проектное предложение на основе `объекты_MVP_Итоги_года.docx`. Это не исходное требование кейса и может уточняться командой. Главная цель — заранее синхронизировать frontend и backend.

## Границы

Backend принимает выбранный тестовый профиль и год, рассчитывает recap и отдаёт готовые presentation-объекты. Frontend не получает сырые события и не реализует бизнес-расчёты.

В MVP не входят полноценный каталог объявлений, реальные чаты, регистрация пользователей, платежи, публичная загрузка сырых событий и административный конструктор правил.

## Основные решения

- REST API: `/api/v1`, JSON, поля в `snake_case`.
- Внешние идентификаторы: UUID.
- Время: RFC 3339, UTC.
- `POST /recaps` поддерживает `Idempotency-Key`.
- Готовый `Recap` является неизменяемым snapshot.
- Правила и схема признаков имеют версии; сохраняется hash исходной активности.
- Карточки образуют discriminated union по полю `type`.
- Share-card является отдельной privacy-проекцией, а не копией личного recap.
- Единый формат ошибок — `APIError`.

## Карта endpoint

| Метод и путь | Назначение | Результат |
|---|---|---|
| `GET /profiles` | Список тестовых профилей | `ProfileList` |
| `GET /profiles/{id}` | Детали профиля | `Profile` |
| `POST /recaps` | Запуск/получение генерации | `CreateRecapResponse` |
| `GET /recaps/{id}` | Статус или готовый результат | `Recap` |
| `GET /recaps/{id}/explanation` | Причины решений | `RecapExplanation` |
| `GET /recaps/{id}/share` | Безопасная публичная версия | `ShareCard` |
| `POST /recaps/{id}/interactions` | Продуктовое событие | `InteractionResponse` |

## Ключевые объекты

- Профили: `ProfileSummary`, `ProfileList`, `Profile`.
- Генерация: `CreateRecapRequest`, `CreateRecapResponse`, `RecapStatus`, `RecapLinks`.
- Результат: `Recap`, `RecapProfile`, `RecapGeneration`, `RecapTheme`, `RecapCapabilities`.
- Карточки: `BaseRecapCard`, `IntroCard`, `MetricCard`, `CategoryCard`, `ArchetypeCard`, `AchievementCard`, `OpportunityCard`, `FinalCard`.
- Рекомендации: `Recommendation`, `RecommendationTarget`, `RecommendationReason`.
- Объяснимость: `RecapExplanation`, `DecisionExplanation`, `RuleFact`.
- Публичная проекция: `ShareCard`, `ShareFact`, `ShareAchievement`, `ShareVisual`.
- Аналитика: `InteractionRequest`, `InteractionResponse`.
- Ошибки: `APIError`, `ErrorDetail`.

## Жизненный цикл

`RecapStatus`: `queued`, `processing`, `ready`, `failed`.

MVP может синхронно отвечать `201 Created` со статусом `ready`, но контракт сохраняет промежуточные статусы, чтобы позднее перейти к worker/queue без изменения frontend-модели.

## Типы карточек

Каждая карточка содержит `id`, `type`, `position`, `visibility`, `title` и типизированный `data`. Поддерживаются:

- `intro` — вводный экран;
- `metric` — один измеримый показатель;
- `category` — главная категория;
- `archetype` — поведенческий тип;
- `achievement` — ачивка;
- `opportunity` — незавершённый сценарий;
- `final` — завершение и действия.

Видимость: `personal`, `shareable`, `private`. Порядок определяется `position`.

## Приватность ShareCard

В публичную проекцию запрещено включать `profile_id`, `listing_id`, точные покупки и цены, данные собеседников, тексты сообщений, внутренние score и персональные рекомендации. Допускаются только заранее отобранные факты, безопасная ачивка и абстрактная визуальная тема.

## Ошибки

Минимальный набор кодов: `invalid_argument`, `profile_not_found`, `recap_not_found`, `insufficient_activity`, `idempotency_conflict`, `recap_already_generating`, `rate_limit_exceeded`, `dependency_unavailable`, `internal_error`.

## Приоритет реализации

1. P0: профили, запуск генерации, базовый `Recap`, базовые карточки и единые ошибки — один профиль проходит end-to-end.
2. P1: ачивки, opportunity, рекомендация, объяснение и share-card — полная продуктовая история.
3. P2: продуктовые interactions и аналитика.

## Следующие артефакты

- `openapi.yaml`, проходящий lint;
- `recap.example.json` для независимой разработки frontend;
- схема хранения сырых событий, агрегированных признаков и snapshots;
- contract tests frontend/backend;
- ADR с выбором базы данных.

После согласования первого `recap.example.json` существующие поля карточек следует считать замороженными: допустимы новые необязательные поля, но переименование или изменение смысла требует совместного решения frontend и backend.

