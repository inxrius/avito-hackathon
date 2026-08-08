import { PROFILES, RECAPS } from './mock/fixtures';
import type { Profile, Recap } from '@/shared/types/recap';

/**
 * Единственная точка входа за данными. Сейчас источник — локальные фикстуры;
 * когда появится Go-бэкенд, меняются только тела двух функций:
 *   GET /api/v1/profiles        -> Profile[]
 *   GET /api/v1/recap/:profile  -> Recap
 * Формат ответа уже совпадает с типами в shared/types/recap.ts.
 */

/** Генерация «занимает время» намеренно: сборка города — часть шоу, а не лоадер. */
const GENERATION_MS = 2600;

export function fetchProfiles(): Promise<Profile[]> {
  return Promise.resolve(PROFILES);
}

export class RecapNotFoundError extends Error {
  constructor(profileId: string) {
    super(`Профиль «${profileId}» не найден`);
    this.name = 'RecapNotFoundError';
  }
}

/**
 * Уже собранные города. Экран итогов берёт recap отсюда, чтобы переход
 * «генерация → итоги» был мгновенным, а не запускал сборку заново.
 */
const cache = new Map<string, Recap>();

export function getCachedRecap(profileId: string): Recap | undefined {
  return cache.get(profileId);
}

export function generateRecap(profileId: string, signal?: AbortSignal): Promise<Recap> {
  return new Promise((resolve, reject) => {
    const timer = setTimeout(() => {
      const recap = RECAPS[profileId];
      if (recap) {
        cache.set(profileId, recap);
        resolve(recap);
      } else {
        reject(new RecapNotFoundError(profileId));
      }
    }, GENERATION_MS);

    signal?.addEventListener('abort', () => {
      clearTimeout(timer);
      reject(new DOMException('Генерация отменена', 'AbortError'));
    });
  });
}
