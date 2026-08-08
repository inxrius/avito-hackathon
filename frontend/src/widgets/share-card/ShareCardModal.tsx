import { useEffect, useMemo, useRef } from 'react';
import { createRandom } from '@/shared/lib/random';
import { DISTRICTS, pluralize } from '@/shared/lib/plural';
import type { ShareCard } from '@/shared/types/recap';
import './ShareCardModal.css';

interface Props {
  card: ShareCard;
  onClose: () => void;
}

const BARS = 26;

/**
 * Силуэт города строится из сида, а не из активности: высоты домов здесь
 * декоративные и ничего не сообщают о поведении. Тот же город — но без данных.
 */
function useSilhouette(seed: number): number[] {
  return useMemo(() => {
    const rand = createRandom(seed);
    return Array.from({ length: BARS }, (_, i) => {
      const wave = Math.sin((i / BARS) * Math.PI);
      return 0.25 + wave * 0.5 + rand() * 0.25;
    });
  }, [seed]);
}

export function ShareCardModal({ card, onClose }: Props) {
  const bars = useSilhouette(card.silhouetteSeed);
  const closeRef = useRef<HTMLButtonElement | null>(null);

  useEffect(() => {
    closeRef.current?.focus();
    const onKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose();
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [onClose]);

  return (
    <div className="share" role="dialog" aria-modal="true" aria-label="Публичная карточка">
      <div className="share__backdrop" onClick={onClose} />

      <div className="share__body">
        <div className="share__card">
          <p className="share__kicker">Итоги года · 2025</p>
          <h2 className="share__city">{card.cityName}</h2>

          <svg className="share__silhouette" viewBox="0 0 260 90" role="img" aria-hidden="true">
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

          <div className="share__traits">
            <span className="share__trait">{card.roleTitle}</span>
            <span className="share__trait share__trait--soft">{card.styleTitle}</span>
          </div>
          <p className="share__meta">{pluralize(card.districts, DISTRICTS)} застроено</p>
        </div>

        <div className="share__side">
          <h3 className="share__side-title">Что уходит в публичную карточку</h3>
          <ul className="share__list share__list--ok">
            <li>Название города и роль</li>
            <li>Стиль поведения</li>
            <li>Количество районов</li>
          </ul>

          <h3 className="share__side-title">Чего в ней нет</h3>
          <ul className="share__list share__list--no">
            <li>Категорий и счётчиков по ним</li>
            <li>Цен, сделок и объявлений</li>
            <li>Переписок и других пользователей</li>
          </ul>

          <p className="share__note">
            Приватный recap подробнее публичной карточки намеренно: по силуэту нельзя восстановить,
            что человек искал и покупал.
          </p>

          <button type="button" className="btn btn--primary share__copy">
            Скопировать ссылку
          </button>
          <button ref={closeRef} type="button" className="btn btn--ghost" onClick={onClose}>
            Закрыть
          </button>
        </div>
      </div>
    </div>
  );
}
