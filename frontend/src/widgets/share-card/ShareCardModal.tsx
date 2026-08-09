import { useEffect, useMemo, useRef, useState } from 'react';
import { loadShareCard } from '@/shared/api/recap';
import { createRandom } from '@/shared/lib/random';
import type { ShareCardDTO } from '@/shared/api/dto';
import './ShareCardModal.css';

interface Props {
  recapId: string;
  onClose: () => void;
}

const BARS = 26;

/**
 * Силуэт города строится из id recap, а не из активности: высоты домов здесь
 * декоративные и ничего не сообщают о поведении. Тот же город — но без данных.
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

export function ShareCardModal({ recapId, onClose }: Props) {
  const bars = useSilhouette(recapId);
  const closeRef = useRef<HTMLButtonElement | null>(null);
  const [card, setCard] = useState<ShareCardDTO | null>(null);
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    closeRef.current?.focus();
    const onKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose();
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [onClose]);

  useEffect(() => {
    let active = true;
    loadShareCard(recapId)
      .then((value) => {
        if (active) setCard(value);
      })
      .catch(() => {
        if (active) setFailed(true);
      });
    return () => {
      active = false;
    };
  }, [recapId]);

  return (
    <div className="share" role="dialog" aria-modal="true" aria-label="Публичная карточка">
      <div className="share__backdrop" onClick={onClose} />

      <div className="share__body">
        <div className="share__card">
          {card ? (
            <>
              <p className="share__kicker">Итоги года · {card.year}</p>
              <h2 className="share__city">{card.title}</h2>
            </>
          ) : (
            <p className="share__kicker">{failed ? 'Карточка недоступна' : 'Готовим карточку…'}</p>
          )}

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

          {card && (
            <>
              <div className="share__traits">
                <span className="share__trait">{card.main_district.title}</span>
                <span className="share__trait share__trait--soft">{card.subtitle}</span>
              </div>
              <ul className="share__facts">
                {card.facts.map((fact) => (
                  <li key={fact.kind}>
                    <b>{fact.value}</b> · {fact.label}
                  </li>
                ))}
              </ul>
            </>
          )}
        </div>

        <div className="share__side">
          <h3 className="share__side-title">Что уходит в публичную карточку</h3>
          <ul className="share__list share__list--ok">
            <li>Главный район и звания</li>
            <li>Роль и стиль</li>
            <li>Дни активности</li>
          </ul>

          <h3 className="share__side-title">Чего в ней нет</h3>
          <ul className="share__list share__list--no">
            <li>Идентификаторов пользователя</li>
            <li>Цен, сделок и объявлений</li>
            <li>Переписок и других пользователей</li>
          </ul>

          <p className="share__note">
            Карточку собирает бэкенд отдельным запросом: во фронт не попадает ничего, чего
            в ней быть не должно.
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
