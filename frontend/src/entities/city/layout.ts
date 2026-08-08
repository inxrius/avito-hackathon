import { createRandom } from '@/shared/lib/random';
import type { District, DistrictId } from '@/shared/types/recap';

export const GRID = 14;
/** Потолок этажности. Выше — небоскрёбы начинают перекрывать заголовок сцены. */
const MAX_LEVELS = 9;

export type CellKind = 'building' | 'park' | 'home';

export interface Cell {
  gx: number;
  gy: number;
  districtId: DistrictId;
  /** Этажность, 1..MAX_LEVELS. Растёт с активностью в категории. */
  levels: number;
  kind: CellKind;
}

export interface Site {
  gx: number;
  gy: number;
}

export interface CityLayout {
  cells: Cell[];
  /** Незавершённые сценарии — стройки по краю карты. */
  sites: Site[];
}

interface Rect {
  x: number;
  y: number;
  w: number;
  h: number;
}

/**
 * Slice-and-dice: рекурсивно режем квадрат на прямоугольники, площадь которых
 * пропорциональна доле района. Получаются кварталы, а не полосы, и каждый
 * район остаётся цельным куском — по нему можно кликнуть.
 */
function carve(items: District[], rect: Rect, out: Map<DistrictId, Rect>): void {
  if (items.length === 0) return;
  if (items.length === 1) {
    out.set(items[0].id, rect);
    return;
  }

  const total = items.reduce((sum, d) => sum + d.share, 0);
  // Набираем первую группу примерно до половины веса — так дерево остаётся сбалансированным.
  let acc = 0;
  let split = 0;
  while (split < items.length - 1 && acc + items[split].share <= total / 2) {
    acc += items[split].share;
    split += 1;
  }
  if (split === 0) split = 1;

  const ratio = acc / total;
  const head = items.slice(0, split);
  const tail = items.slice(split);

  // Режем всегда поперёк длинной стороны, иначе кварталы вырождаются в нитки.
  if (rect.w >= rect.h) {
    const cut = Math.max(1, Math.min(rect.w - 1, Math.round(rect.w * ratio)));
    carve(head, { ...rect, w: cut }, out);
    carve(tail, { x: rect.x + cut, y: rect.y, w: rect.w - cut, h: rect.h }, out);
  } else {
    const cut = Math.max(1, Math.min(rect.h - 1, Math.round(rect.h * ratio)));
    carve(head, { ...rect, h: cut }, out);
    carve(tail, { x: rect.x, y: rect.y + cut, w: rect.w, h: rect.h - cut }, out);
  }
}

export function buildCityLayout(
  districts: District[],
  seed: number,
  siteCount: number,
): CityLayout {
  const rand = createRandom(seed);
  const ordered = [...districts].sort((a, b) => b.share - a.share || a.id.localeCompare(b.id));

  const rects = new Map<DistrictId, Rect>();
  carve(ordered, { x: 0, y: 0, w: GRID, h: GRID }, rects);

  const maxActions = Math.max(...districts.map((d) => d.actions), 1);
  const cells: Cell[] = [];

  // Дом — в самом крупном районе года: визуальный якорь «здесь ты жил».
  const homeDistrict = ordered[0]?.id;
  let homePlaced = false;

  for (const district of ordered) {
    const rect = rects.get(district.id);
    if (!rect) continue;

    const intensity = district.actions / maxActions;
    const centerX = rect.x + rect.w / 2;
    const centerY = rect.y + rect.h / 2;

    for (let gy = rect.y; gy < rect.y + rect.h; gy += 1) {
      for (let gx = rect.x; gx < rect.x + rect.w; gx += 1) {
        const roll = rand();

        // Немного зелени, чтобы кварталы не выглядели сплошной стеной.
        if (roll < 0.1) {
          cells.push({ gx, gy, districtId: district.id, levels: 1, kind: 'park' });
          continue;
        }

        // Этажность — главная метафора: чем больше действий в категории, тем выше дома.
        // К центру квартала добавляется рост, чтобы силуэт читался как настоящий район.
        const distance = Math.hypot(gx - centerX + 0.5, gy - centerY + 0.5);
        const falloff = 1 - Math.min(1, distance / (Math.max(rect.w, rect.h) / 2 + 0.5));
        const raw = 1 + intensity * 7 * (0.35 + 0.65 * falloff) + rand() * 1.5;
        const levels = Math.max(1, Math.min(MAX_LEVELS, Math.round(raw)));

        const isHome = !homePlaced && district.id === homeDistrict && roll > 0.9;
        if (isHome) homePlaced = true;

        cells.push({
          gx,
          gy,
          districtId: district.id,
          levels: isHome ? Math.max(2, levels) : levels,
          kind: isHome ? 'home' : 'building',
        });
      }
    }
  }

  // Фолбэк: если бросок кубика так и не выпал, ставим дом в центр главного квартала.
  if (!homePlaced && homeDistrict) {
    const rect = rects.get(homeDistrict);
    if (rect) {
      const gx = rect.x + Math.floor(rect.w / 2);
      const gy = rect.y + Math.floor(rect.h / 2);
      const cell = cells.find((c) => c.gx === gx && c.gy === gy);
      if (cell) cell.kind = 'home';
    }
  }

  // Стройки выносим за границу города: они ещё не часть года, но уже видны.
  const sites: Site[] = [];
  for (let i = 0; i < siteCount; i += 1) {
    const along = Math.floor(rand() * GRID);
    sites.push(i % 2 === 0 ? { gx: GRID + 1, gy: along } : { gx: along, gy: GRID + 1 });
  }

  return { cells, sites };
}

/** Сортировка «от дальнего к ближнему»: в изометрии это единственный корректный порядок отрисовки. */
export function sortForPainting(cells: Cell[]): Cell[] {
  return [...cells].sort((a, b) => a.gx + a.gy - (b.gx + b.gy) || a.gx - b.gx);
}
