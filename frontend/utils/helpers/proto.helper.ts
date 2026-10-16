import { create } from '@bufbuild/protobuf';
import { timestampDate, timestampFromDate, type Timestamp } from '@bufbuild/protobuf/wkt';
import { DateSchema, type Date as ProtoDate } from '#genproto/google/type/date_pb';

export function fromProtoDate(value: ProtoDate | undefined): Date | null {
  if (!value || !value.year || !value.month || !value.day) return null;
  return new Date(value.year, value.month - 1, value.day);
}

export function toProtoDate(value: Date): ProtoDate {
  return create(DateSchema, { year: value.getFullYear(), month: value.getMonth() + 1, day: value.getDate() });
}

export function fromTimestamp(value: Timestamp | undefined): Date | null {
  return value ? timestampDate(value) : null;
}

export function fromBigInt(value: bigint | undefined): number {
  return value === undefined ? 0 : Number(value);
}

export function toTimestamp(value: Date): Timestamp {
  return timestampFromDate(value);
}

export function toApiPage(pagination: { limit: number; offset: number }): { limit: number; offset: bigint } {
  return { limit: pagination.limit, offset: BigInt(pagination.offset) };
}

export function createRequestId(): string {
  return crypto.randomUUID();
}
