import { CATEGORY_META, DISTRICT_ORDER } from '@/shared/api/categories';
import type {
  Badge,
  Chapter,
  District,
  DistrictId,
  Profile,
  Recap,
  Trait,
  Unfinished,
} from '@/shared/types/recap';

/** Сырая активность профиля: действия, объявления, диалоги по категориям. */
type Activity = Record<DistrictId, readonly [actions: number, listings: number, dialogs: number]>;

interface Fixture {
  profile: Profile;
  cityName: string;
  seed: number;
  role: Trait;
  style: Trait;
  activity: Activity;
  badges: Badge[];
  chapters: Chapter[];
  unfinished: Unfinished[];
}

/**
 * Собирает районы из сырой активности. `share` — всегда производная от действий,
 * поэтому площадь квартала и цифра в легенде не могут разойтись.
 */
function toDistricts(activity: Activity): District[] {
  const total = DISTRICT_ORDER.reduce((sum, id) => sum + activity[id][0], 0);
  return DISTRICT_ORDER.map((id) => {
    const [actions, listings, dialogs] = activity[id];
    return {
      id,
      title: CATEGORY_META[id].title,
      tone: CATEGORY_META[id].tone,
      shade: CATEGORY_META[id].shade,
      actions,
      listings,
      dialogs,
      share: total === 0 ? 0 : actions / total,
    };
  });
}

const FIXTURES: Fixture[] = [
  {
    profile: {
      id: 'marina',
      name: 'Марина',
      tagline: 'Искала квартиру весь год',
      hint: 'Много просмотров, мало объявлений',
      tone: 'blue',
    },
    cityName: 'Город Марины',
    seed: 20250114,
    role: {
      title: 'Охотник за домом',
      reason:
        'Недвижимость забрала 31% всех действий за год — больше, чем следующие две категории вместе.',
    },
    style: {
      title: 'Дотошный исследователь',
      reason:
        'На одно избранное приходится 14 просмотров: решения принимаются после долгого сравнения, а не с первого взгляда.',
    },
    activity: {
      realty: [312, 0, 18],
      home: [204, 2, 9],
      electronics: [134, 12, 4],
      work: [96, 1, 6],
      personal: [88, 7, 5],
      hobby: [61, 0, 2],
      auto: [58, 0, 1],
      services: [44, 0, 3],
      pets: [19, 0, 0],
    },
    badges: [
      {
        id: 'early-bird',
        group: 'time',
        groupTitle: 'Время',
        title: 'Ранняя пташка',
        reason: 'Больше половины действий приходится на утро — до начала рабочего дня.',
        facts: [
          '68% просмотров — до 9:00',
          'Самый активный час — 7:40',
          '43 дня подряд с утренним заходом',
        ],
      },
      {
        id: 'deal-master',
        group: 'social',
        groupTitle: 'Общение',
        title: 'Мастер торга',
        reason:
          'Почти каждый диалог доходил до обсуждения условий, а не обрывался на первом сообщении.',
        facts: [
          '48 диалогов за год',
          'В 31 из них — больше 5 сообщений',
          'Среднее время ответа — 12 минут',
        ],
      },
      {
        id: 'settler',
        group: 'result',
        groupTitle: 'Итог',
        title: 'Новосёл',
        reason: 'Год закрылся сделкой в той категории, которая весь год была главной.',
        facts: [
          'Поиск длился 9 месяцев',
          'Просмотрено 312 объектов',
          'Финальный выбор — в декабре',
        ],
      },
    ],
    chapters: [
      {
        index: 1,
        title: 'Всё началось с одного района',
        stat: { value: '312', label: 'объявлений о недвижимости' },
        narrative:
          'Год ты начала с поиска квартиры — и он задал ритм всему остальному. Этот квартал вырос первым и остался самым большим.',
        districtId: 'realty',
      },
      {
        index: 2,
        title: 'Твой город просыпался рано',
        stat: { value: '7:40', label: 'самый активный час' },
        narrative:
          'Пока остальные ещё спали, ты уже листала новые объявления. Утро — время, когда решения давались легче всего.',
        districtId: 'home',
        badgeId: 'early-bird',
      },
      {
        index: 3,
        title: 'Рядом выросла техника',
        stat: { value: '12', label: 'объявлений ты разместила сама' },
        narrative:
          'Освобождая место под переезд, ты стала продавцом. Небольшой, но самый деловой квартал города.',
        districtId: 'electronics',
      },
      {
        index: 4,
        title: 'Ты умеешь договариваться',
        stat: { value: '48', label: 'диалогов за год' },
        narrative:
          'Ты редко соглашалась на первую цену — и почти всегда доводила разговор до конца.',
        districtId: 'work',
        badgeId: 'deal-master',
      },
      {
        index: 5,
        title: 'Год закрылся сделкой',
        stat: { value: '9', label: 'месяцев поиска' },
        narrative:
          'Длинная дорога закончилась там же, где началась — в районе недвижимости. Твой дом теперь отмечен на карте.',
        districtId: 'personal',
        badgeId: 'settler',
      },
      {
        index: 6,
        title: 'Город собран',
        stat: { value: '1 016', label: 'действий за год' },
        narrative:
          'Девять районов, три стройки на окраине. Покрути город, загляни в каждый квартал — и посмотри, что осталось недостроенным.',
      },
    ],
    unfinished: [
      {
        id: 'favorites-silent',
        title: 'квартир в избранном без диалога',
        count: 8,
        ctaLabel: 'Написать продавцам',
      },
      {
        id: 'draft',
        title: 'черновик объявления с апреля',
        count: 1,
        ctaLabel: 'Дописать и опубликовать',
      },
      {
        id: 'saved-search',
        title: 'сохранённых поиска без новых заходов',
        count: 2,
        ctaLabel: 'Посмотреть, что появилось',
      },
    ],
  },

  {
    profile: {
      id: 'artem',
      name: 'Артём',
      tagline: 'Купил, починил, продал',
      hint: 'Продаж больше, чем покупок',
      tone: 'purple',
    },
    cityName: 'Город Артёма',
    seed: 20250422,
    role: {
      title: 'Мастер апгрейда',
      reason: 'На каждую покупку приходится 2,4 продажи: техника у тебя не оседает, а идёт дальше.',
    },
    style: {
      title: 'Ночной торговец',
      reason: '61% действий — после 22:00. Сделки закрывались тогда, когда город уже спал.',
    },
    activity: {
      electronics: [480, 64, 41],
      auto: [190, 9, 12],
      personal: [150, 28, 17],
      hobby: [90, 6, 5],
      services: [60, 0, 4],
      home: [40, 3, 2],
      work: [30, 0, 1],
      realty: [20, 0, 0],
      pets: [8, 0, 0],
    },
    badges: [
      {
        id: 'night-owl',
        group: 'time',
        groupTitle: 'Время',
        title: 'Ночной смотритель',
        reason: 'Основная активность приходится на часы после полуночи.',
        facts: ['61% действий — после 22:00', 'Пик — 00:20', '17 сделок закрыто ночью'],
      },
      {
        id: 'fast-reply',
        group: 'social',
        groupTitle: 'Общение',
        title: 'Отвечает за минуту',
        reason: 'Скорость ответа — главное конкурентное преимущество твоих объявлений.',
        facts: [
          '82 диалога за год',
          'Среднее время ответа — 74 секунды',
          'Ни одного диалога без ответа',
        ],
      },
      {
        id: 'top-seller',
        group: 'result',
        groupTitle: 'Итог',
        title: 'Оборот года',
        reason: 'Больше сотни объявлений за год — это уже не разбор шкафа, а система.',
        facts: [
          '110 объявлений размещено',
          '76 закрыто продажей',
          'Медианный срок продажи — 4 дня',
        ],
      },
    ],
    chapters: [
      {
        index: 1,
        title: 'Целый квартал под технику',
        stat: { value: '480', label: 'действий с электроникой' },
        narrative:
          'Почти половина твоего года — один район. Он вырос выше всех и виден с любой точки города.',
        districtId: 'electronics',
      },
      {
        index: 2,
        title: 'Твой город не спал',
        stat: { value: '00:20', label: 'час пиковой активности' },
        narrative:
          'Пока лента затихала, ты только разгонялся. Ночью конкуренция за ответ продавцу минимальная — и ты этим пользовался.',
        districtId: 'auto',
        badgeId: 'night-owl',
      },
      {
        index: 3,
        title: 'Вещи не задерживались',
        stat: { value: '4', label: 'дня — медианный срок продажи' },
        narrative:
          'Купленное почти сразу уходило дальше. Твой город — это не склад, а перевалочный пункт.',
        districtId: 'personal',
      },
      {
        index: 4,
        title: 'Ты отвечал быстрее всех',
        stat: { value: '74', label: 'секунды — среднее время ответа' },
        narrative:
          'Ни один диалог не остался без ответа. Для покупателя это часто важнее, чем цена.',
        districtId: 'hobby',
        badgeId: 'fast-reply',
      },
      {
        index: 5,
        title: 'Сто десять объявлений',
        stat: { value: '110', label: 'объявлений за год' },
        narrative:
          'Это уже не случайная распродажа старого — это выстроенный процесс со своим ритмом.',
        districtId: 'services',
        badgeId: 'top-seller',
      },
      {
        index: 6,
        title: 'Город собран',
        stat: { value: '1 068', label: 'действий за год' },
        narrative:
          'Девять районов и две стройки на окраине. Загляни в кварталы — и добей то, что осталось висеть.',
      },
    ],
    unfinished: [
      {
        id: 'unsold',
        title: 'объявления висят дольше 60 дней',
        count: 6,
        ctaLabel: 'Обновить цену',
      },
      {
        id: 'no-delivery',
        title: 'объявления без Авито Доставки',
        count: 14,
        ctaLabel: 'Подключить доставку',
      },
    ],
  },

  {
    profile: {
      id: 'lena',
      name: 'Лена',
      tagline: 'Дача, сад и кот',
      hint: 'Ярко выраженная сезонность',
      tone: 'green',
    },
    cityName: 'Город Лены',
    seed: 20250607,
    role: {
      title: 'Хозяйка сезона',
      reason:
        'Дом и дача плюс животные — 56% года. Активность включается весной и гаснет к ноябрю.',
    },
    style: {
      title: 'Плановый покупатель',
      reason:
        'Между добавлением в избранное и покупкой в среднем 11 дней: ты возвращаешься к решению, а не берёшь сразу.',
    },
    activity: {
      home: [340, 14, 21],
      pets: [180, 5, 16],
      hobby: [120, 3, 7],
      personal: [95, 11, 8],
      services: [80, 0, 9],
      realty: [40, 1, 2],
      electronics: [35, 2, 1],
      work: [20, 0, 1],
      auto: [12, 0, 0],
    },
    badges: [
      {
        id: 'spring-wave',
        group: 'time',
        groupTitle: 'Время',
        title: 'Живёт по сезону',
        reason: 'Активность распределена по году крайне неравномерно — с чётким весенним пиком.',
        facts: [
          '64% действий — с апреля по июль',
          'Пик — вторая неделя мая',
          'В декабре — 4 действия',
        ],
      },
      {
        id: 'good-neighbour',
        group: 'social',
        groupTitle: 'Общение',
        title: 'Добрый сосед',
        reason:
          'Ты почти всегда оставляла отзыв после сделки — так растёт доверие ко всей площадке.',
        facts: ['64 диалога за год', '19 отзывов оставлено', 'Средняя оценка от продавцов — 4,9'],
      },
      {
        id: 'garden-done',
        group: 'result',
        groupTitle: 'Итог',
        title: 'Участок укомплектован',
        reason: 'К концу сезона закрыты все категории, с которых год начинался.',
        facts: ['Куплено 23 вещи для дачи', '14 своих объявлений', 'Сезон закрыт в октябре'],
      },
    ],
    chapters: [
      {
        index: 1,
        title: 'Весной город ожил',
        stat: { value: '340', label: 'действий в разделе «Дом и дача»' },
        narrative: 'Твой год начался не в январе, а в апреле. Главный квартал вырос за одну весну.',
        districtId: 'home',
      },
      {
        index: 2,
        title: 'Май был громче всех',
        stat: { value: '64%', label: 'года уместилось в четыре месяца' },
        narrative:
          'С апреля по июль ты заходила почти каждый день, а к декабрю город затих. Так выглядит сезонная жизнь.',
        districtId: 'pets',
        badgeId: 'spring-wave',
      },
      {
        index: 3,
        title: 'Появился район с лапами',
        stat: { value: '16', label: 'диалогов про животных' },
        narrative: 'Корм, переноска, ветеринар — целый квартал, который вырос вокруг одного кота.',
        districtId: 'hobby',
      },
      {
        index: 4,
        title: 'После тебя оставались отзывы',
        stat: { value: '19', label: 'отзывов за год' },
        narrative: 'Ты почти всегда возвращалась, чтобы оценить сделку. Продавцы это запоминают.',
        districtId: 'personal',
        badgeId: 'good-neighbour',
      },
      {
        index: 5,
        title: 'Сезон закрыт',
        stat: { value: '23', label: 'вещи для участка' },
        narrative: 'К октябрю всё, что планировалось весной, было куплено. Город ушёл на зимовку.',
        districtId: 'services',
        badgeId: 'garden-done',
      },
      {
        index: 6,
        title: 'Город собран',
        stat: { value: '922', label: 'действия за год' },
        narrative:
          'Девять районов и три стройки на окраине. Загляни в кварталы — и посмотри, что стоит достроить к следующей весне.',
      },
    ],
    unfinished: [
      {
        id: 'season-search',
        title: 'сохранённых поиска ждут весны',
        count: 3,
        ctaLabel: 'Включить уведомления',
      },
      {
        id: 'cart-left',
        title: 'товара в избранном с июля',
        count: 5,
        ctaLabel: 'Вернуться к избранному',
      },
      {
        id: 'draft-lena',
        title: 'черновик объявления о рассаде',
        count: 1,
        ctaLabel: 'Опубликовать',
      },
    ],
  },

  {
    profile: {
      id: 'kirill',
      name: 'Кирилл',
      tagline: 'Искал работу, а нашёл стол',
      hint: 'Небольшой город — показывает нижнюю границу',
      tone: 'red',
    },
    cityName: 'Город Кирилла',
    seed: 20250903,
    role: {
      title: 'Соискатель',
      reason: 'Работа — 40% всех действий, при том что год в целом был спокойным.',
    },
    style: {
      title: 'Быстрые решения',
      reason:
        'На одно избранное приходится 3 просмотра — ты почти не сравниваешь и берёшь то, что подошло сразу.',
    },
    activity: {
      work: [120, 2, 11],
      personal: [55, 9, 6],
      electronics: [40, 1, 3],
      realty: [28, 0, 2],
      hobby: [20, 0, 1],
      services: [14, 0, 2],
      home: [10, 1, 1],
      auto: [8, 0, 0],
      pets: [3, 0, 0],
    },
    badges: [
      {
        id: 'weekend-only',
        group: 'time',
        groupTitle: 'Время',
        title: 'Житель выходных',
        reason: 'Почти вся активность приходится на субботу и воскресенье.',
        facts: [
          '71% действий — в выходные',
          'Пик — воскресенье, 14:00',
          'В будни в среднем 0,4 действия в день',
        ],
      },
      {
        id: 'to-the-point',
        group: 'social',
        groupTitle: 'Общение',
        title: 'Строго по делу',
        reason: 'Диалоги короткие и заканчиваются решением, а не затухают.',
        facts: ['26 диалогов за год', 'Медиана — 3 сообщения', '19 диалогов закончились сделкой'],
      },
      {
        id: 'new-desk',
        group: 'result',
        groupTitle: 'Итог',
        title: 'Рабочее место собрано',
        reason: 'Смена работы потянула за собой покупки, которые сложились в один сценарий.',
        facts: [
          'Откликов на вакансии — 34',
          'Куплено 6 вещей для дома',
          'Всё уместилось в 3 недели октября',
        ],
      },
    ],
    chapters: [
      {
        index: 1,
        title: 'Небольшой, но твой',
        stat: { value: '298', label: 'действий за год' },
        narrative:
          'Твой город компактный — и это нормально. Здесь важен не размер, а то, что каждый квартал появился не случайно.',
        districtId: 'work',
      },
      {
        index: 2,
        title: 'Ты приходил по выходным',
        stat: { value: '71%', label: 'действий — в субботу и воскресенье' },
        narrative: 'Будни принадлежали работе, а Авито — выходным. Город оживал два дня в неделю.',
        districtId: 'personal',
        badgeId: 'weekend-only',
      },
      {
        index: 3,
        title: 'Октябрь всё изменил',
        stat: { value: '34', label: 'отклика на вакансии' },
        narrative: 'Три недели плотного поиска — и главный район города вырос почти мгновенно.',
        districtId: 'electronics',
      },
      {
        index: 4,
        title: 'Ты не тратил слов',
        stat: { value: '3', label: 'сообщения — медиана диалога' },
        narrative: 'Короткие разговоры, быстрые решения. Из 26 диалогов 19 закончились сделкой.',
        districtId: 'realty',
        badgeId: 'to-the-point',
      },
      {
        index: 5,
        title: 'Новое место собралось само',
        stat: { value: '6', label: 'покупок для дома' },
        narrative: 'Стол, кресло, лампа — за новой работой пришёл целый маленький квартал.',
        districtId: 'hobby',
        badgeId: 'new-desk',
      },
      {
        index: 6,
        title: 'Город собран',
        stat: { value: '9', label: 'районов' },
        narrative:
          'Даже спокойный год складывается в карту. Загляни в кварталы — половину города ещё можно достроить.',
      },
    ],
    unfinished: [
      {
        id: 'resume',
        title: 'резюме без обновления с ноября',
        count: 1,
        ctaLabel: 'Обновить резюме',
      },
      {
        id: 'kirill-fav',
        title: 'вакансии в избранном без отклика',
        count: 4,
        ctaLabel: 'Откликнуться',
      },
    ],
  },
];

const RULES_VERSION = 'recap-rules-1.0.0';

function toRecap(fixture: Fixture): Recap {
  const districts = toDistricts(fixture.activity);
  const actions = districts.reduce((sum, d) => sum + d.actions, 0);

  return {
    profileId: fixture.profile.id,
    year: 2025,
    rulesVersion: RULES_VERSION,
    seed: fixture.seed,
    cityName: fixture.cityName,
    totals: {
      actions,
      districts: districts.filter((d) => d.actions > 0).length,
      sites: fixture.unfinished.length,
    },
    role: fixture.role,
    style: fixture.style,
    districts,
    chapters: fixture.chapters,
    badges: fixture.badges,
    unfinished: fixture.unfinished,
    shareCard: {
      cityName: fixture.cityName,
      roleTitle: fixture.role.title,
      styleTitle: fixture.style.title,
      districts: districts.filter((d) => d.actions > 0).length,
      silhouetteSeed: fixture.seed,
    },
  };
}

export const PROFILES: Profile[] = FIXTURES.map((f) => f.profile);

export const RECAPS: Record<string, Recap> = Object.fromEntries(
  FIXTURES.map((f) => [f.profile.id, toRecap(f)]),
);
