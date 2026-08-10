# История разработки и ключевых коммитов

Этот документ дополняет основной `README.md` и фиксирует ключевые этапы командной разработки Avito Recap City. Полная техническая история остаётся доступна через Git:

```bash
git log --oneline --graph --decorate --all
```

Разработка велась через персональные feature/fix-ветки и pull requests в интеграционную ветку `dev`. Ниже перечислены ключевые содержательные коммиты и этапы; merge-коммиты опущены там, где они не добавляют отдельной реализации.

## Участники

| Участник | GitHub | Основные зоны |
|---|---|---|
| Александр Цыков | `At0-m` | Архитектура решения, OpenAPI, recap domain/analytics, personalization, narrative/Mistral, pipeline, backend-интеграция, E2E-стабилизация, финальная frontend-полировка, тесты, CI, Docker |
| Станислав Шегай | `inxrius` | Базовая структура проекта, backend models, repository/service/HTTP слои, интеграция generator pipeline |
| Дмитрий | `DemiusHTTV` | Frontend, пользовательский flow, Canvas-город, recap UI, интеграция frontend с актуальным backend contract |
| Роман Карелин | `mamooin` | Восстановление и структурирование требований, ранние архитектурные решения, OpenAPI и документация модели данных |

> Зоны ответственности пересекались: интеграция контрактов, review, исправления и финальная стабилизация выполнялись совместно.

## 1. Инициализация проекта и требований

### `inxrius` - Станислав Шегай

- `d0db524` - `Initial commit`
- `370d04c`, `f772114`, `f156d45` - базовая структура репозитория и проекта.
- `aa4b548` - начальные backend data structures.
- `51c1d7b` - migrations и seed mock data.
- `cd61965` - начальная структура backend и Docker Compose.

### `mamooin` - Роман Карелин

- `a235288` - `docs: restore project requirements and context`
- `5e65d36` - `add openAPI and db models docs`

На этом этапе были восстановлены требования, зафиксированы основные контракты и модель данных, на которых затем строилась реализация.

## 2. Backend foundation

### `inxrius` - Станислав Шегай

- `c1f76fd` - `feat(backend): implement repository layer`
- `36ce6ea` - `feat(backend): implement service layer`
- `1d6656e` - `feat(backend): implement HTTP handlers`
- `7607842` - `chore: init project structure and docker-compose setup`
- `9779a31` - `feat: integrate colleague's recap generator pipeline`
- `87534e1` - обновление personalization/event/category integration.
- `272aa2d` - `feat: add year filtering for activities`
- `1e872bd` - alignment backend models/handlers with OpenAPI.
- `a95b851` - migrations and seed data.
- `edb9c9d` - database connection package.

Эти коммиты сформировали application shell вокруг recap-домена: HTTP, service и repository слои, подключение хранилищ и маршрутизацию.

## 3. Recap domain, analytics и personalization

### `At0-m` - Александр Цыков

- `e073d76`, `1f2b37d` - изменения OpenAPI под актуальный recap contract.
- `ea70ae3` - `feat(recap): add analytics and recap domain`
- `7a9d436` - `feat(recap): add personalization rules`
- `28e3b88` - `feat(recap): add narrative generation`
- `b56e771` - `feat(recap): assemble personalization pipeline`

На этом этапе появились детерминированный подсчёт признаков, выбор роли/стиля/достижений, narrative layer и сборка итогового recap.

## 4. Интеграция с PostgreSQL и ClickHouse

### `At0-m` - Александр Цыков

После первого E2E-прогона были исправлены реальные интеграционные проблемы между доменным pipeline и хранилищами:

- `68655f2` - `fix(infra): complete ClickHouse and backend configuration`
- `da52ee3` - `fix(recap): integrate generator with backend services`
- `460230c` - `fix(api): align recap models with integration`
- `d51fc1d` - `fix(recap): correct personalization and Mistral flow`
- `cc13bbd` - `fix(recap): align integration ports`
- `f836e97` - `db: add recap storage and ClickHouse seed data`
- `45bc6e3` - update backend development commands.
- `d519664` - Docker Compose/backend integration fixes.

В эту фазу вошли исправления ClickHouse UUID/String compatibility, `DateTime64` decoding, interaction queries, persistence snapshot'ов, Mistral structured output и template fallback.

## 5. Frontend и переход с mock contract на реальный API

### `DemiusHTTV` - Дмитрий

- `7182471` - `fix(frontend): align frontend with backend contract`
- `cc692df` - `fix(frontend): complete recap flow on real backend contract`
- `94a9878` - дополнительные fixes к backend contract.

Frontend был переведён с fixtures на настоящий API, сохранив существующий пользовательский flow и Canvas-представление города.

### `At0-m` - Александр Цыков

- `1ada31c` - `fix(frontend): finalize recap experience`

Финальная полировка recap-flow: отображение backend summary, достижений и согласование визуального сценария с фактическими backend-данными.

## 6. Финальная интеграция

Ключевые merge-этапы в `dev`:

- PR #12 - backend integration fixes и Docker Compose updates.
- PR #14 - frontend/backend contract alignment.
- PR #15 - final recap experience.

Финальный E2E-путь:

```text
Browser
  -> React frontend
  -> Go/Gin API
  -> ClickHouse activity events
  -> analytics + personalization
  -> Mistral or deterministic template fallback
  -> PostgreSQL immutable recap snapshot
  -> explanation/share projections
```

## 7. Quality gate перед сдачей

В финальной фазе проект был дополнен:

- frontend API/integration tests;
- обновлёнными backend tests;
- Go formatting/vet checks;
- ESLint и production frontend build;
- GitHub Actions CI;
- Dockerfile для frontend;
- Nginx `/api` proxy;
- единым `docker compose up -d --build` для всего стека.

Это позволяет проверяющему поднять проект без ручной сборки отдельных частей.

## 8. Что сознательно осталось за рамками MVP

MVP визуально акцентирует главный район пользователя. Следующее продуктово значимое расширение - безопасно отдавать распределение активности сразу по нескольким вертикалям и строить полноценный многорайонный CityCanvas.

Архитектурный roadmap также рассматривает асинхронную генерацию, broker/outbox/inbox, выделенный Analytics Service, materialized views и воспроизводимые benchmark/fault tests. Эти компоненты предполагается добавлять только после измерений, подтверждающих их необходимость.
