import type { ICreateUserDto } from '~/core/users/services/users/users.types';

export const validateUserDto = (data: ICreateUserDto): IValidationResult => {
  const messages: string[] = [];

  if (!/^\S+@\S+\.\S+$/.test(data.email)) messages.push('Введите корректный Email');
  if (!data.firstName.length) messages.push('Введите имя пользователя');

  return { isValid: !messages.length, message: messages.join(', ') };
};
