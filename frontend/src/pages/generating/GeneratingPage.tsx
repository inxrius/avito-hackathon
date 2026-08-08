import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { generateRecap } from '@/shared/api/recap';
import './GeneratingPage.css';

/**
 * Шаги проговаривают, что именно система посчитала важным. Это требование ТЗ
 * («какие действия система посчитала важными») закрыто ещё до первого экрана итогов.
 */
const STEPS = [
  'Читаем действия за год',
  'Раскладываем их по категориям',
  'Ищем повторяющиеся сценарии',
  'Назначаем роль и ачивки',
  'Строим город',
] as const;

const STEP_MS = 520;

export function GeneratingPage() {
  const { profileId = '' } = useParams();
  const navigate = useNavigate();
  const [step, setStep] = useState(0);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const controller = new AbortController();

    const ticker = setInterval(() => {
      setStep((current) => Math.min(current + 1, STEPS.length - 1));
    }, STEP_MS);

    void generateRecap(profileId, controller.signal)
      .then(() => {
        void navigate(`/recap/${profileId}`, { replace: true });
      })
      .catch((cause: unknown) => {
        if (cause instanceof DOMException && cause.name === 'AbortError') return;
        setError(cause instanceof Error ? cause.message : 'Не удалось собрать итоги');
      });

    return () => {
      clearInterval(ticker);
      controller.abort();
    };
  }, [navigate, profileId]);

  if (error) {
    return (
      <main className="generating">
        <div className="generating__inner">
          <p className="kicker">Ошибка</p>
          <h1 className="generating__title">{error}</h1>
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
