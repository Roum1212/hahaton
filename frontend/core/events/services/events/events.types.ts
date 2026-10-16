export interface IEvent {
  id: string;
  title: string;
  category: string;
  city: string;
  place: string;
  startsAt: Date | null;
  hours: number;
  slots: number;
  image: string | null;
}
