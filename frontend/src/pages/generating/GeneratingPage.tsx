import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { generateRecap } from '@/shared/api/recap';
import { describeFailure, type FailureView } from '@/shared/api/errors';
import { cacheRecap } from '@/shared/api/cache';
import './GeneratingPage.css';

/**
 * Шаги проговаривают, что именно система посчитала важным. Это требование ТЗ
 * («какие действия система посчитала важными») закрыто ещё до первого экрана итогов.
 */
const STEPS = [
  'Читаем действия за год',
  'Раскладываем их по районам',
  'Ищем повторяющиеся сценарии',
  'Назначаем роль и звания',
  'Строим город',
] as const;

const STEP_MS = 420;
/** Минимальный показ анимации, чтобы сборка не мигала, если бэкенд ответил мгновенно. */
const MIN_VISIBLE_MS = STEP_MS * STEPS.length;

export function GeneratingPage() {
  const { profileId = '', year = '' } = useParams();
  const navigate = useNavigate();
  const [step, setStep] = useState(0);
  const [failure, setFailure] = useState<FailureView | null>(null);

  useEffect(() => {
    let active = true;

    const ticker = setInterval(() => {
      setStep((current) => Math.min(current + 1, STEPS.length - 1));
    }, STEP_MS);

    // Запрос уходит сразу; ждём только оставшуюся часть минимального показа.
    const startedAt = Date.now();
    const holdRemaining = () =>
      new Promise<void>((resolve) =>
        setTimeout(resolve, Math.max(0, MIN_VISIBLE_MS - (Date.now() - startedAt))),
      );

    generateRecap(profileId, Number(year))
      .then(async (recap) => {
        await holdRemaining();
        if (!active) return;
        cacheRecap(recap);
        void navigate(`/recap/${recap.recapId}`, { replace: true });
      })
      .catch((cause: unknown) => {
        if (!active) return;
        // 422 — не сбой, а штатный ответ: за год слишком мало значимых действий.
        setFailure(describeFailure(cause));
      })
      .finally(() => clearInterval(ticker));

    return () => {
      active = false;
      clearInterval(ticker);
    };
  }, [navigate, profileId, year]);

  if (failure) {
    return (
      <main className="generating">
        <div className="generating__inner">
          <p className="kicker">Города пока нет</p>
          <h1 className="generating__title">{failure.title}</h1>
          {failure.hint && <p className="generating__hint">{failure.hint}</p>}
          <div className="generating__actions">
            {failure.retryable && (
              <button
                type="button"
                className="btn btn--primary"
                onClick={() => window.location.reload()}
              >
                Повторить
              </button>
            )}
            <button
              type="button"
              className={failure.retryable ? 'btn btn--ghost' : 'btn btn--primary'}
              onClick={() => void navigate('/')}
            >
              Выбрать другой профиль
            </button>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="generating">
      <div className="generating__inner">
        <p className="kicker">Собираем итоги года</p>
        <h1 className="generating__title">Город строится</h1>

        <ol className="generating__steps" aria-live="polite">
          {STEPS.map((label, index) => (
            <li
              key={label}
              className={
                'generating__step' +
                (index < step ? ' generating__step--done' : '') +
                (index === step ? ' generating__step--active' : '')
              }
            >
              <span className="generating__marker" aria-hidden="true" />
              {label}
            </li>
          ))}
        </ol>

        <div className="generating__bar">
          <div
            className="generating__fill"
            style={{ width: `${((step + 1) / STEPS.length) * 100}%` }}
          />
        </div>
      </div>
    </main>
  );
}
