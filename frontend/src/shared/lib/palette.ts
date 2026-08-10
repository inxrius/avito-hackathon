import {
  OTHER_DISTRICT,
  type BadgeGroup,
  type BrandTone,
  type District,
} from '@/shared/types/recap';

/**
 * Девять районов на четырёх брендовых цветах. Каждый тон разведён на три
 * светлоты — так категории различимы, но город остаётся в палитре Авито.
 */
const TONES: Record<BrandTone, readonly [string, string, string]> = {
  blue: ['#00aaff', '#4cc4ff', '#8ad8ff'],
  green: ['#05e061', '#4aeb8d', '#8af3b8'],
  red: ['#ff4154', '#ff7583', '#ffa3ac'],
  purple: ['#9c61e7', '#b68bee', '#cfb2f5'],
};

export function toneColor(tone: BrandTone, shade: 0 | 1 | 2): string {
  return TONES[tone][shade];
}

/**
 * «Остальной город» намеренно серый: бэкенд не говорит, из чего он состоит,
 * и красить неизвестное в брендовый цвет — значит делать вид, что мы знаем.
 */
const UNKNOWN_DISTRICT_COLOR = '#4d525f';

export function districtColor(district: Pick<District, 'id' | 'tone' | 'shade'>): string {
  if (district.id === OTHER_DISTRICT) return UNKNOWN_DISTRICT_COLOR;
  return toneColor(district.tone, district.shade);
}

/**
 * Цвет уровня звания. Один и тот же на экране глав и в панели итогов,
 * чтобы «Эксперт» везде выглядел одинаково.
 */
export const LEVEL_TONE: Record<BadgeGroup, string> = {
  newcomer: 'var(--avito-blue)',
  local: 'var(--avito-purple)',
  expert: 'var(--avito-green)',
  guru: 'var(--avito-red)',
};

/** Осветление/затемнение для граней изометрических блоков. */
export function shift(hex: string, amount: number): string {
  const n = parseInt(hex.slice(1), 16);
  const clamp = (v: number) => Math.max(0, Math.min(255, Math.round(v)));
  const r = clamp(((n >> 16) & 0xff) * amount);
  const g = clamp(((n >> 8) & 0xff) * amount);
  const b = clamp((n & 0xff) * amount);
  return `rgb(${r},${g},${b})`;
}
