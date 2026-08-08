import { useCallback, useEffect, useMemo, useRef, useState, type MouseEvent } from 'react';
import { buildCityLayout, sortForPainting, type Site } from './layout';
import { shift, toneColor } from '@/shared/lib/palette';
import { ACTIONS, DIALOGS, DISTRICTS, LISTINGS, pluralize } from '@/shared/lib/plural';
import type { District, DistrictId, Recap } from '@/shared/types/recap';
import './CityCanvas.css';

const CANVAS_W = 1040;
const CANVAS_H = 490;
/** Тайл 2:1 — классическая изометрия, читается без искажений. */
const TILE_W = 46;
const TILE_H = 23;
const LEVEL_H = 11;
const ORIGIN_Y = 120;
/** Здание уже своей клетки: между домами остаётся зазор, и район не выглядит монолитом. */
const BUILDING_SCALE = 0.84;
const PARK_SCALE = 0.6;

interface Props {
  recap: Recap;
  /** Районы, которые уже построены. По ходу глав множество растёт. */
  revealed: ReadonlySet<DistrictId>;
  /** Район в фокусе: подсвечен, остальные приглушены. */
  focus?: DistrictId | null;
  showSites?: boolean;
  interactive?: boolean;
  onSelectDistrict?: (id: DistrictId) => void;
}

interface Tooltip {
  district: District;
  x: number;
  y: number;
}

function project(gx: number, gy: number): { sx: number; sy: number } {
  return {
    sx: CANVAS_W / 2 + ((gx - gy) * TILE_W) / 2,
    sy: ORIGIN_Y + ((gx + gy) * TILE_H) / 2,
  };
}

function fillPath(ctx: CanvasRenderingContext2D, points: number[][], color: string): void {
  ctx.beginPath();
  ctx.moveTo(points[0][0], points[0][1]);
  for (let i = 1; i < points.length; i += 1) ctx.lineTo(points[i][0], points[i][1]);
  ctx.closePath();
  ctx.fillStyle = color;
  ctx.fill();
}

/** Плоская плитка земли под кварталом: по ней читаются границы районов. */
function drawGround(ctx: CanvasRenderingContext2D, sx: number, sy: number, color: string): void {
  const hw = TILE_W / 2;
  const hh = TILE_H / 2;
  fillPath(
    ctx,
    [
      [sx, sy - hh],
      [sx + hw, sy],
      [sx, sy + hh],
      [sx - hw, sy],
    ],
    color,
  );
}

/** Изометрический параллелепипед: верх, левая и правая грани с разной светлотой. */
function drawBox(
  ctx: CanvasRenderingContext2D,
  sx: number,
  sy: number,
  height: number,
  color: string,
  flat: boolean,
  scale = 1,
): void {
  const hw = (TILE_W / 2) * scale;
  const hh = (TILE_H / 2) * scale;
  const top = sy - height;

  if (flat) {
    // Режим пикинга: грани одного цвета, иначе тени попадут в буфер как чужой район.
    fillPath(
      ctx,
      [
        [sx, top - hh],
        [sx + hw, top],
        [sx + hw, sy],
        [sx, sy + hh],
        [sx - hw, sy],
        [sx - hw, top],
      ],
      color,
    );
    return;
  }

  fillPath(
    ctx,
    [
      [sx - hw, top],
      [sx, top + hh],
      [sx, sy + hh],
      [sx - hw, sy],
    ],
    shift(color, 0.55),
  );
  fillPath(
    ctx,
    [
      [sx, top + hh],
      [sx + hw, top],
      [sx + hw, sy],
      [sx, sy + hh],
    ],
    shift(color, 0.78),
  );
  fillPath(
    ctx,
    [
      [sx, top - hh],
      [sx + hw, top],
      [sx, top + hh],
      [sx - hw, top],
    ],
    color,
  );
}

function drawSite(ctx: CanvasRenderingContext2D, site: Site): void {
  const { sx, sy } = project(site.gx, site.gy);
  const hw = TILE_W / 2;
  const hh = TILE_H / 2;
  const height = LEVEL_H * 2.5;

  fillPath(
    ctx,
    [
      [sx, sy - hh],
      [sx + hw, sy],
      [sx, sy + hh],
      [sx - hw, sy],
    ],
    'rgba(255,65,84,0.14)',
  );

  // Каркас без заливки: стройка выглядит незаконченной буквально.
  ctx.strokeStyle = 'rgba(255,65,84,0.85)';
  ctx.lineWidth = 1.5;
  ctx.setLineDash([4, 3]);
  ctx.beginPath();
  ctx.moveTo(sx - hw, sy);
  ctx.lineTo(sx - hw, sy - height);
  ctx.lineTo(sx, sy - height + hh);
  ctx.lineTo(sx, sy + hh);
  ctx.moveTo(sx, sy - height + hh);
  ctx.lineTo(sx + hw, sy - height);
  ctx.lineTo(sx + hw, sy);
  ctx.moveTo(sx + hw, sy - height);
  ctx.lineTo(sx, sy - height - hh);
  ctx.lineTo(sx - hw, sy - height);
  ctx.stroke();
  ctx.setLineDash([]);
}

export function CityCanvas({
  recap,
  revealed,
  focus = null,
  showSites = false,
  interactive = false,
  onSelectDistrict,
}: Props) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const pickRef = useRef<HTMLCanvasElement | null>(null);
  const growthRef = useRef<Map<DistrictId, number>>(new Map());
  const frameRef = useRef<number>(0);
  const [tooltip, setTooltip] = useState<Tooltip | null>(null);
  const [hovered, setHovered] = useState<DistrictId | null>(null);

  const layout = useMemo(
    () => buildCityLayout(recap.districts, recap.seed, recap.totals.sites),
    [recap],
  );
  const painted = useMemo(() => sortForPainting(layout.cells), [layout]);
  const byId = useMemo(() => new Map(recap.districts.map((d) => [d.id, d])), [recap.districts]);
  // Индекс района кодируется в красный канал буфера пикинга, поэтому важен порядок.
  const order = useMemo(() => recap.districts.map((d) => d.id), [recap.districts]);

  const paint = useCallback(() => {
    const canvas = canvasRef.current;
    const pick = pickRef.current;
    if (!canvas || !pick) return;

    const ctx = canvas.getContext('2d');
    const pickCtx = pick.getContext('2d', { willReadFrequently: true });
    if (!ctx || !pickCtx) return;

    const dpr = Math.min(window.devicePixelRatio || 1, 2);
    if (canvas.width !== CANVAS_W * dpr) {
      canvas.width = CANVAS_W * dpr;
      canvas.height = CANVAS_H * dpr;
      pick.width = CANVAS_W;
      pick.height = CANVAS_H;
    }

    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.clearRect(0, 0, CANVAS_W, CANVAS_H);
    pickCtx.clearRect(0, 0, CANVAS_W, CANVAS_H);

    const highlight = focus ?? hovered;

    // Земля рисуется отдельным проходом: плоская плоскость не может перекрыть
    // высокое здание соседней клетки, поэтому мешать её со зданиями нельзя.
    for (const cell of painted) {
      const grow = growthRef.current.get(cell.districtId) ?? 0;
      if (grow <= 0.001) continue;

      const district = byId.get(cell.districtId);
      if (!district) continue;

      const dimmed = highlight !== null && highlight !== cell.districtId;
      const { sx, sy } = project(cell.gx, cell.gy);
      ctx.globalAlpha = (dimmed ? 0.22 : 1) * grow;
      drawGround(ctx, sx, sy, shift(toneColor(district.tone, district.shade), 0.3));
    }

    ctx.globalAlpha = 1;

    for (const cell of painted) {
      const grow = growthRef.current.get(cell.districtId) ?? 0;
      if (grow <= 0.001) continue;

      const district = byId.get(cell.districtId);
      if (!district) continue;

      const { sx, sy } = project(cell.gx, cell.gy);
      const base = toneColor(district.tone, district.shade);
      const dimmed = highlight !== null && highlight !== cell.districtId;
      const isPark = cell.kind === 'park';
      const height = (isPark ? LEVEL_H * 0.7 : cell.levels * LEVEL_H) * grow;

      ctx.globalAlpha = dimmed ? 0.22 : 1;
      drawBox(
        ctx,
        sx,
        sy,
        height,
        isPark ? shift(base, 0.45) : base,
        false,
        isPark ? PARK_SCALE : BUILDING_SCALE,
      );

      if (cell.kind === 'home' && grow > 0.85) {
        ctx.globalAlpha = dimmed ? 0.3 : 1;
        ctx.fillStyle = '#05e061';
        ctx.beginPath();
        ctx.arc(sx, sy - height - TILE_H / 2 - 9, 4, 0, Math.PI * 2);
        ctx.fill();
        ctx.font = '700 11px Manrope, sans-serif';
        ctx.textAlign = 'center';
        ctx.fillText('твой дом', sx, sy - height - TILE_H / 2 - 18);
      }

      ctx.globalAlpha = 1;

      if (interactive) {
        // В буфер пикинга клетка идёт целиком: между домами не должно быть «дырок».
        const index = order.indexOf(cell.districtId) + 1;
        drawBox(pickCtx, sx, sy, height, `rgb(${index},0,0)`, true);
      }
    }

    if (showSites) {
      for (const site of layout.sites) drawSite(ctx, site);
    }
  }, [byId, focus, hovered, interactive, layout, order, painted, showSites]);

  // Рост города: каждый район подтягивается к своей цели, кадр за кадром.
  useEffect(() => {
    const animate = () => {
      let moving = false;
      for (const district of recap.districts) {
        const target = revealed.has(district.id) ? 1 : 0;
        const current = growthRef.current.get(district.id) ?? 0;
        const next = current + (target - current) * 0.09;
        if (Math.abs(target - next) > 0.004) {
          moving = true;
          growthRef.current.set(district.id, next);
        } else {
          growthRef.current.set(district.id, target);
        }
      }
      paint();
      if (moving) frameRef.current = requestAnimationFrame(animate);
    };

    cancelAnimationFrame(frameRef.current);
    frameRef.current = requestAnimationFrame(animate);
    return () => cancelAnimationFrame(frameRef.current);
  }, [paint, recap.districts, revealed]);

  const districtAt = useCallback(
    (clientX: number, clientY: number): DistrictId | null => {
      const canvas = canvasRef.current;
      const pick = pickRef.current;
      if (!canvas || !pick) return null;
      const pickCtx = pick.getContext('2d', { willReadFrequently: true });
      if (!pickCtx) return null;

      const rect = canvas.getBoundingClientRect();
      const x = Math.round(((clientX - rect.left) / rect.width) * CANVAS_W);
      const y = Math.round(((clientY - rect.top) / rect.height) * CANVAS_H);
      if (x < 0 || y < 0 || x >= CANVAS_W || y >= CANVAS_H) return null;

      const [r, , , a] = pickCtx.getImageData(x, y, 1, 1).data;
      if (a === 0 || r === 0) return null;
      return order[r - 1] ?? null;
    },
    [order],
  );

  const handleMove = (event: MouseEvent<HTMLCanvasElement>) => {
    if (!interactive) return;
    const id = districtAt(event.clientX, event.clientY);
    setHovered(id);

    if (!id) {
      setTooltip(null);
      return;
    }
    const district = byId.get(id);
    if (!district) return;

    const rect = event.currentTarget.getBoundingClientRect();
    setTooltip({
      district,
      x: event.clientX - rect.left,
      y: event.clientY - rect.top,
    });
  };

  const handleClick = (event: MouseEvent<HTMLCanvasElement>) => {
    if (!interactive || !onSelectDistrict) return;
    const id = districtAt(event.clientX, event.clientY);
    if (id) onSelectDistrict(id);
  };

  return (
    <div className="city">
      <canvas
        ref={canvasRef}
        className="city__canvas"
        style={{
          aspectRatio: `${CANVAS_W} / ${CANVAS_H}`,
          cursor: hovered ? 'pointer' : 'default',
        }}
        onMouseMove={handleMove}
        onMouseLeave={() => {
          setHovered(null);
          setTooltip(null);
        }}
        onClick={handleClick}
        role="img"
        aria-label={`Город «${recap.cityName}»: ${pluralize(recap.totals.districts, DISTRICTS)}, ${pluralize(recap.totals.actions, ACTIONS)} за год`}
      />
      {/* Внеэкранный буфер: цвет пикселя = индекс района под курсором. */}
      <canvas ref={pickRef} className="city__pick" aria-hidden="true" />

      {tooltip && (
        <div
          className="city__tooltip"
          style={{ left: tooltip.x + 16, top: tooltip.y - 12 }}
          role="status"
        >
          <div className="city__tooltip-title">
            <span
              className="city__dot"
              style={{ background: toneColor(tooltip.district.tone, tooltip.district.shade) }}
            />
            {tooltip.district.title}
          </div>
          <div className="city__tooltip-stats">
            {pluralize(tooltip.district.actions, ACTIONS)} ·{' '}
            {pluralize(tooltip.district.listings, LISTINGS)} ·{' '}
            {pluralize(tooltip.district.dialogs, DIALOGS)}
          </div>
        </div>
      )}
    </div>
  );
}
