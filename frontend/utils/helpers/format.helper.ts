import { format as formatInTimeZone, toZonedTime } from 'date-fns-tz';
import { ru } from 'date-fns/locale';
export function formatDate(date: Date | string | null | undefined, format = 'yyyy-MM-dd HH:mm:ss'): string {
  if (!date) return '-';
  const moscowTime = toZonedTime(date, 'Europe/Moscow');
  return formatInTimeZone(moscowTime, format, {
    timeZone: 'Europe/Moscow',
    locale: ru,
  });
}
export function formatPrice(price: string | number | undefined, symbol = '₽') {
  if (price === undefined || price === null) return '-';
  const parts = price.toString().split('.');
  const integerPart = parts[0];
  if (!integerPart) return '-';
  const fractionalPart = parts.length > 1 ? '.' + parts[1] : '';
  const formattedInteger = integerPart.replace(/\B(?=(\d{3})+(?!\d))/g, ' ');
  return `${formattedInteger}${fractionalPart} ${symbol}`;
}
export function setLocale() {
  const primevue = usePrimeVue();
  primevue.config.locale = {
    accept: 'Да',
    addRule: 'Добавить правило',
    am: 'AM',
    apply: 'Применить',
    cancel: 'Отменить',
    choose: 'Выбрать',
    chooseDate: 'Выбрать дату',
    chooseMonth: 'Выбрать месяц',
    chooseYear: 'Выбрать год',
    clear: 'Очистить',
    completed: 'Завершено',
    contains: 'Содержит',
    dateAfter: 'Дата после',
    dateBefore: 'Дата до',
    dateFormat: 'dd.mm.yy',
    dateIs: 'Дата равна',
    dateIsNot: 'Дата не равна',
    dayNames: ['Воскресенье', 'Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота'],
    dayNamesMin: ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'],
    dayNamesShort: ['Вос', 'Пон', 'Вто', 'Сре', 'Чет', 'Пят', 'Суб'],
    emptyFilterMessage: 'Результаты не найдены',
    emptyMessage: 'Доступных опций нет',
    emptySearchMessage: 'Результаты не найдены',
    emptySelectionMessage: 'Не выбран элемент',
    endsWith: 'Заканчивается на',
    equals: 'Равно',
    fileSizeTypes: ['Б', 'КБ', 'МБ', 'ГБ', 'ТБ', 'ПБ', 'ЭБ', 'ЗБ', 'ЙБ'],
    firstDayOfWeek: 1,
    gt: 'Больше чем',
    gte: 'Больше или равно',
    lt: 'Меньше чем',
    lte: 'Меньше или равно',
    matchAll: 'Соответствует всем',
    matchAny: 'Соответствует любому',
    medium: 'Средний',
    monthNames: ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'],
    monthNamesShort: ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'],
    nextDecade: 'Следующее десятилетие',
    nextHour: 'Следующий час',
    nextMinute: 'Следующая минута',
    nextMonth: 'Следующий месяц',
    nextSecond: 'Следующая секунда',
    nextYear: 'Следующий год',
    noFilter: 'Без фильтра',
    notContains: 'Не содержит',
    notEquals: 'Не равно',
    passwordPrompt: 'Введите пароль',
    pending: 'В ожидании',
    pm: 'PM',
    prevDecade: 'Предыдущее десятилетие',
    prevHour: 'Предыдущий час',
    prevMinute: 'Предыдущая минута',
    prevMonth: 'Предыдущий месяц',
    prevSecond: 'Предыдущая секунда',
    prevYear: 'Предыдущий год',
    reject: 'Нет',
    removeRule: 'Удалить правило',
    searchMessage: '{0} результатов доступно',
    selectionMessage: '{0} элементов выбрано',
    showMonthAfterYear: false,
    startsWith: 'Начинается с',
    strong: 'Сильный',
    today: 'Сегодня',
    upload: 'Загрузить',
    weak: 'Слабый',
    weekHeader: 'Нед',
  };
}
