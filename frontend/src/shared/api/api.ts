/**
 * Адреса ручек бэкенда. Только пути — без запросов и типов:
 * их подставляет слой, который эти константы использует.
 *
 * Полное описание ответов — в docs/openapi.yaml на ветке dev.
 */

export const API_BASE = '/api/v1';

export const endpoints = {
  /** GET — список тестовых профилей с доступными годами. */
  profiles: () => `${API_BASE}/profiles`,

  /** GET — один профиль. */
  profile: (profileId: string) => `${API_BASE}/profiles/${profileId}`,

  /**
   * POST — сформировать итоги за год, тело `{ profile_id, year }`.
   * Первый вызов отдаёт 201, повторный для той же пары — 200 с тем же recap.
   */
  recaps: () => `${API_BASE}/recaps`,

  /** GET — ранее сформированные итоги по id самого recap, не профиля. */
  recap: (recapId: string) => `${API_BASE}/recaps/${recapId}`,

  /** GET — обоснования роли, стиля и ачивок. 409, если выключено в capabilities. */
  recapExplanation: (recapId: string) => `${API_BASE}/recaps/${recapId}/explanation`,

  /** GET — публичная карточка. 409, если выключена в capabilities. */
  recapShare: (recapId: string) => `${API_BASE}/recaps/${recapId}/share`,

  /** POST — продуктовое событие просмотра итогов. */
  recapInteractions: (recapId: string) => `${API_BASE}/recaps/${recapId}/interactions`,
} as const;
