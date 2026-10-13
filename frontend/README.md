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

| Команда                                        | Назначение                                      |
| ---------------------------------------------- | ----------------------------------------------- |
| `pnpm dev`                                     | dev-сервер на `http://localhost:3000`           |
| `pnpm build` / `pnpm preview`                  | production-сборка и её просмотр                 |
| `pnpm start` / `pnpm stop`                     | запуск/остановка через pm2 (`build.config.cjs`) |
| `pnpm proto:gen`                               | генерация TS-кода из `proto/` в `gen/`          |
| `pnpm proto:lint`                              | `buf lint` для proto-контрактов                 |
| `pnpm mock`                                    | mock Connect RPC сервер                         |
| `pnpm typecheck` / `pnpm lint` / `pnpm format` | проверки и форматирование                       |

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
proto/                  контракты (*.proto), источник правды для API
gen/                    сгенерированный код protobuf-es (не редактировать руками)
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

## Тема и оболочка

Светлая оболочка в стиле kms-frontend: колонка разделов с lucide-иконками, панель меню с группами и пользователем внизу, топбар с крошками «раздел › страница», страницы `AppPage` с заголовком, подзаголовком и действиями.

- `theme.settings.js` — пресет PrimeVue: оранжевый `primary` и нейтральные `surface` для светлой и тёмной схем.
- `assets/scss/_base.scss` — токены `--app-*` (поверхности, границы, текст, акцент, размеры колонки и панели).
- `assets/scss/_cabinet.scss` — блоки `app-shell`, `app-rail`, `app-panel`, `app-nav`, `app-user`, `app-topbar`, `app-page`, `app-card`, `app-grid`, `app-field`, `app-form-grid`, `app-list`, `app-table`, `app-empty`, `app-settings-*`.
- `config/menu.config.ts` — разделы (`IMenuSection`: `lucide`-иконка и группы ссылок).
- Кнопка в топбаре сворачивает панель до колонки иконок; на мобильных сайдбар выезжает поверх контента.

Для таблиц используйте `UiTable` (toolbar, пустое состояние, «Показать ещё», выделение строк с панелью массовых действий). Пример — `components/users/List.vue`.

## Как добавить новый модуль

1. Описать сервис в `proto/<module>/v1/<module>.proto` и выполнить `pnpm proto:gen`.
2. Создать `core/<module>/api/<entity>/<entity>.api.ts` — класс `extends BaseRpcApi<typeof XxxService>`, и `<entity>.api.types.ts` с типами на основе `MessageInitShape`.
3. Создать `core/<module>/services/<entity>/` — сервис с маппингом в доменные интерфейсы (`MapperService`) и `*.types.ts`.
4. Создать `core/<module>/index.ts` с `<MODULE>_API()` и `<MODULE>_SERVICES()`; зарегистрировать сервисы в `composables/helpers/useServices.ts`.
5. Добавить стор в `core/<module>/stores/` (подхватывается автоимпортом).
6. Добавить composables, компоненты, страницы и раздел или пункты в `config/menu.config.ts`.

Пример полного модуля — `core/users` + `composables/users` + `components/users` + `pages/cabinet/settings/users`.
