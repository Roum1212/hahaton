export function getWorkingDaysByNumber(days: (1 | 2 | 3 | 4 | 5 | 6 | 7)[]): string {
  const daysEnum = {
    1: 'Пн',
    2: 'Вт',
    3: 'Ср',
    4: 'Чт',
    5: 'Пт',
    6: 'Сб',
    7: 'Вс',
  };

  return days.map((d) => daysEnum[d]).join(', ');
}
export function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function isEmpty(obj: object): boolean {
  return JSON.stringify(obj) === '{}';
}
