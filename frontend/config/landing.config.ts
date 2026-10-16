import {
  CalendarCheck,
  ClipboardCheck,
  Coins,
  FileWarning,
  Gift,
  History,
  Puzzle,
  Repeat,
  Search,
  Send,
  ShieldCheck,
  ShoppingBag,
  UserX,
  Wallet,
} from 'lucide-vue-next';
import { markRaw, type Component } from 'vue';

export const COINS_PER_HOUR = 10;

export const COINS_DISCLAIMER =
  'Койны — внутренние баллы демонстрационной системы. Они не являются денежным средством и не дают гарантированного права на награду.';

export interface ILandingPhoto {
  src: string;
  alt: string;
  source?: string;
  href?: string;
  fit?: 'cover' | 'contain';
}

export interface ILandingNavItem {
  label: string;
  anchor: string;
}

export interface ILandingStat {
  value: string;
  label: string;
  link: string;
  anchor: string;
  hint?: string;
}

export interface ILandingCard {
  title: string;
  text: string;
  tags: string[];
  icon: Component;
  photo?: ILandingPhoto;
}

export interface ILandingStep {
  title: string;
  text: string;
  chip: string;
  icon: Component;
  photo?: ILandingPhoto;
}

export interface ILandingAudience {
  key: string;
  label: string;
  title: string;
  points: string[];
  action: string;
  photo: ILandingPhoto;
}

export interface ILandingFaq {
  question: string;
  answer: string;
}

const YANDEX_SOURCE = 'Источник: Yandex Images';
const ILLUSTRATION_SOURCE = 'Иллюстрация: команда ВОЛОНТЕРИМ';

export const LANDING_PHOTOS = {
  team: {
    src: '/images/landing/team.webp',
    alt: 'Волонтёры в толстовках «Меняю мир добротой»',
    source: YANDEX_SOURCE,
  },
  words: {
    src: '/images/landing/words.webp',
    alt: 'Слова о волонтёрстве',
    source: YANDEX_SOURCE,
  },
  hands: {
    src: '/images/landing/hands.webp',
    alt: 'Поднятые ладони с сердцами',
    source: YANDEX_SOURCE,
  },
  dobrovolec: {
    src: '/images/landing/dobrovolec.webp',
    alt: 'Логотип «Доброволец России»',
    source: 'Изображение: Всероссийский конкурс «Доброволец России»',
    fit: 'contain',
  },
  coordinator: {
    src: '/images/landing/coordinator.svg',
    alt: 'Начисление койнов координатором',
    source: ILLUSTRATION_SOURCE,
  },
  community: {
    src: '/images/landing/community.svg',
    alt: 'Молодёжное сообщество',
    source: ILLUSTRATION_SOURCE,
  },
  banner: {
    src: '/images/landing/banner-volunteers.jpg',
    alt: 'Волонтёры убирают листья в парке',
  },
  eventsEmpty: {
    src: '/images/landing/events-empty.jpg',
    alt: 'Пустой календарь с лупой',
  },
  featureTracking: {
    src: '/images/landing/feature-tracking.webp',
    alt: 'Планшет с отметками участия',
  },
  featureCoins: {
    src: '/images/landing/feature-coins.webp',
    alt: 'Койны с сердцем и кубок',
  },
  featureShield: {
    src: '/images/landing/feature-shield.webp',
    alt: 'Щит с галочкой и замок',
  },
  featureRewards: {
    src: '/images/landing/feature-rewards.webp',
    alt: 'Подарочная коробка с поощрениями',
  },
  stepCatalog: {
    src: '/images/landing/step-catalog.jpg',
    alt: 'Карточки мероприятий и лупа',
  },
  stepRequest: {
    src: '/images/landing/step-request.jpg',
    alt: 'Бумажный самолётик вылетает из телефона',
  },
  stepResults: {
    src: '/images/landing/step-results.jpg',
    alt: 'Планшет с оценками участников',
  },
  stepAccrual: {
    src: '/images/landing/step-accrual.jpg',
    alt: 'Койн падает в кошелёк',
  },
  stepReward: {
    src: '/images/landing/step-reward.jpg',
    alt: 'Сумка с койнами',
  },
  stepHandover: {
    src: '/images/landing/step-handover.jpg',
    alt: 'Передача подарка из рук в руки',
  },
  problemScattered: {
    src: '/images/landing/problem-scattered.jpg',
    alt: 'Разрозненные кусочки пазла',
  },
  problemVerify: {
    src: '/images/landing/problem-verify.jpg',
    alt: 'Длинный список под лупой',
  },
  problemUnclear: {
    src: '/images/landing/problem-unclear.jpg',
    alt: 'Койн и знак вопроса',
  },
  problemCheating: {
    src: '/images/landing/problem-cheating.jpg',
    alt: 'Дублированные койны и знак предупреждения',
  },
  problemHistory: {
    src: '/images/landing/problem-history.jpg',
    alt: 'Пустой блокнот и часы',
  },
  problemMotivation: {
    src: '/images/landing/problem-motivation.jpg',
    alt: 'Увядший цветок в горшке',
  },
} satisfies Record<string, ILandingPhoto>;

export const LANDING_NAV: ILandingNavItem[] = [
  { label: 'О проекте', anchor: 'about' },
  { label: 'Как это работает', anchor: 'steps' },
  { label: 'Мероприятия', anchor: 'events' },
  { label: 'Для кого', anchor: 'audience' },
  { label: 'Вопросы', anchor: 'faq' },
];

export const LANDING_STATS: ILandingStat[] = [
  { value: '1', label: 'город в пилоте', link: 'Выбрать мероприятие', anchor: 'events' },
  { value: `${COINS_PER_HOUR}`, label: 'койнов за 1 час участия', link: 'Как начисляются', anchor: 'steps', hint: COINS_DISCLAIMER },
  { value: '3', label: 'поощрения в магазине', link: 'Как обменять', anchor: 'steps' },
  { value: '2', label: 'роли: волонтёр и координатор', link: 'Для кого сервис', anchor: 'audience' },
];

export const LANDING_PROBLEMS: ILandingCard[] = [
  {
    title: 'Данные в разных местах',
    text: 'Мероприятия, часы и поощрения приходится собирать по чатам, таблицам и бумажным спискам.',
    tags: ['Потерянные часы', 'Ручной учёт'],
    icon: markRaw(Puzzle),
    photo: LANDING_PHOTOS.problemScattered,
  },
  {
    title: 'Сложно проверить активность',
    text: 'Координатор тратит время на сверку списков и не видит, кто реально был на мероприятии.',
    tags: ['Сверка вручную', 'Ошибки'],
    icon: markRaw(Search),
    photo: LANDING_PHOTOS.problemVerify,
  },
  {
    title: 'Непонятно, откуда баллы',
    text: 'Участник не видит, за что и когда ему начислили баллы, — появляются вопросы и споры.',
    tags: ['Недоверие', 'Споры'],
    icon: markRaw(FileWarning),
    photo: LANDING_PHOTOS.problemUnclear,
  },
  {
    title: 'Риск накрутки',
    text: 'Без ролей и проверок баллы можно начислить дважды или приписать самому себе.',
    tags: ['Злоупотребления', 'Двойные начисления'],
    icon: markRaw(UserX),
    photo: LANDING_PHOTOS.problemCheating,
  },
  {
    title: 'Нет истории участия',
    text: 'Добрые дела не складываются в понятную историю, которую можно показать.',
    tags: ['Нет портфолио', 'Нет признания'],
    icon: markRaw(History),
    photo: LANDING_PHOTOS.problemHistory,
  },
  {
    title: 'Нет мотивации возвращаться',
    text: 'Без обратной связи и поощрений волонтёры реже приходят на следующие акции.',
    tags: ['Отток', 'Низкая вовлечённость'],
    icon: markRaw(Repeat),
    photo: LANDING_PHOTOS.problemMotivation,
  },
];

export const LANDING_FEATURES: ILandingCard[] = [
  {
    title: 'Прозрачный учёт участия',
    text: 'Каждая заявка проходит понятные статусы: подана, подтверждена или отклонена координатором.',
    tags: ['Статусы заявок', 'История участия'],
    icon: markRaw(ClipboardCheck),
    photo: LANDING_PHOTOS.featureTracking,
  },
  {
    title: 'Койны по итогам мероприятия',
    text: `После завершения мероприятия координатор подтверждает часы участников, а койны начисляются автоматически: 1 час = ${COINS_PER_HOUR} койнов.`,
    tags: [`1 час = ${COINS_PER_HOUR} койнов`, 'Подтверждение координатором'],
    icon: markRaw(Coins),
    photo: LANDING_PHOTOS.featureCoins,
  },
  {
    title: 'Защита от накрутки',
    text: 'Заявку нельзя подтвердить дважды, списать больше баланса или начислить койны себе без роли координатора.',
    tags: ['Роли', 'Проверки'],
    icon: markRaw(ShieldCheck),
    photo: LANDING_PHOTOS.featureShield,
  },
  {
    title: 'Магазин поощрений',
    text: 'Обменивайте койны на учебные сувениры и нематериальные поощрения — без реальной оплаты.',
    tags: ['Сувениры', 'Благодарности'],
    icon: markRaw(Gift),
    photo: LANDING_PHOTOS.featureRewards,
  },
];

export const LANDING_STEPS: ILandingStep[] = [
  {
    title: 'Выбор мероприятия',
    text: 'Волонтёр открывает каталог мероприятий своего города.',
    chip: 'Ваш город',
    icon: markRaw(Search),
    photo: LANDING_PHOTOS.stepCatalog,
  },
  {
    title: 'Заявка',
    text: 'Подаёт заявку на участие в один клик.',
    chip: 'Заявка подана',
    icon: markRaw(Send),
    photo: LANDING_PHOTOS.stepRequest,
  },
  {
    title: 'Итоги мероприятия',
    text: 'Координатор подтверждает, сколько часов отработал каждый участник.',
    chip: '3 часа',
    icon: markRaw(CalendarCheck),
    photo: LANDING_PHOTOS.stepResults,
  },
  {
    title: 'Начисление',
    text: `Койны начисляются по курсу 1 час = ${COINS_PER_HOUR} койнов и сразу попадают в историю операций.`,
    chip: '+30 койнов',
    icon: markRaw(Coins),
    photo: LANDING_PHOTOS.stepAccrual,
  },
  {
    title: 'Поощрение',
    text: 'Волонтёр выбирает учебный сувенир и оформляет запрос.',
    chip: '−20 койнов',
    icon: markRaw(ShoppingBag),
    photo: LANDING_PHOTOS.stepReward,
  },
  {
    title: 'Выдача',
    text: 'Координатор видит запрос на выдачу, у волонтёра остаётся баланс.',
    chip: 'Остаток 10',
    icon: markRaw(Wallet),
    photo: LANDING_PHOTOS.stepHandover,
  },
];

export const LANDING_AUDIENCE: ILandingAudience[] = [
  {
    key: 'volunteers',
    label: 'Добровольцам',
    title: 'Делайте добрые дела и следите за своим вкладом',
    points: [
      'Каталог городских мероприятий',
      'Заявка на участие в один клик',
      'Баланс койнов и история часов',
      'Обмен койнов на поощрения',
    ],
    action: 'Стать волонтёром',
    photo: LANDING_PHOTOS.dobrovolec,
  },
  {
    key: 'coordinators',
    label: 'Координаторам',
    title: 'Начисляйте койны без таблиц и сверок',
    points: [
      'Подтверждение часов участников',
      'Защита от повторного начисления',
      'Очередь запросов на выдачу поощрений',
      'Понятная история по каждому участнику',
    ],
    action: 'Войти как координатор',
    photo: LANDING_PHOTOS.coordinator,
  },
  {
    key: 'center',
    label: 'Добро.Центру Обнинска',
    title: 'Единая картина волонтёрской активности города',
    points: ['Все мероприятия в одном каталоге', 'Прозрачная история начислений', 'Правила поощрений в одном месте'],
    action: 'Подключиться',
    photo: LANDING_PHOTOS.hands,
  },
  {
    key: 'youth',
    label: 'Молодёжным организациям',
    title: 'Вовлекайте участников и отмечайте их вклад',
    points: ['Публикация своих мероприятий', 'Учёт часов участников', 'Поощрения без бумажной волокиты'],
    action: 'Присоединиться',
    photo: LANDING_PHOTOS.community,
  },
];

export interface ILandingRule {
  text: string;
  icon: string;
}

export const LANDING_RULES: ILandingRule[] = [
  { text: '1 заявка — 1 начисление', icon: 'pi pi-check-circle' },
  { text: 'Баланс не уходит в минус', icon: 'pi pi-wallet' },
  { text: 'Койны начисляет только координатор', icon: 'pi pi-shield' },
];

export const LANDING_FAQ: ILandingFaq[] = [
  {
    question: 'Что такое койны?',
    answer: `Внутренние баллы сервиса. За каждый подтверждённый координатором час участия начисляется ${COINS_PER_HOUR} койнов. ${COINS_DISCLAIMER}`,
  },
  {
    question: 'Кто начисляет койны?',
    answer: `Координатор мероприятия. После завершения акции он подтверждает часы каждого участника, и койны начисляются автоматически: 1 час = ${COINS_PER_HOUR} койнов. Волонтёр не может начислить койны самому себе.`,
  },
  {
    question: 'Что будет, если заявку отклонят?',
    answer: 'Койны не начисляются, а в истории заявки отображается статус «Отклонено». Повторно подтвердить ту же заявку нельзя.',
  },
  {
    question: 'Можно ли обменять койны на деньги?',
    answer: 'Нет. Койны обмениваются только на учебные поощрения из магазина. Реальная оплата и выдача призов в пилоте не подключаются.',
  },
  {
    question: 'Как стать координатором?',
    answer: 'Роль координатора назначает администратор сервиса. Начислять койны без этой роли невозможно.',
  },
];
