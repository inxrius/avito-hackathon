import { DISTRICTS, pluralize } from '@/shared/lib/plural';
import type { Unfinished } from '@/shared/types/recap';

interface Props {
  items: Unfinished[];
}

/**
 * Стройки на окраине. Это одновременно «какие возможности пользователь упустил»
 * и точка возврата в продукт: у каждой строки есть конкретное следующее действие.
 */
export function UnfinishedPanel({ items }: Props) {
  return (
    <section className="panel recap-panel">
      <h2 className="kicker">Что не достроено</h2>

      <ul className="unfinished">
        {items.map((item) => (
          <li key={item.id} className="unfinished__item">
            <span className="unfinished__icon" aria-hidden="true" />
            <span className="unfinished__text">
              <b>{item.count}</b> {item.title}
            </span>
            <span className="unfinished__cta">{item.ctaLabel}</span>
          </li>
        ))}
      </ul>

      <button type="button" className="btn btn--primary unfinished__button">
        Достроить {pluralize(items.length, DISTRICTS)} →
      </button>
    </section>
  );
}
