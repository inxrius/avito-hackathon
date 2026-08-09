import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { CityCanvas } from '@/entities/city/CityCanvas';
import { cacheRecap, getCachedRecap } from '@/shared/api/cache';
import { loadExplanation, loadRecap } from '@/shared/api/recap';
import { ShareCardModal } from '@/widgets/share-card/ShareCardModal';
import { DISTRICTS, pluralize } from '@/shared/lib/plural';
import { BadgesPanel } from './components/BadgesPanel';
import { DistrictLegend } from './components/DistrictLegend';
import { TraitsPanel } from './components/TraitsPanel';
import { UnfinishedPanel } from './components/UnfinishedPanel';
import type { DistrictId, Recap } from '@/shared/types/recap';
import './RecapPage.css';

export function RecapPage() {
  const { recapId = '' } = useParams();
  const navigate = useNavigate();

  const [recap, setRecap] = useState<Recap | null>(() => getCachedRecap(recapId) ?? null);
  const [failed, setFailed] = useState(false);
  const [step, setStep] = useState(0);
  const [focus, setFocus] = useState<DistrictId | null>(null);
  const [shareOpen, setShareOpen] = useState(false);
  // Ref, а не state: это защёлка «уже запросили», перерисовывать от неё нечего.
  const explained = useRef(false);

  // Прямая ссылка на итоги: snapshot неизменяем, поэтому просто перечитываем его.
  useEffect(() => {
    if (recap) return;
    let active = true;

    loadRecap(recapId)
      .then((value) => {
        if (!active) return;
        cacheRecap(value);
        setRecap(value);
      })
      .catch(() => {
        if (active) setFailed(true);
      });

    return () => {
      active = false;
    };
  }, [recap, recapId]);

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

  // Обоснования — отдельный запрос, поэтому тянем их только когда дошли до финала.
  useEffect(() => {
    if (!recap || !isFinal || explained.current || !recap.capabilities.explanationAvailable) return;
    let active = true;
    explained.current = true;

    loadExplanation(recap)
      .then((value) => {
        if (!active) return;
        cacheRecap(value);
        setRecap(value);
      })
      .catch(() => undefined);

    return () => {
      active = false;
    };
  }, [isFinal, recap]);

  const revealed = useMemo(() => {
    if (!recap) return new Set<DistrictId>();
    if (isFinal) return new Set(recap.districts.map((d) => d.id));
    const ids = recap.chapters
      .slice(0, step + 1)
      .map((chapter) => chapter.districtId)
      .filter((id): id is DistrictId => Boolean(id));
    return new Set(ids);
  }, [isFinal, recap, step]);

  if (failed) {
    return (
      <main className="recap recap--message">
        <div>
          <p className="kicker">Итоги не найдены</p>
          <h1 className="recap__chapter-title">Такого города нет</h1>
          <button type="button" className="btn btn--primary" onClick={() => void navigate('/')}>
            К выбору профиля
          </button>
        </div>
      </main>
    );
  }

  if (!recap) {
    return (
      <main className="recap recap--message">
        <p className="kicker">Загружаем итоги…</p>
      </main>
    );
  }

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
          {recap.chapters.map((item) => (
            <li
              key={item.index}
              className={'recap__tick' + (item.index <= step + 1 ? ' recap__tick--done' : '')}
            />
          ))}
        </ol>

        <div className="recap__header-actions">
          {/* Кнопку прячем по capabilities, а не ловим 409 после клика. */}
          {isFinal && recap.capabilities.shareAvailable && (
            <button type="button" className="btn btn--ghost" onClick={() => setShareOpen(true)}>
              Поделиться городом
            </button>
          )}
        </div>
      </header>

      <section className="recap__stage">
        <div className="recap__intro">
          <p className="kicker">
            Глава {step + 1} из {total}
            {isFinal ? ' · Город собран' : ''}
          </p>

          {isFinal ? (
            <>
              <h1 className="recap__city-name">{recap.cityName}</h1>
              <p className="recap__totals">
                {recap.totals.activeDays !== undefined && `${recap.totals.activeDays} активных дней · `}
                {pluralize(recap.totals.districts, DISTRICTS)}
              </p>
            </>
          ) : (
            <>
              <h1 className="recap__chapter-title">{chapter.title}</h1>
              {chapter.stat && (
                <p className="recap__stat">
                  <span className="recap__stat-value">{chapter.stat.value}</span>
                  <span className="recap__stat-label">{chapter.stat.label}</span>
                </p>
              )}
              {chapter.narrative && <p className="recap__narrative">{chapter.narrative}</p>}
              {chapterBadge && (
                <p className="recap__badge-toast">
                  <span className="recap__badge-mark" aria-hidden="true" />
                  Новое звание: <b>{chapterBadge.title}</b>
                </p>
              )}
            </>
          )}
        </div>

        <CityCanvas
          recap={recap}
          revealed={revealed}
          focus={isFinal ? focus : (chapter.districtId ?? null)}
          showSites={false}
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
            {/* Панель появится, когда бэкенд начнёт отдавать незавершённые сценарии. */}
            {recap.unfinished.length > 0 && <UnfinishedPanel items={recap.unfinished} />}
          </div>

          <footer className="recap__footer">
            <p className="recap__privacy">
              В городе нет переписок, цен и данных других людей — только форма твоего года.
            </p>
            <p className="recap__version">
              Правила: {recap.rulesVersion} · текст:{' '}
              {recap.narrativeSource === 'mistral' ? 'сгенерирован ИИ' : 'шаблон'}. Тот же профиль
              всегда даёт тот же город.
            </p>
          </footer>
        </section>
      )}

      {shareOpen && <ShareCardModal recapId={recap.recapId} onClose={() => setShareOpen(false)} />}
    </main>
  );
}
