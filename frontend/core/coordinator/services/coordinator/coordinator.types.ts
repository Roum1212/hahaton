export interface ICoordinator {
  id: string;
  displayName: string;
  description: string;
  avatarUrl: string | null;
}

export interface IUpdateCoordinatorDto {
  description?: string;
  avatarPath?: string;
}

export interface ICoordinatorSettings {
  cooldownSeconds: number;
}
