import type { BrandTone } from '@/shared/types/recap';

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

/** Осветление/затемнение для граней изометрических блоков. */
export function shift(hex: string, amount: number): string {
  const n = parseInt(hex.slice(1), 16);
  const clamp = (v: number) => Math.max(0, Math.min(255, Math.round(v)));
  const r = clamp(((n >> 16) & 0xff) * amount);
  const g = clamp(((n >> 8) & 0xff) * amount);
  const b = clamp((n & 0xff) * amount);
  return `rgb(${r},${g},${b})`;
}
