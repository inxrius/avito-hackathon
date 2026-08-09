import { adaptProfiles, adaptRecap, applyExplanation } from './adapter';
import { APIError, createRecap, getExplanation, getProfiles, getRecap, getShareCard } from './client';
import type { ShareCardDTO } from './dto';
import type { Profile, Recap } from '@/shared/types/recap';

/**
 * Прикладной слой: страницы ходят сюда и получают уже модель экранов.
 * Про HTTP и DTO знают client.ts и adapter.ts, компоненты — нет.
 */

export function fetchProfiles(): Promise<Profile[]> {
  return getProfiles().then(adaptProfiles);
}

/**
 * Собирает итоги. Повторный вызов для той же пары профиль–год бэкенд
 * не пересчитывает: возвращает тот же snapshot с 200 вместо 201.
 */
export function generateRecap(profileId: string, year: number): Promise<Recap> {
  return createRecap(profileId, year).then(adaptRecap);
}

export function loadRecap(recapId: string): Promise<Recap> {
  return getRecap(recapId).then(adaptRecap);
}

/**
 * Догружает обоснования. Вызывается лениво — только когда пользователь
 * раскрывает «почему», и только если это разрешено в capabilities.
 */
export function loadExplanation(recap: Recap): Promise<Recap> {
  if (!recap.capabilities.explanationAvailable) return Promise.resolve(recap);
  return getExplanation(recap.recapId).then((dto) => applyExplanation(recap, dto));
}

export function loadShareCard(recapId: string): Promise<ShareCardDTO> {
  return getShareCard(recapId);
}

/** Профиль без достаточной активности — штатный ответ бэкенда, не сбой. */
export function isInsufficientActivity(error: unknown): boolean {
  return error instanceof APIError && error.code === 'insufficient_activity';
}
