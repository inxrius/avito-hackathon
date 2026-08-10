import { endpoints } from './api.ts';
import type {
  APIErrorCode,
  APIErrorDTO,
  ProfileDTO,
  RecapDTO,
  RecapExplanationDTO,
  ShareCardDTO,
} from './dto.ts';

/**
 * HTTP-слой. Ничего не знает про экраны: отдаёт DTO как есть
 * и превращает ответ бэкенда об ошибке в типизированное исключение.
 */

/** Ошибка бэкенда в его формате. `code` — то, по чему ветвится UI. */
export class APIError extends Error {
  readonly code: APIErrorCode;
  readonly status: number;
  readonly requestId: string;

  constructor(status: number, payload: APIErrorDTO) {
    super(payload.message);
    this.name = 'APIError';
    this.status = status;
    this.code = payload.code;
    this.requestId = payload.request_id;
  }
}

/** Сеть недоступна или ответ нечитаем — это не ошибка бизнес-логики. */
export class NetworkError extends Error {
  constructor(cause: unknown) {
    super('Сервис недоступен');
    this.name = 'NetworkError';
    this.cause = cause;
  }
}

async function request<T>(url: string, init?: RequestInit): Promise<T> {
  let response: Response;
  try {
    response = await fetch(url, {
      ...init,
      headers: { 'Content-Type': 'application/json', ...init?.headers },
    });
  } catch (cause) {
    throw new NetworkError(cause);
  }

  const body: unknown = await response.json().catch(() => null);

  if (!response.ok) {
    if (body && typeof body === 'object' && 'code' in body) {
      throw new APIError(response.status, body as APIErrorDTO);
    }
    throw new NetworkError(new Error(`HTTP ${response.status}`));
  }

  return body as T;
}

export function getProfiles(): Promise<ProfileDTO[]> {
  return request<ProfileDTO[]>(endpoints.profiles());
}

/**
 * Формирует итоги. Бэкенд отвечает 201 при первом вызове и 200, если snapshot
 * уже существует, — оба статуса успешные и тело у них одинаковое.
 */
export function createRecap(profileId: string, year: number): Promise<RecapDTO> {
  return request<RecapDTO>(endpoints.recaps(), {
    method: 'POST',
    body: JSON.stringify({ profile_id: profileId, year }),
  });
}

export function getRecap(recapId: string): Promise<RecapDTO> {
  return request<RecapDTO>(endpoints.recap(recapId));
}

export function getExplanation(recapId: string): Promise<RecapExplanationDTO> {
  return request<RecapExplanationDTO>(endpoints.recapExplanation(recapId));
}

export function getShareCard(recapId: string): Promise<ShareCardDTO> {
  return request<ShareCardDTO>(endpoints.recapShare(recapId));
}
