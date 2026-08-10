/** Формы слова для русского счёта: [1 район, 2 района, 5 районов]. */
export type PluralForms = readonly [one: string, few: string, many: string];

/**
 * Русское склонение по числу. Без него в интерфейсе появляются
 * «2 объявлений» и «11 диалога» — мелочь, которая сразу видна.
 */
export function plural(count: number, forms: PluralForms): string {
  const abs = Math.abs(count) % 100;
  const tail = abs % 10;

  if (abs > 10 && abs < 20) return forms[2];
  if (tail === 1) return forms[0];
  if (tail >= 2 && tail <= 4) return forms[1];
  return forms[2];
}

/** Число вместе со склонённым словом: `12 объявлений`. */
export function pluralize(count: number, forms: PluralForms): string {
  return `${count.toLocaleString('ru-RU')} ${plural(count, forms)}`;
}

export const ACTIONS: PluralForms = ['действие', 'действия', 'действий'];
export const LISTINGS: PluralForms = ['объявление', 'объявления', 'объявлений'];
export const DIALOGS: PluralForms = ['диалог', 'диалога', 'диалогов'];
export const DISTRICTS: PluralForms = ['район', 'района', 'районов'];
export const SITES: PluralForms = ['стройка', 'стройки', 'строек'];
