import { useEffect, useState, type CSSProperties } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchProfiles } from '@/shared/api/recap';
import { describeFailure, type FailureView } from '@/shared/api/errors';
import { toneColor } from '@/shared/lib/palette';
import type { Profile } from '@/shared/types/recap';
import './ProfilesPage.css';

export function ProfilesPage() {
  const navigate = useNavigate();
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [state, setState] = useState<'loading' | 'ready' | 'failed'>('loading');
  const [failure, setFailure] = useState<FailureView | null>(null);
  const [selected, setSelected] = useState<Profile | null>(null);

  useEffect(() => {
    let active = true;

    fetchProfiles()
      .then((list) => {
        if (!active) return;
        setProfiles(list);
        setState('ready');
      })
      .catch((cause: unknown) => {
        // Пустой список и молча выключенная кнопка — худший из вариантов:
        // пользователь не понимает, почему ничего не происходит.
        if (!active) return;
        setFailure(describeFailure(cause));
        setState('failed');
      });

    return () => {
      active = false;
    };
  }, []);

  // Год берём из профиля: хардкода нет, а без года генерировать нечего.
  const year = selected?.availableYears[0];

  const start = () => {
    if (!selected || year === undefined) return;
    void navigate(`/generate/${selected.id}/${year}`);
  };

  return (
    <main className="profiles">
      <div className="profiles__inner">
        <p className="kicker">Итоги года</p>
        <h1 className="profiles__title">
          Год на Авито —
          <br />
          это целый город
        </h1>
        <p className="profiles__lead">
          Каждая вертикаль, в которой ты что-то делал, становится районом. Чем больше действий —
          тем выше дома. Выбери профиль и посмотри, что построилось за год.
        </p>

        {state === 'loading' && <p className="profiles__status">Загружаем профили…</p>}

        {state === 'failed' && (
          <div className="profiles__status profiles__status--error">
            <p>
              <b>{failure?.title ?? 'Не удалось загрузить профили'}</b>
              {failure?.hint ? ` ${failure.hint}` : ''}
            </p>
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => window.location.reload()}
            >
              Повторить
            </button>
          </div>
        )}

        {state === 'ready' && (
          <ul className="profiles__list">
            {profiles.map((profile) => (
              <li key={profile.id}>
                <button
                  type="button"
                  className={`profile-card${selected?.id === profile.id ? ' profile-card--active' : ''}`}
                  style={{ '--tone': toneColor(profile.tone, 0) } as CSSProperties}
                  onClick={() => setSelected(profile)}
                  aria-pressed={selected?.id === profile.id}
                >
                  <span className="profile-card__avatar">{profile.name.slice(0, 1)}</span>
                  <span className="profile-card__body">
                    <span className="profile-card__name">{profile.name}</span>
                    <span className="profile-card__tagline">{profile.description}</span>
                    <span className="profile-card__hint">
                      {profile.availableYears.length > 0
                        ? `Доступный год: ${profile.availableYears.join(', ')}`
                        : 'Нет годов с данными'}
                    </span>
                  </span>
                </button>
              </li>
            ))}
          </ul>
        )}

        <div className="profiles__footer">
          <button
            type="button"
            className="btn btn--primary profiles__cta"
            disabled={!selected || year === undefined}
            onClick={start}
          >
            Построить город →
          </button>
          <p className="profiles__note">
            Профили тестовые. Данные активности подготовлены заранее — переписок и реальных
            пользователей здесь нет.
          </p>
        </div>
      </div>
    </main>
  );
}
