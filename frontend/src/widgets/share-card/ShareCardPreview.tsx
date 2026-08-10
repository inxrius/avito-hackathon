import { useMemo } from 'react';
import { createRandom } from '@/shared/lib/random';
import type { ShareCardDTO } from '@/shared/api/dto';
import './ShareCardPreview.css';

interface Props {
  card: ShareCardDTO;
  /** Строка, из которой выводится силуэт. Обычно recapId. */
  seed: string;
}

const BARS = 26;

/**
 * Силуэт города строится из сида, а не из активности: высоты домов здесь
 * декоративные и ничего не сообщают о поведении.
 */
function useSilhouette(seed: string): number[] {
  return useMemo(() => {
    let numeric = 0;
    for (const char of seed) numeric = (numeric * 31 + char.charCodeAt(0)) >>> 0;
    const rand = createRandom(numeric);
    return Array.from({ length: BARS }, (_, i) => {
      const wave = Math.sin((i / BARS) * Math.PI);
      return 0.25 + wave * 0.5 + rand() * 0.25;
    });
  }, [seed]);
}

/**
 * Публичная карточка. Один и тот же компонент показывается в модалке автору
 * и на странице по внешней ссылке — так автор видит ровно то, что увидят другие.
 *
 * Источник данных только `ShareCardDTO`: приватную модель recap сюда не передаём.
 */
export function ShareCardPreview({ card, seed }: Props) {
  const bars = useSilhouette(seed);

  return (
    <div className="share-card">
      <p className="share-card__kicker">Итоги года · {card.year}</p>
      <h2 className="share-card__title">{card.title}</h2>

      <svg className="share-card__silhouette" viewBox="0 0 260 90" role="img" aria-hidden="true">
        {bars.map((height, index) => {
          const w = 260 / BARS;
          const h = height * 78;
          return (
            <rect
              key={index}
              x={index * w + 1}
              y={90 - h}
              width={w - 2}
              height={h}
              rx={1.5}
              fill="#05221a"
              opacity={0.35 + height * 0.5}
            />
          );
        })}
      </svg>

      <div className="share-card__traits">
        <span className="share-card__trait">{card.main_district.title}</span>
        <span className="share-card__trait share-card__trait--soft">{card.subtitle}</span>
      </div>

      <ul className="share-card__facts">
        {card.facts.map((fact) => (
          <li key={fact.kind}>
            <b>{fact.value}</b> · {fact.label}
          </li>
        ))}
      </ul>

      {card.achievements.length > 0 && (
        <ul className="share-card__awards">
          {card.achievements.map((achievement) => (
            <li key={achievement.code}>{achievement.title}</li>
          ))}
        </ul>
      )}
    </div>
  );
}
