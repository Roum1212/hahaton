# Backend

Go backend проекта. Модуль объявлен в `go.mod`, сгенерированные ConnectRPC-пакеты находятся в `genproto/`.

Исполняемые приложения размещаются в `cmd/`, внутренняя реализация — в `internal/`.

Реализован модуль Session: агрегат, сервисы `app`, Connect RPC,
аутентификация сессии, PostgreSQL-репозиторий на pgx/Squirrel и миграция `migrations/00001_sessions.sql`.
Пользователь в сессии представлен общим `user_id`; роли в ней не хранятся.
Логика UUIDv7 перенесена из Flida в локальный пакет `internal/pkg/types/uuid-v7`.
ID сессии и пользователя используют этот пакет.

Каждая доменная модель размещается в собственной директории: `session-id`,
`user-id`. Поля приватные, чтение — через методы.
Конструкторы моделей вызывают `Validate`; доменные ошибки находятся в агрегате,
интерфейс репозитория и его ошибки — в `domain/aggregate/session/repository.go`.
Агрегат хранит отдельные значения и изменяется командами открытия и авторизации.
Сервис закрытия удаляет запись через репозиторий; повторное удаление успешно.
Состояние закрытия, дата закрытия и отдельный ключ сессии не хранятся.
Для каждого типа используется один вид receiver; имя receiver — `x`.
После типа идут его методы, затем обработчики команд; конструктор располагается в конце.
Доменные ошибки объявляются сразу после импортов.
Тексты ошибок пишутся в нижнем регистре, включая сокращения `id`, `uuid` и `sql`.
При обработке ошибок каждой ошибке соответствует отдельный `case`, даже если результат одинаков.
`NewInvalidArgumentError` всегда получает заполненный `BadRequest` с полем и причиной ошибки.
Входные параметры функций и методов всегда именуются, включая методы интерфейсов.
У каждого параметра явно указывается тип; в многострочных сигнатурах каждый параметр идёт отдельной строкой.
В литералах структур с несколькими полями каждое поле записывается на отдельной строке.
Сервисы принимают `ctx, r Request`; `Request.Validate` вызывается в начале сервиса.
Конструкторы запросов только собирают DTO. Поля запросов приватные, чтение — через геттеры.

Команды — контейнеры валидных входных данных с приватными полями и геттерами.
Их конструкторы только собирают значения, не вызывают `Validate` и не возвращают ошибку.
Агрегат принимает команды через публичные `Open` и `Authenticate`, без общего `Execute`.
Эти методы проверяют допустимость перехода из текущего состояния сессии.
Конструкторы агрегата принимают валидные модели, только собирают объект и не возвращают ошибку.
`NewSessionFromFields` восстанавливает агрегат из полей и размещается после `NewSession`.
У агрегата и команд нет метода `Validate`: этот метод используется у моделей и запросов сервисов.
SQL отвечает за сохранение и чтение, без бизнес-проверок `CHECK`.

Механизм context и cookie повторяет Flida: `authn.SetInfo/GetInfo`,
`go.flida.dev/cookies`, cookie `__Secure-SID` и шифрование PASETO v4.local.
В cookie находится защищённый `session_id`, отдельного ключа сессии нет.
PASETO проверяет `iat`, `nbf` и `exp` относительно текущего времени при каждом чтении.
Работа с токенами перенесена локально из Flida: `internal/pkg/token` содержит generic-интерфейсы
`Service[T]`, `Encryptor[T]`, `Decryptor[T]` и модель `EncryptedToken` с приватным значением.
`internal/pkg/token/paseto` реализует шифрование; `internal/pkg/session/token` определяет
payload с ID сессии и датами, преобразование в PASETO и фабрику `PASETO`.
В cookie передаётся `EncryptedToken.String()`, при чтении — `NewEncryptedToken(cookie.Value)`.
В context передаётся `*aggregate_session.Session`, без отдельной проекции и CQRS. Состояние и срок сессии
проверяются по PostgreSQL через `Authentication` в `internal/pkg/session/authentication/interceptor`.
Как во Flida, его `Authenticate(ctx, r *http.Request) (any, error)` передаётся в `authn.NewMiddleware`.
Middleware помещает возвращённый агрегат в context; RPC проверяет наличие применимой сессии.
Режимы RPC перечислены в `authModeByProcedure`: анонимная, авторизованная или любая действующая сессия.
Open/Close допускают отсутствие сессии; Get требует действующую, Authenticate — анонимную,
UploadImage — авторизованную. Неизвестные RPC и AdminService через этот listener получают NotFound.
Раздача `/storage/` публичная. Проверки `Origin`, Bearer-входа и фонового опроса сессий нет.
Middleware проверяет стрим при открытии; дальнейший отзыв доступа реализуется в будущем чат-сервисе.

API следует структуре Flida: фасад `Service`, отдельный `RPC` для каждой ручки и пакет `serde`.
Сессионный агрегат доступен через `session_context.GetSession`; RPC преобразует ошибки своего сервиса.
`OpenSession` и `AuthenticateSession` возвращают `google.protobuf.Empty` и устанавливают cookie.
`OpenSession`, `GetSession` и `CloseSession` принимают `google.protobuf.Empty`.
Их `Request` в `app` также пустые; сервисы сами получают текущий агрегат из context.
Повторно загружать сессию из БД для Open/Get не требуется: её уже проверил Authentication.
Open/Authenticate возвращают `token.EncryptedToken`; RPC передаёт его `.String()` в `cookies.Service.SetCookie`.
Пакет `session/cookie` содержит только имя `__Secure-SID`, обёртки Manager нет.
Действующая сессия при `OpenSession` сохраняется. Данные кабинета читает `GetSession`;
наличие профиля координатора вычисляется заново, срок сессии не продлевается.
Признак профиля координатора находится в DTO `service_get_session.Response`, а не в агрегате Session.
`Request` и `Response` каждого сервиса размещаются в его `dto.go`.
Интерфейсы зависимостей и ошибки объявляются в сервисах, которые их используют; общего `service_session` нет.

Авторизация создаёт агрегат с новым `session_id`, а прежнюю анонимную запись удаляет
в одной транзакции PostgreSQL через `Client.InTx(ctx, func(ctx context.Context) error)`.
Сетевая проверка личности выполняется до транзакции. Повторной загрузки сессии и `FOR UPDATE` нет:
для демо не реализуются специальные сценарии параллельных запросов.
Репозиторий выполняет только `INSERT`, `SELECT`, `UPDATE`, `DELETE` и преобразование данных.
Решения о допустимости входа, сроке сессии и удалении записи принимаются в `app`.

Клиент `internal/infra/client/postgresql` оборачивает `pgxpool.Pool`.
Его `Exec`, `Query`, `QueryRow` используют транзакцию из context, если она есть, иначе пул.
`InTx` передаёт транзакцию через context; commit и rollback выполняет `pgx.BeginFunc`.
При сборке зависимостей один клиент передаётся в `NewSessions` и в сервис авторизации как `Transactor`.
Внутри callback репозиторий вызывается с полученным context. Вложенные транзакции не используются.

Регистрация пользователей и проверка Flida пока представлены портом `service_authenticate_session.Authenticator`.
Его реализация должна проверять подпись, `iss`, `aud`, `exp` и применимые OIDC claims,
находить участника по `(iss, sub)` или регистрировать его вместе с нулевым кошельком.
Неполный профиль при первой регистрации возвращает `ErrIncompleteProfile`, неверный токен — `ErrInvalidIDToken`.
Наличие координаторского профиля предоставляет порт `service_get_session.CoordinatorReader`.
Другие агрегаты и запуск сервера пока не реализованы.

Собранный API регистрируется через сгенерированный `session_v1connect.NewSessionServiceHandler`.
Как во Flida, HTTP router оборачивается в библиотечный middleware:

```go
parser := paseto.MakeParser([]paseto.Rule{paseto_rules.ValidNow()})
tokens := session_token.PASETO(key, parser, sessionTTL)
authentication := session_authentication_interceptor.NewAuthentication(cookiesService, repository, tokens)
handler := authn.NewMiddleware(authentication.Authenticate).Wrap(mux)
```

`paseto_rules` — локальный пакет `internal/pkg/paseto`, как `pkg/paseto` во Flida.
Один `token.Service[session_token.Token]` передаётся сервисам Open/Authenticate как Encryptor,
а Authentication — как Decryptor. Payload берёт даты из агрегата: TTL общего PASETO-сервиса
не продлевает существующую сессию.

Для Open/Authenticate передаётся `cookies.Service` с `WithHTTPOnly`, `WithSecure`, `WithPath("/")`,
`WithSameSite` и `WithMaxAge` из конфигурации. Срок сессии всё равно ограничивается expires_at
в PostgreSQL и PASETO: повторный Open не продлевает серверный срок.
Close получает отдельный `expiredCookiesService` с теми же параметрами имени, Domain и Path,
но `WithMaxAge(-time.Second)`, чтобы удалить cookie через `SetCookie(header, "")`.

## Изображения

Добавлен `ImageService.UploadImage`: авторизованный участник загружает JPEG/PNG до 5 МиБ
и до 4096 пикселей по каждой стороне. `ImageData` проверяет содержимое и размеры;
`ImagePath` содержит UUIDv7 владельца и файла, расширение определяется сервером.
`service_upload_image` сохраняет файл через порт `Storage`, его локальный адаптер —
`internal/infra/client/image-storage`. Новых агрегатов и таблиц для файлов нет.

Папка задаётся в `client_image_storage.NewClient("./storage")`. Файлы размещаются
в `storage/{user_id}/{image_id}.jpg` или `.png`. Папка исключена из Git и должна
сохраняться между запусками; в контейнере ей нужен volume.

При сборке HTTP listener зарегистрировать сгенерированный `NewImageServiceHandler`
с `connect.WithReadMaxBytes(8 << 20)` на router, обёрнутом в общий `authn.NewMiddleware`.
Зависимости: локальный клиент → `service_upload_image.NewService` → `NewRPC` → API `NewService`.
Публичную раздачу той же папки подключить к mux:

```go
mux.Handle("/storage/", http.StripPrefix("/storage/", http.FileServer(http.Dir("./storage"))))
```

Ответ загрузки содержит `image_path` для сохранения в карточке и `image_url` вида `/storage/...`
для отображения. Новые контракты Volunteer/Coordinator/Reward принимают путь своей загруженной картинки,
а в ответах возвращают вычисляемые `avatar_url`/`image_url`.
Сами агрегаты и сервисы карточек ещё не реализованы. При их реализации нужно проверять владельца
из `ImagePath.GetUserID()` и существование файла перед сохранением привязки.
Замена/очистка привязки файл не удаляет; фоновой очистки в демо нет.
Файлы публичные: для демо используются только учебные изображения.

Фронтовой helper: `frontend/utils/helpers/image.helper.ts`. Он отправляет байты с cookie,
возвращает путь для Create/Update и полный URL картинки, учитывая отдельный origin бэкенда.

Проверка сборки: `go build ./...`. Форматирование и линтер: `golangci-lint fmt` и
`golangci-lint run ./...`. Конфигурация взята из Flida и адаптирована под модуль
и ограничение строк в 120 символов. Unit-тесты пока не добавляются.
Смешанные pointer/value receivers проверяет `recvcheck`.
Имена параметров интерфейсов проверяет `inamedparam`.
