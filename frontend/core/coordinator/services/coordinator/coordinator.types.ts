export interface ICoordinator {
  id: string;
  displayName: string;
  description: string;
  avatarUrl: string;
}

export interface IUpdateCoordinatorDto {
  description: string;
  avatarUrl: string;
}

export interface ICoordinatorSettings {
  cooldownSeconds: number;
}
