import { useState } from 'react';
import type { Trait } from '@/shared/types/recap';

interface Props {
  role: Trait;
  style: Trait;
}

/**
 * Роль и стиль — главный вывод о пользователе. Обоснование спрятано под кнопку:
 * не мешает читать итог, но всегда доступно и отвечает на «почему именно так».
 */
export function TraitsPanel({ role, style }: Props) {
  const [open, setOpen] = useState(false);

  return (
    <section className="panel recap-panel">
      <h2 className="kicker">Роль и стиль</h2>

      <div className="trait">
        <p className="trait__label">Роль</p>
        <p className="trait__value trait__value--accent">{role.title}</p>
      </div>

      <div className="trait">
        <p className="trait__label">Стиль</p>
        <p className="trait__value">{style.title}</p>
      </div>

      <button
        type="button"
        className="why-toggle"
        onClick={() => setOpen((value) => !value)}
        aria-expanded={open}
      >
        {open ? 'Скрыть' : 'Почему я такой?'}
      </button>

      {open && (
        <div className="why">
          <p className="why__line">{role.reason}</p>
          <p className="why__line">{style.reason}</p>
        </div>
      )}
    </section>
  );
}
