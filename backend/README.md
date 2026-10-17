# Backend

Go-монолит с ConnectRPC, PostgreSQL, pgx и Squirrel. Реализованы все методы текущих proto:
сессии и вход через Flida, профили, административное создание координатора,
сотрудничество и приглашения, мероприятия и заявки, кошелёк, поощрения и заказы,
чат, геопоиск и загрузка изображений.

Правила предметной области: [domain-design](../docs/domain-design.md).
Состояние реализации и проверки: [backend-implementation](../docs/backend-implementation.md).

## Запуск

Из каталога `backend`:

```sh
cp .env.example .env
sh scripts/generate-keys.sh
sh scripts/fetch-cities.sh
docker compose up -d postgres
```

Перенесите сгенерированные `SESSION_TOKEN_V4_SYMMETRIC_KEY`, `DATA_ENCRYPTION_KEY`, `ADMIN_API_KEY` в `.env`.
Укажите `FLIDA_ISSUER` и `FLIDA_CLIENT_ID` своего приложения. Flida должна отдавать
OIDC discovery/JWKS и ID token с данными профиля. Ключи DaData и администратора
остаются только на сервере.

Конфигурация находится в `internal/config`, по подходу Flida: `caarlos0/env`,
разделы с `envPrefix`, `ParseOptions` для типизированных ключей PASETO/AES.
Переменные используют префиксы `POSTGRESQL_`, `HTTP_SERVER_`, `ADMIN_HTTP_SERVER_`,
`FLIDA_`, `SESSION_TOKEN_`, `STORAGE_`, `CITY_CATALOG_`, `DADATA_`; актуальные имена — в `.env.example`.
Поля Config и вложенных секций публичные, теги `env`/`envPrefix` стоят прямо на них.
Точки входа заполняют Config через `env.ParseWithOptions(&cfg, config.ParseOptions())`,
затем вызывают Validate. Промежуточного environment, геттеров и NewConfig нет.
Сервер ожидает переменные окружения; `.env` сам не читает:

```sh
set -a
. ./.env
set +a
goose -dir migrations postgres "$POSTGRESQL_DSN" up
go run ./cmd
```

Миграции применяются вручную через установленный goose, до запуска сервера или seed:

```sh
goose -dir migrations postgres "$POSTGRESQL_DSN" up
goose -dir migrations postgres "$POSTGRESQL_DSN" status
```

Одна таблица — одна миграция; версии хранит `goose_db_version`.
Сервер и seed не меняют схему автоматически. Рабочие запросы выполняются через pgx/Squirrel.

Клиентский listener по умолчанию `:8080`, административный — `127.0.0.1:8081`.
`GET /healthz` возвращает успешный ответ после сборки приложения.
Для браузера нужны HTTPS и `credentials: 'include'`: cookie `__Secure-SID` всегда Secure/HttpOnly.
TLS можно завершать на прокси либо задать `HTTP_SERVER_TLS_CERTIFICATE` и `HTTP_SERVER_TLS_KEY` самому серверу.
`FRONTEND_URL` задаёт адрес фронта для ссылок приглашения.
CORS разрешает любой origin с credentials: сервер отражает origin запроса;
методы и заголовки Connect/gRPC-Web берутся из `connectrpc.com/cors`.
Cookie имеет `SameSite=None; Secure; HttpOnly`, на фронте нужен `credentials: 'include'`.
Клиентский listener не предоставляет AdminService.

Для контейнерного запуска заполните `.env`, загрузите справочник и выполните:

```sh
docker compose --profile app up --build -d
```

Серверный образ собирается из этого каталога. PostgreSQL и storage используют volumes;
`data` подключается как каталог только для чтения. HTTPS обеспечьте прокси перед приложением.

## Первый координатор и демонстрационные данные

Сначала два учебных участника входят через Flida. Вызовите AdminService.FindVolunteers
по ID, телефону или почте и выберите ID будущего координатора. Служебные запросы требуют
`Authorization: Bearer <ADMIN_API_KEY>` на административном listener:

```sh
curl --fail-with-body http://127.0.0.1:8081/admin.v1.AdminService/CreateCoordinator \
  -H "Authorization: Bearer $ADMIN_API_KEY" \
  -H 'Content-Type: application/json' \
  --data '{"volunteer_id":"<ID>","cooperation_request_cooldown_seconds":"60"}'
```

Для подготовки защиты:

```sh
go run ./cmd seed -coordinator-id '<ID координатора>' -volunteer-id '<ID другого участника>'
```

Seed использует существующие профили, создаёт активное сотрудничество, три учебных
мероприятия и три поощрения по 20 коинов. Для прошедшего мероприятия уже есть принятая
заявка: можно подтвердить 180 минут и купить поощрение. Повторный seed не сбрасывает баланс,
результаты и историю. Подготовленное прошлое событие не выдавайте за событие, созданное на защите:
API завершает мероприятия только после `ends_at`.

## Сборка приложения

`cmd/main.go` — единая точка входа. В `cmd/hahaton` сборка разделена как во Flida:
`init_database_clients`, `init_external_clients`, `init_repositories`, `init_tokens`,
`init_support_services`, `init_services`, `init_rpc_services`, `init_http_handlers`,
`init_http_server`. Сервисы приложения не получают зависимости друг на друга.
Команда `go run ./cmd seed ...` готовит только учебные данные и не применяет миграции.

## Слои и транзакции

- `api`: фасады сервисов, отдельные RPC, serde и преобразование ошибок.
- `app`: сервисы `ctx, r Request`, проверка запросов и координация агрегатов.
- `domain/aggregate`: состояние, команды, репозиторные контракты и доменные ошибки.
- `domain/model`: отдельная директория на модель, приватные поля, конструктор с Validate.
- `infra`: PostgreSQL-репозитории, Flida/DaData, каталог и локальное хранилище файлов.

Optional-поля app-запросов оформлены по образцу Flida: RPC собирает `RequestOption` через `With…`,
DTO хранит `Optional[T]`, сервис проверяет значение и presence через `Get()`.
Отсутствие поля, пустая строка и явный ноль различаются и при обновлениях, и в фильтрах.

Сервисы не вызывают другие сервисы. Последовательность регистрации и авторизации
собирается в RPC. Репозиторий передаёт callback транзакции; вложенные обращения используют
тот же pgx transaction через context. SQL только читает и сохраняет состояние.
Версий агрегатов, event sourcing, CQRS и специальных механизмов редких гонок нет.

Сессия в context — агрегат Session. Authentication из `internal/pkg/session/authentication/interceptor`
подключён через `authn.NewMiddleware`, как во Flida. ID находится в зашифрованной PASETO-cookie.
Open/Authenticate только устанавливают cookie, Close удаляет текущую сессию из БД.
Volunteer и Coordinator используют общий ID, а наличие координаторского профиля проверяется отдельно.

## Города, изображения и чат

`ListCities`/`GetCity` читают таблицу `cities`. На пустой БД справочник импортируется
из `CITY_CATALOG_CSV_PATH` с версией `CITY_CATALOG_VERSION`; после заполнения файл не нужен для чтения.
DaData используется отдельно для адресных подсказок. Частичный ответ подсказок
не считается полным справочником. Формат CSV и необязательных границ: [data/README](data/README.md).
Без проверенного полигона поиск по городу возвращает FailedPrecondition; клиент явно выбирает радиус.

UploadImage принимает JPEG/PNG до 5 МиБ и 4096 пикселей на сторону. Ответ содержит
`image_path` для сохранения в карточке и `image_url` для показа. Файлы размещаются в
`storage/{user_id}/{image_id}.jpg|png`, читаются публично через `/storage/`.
Профиль, мероприятие и поощрение проверяют владельца и наличие файла.
Пустой optional-путь очищает картинку, отсутствие поля сохраняет прежнюю.
Helper фронта: `frontend/utils/helpers/image.helper.ts`; привязка полей выбора к формам остаётся задачей UI.

Сообщения чата и секреты приглашений шифруются AES-GCM ключом `DATA_ENCRYPTION_KEY`. Храните ключ между перезапусками:
его замена без переноса данных сделает старую переписку и приглашения нечитаемыми.
Стрим проверяет новые сообщения и право доступа каждые 500 мс. После закрытия/истечения
сессии или отмены участия доступ отзывается; после завершения события остаётся чтение.
При первом и повторном вступлении видны только новые сообщения.

## Проверки и ограничения демо

Интеграционные тесты репозиториев устроены по образцу Flida: testify/suite и require,
`tests/integration/repository/<агрегат>/postgresql`, контейнерный helper в
`pkg/testcontainers/postgresql`. Нужен запущенный Docker; отдельный PostgreSQL/DSN не требуется.
Каждый из 12 suites поднимает PostgreSQL 17, применяет Goose и очищает свои таблицы перед тестом.
Контейнер удаляется автоматически. Рабочая БД не используется.

```sh
go test -tags=integration -count=1 -p 2 ./tests/integration/repository/...
go test -tags=integration -count=1 -p 2 \
  -coverpkg=./internal/infra/repository/... -coverprofile=repository.cover \
  ./tests/integration/repository/...
go tool cover -func=repository.cover
```


```sh
go build ./...
golangci-lint fmt
golangci-lint run ./...
```

Линтер перенесён из Flida, длина строк — 120. Unit-тесты по договорённости не добавлены.
Для воспроизводимой проверки всех 58 RPC через HTTP с PostgreSQL:

```sh
INTEGRATION_POSTGRESQL_DSN="$POSTGRESQL_DSN" go test -count=1 -tags=integration ./integration
```

Проверки создают отдельную временную схему и удаляют её после завершения;
к рабочим таблицам не обращаются. OIDC/JWKS и DaData подменяются учебными серверами/ответами.
Без `INTEGRATION_POSTGRESQL_DSN` сценарии пропускаются. Это интеграционные проверки;
unit-тесты не добавлены.

На отдельной PostgreSQL 17 проверены goose Up/Down/повторный Up, полный сценарий MVP,
повторы начислений и покупок, нехватка баланса, возвраты, история выданного,
каскады прекращения сотрудничества, доступ к чату, шифрование и cookie-аутентификация Connect JSON.
Отдельно проверены OIDC discovery/JWKS с подписанным учебным токеном, загрузка и привязка
изображения, admin-изоляция, полный импорт 1103 записей CSV и чтение без исходного файла,
Connect-стрим с историей, успешный EOF и закрытие при отзыве сессии.

Списки для демо фильтруются в памяти перед пагинацией; геометрия проверяется без PostGIS.
Для нескольких экземпляров потребуются общее файловое хранилище и пересмотр синхронизации чата.
Локальная демонстрация не требует DaData для уже сохранённых точек и справочника,
но первый вход требует доступной Flida. Реальных платежей и призов нет.
