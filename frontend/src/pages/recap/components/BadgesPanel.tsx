import { useState } from 'react';
import { LEVEL_TONE } from '@/shared/lib/palette';
import type { Badge } from '@/shared/types/recap';

interface Props {
  badges: Badge[];
}

export function BadgesPanel({ badges }: Props) {
  const [openId, setOpenId] = useState<string | null>(null);

  return (
    <section className="panel recap-panel">
      <h2 className="kicker">Городские звания</h2>

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
                <span className="badge__dot" style={{ background: LEVEL_TONE[badge.group] }} />
                <span className="badge__group">{badge.groupTitle}</span>
                <span className="badge__title">{badge.title}</span>
              </button>

              {open && (
                <div className="why">
                  {badge.reason && <p className="why__line">{badge.reason}</p>}
                  {/* Факты приходят из /explanation отдельным запросом. */}
                  {badge.facts.length > 0 && (
                    <ul className="why__facts">
                      {badge.facts.map((fact) => (
                        <li key={fact}>{fact}</li>
                      ))}
                    </ul>
                  )}
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
