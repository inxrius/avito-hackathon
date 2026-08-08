/**
 * Контракт «Итогов года». Это же — форма ответа `GET /api/v1/recap/:profileId`
 * на бэкенде: фронт ничего не досчитывает, он только рисует то, что пришло.
 * Любая цифра на экране должна иметь источник в этом объекте.
 */

/** Опорный брендовый оттенок района. Девять категорий разведены по четырём
 *  цветам Авито через `shade`, чтобы город не превращался в радугу. */
export type BrandTone = 'blue' | 'green' | 'red' | 'purple';

export type DistrictId =
  'realty' | 'home' | 'electronics' | 'work' | 'personal' | 'hobby' | 'auto' | 'services' | 'pets';

export interface District {
  id: DistrictId;
  title: string;
  /** Всего значимых действий в категории за год. */
  actions: number;
  listings: number;
  dialogs: number;
  /** Доля района в городе, 0..1. Определяет площадь застройки. */
  share: number;
  tone: BrandTone;
  /** 0 — базовый оттенок, 1..2 — светлее. Разводит категории одного тона. */
  shade: 0 | 1 | 2;
}

/** Группы ачивок: по одной из каждой, чтобы итог не превращался в список из 20 бейджей. */
export type BadgeGroup = 'time' | 'social' | 'result';

export interface Badge {
  id: string;
  group: BadgeGroup;
  groupTitle: string;
  title: string;
  /** Ответ на вопрос ТЗ «почему этот пользователь получил именно такой recap». */
  reason: string;
  /** Конкретные факты из активности, которые зажгли ачивку. */
  facts: string[];
}

export interface Chapter {
  index: number;
  title: string;
  /** Крупная цифра главы — то, ради чего глава существует. */
  stat: { value: string; label: string };
  narrative: string;
  /** Район, который достраивается в этой главе. */
  districtId?: DistrictId;
  /** Ачивка, которая вручается в этой главе. */
  badgeId?: string;
}

/** Незавершённый сценарий = стройка на окраине города и повод вернуться в продукт. */
export interface Unfinished {
  id: string;
  title: string;
  count: number;
  ctaLabel: string;
}

export interface Trait {
  title: string;
  reason: string;
}

/**
 * Публичная карточка. Намеренно беднее приватного recap:
 * ни категорий, ни цен, ни счётчиков по районам — только силуэт и роль.
 */
export interface ShareCard {
  cityName: string;
  roleTitle: string;
  styleTitle: string;
  districts: number;
  /** Сид силуэта — тот же город, но без единой цифры о поведении. */
  silhouetteSeed: number;
}

export interface Recap {
  profileId: string;
  year: number;
  /** Версия правил генерации: тот же профиль + та же версия = тот же город. */
  rulesVersion: string;
  seed: number;
  cityName: string;
  totals: {
    actions: number;
    districts: number;
    sites: number;
  };
  role: Trait;
  style: Trait;
  districts: District[];
  chapters: Chapter[];
  badges: Badge[];
  unfinished: Unfinished[];
  shareCard: ShareCard;
}

export interface Profile {
  id: string;
  name: string;
  tagline: string;
  /** Короткий намёк на архетип — чтобы выбор профиля не был выбором вслепую. */
  hint: string;
  tone: BrandTone;
}
