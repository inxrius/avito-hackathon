/**
 * Модель экранов. Это НЕ контракт бэкенда — сырой контракт лежит в api/dto.ts,
 * приведение одного к другому в api/adapter.ts.
 *
 * Правило то же, что и было: каждая цифра на экране обязана иметь источник
 * в ответе API. Поля, которых бэкенд не отдаёт, помечены необязательными —
 * компоненты обязаны уметь их не показывать, а не подставлять нули.
 */

export type BrandTone = 'blue' | 'green' | 'red' | 'purple';

/** Код вертикали бэкенда (`goods`, `transport`, …) либо служебный `other`. */
export type DistrictId = string;

/** Синтетический район «весь остальной город»: доля известна, состав — нет. */
export const OTHER_DISTRICT: DistrictId = 'other';

export interface District {
  id: DistrictId;
  title: string;
  /** Доля района в году, 0..1. Приходит из `activity_share`, не пересчитывается. */
  share: number;
  tone: BrandTone;
  shade: 0 | 1 | 2;
  /** Улица — топовая категория внутри района. Есть только у главного. */
  topCategory?: string;
  /** Счётчиков по районам бэкенд пока не отдаёт, поэтому необязательные. */
  actions?: number;
  listings?: number;
  dialogs?: number;
}

/**
 * Групп у ачивок бэкенд не знает — вместо них приходит уровень.
 * Он и стал группировкой: честнее, чем выдумывать свои категории.
 */
export type BadgeGroup = 'newcomer' | 'local' | 'expert' | 'guru';

export interface Badge {
  id: string;
  group: BadgeGroup;
  groupTitle: string;
  title: string;
  /** Из `/explanation`; пока объяснение не загружено — пусто. */
  reason?: string;
  facts: string[];
}

/** Тип карточки бэкенда — определяет, чем глава наполняется на экране. */
export type ChapterKind =
  | 'intro'
  | 'metric'
  | 'district'
  | 'archetype'
  | 'achievements'
  | 'summary'
  | 'final';

export interface Chapter {
  index: number;
  kind: ChapterKind;
  /** Надзаголовок карточки бэкенда: у роли это «Твоя роль в городе». */
  eyebrow?: string;
  title: string;
  /** Есть не у всех карточек: у intro и archetype крупной цифры нет. */
  stat?: { value: string; label: string };
  narrative: string;
  districtId?: DistrictId;
}

export interface Unfinished {
  id: string;
  title: string;
  count: number;
  ctaLabel: string;
}

export interface Trait {
  title: string;
  reason?: string;
}

export interface Recap {
  /** id самого recap — им ходим в /explanation, /share, /interactions. */
  recapId: string;
  profileId: string;
  year: number;
  /** `algorithm_version` бэкенда. */
  rulesVersion: string;
  /** Выведен из `activity_hash`: та же активность — тот же город. */
  seed: number;
  /** Презентационное имя, бэкенд его не отдаёт. */
  cityName: string;
  /** Готовая персональная суммаризация из summary-карточки. Свою не пишем. */
  summaryText?: string;
  totals: {
    activeDays?: number;
    districts: number;
    sites: number;
  };
  role: Trait;
  style: Trait;
  districts: District[];
  chapters: Chapter[];
  badges: Badge[];
  unfinished: Unfinished[];
  capabilities: {
    shareAvailable: boolean;
    explanationAvailable: boolean;
    feedbackAvailable: boolean;
  };
  /** `mistral` или `template` — писал ли текст ИИ. */
  narrativeSource: 'mistral' | 'template';
}

export interface Profile {
  id: string;
  name: string;
  description: string;
  availableYears: number[];
  avatarUrl?: string | null;
  /** Только оформление: бэкенд цвет не присылает, он выводится из позиции. */
  tone: BrandTone;
}
