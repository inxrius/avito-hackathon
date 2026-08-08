import type { BrandTone, DistrictId } from '@/shared/types/recap';

interface CategoryMeta {
  title: string;
  tone: BrandTone;
  shade: 0 | 1 | 2;
}

/**
 * Единая таблица категорий. Девять районов разложены по четырём брендовым
 * цветам Авито — ни одного оттенка «на глаз» вне палитры.
 * Порядок ключей фиксирован: он же задаёт порядок легенды и индексы пикинга.
 */
export const CATEGORY_META: Record<DistrictId, CategoryMeta> = {
  realty: { title: 'Недвижимость', tone: 'blue', shade: 0 },
  home: { title: 'Дом и дача', tone: 'green', shade: 1 },
  electronics: { title: 'Электроника', tone: 'purple', shade: 0 },
  work: { title: 'Работа', tone: 'green', shade: 0 },
  personal: { title: 'Личные вещи', tone: 'red', shade: 1 },
  hobby: { title: 'Хобби и отдых', tone: 'red', shade: 0 },
  auto: { title: 'Авто', tone: 'blue', shade: 2 },
  services: { title: 'Услуги', tone: 'purple', shade: 1 },
  pets: { title: 'Животные', tone: 'green', shade: 2 },
};

export const DISTRICT_ORDER = Object.keys(CATEGORY_META) as DistrictId[];
