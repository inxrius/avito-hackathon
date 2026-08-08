import { toneColor } from '@/shared/lib/palette';
import type { District, DistrictId } from '@/shared/types/recap';

interface Props {
  districts: District[];
  revealed: ReadonlySet<DistrictId>;
  active: DistrictId | null;
  onSelect: (id: DistrictId) => void;
}

export function DistrictLegend({ districts, revealed, active, onSelect }: Props) {
  const max = Math.max(...districts.map((d) => d.actions), 1);

  return (
    <ul className="legend">
      {districts.map((district) => {
        const built = revealed.has(district.id);
        const color = toneColor(district.tone, district.shade);

        return (
          <li key={district.id}>
            <button
              type="button"
              className={
                'legend__item' +
                (active === district.id ? ' legend__item--active' : '') +
                (built ? '' : ' legend__item--pending')
              }
              onClick={() => onSelect(district.id)}
              aria-pressed={active === district.id}
            >
              <span className="legend__head">
                <span className="legend__dot" style={{ background: color }} />
                <span className="legend__title">{district.title}</span>
                <span className="legend__value">{district.actions}</span>
              </span>
              <span className="legend__track">
                <span
                  className="legend__bar"
                  style={{
                    width: built ? `${(district.actions / max) * 100}%` : '0%',
                    background: color,
                  }}
                />
              </span>
            </button>
          </li>
        );
      })}
    </ul>
  );
}
