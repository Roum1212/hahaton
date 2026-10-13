import { UserRole } from '~/gen/user/v1/user_pb';

export { UserRole };

export interface IUser {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  fullName: string;
  role: UserRole;
  active: boolean;
  createdAt: Date | null;
}

export interface ICreateUserDto {
  email: string;
  firstName: string;
  lastName: string;
  role: UserRole;
  active?: boolean;
}

export const USER_ROLE_LABELS: Record<UserRole, string> = {
  [UserRole.UNSPECIFIED]: 'Не указана',
  [UserRole.USER]: 'Пользователь',
  [UserRole.ADMIN]: 'Администратор',
};
