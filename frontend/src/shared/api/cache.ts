import type { Recap } from '@/shared/types/recap';

/**
 * Кэш на время сессии. Нужен только чтобы переход «генерация → итоги» был
 * мгновенным: сам snapshot неизменяем и живёт на бэкенде, так что при прямой
 * ссылке или перезагрузке страница просто перезапрашивает его по recapId.
 */
const recaps = new Map<string, Recap>();

export function cacheRecap(recap: Recap): void {
  recaps.set(recap.recapId, recap);
}

export function getCachedRecap(recapId: string): Recap | undefined {
  return recaps.get(recapId);
}
