import { useEffect, useRef, useState } from 'react';
import { loadShareCard } from '@/shared/api/recap';
import { describeFailure, type FailureView } from '@/shared/api/errors';
import { ShareCardPreview } from './ShareCardPreview';
import type { ShareCardDTO } from '@/shared/api/dto';
import './ShareCardModal.css';

interface Props {
  recapId: string;
  onClose: () => void;
}

type CopyState = 'idle' | 'copied' | 'failed';

export function ShareCardModal({ recapId, onClose }: Props) {
  const closeRef = useRef<HTMLButtonElement | null>(null);
  const [card, setCard] = useState<ShareCardDTO | null>(null);
  const [failure, setFailure] = useState<FailureView | null>(null);
  const [copyState, setCopyState] = useState<CopyState>('idle');

  // Ссылка ведёт на публичный просмотр, а не на приватные итоги.
  const publicUrl = `${window.location.origin}/share/${encodeURIComponent(recapId)}`;

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
      .catch((cause: unknown) => {
        if (active) setFailure(describeFailure(cause));
      });

    return () => {
      active = false;
    };
  }, [recapId]);

  const copyPublicLink = async () => {
    try {
      await navigator.clipboard.writeText(publicUrl);
      setCopyState('copied');
    } catch {
      // Буфер может быть недоступен без https или без разрешения —
      // тогда показываем ссылку, чтобы её можно было выделить руками.
      setCopyState('failed');
    }
  };

  return (
    <div className="share" role="dialog" aria-modal="true" aria-label="Поделиться итогами">
      <div className="share__backdrop" onClick={onClose} />

      <div className="share__body">
        <div className="share__preview">
          {card ? (
            <ShareCardPreview card={card} seed={recapId} />
          ) : (
            <div className="share__placeholder">
              {failure ? failure.title : 'Готовим карточку…'}
            </div>
          )}
        </div>

        <div className="share__side">
          <h2 className="share__heading">Поделиться итогами</h2>
          <p className="share__lead">
            Так карточку увидят другие. В ней только итоговые факты — без переписок и личных данных.
          </p>

          {copyState === 'failed' && (
            <label className="share__fallback">
              Скопируй ссылку вручную
              <input
                className="share__url"
                readOnly
                value={publicUrl}
                onFocus={(event) => event.currentTarget.select()}
              />
            </label>
          )}

          <div className="share__actions">
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => void copyPublicLink()}
              disabled={!card}
            >
              {copyState === 'copied' ? 'Ссылка скопирована ✓' : 'Скопировать ссылку'}
            </button>
            <button ref={closeRef} type="button" className="btn btn--ghost" onClick={onClose}>
              Закрыть
            </button>
          </div>

          <p className="share__note">Без переписок, цен и идентификаторов профиля.</p>
        </div>
      </div>
    </div>
  );
}
