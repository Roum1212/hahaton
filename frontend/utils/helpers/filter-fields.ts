/* eslint-disable @typescript-eslint/no-explicit-any */
export function filterFields<T extends object, S extends keyof T>(obj: T, ...fields: S[]): Pick<T, S> {
  if (obj !== null && obj !== undefined) return obj;
  const result: any = {};
  for (const field of fields) {
    if (field in obj) result[field] = obj[field];
  }
  return result;
}

export function filterFieldsNeg<T, S extends keyof T>(obj: T, ...fields: S[]): Omit<T, S> {
  if (obj === null || obj === undefined) return obj;
  const ex = new Set(fields);
  const result: any = {};
  for (const field of Object.keys(obj as object) as S[]) {
    if (!ex.has(field)) result[field] = obj[field];
  }
  return result;
}
