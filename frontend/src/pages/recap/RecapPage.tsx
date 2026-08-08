import { useCallback, useEffect, useMemo, useState } from 'react';
import { Navigate, useNavigate, useParams } from 'react-router-dom';
import { CityCanvas } from '@/entities/city/CityCanvas';
import { getCachedRecap } from '@/shared/api/recap';
import { ShareCardModal } from '@/widgets/share-card/ShareCardModal';
import { ACTIONS, DISTRICTS, SITES, pluralize } from '@/shared/lib/plural';
import { BadgesPanel } from './components/BadgesPanel';
import { DistrictLegend } from './components/DistrictLegend';
import { TraitsPanel } from './components/TraitsPanel';
import { UnfinishedPanel } from './components/UnfinishedPanel';
import type { DistrictId } from '@/shared/types/recap';
import './RecapPage.css';

export function RecapPage() {
  const { profileId = '' } = useParams();
  const navigate = useNavigate();
  const recap = getCachedRecap(profileId);

  const [step, setStep] = useState(0);
  const [focus, setFocus] = useState<DistrictId | null>(null);
  const [shareOpen, setShareOpen] = useState(false);

  const total = recap?.chapters.length ?? 0;
  const isFinal = total > 0 && step === total - 1;

  const next = useCallback(() => setStep((v) => Math.min(v + 1, total - 1)), [total]);
  const back = useCallback(() => setStep((v) => Math.max(v - 1, 0)), []);

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      if (event.key === 'ArrowRight' || event.key === ' ') {
        event.preventDefault();
        next();
      }
      if (event.key === 'ArrowLeft') back();
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [back, next]);

  /**
   * Город растёт по главам: каждая достраивает свой район, финальная —
   * раскрывает всё остальное разом. Это и есть «город собран».
   */
  const revealed = useMemo(() => {
    if (!recap) return new Set<DistrictId>();
    if (isFinal) return new Set(recap.districts.map((d) => d.id));
    const ids = recap.chapters
      .slice(0, step + 1)
      .map((chapter) => chapter.districtId)
      .filter((id): id is DistrictId => Boolean(id));
    return new Set(ids);
  }, [isFinal, recap, step]);

  // Прямая ссылка на итоги без сборки: отправляем на генерацию, а не показываем пустой экран.
  if (!recap) return <Navigate to={`/generate/${profileId}`} replace />;

  const chapter = recap.chapters[step];
  const chapterBadge = chapter.badgeId
    ? recap.badges.find((badge) => badge.id === chapter.badgeId)
    : undefined;

  return (
    <main className="recap">
      <header className="recap__header">
        <button type="button" className="recap__logo" onClick={() => void navigate('/')}>
          Итоги года · {recap.year}
        </button>

        <ol className="recap__progress" aria-label={`Глава ${step + 1} из ${total}`}>
          {recap.chapters.map((item, index) => (
            <li
              key={item.index}
              className={'recap__tick' + (index <= step ? ' recap__tick--done' : '')}
            />
          ))}
        </ol>

        <div className="recap__header-actions">
          {isFinal && (
            <button type="button" className="btn btn--ghost" onClick={() => setShareOpen(true)}>
              Поделиться городом
            </button>
          )}
        </div>
      </header>

      <section className="recap__stage">
        {isFinal ? (
          <div className="recap__intro">
            <p className="kicker">
              Глава {step + 1} из {total} · Город собран
            </p>
            <h1 className="recap__city-name">{recap.cityName}</h1>
            <p className="recap__totals">
              {pluralize(recap.totals.actions, ACTIONS)} ·{' '}
              {pluralize(recap.totals.districts, DISTRICTS)} ·{' '}
              {pluralize(recap.totals.sites, SITES)} на окраине
            </p>
          </div>
        ) : (
          <div className="recap__intro">
            <p className="kicker">
              Глава {step + 1} из {total}
            </p>
            <h1 className="recap__chapter-title">{chapter.title}</h1>
            <p className="recap__stat">
              <span className="recap__stat-value">{chapter.stat.value}</span>
              <span className="recap__stat-label">{chapter.stat.label}</span>
            </p>
            <p className="recap__narrative">{chapter.narrative}</p>
            {chapterBadge && (
              <p className="recap__badge-toast">
                <span className="recap__badge-mark" aria-hidden="true" />
                Новая ачивка: <b>{chapterBadge.title}</b>
              </p>
            )}
          </div>
        )}

        <CityCanvas
          recap={recap}
          revealed={revealed}
          focus={isFinal ? focus : (chapter.districtId ?? null)}
          showSites={isFinal}
          interactive={isFinal}
          onSelectDistrict={(id) => setFocus((current) => (current === id ? null : id))}
        />

        {!isFinal && (
          <div className="recap__controls">
            <button type="button" className="btn btn--ghost" onClick={back} disabled={step === 0}>
              Назад
            </button>
            <button type="button" className="btn btn--primary" onClick={next}>
              Дальше →
            </button>
          </div>
        )}
      </section>

      {isFinal && (
        <section className="recap__final">
          <DistrictLegend
            districts={recap.districts}
            revealed={revealed}
            active={focus}
            onSelect={(id) => setFocus((current) => (current === id ? null : id))}
          />

          <div className="recap__panels">
            <TraitsPanel role={recap.role} style={recap.style} />
            <BadgesPanel badges={recap.badges} />
            <UnfinishedPanel items={recap.unfinished} />
          </div>

          <footer className="recap__footer">
            <p className="recap__privacy">
              В городе нет переписок, цен и данных других людей — только форма твоего года.
            </p>
            <p className="recap__version">
              Правила генерации: {recap.rulesVersion} · сид {recap.seed}. Тот же профиль всегда даёт
              тот же город.
            </p>
          </footer>
        </section>
      )}

      {shareOpen && <ShareCardModal card={recap.shareCard} onClose={() => setShareOpen(false)} />}
    </main>
  );
}
