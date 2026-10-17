# Template Frontend

Базовый шаблон фронтенда: Nuxt 4 (SPA) + PrimeVue 4 + Tailwind 4 + Pinia. Связь с бэкендом — [Connect RPC](https://connectrpc.com/) по контрактам из `proto/`, вход — через [Флида ID](https://docs.flida.ru/id/sdk/web).

## Быстрый старт

```bash
pnpm install
cp .env.example .env
pnpm mock
pnpm dev
```

`pnpm mock` поднимает тестовый Connect-сервер на `http://localhost:8080` с данными в памяти. Вход в mock проверяет токен через UserInfo Флида ID; пользователь с подтверждённым email из `MOCK_ADMIN_EMAIL` (по умолчанию `admin@example.com`) получает роль администратора.

## Команды

| Команда                                        | Назначение                                               |
| ---------------------------------------------- | -------------------------------------------------------- |
| `pnpm dev`                                     | dev-сервер на `http://localhost:3000`                    |
| `pnpm build` / `pnpm preview`                  | production-сборка и её просмотр                          |
| `pnpm start` / `pnpm stop`                     | запуск/остановка через pm2 (`build.config.cjs`)          |
| `pnpm proto:gen`                               | генерация TS-кода из `../proto/` в `genproto/`          |
| `pnpm proto:lint`                              | `buf lint` для proto-контрактов                          |
| `pnpm mock`                                    | mock Connect RPC сервер                                  |
| `pnpm typecheck` / `pnpm lint` / `pnpm format` | проверки и форматирование                                |

Из корня репозитория `make proto-generate` генерирует одновременно Go-код в `../backend/genproto/` и TypeScript-код в `genproto/`. Для раздельной генерации доступны `make proto-generate-go` и `make proto-generate-ts`.

## Переменные окружения

| Переменная                           | Описание                                                                                                           |
| ------------------------------------ | ------------------------------------------------------------------------------------------------------------------ |
| `APP_RPC_URL`                        | базовый URL Connect RPC бэкенда                                                                                    |
| `APP_IMAGES_URL`                     | базовый URL для изображений (`UiImage`, `UiFile`)                                                                  |
| `APP_PUBLIC_URL`                     | публичный URL фронтенда                                                                                            |
| `APP_PUBLIC_SITE_NAME`               | название проекта (title, сайдбар, страница входа)                                                                  |
| `APP_PUBLIC_SITE_DESCRIPTION`        | meta description                                                                                                   |
| `APP_FLIDA_CLIENT_ID`                | Client ID проекта в консоли Флида ID                                                                               |
| `APP_FLIDA_REDIRECT_URI`             | адрес возврата, по умолчанию `<APP_PUBLIC_URL>/auth/callback`; должен совпадать с Redirect URL в «Платформы → Web» |
| `APP_FLIDA_MODE`                     | OIDC-режим: `redirect` (по умолчанию) или `popup`                                                                  |
| `APP_FLIDA_FEDCM`                    | режимы FedCM через запятую: `active,passive` (по умолчанию), `active`, `passive` или пусто                         |
| `APP_FLIDA_POST_LOGOUT_REDIRECT_URI` | если задан, после выхода пользователь также выходит из Флида ID и возвращается на этот адрес                       |

## Структура

```
../proto/               общие контракты (*.proto), источник правды для API
genproto/               сгенерированный код protobuf-es (не редактировать руками)
plugins/connect.ts      Connect-транспорты: $transport и $publicTransport
core/
  app/rpc/              BaseRpcApi, интерсепторы (auth, refresh, errors), хранилище токена
  app/stores/           общие сторы (настройки таблиц, фильтры, списки)
  app/types/            общие типы (MapperService, ApiPagination)
  auth/                 вход через Флида ID и сессия (api/flida — обёртка Web SDK)
  users/                пример модуля: пользователи
  <module>/
    index.ts            <MODULE>_API() и <MODULE>_SERVICES()
    api/<entity>/       обёртка над RPC-клиентом + типы запросов
    services/<entity>/  маппинг protobuf-сообщений в доменные интерфейсы
    stores/             Pinia-сторы модуля
composables/
  app/                  useApp, useAppLayout, useAppMenu, useDialogManager, useLoader, usePagination ...
  auth/                 useAuthActions, useFlidaAuth
  helpers/useServices   единая точка доступа к сервисам всех модулей
  <module>/             form/actions/table/dialogs-хелперы фичи
components/
  layout/               лоадер, модальный хост, боковые панели
  ui/                   UI-кит (UiTable, UiTableV2, UiButton, UiCard, UiSearch ...)
  <module>/             компоненты фичи
layouts/cabinet/        оболочка кабинета: колонка разделов + панель меню, топбар, AppPage
assets/scss/            _base (токены --app-*), _cabinet (оболочка и блоки app-*), _mobile
middleware/             0.auth → 1.permissions → 2.menu
config/                 меню и маршруты
mock/                   mock Connect RPC сервер
```

## Поток данных

```
Page / Component
  → composable фичи (useUsersActions, useUsersDialogs)
    → Pinia store (useUsersStore)
      → Service (UsersService: DTO ↔ protobuf, доменные модели)
        → Api (UsersApi extends BaseRpcApi<typeof UserService>)
          → Connect transport (interceptors) → бэкенд
```

## Авторизация

Вход выполняется через Флида ID Web SDK. SDK подключается с CDN (`https://id.flida.ru/sdk/v1/sdk.es.js`) в `core/auth/api/flida/flida.api.ts`, типы описаны в `flida.types.ts`.

```
/auth/login ── кнопка ──▶ FedCM (active) или OIDC redirect/popup ──▶ токены Флида ID
/auth/callback ◀── redirect ── SDK сам проверяет state/PKCE и обменивает code
        │
        ▼
AuthService.LoginWithFlida({ flidaAccessToken }) ──▶ бэкенд: GET id.flida.ru/api/oidc/userinfo,
                                                    поиск/создание пользователя по sub
        ◀── access-токен проекта + refresh-кука
```

- На странице входа SDK сразу запускает фоновую попытку (passive FedCM); её ошибки не показываются.
- После явного выхода passive FedCM отключается до конца сессии вкладки, чтобы пользователя не залогинило обратно автоматически.
- Токены Флида ID на фронтенде не хранятся: они нужны только для обмена на сессию проекта.
- Сессия проекта: access-токен в `localStorage` уходит в `Authorization: Bearer <token>`, refresh-токен бэкенд ставит httpOnly-кукой; запросы идут с `credentials: 'include'`.
- На ответ `Unauthenticated` интерсептор один раз вызывает `AuthService.RefreshToken` (параллельные запросы ждут один общий refresh) и повторяет исходный запрос.
- Ошибки RPC показываются тостом: текст ошибки сервера или сообщение по коду (`utils/rpc.utils.ts`).
- Права: `definePageMeta({ permissions: ['admin'] })` на странице и `permissions` у разделов и пунктов меню.

Что нужно от бэкенда:

1. `LoginWithFlida`: вызвать `GET https://id.flida.ru/api/oidc/userinfo` с `Authorization: Bearer <flidaAccessToken>`; на `401` вернуть `Unauthenticated`.
2. Идентифицировать пользователя по `sub`. Привязывать существующего пользователя по email — только если `email_verified: true`.
3. Разрешить CORS для origin фронтенда с `Access-Control-Allow-Credentials: true`.

Настройка в консоли Флида ID: в «Платформы → Web» добавить Redirect URL, совпадающий с `APP_FLIDA_REDIRECT_URI` (например `http://localhost:3000/auth/callback`).

## Загрузка аватаров и изображений поощрений

Загрузка идёт через домен `core/image` (ImageApi → ImageService → `useImagesStore`) по общему
Connect transport с cookie-сессией и перехватчиком ошибок:

```ts
const image = await useImagesStore().upload(file);
// image.url — полный URL для превью, image.path — для сохранения карточки.
// UpdateVolunteer/UpdateCoordinator: { avatarPath: image.path }.
// CreateReward/UpdateReward: { imagePath: image.path }.
```

Готовый компонент выбора файла — `AppImageUpload` (`components/app/ImageUpload.vue`):
`accept="image/jpeg,image/png"`, проверка размера до 5 МиБ на клиенте, превью и кнопка удаления.
Относительные `/storage/...` из ответов бэкенда превращаются в полный URL через `resolveImageUrl`
(`utils/helpers/image.helper.ts`) в сервисах волонтёра, координатора и поощрений.
Пустой путь удаляет картинку, отсутствующее поле сохраняет прежнюю.
Аватары профиля волонтёра и координатора независимы.
Mock-сервер хранит загруженные файлы в памяти и раздаёт их по `/storage/...`.

## Docker

Полный запуск описан в [корневой инструкции](../README.md). `frontend/Dockerfile`
собирает Nuxt с preset `node-server` и `nitro.noExternals`, чтобы зависимости PrimeVue
попали в серверный bundle целиком. Runtime содержит только `.output` и Node.js,
запускается от пользователя `node`, слушает HTTP на `0.0.0.0:3000` и имеет healthcheck.
Nginx и сертификаты настраиваются отдельно. В образ не попадают `.env`, node_modules
разработчика, mock-сервер и отчёты браузерных тестов.

Production-настройки передаются при запуске через стандартные переменные Nuxt:
`NUXT_PUBLIC_RPC_URL`, `NUXT_PUBLIC_SITE_URL`, `NUXT_PUBLIC_FLIDA_CLIENT_ID`,
`NUXT_PUBLIC_FLIDA_REDIRECT_URI` и остальные `NUXT_PUBLIC_*` из корневого Compose.
Это публичные значения для браузера: RPC URL должен быть внешним HTTPS-адресом,
а не внутренним `http://backend:8080`. Настройки API/входа меняются без пересборки.
Название и meta-description задаются также build arguments для статического HTML head;
при их изменении пересоберите frontend.
Серверные ключи, пароль PostgreSQL и ADMIN_API_KEY frontend не передаются.
