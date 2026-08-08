import { useState } from 'react';
import type { Badge, BadgeGroup } from '@/shared/types/recap';

interface Props {
  badges: Badge[];
}

/** Цвет группы — смысловой маркер, а не украшение: время, общение, результат. */
const GROUP_TONE: Record<BadgeGroup, string> = {
  time: 'var(--avito-blue)',
  social: 'var(--avito-purple)',
  result: 'var(--avito-green)',
};

export function BadgesPanel({ badges }: Props) {
  const [openId, setOpenId] = useState<string | null>(null);

  return (
    <section className="panel recap-panel">
      <h2 className="kicker">Достижения · по одной из группы</h2>

      <ul className="badges">
        {badges.map((badge) => {
          const open = openId === badge.id;

          return (
            <li key={badge.id} className="badges__item">
              <button
                type="button"
                className="badge"
                onClick={() => setOpenId(open ? null : badge.id)}
                aria-expanded={open}
              >
                <span className="badge__dot" style={{ background: GROUP_TONE[badge.group] }} />
                <span className="badge__group">{badge.groupTitle}</span>
                <span className="badge__title">{badge.title}</span>
              </button>

              {open && (
                <div className="why">
                  <p className="why__line">{badge.reason}</p>
                  <ul className="why__facts">
                    {badge.facts.map((fact) => (
                      <li key={fact}>{fact}</li>
                    ))}
                  </ul>
                </div>
              )}
            </li>
          );
        })}
      </ul>

      <p className="badges__hint">Нажми на ачивку, чтобы увидеть, за что она выдана.</p>
    </section>
  );
}
