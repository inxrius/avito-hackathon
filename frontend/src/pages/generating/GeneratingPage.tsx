import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { generateRecap, isInsufficientActivity } from '@/shared/api/recap';
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

type Failure = { kind: 'insufficient' } | { kind: 'error'; message: string };

export function GeneratingPage() {
  const { profileId = '', year = '' } = useParams();
  const navigate = useNavigate();
  const [step, setStep] = useState(0);
  const [failure, setFailure] = useState<Failure | null>(null);

  useEffect(() => {
    let active = true;

    const ticker = setInterval(() => {
      setStep((current) => Math.min(current + 1, STEPS.length - 1));
    }, STEP_MS);

    generateRecap(profileId, Number(year))
      .then((recap) => {
        if (!active) return;
        cacheRecap(recap);
        void navigate(`/recap/${recap.recapId}`, { replace: true });
      })
      .catch((cause: unknown) => {
        if (!active) return;
        // 422 — не сбой, а штатный ответ: за год слишком мало значимых действий.
        if (isInsufficientActivity(cause)) {
          setFailure({ kind: 'insufficient' });
          return;
        }
        setFailure({
          kind: 'error',
          message: cause instanceof Error ? cause.message : 'Не удалось собрать итоги',
        });
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
          {failure.kind === 'insufficient' ? (
            <>
              <p className="kicker">Города пока нет</p>
              <h1 className="generating__title">За {year} год слишком мало действий</h1>
              <p className="generating__hint">
                Чтобы собрать итоги, нужно хотя бы несколько недель активности: просмотры,
                избранное, диалоги. Пока их не хватает даже на один квартал.
              </p>
            </>
          ) : (
            <>
              <p className="kicker">Ошибка</p>
              <h1 className="generating__title">{failure.message}</h1>
            </>
          )}
          <button type="button" className="btn btn--primary" onClick={() => void navigate('/')}>
            Выбрать другой профиль
          </button>
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
