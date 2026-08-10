import * as assert from 'node:assert/strict';
import { test } from 'node:test';

import { adaptProfiles, adaptRecap, applyExplanation } from './adapter.ts';
import type { ProfileDTO, RecapDTO, RecapExplanationDTO } from './dto.ts';

function makeRecap(): RecapDTO {
  return {
    schema_version: '2.0',
    id: 'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa',
    profile_id: '11111111-1111-1111-1111-111111111111',
    year: 2026,
    profile: {
      id: '11111111-1111-1111-1111-111111111111',
      name: 'Марина',
      avatar_url: 'https://cdn.example/marina.png',
    },
    generation: {
      algorithm_version: 'recap-rules-2026.08.2-spike',
      feature_schema_version: 'features-v1',
      activity_hash: 'sha256:01020304ffffffffffffffffffffffffffffffffffffffffffffffffffffffff',
      generated_at: '2026-08-08T03:36:46Z',
      narrative: {
        source: 'template',
        prompt_version: 'city-summary-v2-spike',
      },
    },
    theme: {
      code: 'city',
      main_district: { code: 'goods', title: 'Товары' },
      accent_token: 'violet',
    },
    cards: [
      {
        id: 'intro',
        type: 'intro',
        position: 0,
        visibility: 'shareable',
        title: 'Твой город за год готов',
        description: 'Посмотрим, каким получился год',
        explainable: false,
        data: { year: 2026 },
      },
      {
        id: 'active-days',
        type: 'metric',
        position: 1,
        visibility: 'shareable',
        title: '18 дней в ритме города',
        explainable: false,
        data: {
          metric_code: 'active_days',
          value: 18,
          unit: 'days',
          secondary_label: 'Активность была заметна в 1 месяце',
        },
      },
      {
        id: 'sales',
        type: 'metric',
        position: 2,
        visibility: 'shareable',
        title: '6 успешных продаж',
        description: 'Твоя витрина приносила завершённые сделки',
        explainable: false,
        data: { metric_code: 'sales_count', value: 6, unit: 'items' },
      },
      {
        id: 'main-district',
        type: 'district',
        position: 3,
        visibility: 'shareable',
        title: 'Твой главный район - Товары',
        description: 'Чаще всего маршрут проходил по улице «Электроника»',
        explainable: false,
        data: {
          activity_share: 1,
          vertical: { code: 'goods', title: 'Товары' },
          top_category: {
            code: 'electronics',
            title: 'Электроника',
            vertical_code: 'goods',
          },
        },
      },
      {
        id: 'archetype',
        type: 'archetype',
        position: 4,
        visibility: 'shareable',
        eyebrow: 'Твоя роль в городе',
        title: 'Хозяин витрины',
        description: 'Твой стиль - Результативный',
        explainable: true,
        data: {
          role: { code: 'showcase_owner', title: 'Хозяин витрины' },
          style: { code: 'result_oriented', title: 'Результативный' },
        },
      },
      {
        id: 'achievements',
        type: 'achievements',
        position: 5,
        visibility: 'shareable',
        title: 'Твои городские звания',
        description: 'Достижения, которые лучше всего описывают твой год',
        explainable: true,
        data: {
          total_count: 2,
          items: [
            {
              code: 'city_rhythm',
              title: 'В ритме города',
              description: 'Ты сохранял серию активности несколько дней подряд.',
              level: 'expert',
              icon: 'city-rhythm-expert',
              metric_code: 'max_activity_streak',
              current_value: 18,
              next_level_threshold: 30,
            },
            {
              code: 'own_showcase',
              title: 'Своя витрина',
              description: 'Ты активно пополнял собственную витрину объявлениями.',
              level: 'expert',
              icon: 'own-showcase-expert',
              metric_code: 'published_listings_count',
              current_value: 8,
              next_level_threshold: 15,
            },
          ],
        },
      },
      {
        id: 'summary',
        type: 'summary',
        position: 6,
        visibility: 'personal',
        title: 'Твой город за год',
        description: 'Главным районом стали «Товары». Год получился результативным.',
        explainable: false,
        data: {
          role_code: 'showcase_owner',
          style_code: 'result_oriented',
          achievement_codes: ['city_rhythm', 'own_showcase'],
        },
      },
      {
        id: 'final',
        type: 'final',
        position: 7,
        visibility: 'shareable',
        title: 'До встречи в городе',
        description: 'Сохрани итоги или поделись городскими званиями',
        explainable: false,
        data: { show_feedback: true, show_share_button: true },
      },
    ],
    capabilities: {
      share_available: true,
      explanation_available: true,
      feedback_available: true,
    },
  };
}

void test('adaptProfiles keeps backend identifiers and supplies presentation tones', () => {
  const profiles: ProfileDTO[] = [
    {
      id: '11111111-1111-1111-1111-111111111111',
      name: 'Марина',
      description: 'Первый профиль',
      available_years: [2026],
    },
    {
      id: '22222222-2222-2222-2222-222222222222',
      name: 'Алексей',
      description: 'Второй профиль',
      available_years: [2025, 2026],
    },
  ];

  const result = adaptProfiles(profiles);
  assert.equal(result[0]?.id, profiles[0]?.id);
  assert.deepEqual(result[0]?.availableYears, [2026]);
  assert.equal(result[0]?.tone, 'blue');
  assert.equal(result[1]?.tone, 'purple');
});

void test('adaptRecap maps the backend contract without losing summary or achievements', () => {
  const result = adaptRecap(makeRecap());

  assert.equal(result.recapId, 'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa');
  assert.equal(result.profileId, '11111111-1111-1111-1111-111111111111');
  assert.equal(result.year, 2026);
  assert.equal(result.summaryText, 'Главным районом стали «Товары». Год получился результативным.');
  assert.equal(result.role.title, 'Хозяин витрины');
  assert.equal(result.style.title, 'Результативный');
  assert.equal(result.totals.activeDays, 18);
  assert.equal(result.districts[0]?.id, 'goods');
  assert.equal(result.districts[0]?.topCategory, 'Электроника');
  assert.equal(result.badges.length, 2);
  assert.equal(result.badges[0]?.title, 'В ритме города');
  assert.equal(result.badges[1]?.title, 'Своя витрина');
  assert.equal(result.chapters.length, 8);
  assert.equal(result.chapters[3]?.districtId, 'goods');
  assert.doesNotMatch(result.chapters[4]?.title ?? '', /\[object Object\]/u);
});

void test('adaptRecap creates a deterministic seed from the activity hash', () => {
  const first = adaptRecap(makeRecap());
  const second = adaptRecap(makeRecap());

  assert.equal(first.seed, second.seed);
  assert.equal(first.seed, 0x01020304);
});

void test('applyExplanation enriches role, style and matching achievements', () => {
  const recap = adaptRecap(makeRecap());
  const explanation: RecapExplanationDTO = {
    recap_id: recap.recapId,
    algorithm_version: recap.rulesVersion,
    activity_hash: 'sha256:01020304ffffffffffffffffffffffffffffffffffffffffffffffffffffffff',
    decisions: [
      {
        card_id: 'archetype',
        kind: 'archetype_role',
        code: 'showcase_owner',
        reason: 'Продавцовская активность была заметнее всего.',
        rule_version: 'archetype-rules-v2-spike',
        facts: [],
      },
      {
        card_id: 'archetype',
        kind: 'archetype_style',
        code: 'result_oriented',
        reason: 'Активность часто приводила к завершённым сделкам.',
        rule_version: 'archetype-rules-v2-spike',
        facts: [],
      },
      {
        card_id: 'achievements',
        kind: 'achievement',
        code: 'city_rhythm',
        reason: 'Достижение получено на уровне «expert».',
        rule_version: 'achievement-rules-v1-spike',
        facts: [
          {
            metric_code: 'max_activity_streak',
            actual: 18,
            operator: 'gte',
            threshold: 14,
            matched: true,
          },
        ],
      },
    ],
  };

  const result = applyExplanation(recap, explanation);
  assert.equal(result.role.reason, 'Продавцовская активность была заметнее всего.');
  assert.equal(result.style.reason, 'Активность часто приводила к завершённым сделкам.');
  assert.equal(result.badges[0]?.reason, 'Достижение получено на уровне «expert».');
  assert.deepEqual(result.badges[0]?.facts, ['Максимальная серия активности: 18 — не меньше 14']);
  assert.equal(result.badges[1]?.facts.length, 0);
});
