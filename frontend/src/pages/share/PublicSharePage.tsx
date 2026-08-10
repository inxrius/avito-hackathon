import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { loadShareCard } from '@/shared/api/recap';
import { describeFailure, type FailureView } from '@/shared/api/errors';
import { ShareCardPreview } from '@/widgets/share-card/ShareCardPreview';
import type { ShareCardDTO } from '@/shared/api/dto';
import './PublicSharePage.css';

/**
 * Страница по внешней ссылке. Единственный источник данных —
 * `GET /recaps/{id}/share`: приватный recap и обоснования здесь не запрашиваются,
 * поэтому по ссылке физически невозможно вытащить лишнее.
 */
export function PublicSharePage() {
  const { recapId = '' } = useParams();
  const navigate = useNavigate();
  const [card, setCard] = useState<ShareCardDTO | null>(null);
  const [failure, setFailure] = useState<FailureView | null>(null);

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

  if (failure) {
    return (
      <main className="public-share">
        <div className="public-share__inner">
          <p className="kicker">Карточка недоступна</p>
          <h1 className="public-share__error-title">{failure.title}</h1>
          {failure.hint && <p className="public-share__error-hint">{failure.hint}</p>}
          <button type="button" className="btn btn--primary" onClick={() => void navigate('/')}>
            Собрать свои итоги
          </button>
        </div>
      </main>
    );
  }

  return (
    <main className="public-share">
      <div className="public-share__inner">
        {card ? (
          <>
            <ShareCardPreview card={card} seed={recapId} />
            <p className="public-share__note">Без переписок, цен и идентификаторов профиля.</p>
            <button type="button" className="btn btn--primary" onClick={() => void navigate('/')}>
              Собрать свои итоги
            </button>
          </>
        ) : (
          <p className="kicker">Загружаем карточку…</p>
        )}
      </div>
    </main>
  );
}
