import { useEffect, useState, type CSSProperties } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchProfiles } from '@/shared/api/recap';
import { toneColor } from '@/shared/lib/palette';
import type { Profile } from '@/shared/types/recap';
import './ProfilesPage.css';

export function ProfilesPage() {
  const navigate = useNavigate();
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [selected, setSelected] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    void fetchProfiles().then((list) => {
      if (active) setProfiles(list);
    });
    return () => {
      active = false;
    };
  }, []);

  return (
    <main className="profiles">
      <div className="profiles__inner">
        <p className="kicker">Итоги года </p>
        <h1 className="profiles__title">
          Год на Авито -
          <br />
          это целый город
        </h1>
        <p className="profiles__lead">
          Каждая категория, в которой ты что-то делал, становится районом. Чем больше действий — тем
          выше дома. Выбери профиль и посмотри, что построилось за год.
        </p>

        <ul className="profiles__list">
          {profiles.map((profile) => (
            <li key={profile.id}>
              <button
                type="button"
                className={`profile-card${selected === profile.id ? ' profile-card--active' : ''}`}
                style={{ '--tone': toneColor(profile.tone, 0) } as CSSProperties}
                onClick={() => setSelected(profile.id)}
                aria-pressed={selected === profile.id}
              >
                <span className="profile-card__avatar">{profile.name.slice(0, 1)}</span>
                <span className="profile-card__body">
                  <span className="profile-card__name">{profile.name}</span>
                  <span className="profile-card__tagline">{profile.tagline}</span>
                  <span className="profile-card__hint">{profile.hint}</span>
                </span>
              </button>
            </li>
          ))}
        </ul>

        <div className="profiles__footer">
          <button
            type="button"
            className="btn btn--primary profiles__cta"
            disabled={!selected}
            onClick={() => {
              if (selected) void navigate(`/generate/${selected}`);
            }}
          >
            Построить город →
          </button>
          <p className="profiles__note">
            Профили тестовые. Данные активности сгенерированы заранее — переписок и реальных
            пользователей здесь нет.
          </p>
        </div>
      </div>
    </main>
  );
}
