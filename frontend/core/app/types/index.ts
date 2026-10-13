export interface MapperService<CREATE_DTO, UPDATE_DTO, CREATE_REQUEST, UPDATE_REQUEST, RESPONSE, MODEL> {
  toApiCreateDto(dto: CREATE_DTO): CREATE_REQUEST;
  toApiUpdateDto(dto: UPDATE_DTO): UPDATE_REQUEST;
  toInterface(data: RESPONSE): MODEL;
}
export interface ApiPagination {
  limit: number;
  offset: number;
}
export type Image = string | null | undefined;

export enum LogTypeEnum {
  SUCCESS = 'success',
  WARNING = 'warning',
  FAIL = 'fail',
  INFO = 'info',
  ERROR = 'error',
}
