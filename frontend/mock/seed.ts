import { CooperationStatus } from '../genproto/cooperation/v1/cooperation_service_pb';
import { EventCandidateRequirement, EventStatus } from '../genproto/event/v1/event_service_pb';
import { ParticipationCancellationReason, ParticipationStatus } from '../genproto/participation/v1/participation_service_pb';
import { RewardOrderCancellationReason, RewardOrderStatus } from '../genproto/reward-order/v1/reward_order_service_pb';
import { VolunteerOperationType } from '../genproto/volunteer/v1/volunteer_service_pb';
import {
  coinsForMinutes,
  cooperationKey,
  db,
  fullName,
  newId,
  newToken,
  type CoordinatorRec,
  type EventRec,
  type ParticipationRec,
  type Point,
  type RewardRec,
  type VolunteerRec,
} from './db';
import { changeParticipation, latestChatSequence, openChatAccess, recordOperation } from './domain';
import { invitationUrl } from './mappers';

const DAY = 24 * 60 * 60 * 1000;
const HOUR = 60 * 60 * 1000;

export const OBNINSK_ID = '05a647ff-b859-4c93-8c8e-3590110a0085';

const at = (days: number, hour: number) => {
  const date = new Date(Date.now() + days * DAY);
  date.setHours(hour, 0, 0, 0);
  return date;
};

const PLACES: Record<string, Point> = {
  belkino: { latitude: 55.1215, longitude: 36.5958 },
  protva: { latitude: 55.0925, longitude: 36.5862 },
  lenina: { latitude: 55.1052, longitude: 36.6101 },
  center: { latitude: 55.0981, longitude: 36.6203 },
  shelter: { latitude: 55.0814, longitude: 36.6402 },
  stadium: { latitude: 55.1102, longitude: 36.6265 },
  kurchatova: { latitude: 55.1013, longitude: 36.6334 },
};

const addCities = () => {
  db.cities.set(OBNINSK_ID, {
    id: OBNINSK_ID,
    name: 'Обнинск',
    regionName: 'Калужская обл',
    center: { latitude: 55.094411, longitude: 36.612196 },
    boundary: [
      { latitude: 55.135, longitude: 36.545 },
      { latitude: 55.142, longitude: 36.605 },
      { latitude: 55.13, longitude: 36.665 },
      { latitude: 55.105, longitude: 36.685 },
      { latitude: 55.075, longitude: 36.67 },
      { latitude: 55.06, longitude: 36.62 },
      { latitude: 55.07, longitude: 36.565 },
      { latitude: 55.095, longitude: 36.54 },
    ],
  });
  db.cities.set('b502ae45-897e-4b6f-9776-6ff49740b537', {
    id: 'b502ae45-897e-4b6f-9776-6ff49740b537',
    name: 'Калуга',
    regionName: 'Калужская обл',
    center: { latitude: 54.513845, longitude: 36.261215 },
  });
  db.cities.set('0c5b2444-70a0-4932-980c-b4dc0d3f02b5', {
    id: '0c5b2444-70a0-4932-980c-b4dc0d3f02b5',
    name: 'Москва',
    regionName: 'г Москва',
    center: { latitude: 55.75396, longitude: 37.620393 },
  });
};

const addVolunteer = (
  firstName: string,
  lastName: string,
  middleName: string,
  email: string,
  phone: string,
  year: number,
): VolunteerRec => {
  const volunteer: VolunteerRec = {
    id: newId(),
    iss: 'seed',
    sub: email,
    firstName,
    lastName,
    middleName,
    email,
    phone,
    dateOfBirth: { year, month: 4, day: 12 },
    coinBalance: 0,
    confirmedMinutes: 0,
    operationSequence: 0,
  };
  db.volunteers.set(volunteer.id, volunteer);
  return volunteer;
};

const addCoordinator = (id: string, description: string, cooldownSeconds = 86400, avatarUrl = ''): CoordinatorRec => {
  const coordinator: CoordinatorRec = { id, description, cooldownSeconds, avatarUrl };
  db.coordinators.set(id, coordinator);
  return coordinator;
};

const setCooperation = (coordinatorId: string, volunteerId: string, status: CooperationStatus, createdAt = at(-30, 12)) => {
  db.cooperations.set(cooperationKey(coordinatorId, volunteerId), {
    coordinatorId,
    volunteerId,
    status,
    retryAfter: status === CooperationStatus.REJECTED ? at(1, 12) : undefined,
    createdAt,
    updatedAt: createdAt,
  });
};

const addInvitation = (coordinatorId: string) => {
  const invitation = { id: newId(), coordinatorId, token: newToken(), createdAt: at(-20, 12) };
  db.invitations.set(invitation.id, invitation);
  return invitation;
};

interface EventInput {
  title: string;
  description: string;
  location: string;
  point: Point;
  startsAt: Date;
  hours: number;
  limit: number;
  requirements?: EventCandidateRequirement[];
}

const addEvent = (coordinatorId: string, input: EventInput): EventRec => {
  const event: EventRec = {
    id: newId(),
    coordinatorId,
    title: input.title,
    description: input.description,
    location: input.location,
    point: input.point,
    startsAt: input.startsAt,
    endsAt: new Date(input.startsAt.getTime() + input.hours * HOUR),
    requirements: input.requirements ?? [],
    volunteerLimit: input.limit,
    status: EventStatus.PUBLISHED,
  };
  db.events.set(event.id, event);
  return event;
};

const apply = (event: EventRec, volunteer: VolunteerRec, createdAt = new Date(event.startsAt.getTime() - 5 * DAY)): ParticipationRec => {
  const participation: ParticipationRec = {
    id: newId(),
    eventId: event.id,
    volunteerId: volunteer.id,
    status: ParticipationStatus.PENDING,
    cancellationReason: ParticipationCancellationReason.UNSPECIFIED,
    createdAt,
    updatedAt: createdAt,
  };
  db.participations.set(participation.id, participation);
  changeParticipation(participation, ParticipationStatus.PENDING, { actor: { case: 'volunteerId', value: volunteer.id }, at: createdAt });
  return participation;
};

const approve = (event: EventRec, participation: ParticipationRec) => {
  changeParticipation(participation, ParticipationStatus.APPROVED, {
    actor: { case: 'coordinatorId', value: event.coordinatorId },
    at: new Date(participation.createdAt.getTime() + HOUR),
  });
  openChatAccess(event.id, participation.volunteerId);
};

const completeHistorical = (event: EventRec, results: { participation: ParticipationRec; minutes?: number }[]) => {
  const completedAt = new Date(event.endsAt.getTime() + 2 * HOUR);
  results.forEach(({ participation, minutes }) => {
    if (minutes) {
      const coins = coinsForMinutes(minutes);
      const operation = recordOperation({
        volunteerId: participation.volunteerId,
        coordinatorId: event.coordinatorId,
        type: VolunteerOperationType.PARTICIPATION_ACCRUAL,
        source: { case: 'participationId', value: participation.id },
        description: `Участие в «${event.title}»`,
        coinDelta: coins,
        minuteDelta: minutes,
        createdAt: completedAt,
      });
      participation.accrualOperationId = operation.id;
      participation.awardedCoins = coins;
      participation.awardedMinutes = minutes;
      changeParticipation(participation, ParticipationStatus.CONFIRMED, {
        actor: { case: 'coordinatorId', value: event.coordinatorId },
        at: completedAt,
      });
    } else {
      changeParticipation(participation, ParticipationStatus.NOT_ATTENDED, {
        actor: { case: 'coordinatorId', value: event.coordinatorId },
        at: completedAt,
      });
    }
  });
  event.status = EventStatus.COMPLETED;
  const requestId = `seed-${event.id}`;
  event.completion = { requestId, payloadKey: requestId };
  db.completionResults.set(`${event.coordinatorId}:${requestId}`, requestId);
};

const addReward = (coordinatorId: string, title: string, description: string, priceCoins: number, quantityLimit: number): RewardRec => {
  const reward: RewardRec = { id: newId(), coordinatorId, title, description, priceCoins, quantityLimit, deleted: false };
  db.rewards.set(reward.id, reward);
  return reward;
};

const buy = (reward: RewardRec, volunteer: VolunteerRec, createdAt: Date, status: RewardOrderStatus) => {
  const id = newId();
  recordOperation({
    volunteerId: volunteer.id,
    coordinatorId: reward.coordinatorId,
    type: VolunteerOperationType.REWARD_PURCHASE,
    source: { case: 'rewardOrderId', value: id },
    description: `Покупка «${reward.title}»`,
    coinDelta: -reward.priceCoins,
    minuteDelta: 0,
    createdAt,
  });
  db.rewardOrders.set(id, {
    id,
    coordinatorId: reward.coordinatorId,
    rewardId: reward.id,
    volunteerId: volunteer.id,
    priceCoins: reward.priceCoins,
    rewardTitle: reward.title,
    status,
    cancellationReason: RewardOrderCancellationReason.UNSPECIFIED,
    requestId: `seed-${id}`,
    createdAt,
    updatedAt: status === RewardOrderStatus.ISSUED ? new Date(createdAt.getTime() + DAY) : createdAt,
  });
};

const say = (event: EventRec, author: VolunteerRec, text: string, createdAt: Date) => {
  const sequence = latestChatSequence(event.id) + 1;
  db.chatSequences.set(event.id, sequence);
  db.chatMessages.push({
    id: newId(),
    eventId: event.id,
    author: event.coordinatorId === author.id ? { case: 'coordinatorId', value: author.id } : { case: 'volunteerId', value: author.id },
    authorDisplayName: fullName(author),
    text,
    sequence,
    requestKey: `seed-${newId()}`,
    createdAt,
  });
};

const seedVolunteers: VolunteerRec[] = [];

export const seed = () => {
  addCities();

  const irina = addVolunteer('Ирина', 'Ковалёва', 'Сергеевна', 'irina.kovaleva@example.com', '+79100000001', 1988);
  const nikolay = addVolunteer('Николай', 'Орлов', 'Андреевич', 'nikolay.orlov@example.com', '+79100000002', 1992);
  const maria = addVolunteer('Мария', 'Иванова', 'Павловна', 'maria.ivanova@example.com', '+79100000003', 2003);
  const dmitry = addVolunteer('Дмитрий', 'Петров', '', 'dmitry.petrov@example.com', '+79100000004', 2001);
  const aleksey = addVolunteer('Алексей', 'Соколов', 'Игоревич', 'aleksey.sokolov@example.com', '+79100000005', 2005);
  const ekaterina = addVolunteer('Екатерина', 'Морозова', '', 'ekaterina.morozova@example.com', '+79100000006', 1999);
  seedVolunteers.push(maria, dmitry, aleksey, ekaterina);

  addCoordinator(irina.id, 'Добро.Центр Обнинска: городские акции, экология и помощь пожилым людям.');
  addCoordinator(nikolay.id, 'Молодёжный центр Обнинска: фестивали, спорт и наставничество.', 0);

  setCooperation(irina.id, maria.id, CooperationStatus.ACTIVE);
  setCooperation(irina.id, dmitry.id, CooperationStatus.ACTIVE);
  setCooperation(irina.id, ekaterina.id, CooperationStatus.PENDING, at(-1, 15));
  setCooperation(nikolay.id, aleksey.id, CooperationStatus.ACTIVE);
  setCooperation(nikolay.id, maria.id, CooperationStatus.ACTIVE);

  const trees = addEvent(irina.id, {
    title: 'Высадка деревьев на набережной Протвы',
    description: 'Высаживаем липы и клёны вдоль набережной. Перчатки и инвентарь выдаём на месте.',
    location: 'Набережная реки Протвы',
    point: PLACES.protva!,
    startsAt: at(-12, 10),
    hours: 4,
    limit: 15,
  });
  const treesMaria = apply(trees, maria, at(-18, 12));
  const treesDmitry = apply(trees, dmitry, at(-17, 12));
  approve(trees, treesMaria);
  approve(trees, treesDmitry);
  say(trees, irina, 'Сбор у лодочной станции в 9:45. Возьмите удобную обувь.', at(-13, 18));
  say(trees, maria, 'Буду вовремя!', at(-13, 19));
  completeHistorical(trees, [{ participation: treesMaria, minutes: 240 }, { participation: treesDmitry }]);

  const shelter = addEvent(irina.id, {
    title: 'Помощь приюту для животных «Друг»',
    description: 'Выгул собак, уборка вольеров и разбор кормов.',
    location: 'Приют «Друг», Киевское шоссе',
    point: PLACES.shelter!,
    startsAt: at(-1, 10),
    hours: 3,
    limit: 6,
  });
  approve(shelter, apply(shelter, maria, at(-6, 12)));
  approve(shelter, apply(shelter, dmitry, at(-5, 12)));

  const subbotnik = addEvent(irina.id, {
    title: 'Субботник в парке усадьбы Белкино',
    description: 'Уборка листьев и мусора в парке. Перчатки, мешки и грабли выдаём.',
    location: 'Парк усадьбы Белкино, главный вход',
    point: PLACES.belkino!,
    startsAt: at(3, 10),
    hours: 3,
    limit: 20,
  });
  approve(subbotnik, apply(subbotnik, maria, at(-2, 12)));
  apply(subbotnik, dmitry, at(-1, 9));
  say(subbotnik, irina, 'Добрый день! Встречаемся у главного входа в парк в 9:50.', at(-1, 11));
  say(subbotnik, maria, 'Спасибо, буду!', at(-1, 12));

  addEvent(irina.id, {
    title: 'Экоакция «Раздельный сбор»',
    description: 'Принимаем от жителей бумагу, пластик и стекло, помогаем сортировать.',
    location: 'Проспект Ленина, у ДК ФЭИ',
    point: PLACES.lenina!,
    startsAt: at(10, 11),
    hours: 4,
    limit: 12,
  });
  addEvent(irina.id, {
    title: 'Зимний праздник для детей',
    description: 'Игры, конкурсы и вручение подарков детям из Добро.Центра.',
    location: 'Добро.Центр Обнинска',
    point: PLACES.center!,
    startsAt: at(40, 12),
    hours: 2,
    limit: 15,
    requirements: [EventCandidateRequirement.AGE_16_PLUS],
  });
  addEvent(nikolay.id, {
    title: 'Турнир по мини-футболу для подростков',
    description: 'Помощь судьям, регистрация команд и работа на площадке.',
    location: 'Стадион «Труд»',
    point: PLACES.stadium!,
    startsAt: at(7, 15),
    hours: 3,
    limit: 8,
  });
  addEvent(nikolay.id, {
    title: 'Школа наставников',
    description: 'Обучающая встреча для будущих наставников подростков.',
    location: 'Молодёжный центр, ул. Курчатова',
    point: PLACES.kurchatova!,
    startsAt: at(15, 18),
    hours: 2,
    limit: 25,
    requirements: [EventCandidateRequirement.AGE_18_PLUS, EventCandidateRequirement.PRELIMINARY_TRAINING],
  });

  const badge = addReward(irina.id, 'Значок «Волонтёр Обнинска»', 'Металлический значок с символикой городского волонтёрства.', 20, 50);
  addReward(irina.id, 'Благодарственное письмо', 'Именное письмо от Добро.Центра о вкладе в городские акции.', 40, 100);
  addReward(irina.id, 'Шоппер с символикой', 'Тканевая сумка с символикой волонтёрского движения.', 60, 10);
  addReward(nikolay.id, 'Билет на мастер-класс', 'Мастер-класс по ораторскому мастерству в Молодёжном центре.', 30, 20);
  buy(badge, maria, at(-10, 13), RewardOrderStatus.ISSUED);

  const invitations = [addInvitation(irina.id), addInvitation(nikolay.id)];
  return invitations.map((item) => ({ coordinator: fullName(db.volunteers.get(item.coordinatorId)!), url: invitationUrl(item.token) }));
};

export const provisionCoordinator = (volunteerId: string, options: Partial<Omit<CoordinatorRec, 'id'>> = {}) => {
  const owner = db.volunteers.get(volunteerId)!;
  addCoordinator(
    volunteerId,
    options.description ?? 'Координатор волонтёрских мероприятий Обнинска.',
    options.cooldownSeconds ?? 86400,
    options.avatarUrl ?? '',
  );
  const [maria, dmitry, aleksey, ekaterina] = seedVolunteers;
  if (!maria || !dmitry || !aleksey || !ekaterina) return;

  setCooperation(volunteerId, maria.id, CooperationStatus.ACTIVE);
  setCooperation(volunteerId, dmitry.id, CooperationStatus.ACTIVE);
  setCooperation(volunteerId, aleksey.id, CooperationStatus.ACTIVE);
  setCooperation(volunteerId, ekaterina.id, CooperationStatus.PENDING, at(0, 9));

  const cleanup = addEvent(volunteerId, {
    title: 'Уборка набережной Протвы',
    description: 'Собираем мусор вдоль набережной и на пляже.',
    location: 'Набережная реки Протвы, лодочная станция',
    point: PLACES.protva!,
    startsAt: at(-2, 10),
    hours: 3,
    limit: 10,
  });
  approve(cleanup, apply(cleanup, maria, at(-6, 12)));
  approve(cleanup, apply(cleanup, dmitry, at(-6, 14)));
  apply(cleanup, aleksey, at(-5, 10));
  say(cleanup, owner, 'Спасибо всем, кто пришёл! Итоги подведу сегодня.', at(-2, 14));

  const yard = addEvent(volunteerId, {
    title: 'Субботник во дворе на улице Курчатова',
    description: 'Покраска лавочек и уборка детской площадки.',
    location: 'Улица Курчатова, двор дома 21',
    point: PLACES.kurchatova!,
    startsAt: at(5, 10),
    hours: 3,
    limit: 8,
  });
  approve(yard, apply(yard, dmitry, at(-1, 12)));
  apply(yard, aleksey, at(-1, 15));

  const badge = addReward(volunteerId, 'Значок «Волонтёр Обнинска»', 'Значок с символикой городского волонтёрства.', 20, 30);
  addReward(volunteerId, 'Футболка волонтёра', 'Футболка с логотипом проекта.', 50, 15);
  addReward(volunteerId, 'Благодарственное письмо', 'Именное письмо координатора.', 40, 100);
  if (maria.coinBalance >= badge.priceCoins) buy(badge, maria, at(0, 9), RewardOrderStatus.PENDING);

  addInvitation(volunteerId);
};
