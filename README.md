# Цифровой волонтёр Обнинска

Docker Compose запускает три постоянных контейнера: Nuxt frontend, Go backend и PostgreSQL 17.
Nginx и HTTPS настраиваются отдельно на сервере. Приложения внутри контейнеров работают по HTTP;
по умолчанию наружу опубликованы только `127.0.0.1:3000` и `127.0.0.1:8000`.
PostgreSQL доступен только во внутренней сети Compose.

## Подготовка

Нужны Docker Engine и Docker Compose v2. Команды выполняются из корня репозитория.

```sh
cp .env.example .env
sh backend/scripts/generate-keys.sh
openssl rand -hex 32
sh backend/scripts/fetch-cities.sh
```

Перенесите три сгенерированных ключа в соответствующие поля `.env`, а последний случайный
секрет — в `POSTGRES_PASSWORD`. Укажите свои значения:

- `FRONTEND_URL` — публичный HTTPS-адрес frontend.
- `BACKEND_URL` — публичный HTTPS-адрес backend, доступный браузеру.
- `FLIDA_ISSUER` и `FLIDA_CLIENT_ID` — issuer и Client ID приложения Flida.
- `FLIDA_DISCOVERY_URL` — отдельный адрес discovery, если он находится не в корне issuer.
  Для текущей Flida: issuer `https://id.flida.ru`, discovery
  `https://id.flida.ru/api/.well-known/openid-configuration`.
- `FLIDA_REDIRECT_URI` — HTTPS-адрес `<FRONTEND_URL>/auth/callback`, разрешённый в Flida.
- `DADATA_API_KEY` и `YANDEX_MAPS_API_KEY` — если используете подсказки адресов и карту.

Не переносите ADMIN_API_KEY, ключи сессии/шифрования и пароль БД в публичные настройки frontend.
Compose передаёт frontend только публичные `NUXT_PUBLIC_*`. API URL и настройки входа меняются
при пересоздании контейнера, пересборка для них не требуется. Название и meta-description сайта
передаются также как build arguments; для их изменения пересоберите frontend.

Города загружаются в PostgreSQL при первом старте backend из подготовленного CSV.
Полигоны необязательны: если они нужны, положите их в `backend/data`, укажите путь
`/app/data/<имя файла>` в `CITY_CATALOG_BOUNDARIES_PATH` и пересоберите backend.

Код работы с cookie и RPC-ошибками перенесён из библиотек Flida во внутренние пакеты
`backend/internal/pkg/cookies` и `backend/internal/pkg/rpc-error`. Отдельных модулей и
`replace` нет; для Docker-сборки не нужны доступ к приватному GitLab Flida и дополнительные токены.

## Первый запуск и миграции

```sh
docker compose config --quiet
docker compose build frontend
docker compose build backend
docker compose up -d --wait postgres
docker compose run --rm --no-deps --entrypoint /app/goose backend up
docker compose up -d --wait
docker compose ps
```

Образы собираются отдельными командами, чтобы сборки не конкурировали за память на небольшой VPS.
Для VPS с 2 ГБ RAM нужен запас памяти на сборку, например swap.

Миграции выполняются вручную Goose из собранного образа; сервер их не применяет автоматически.
Команда `run --rm` создаёт временный контейнер только на время миграции. Постоянно работают три.
Параметры PostgreSQL передаются через `PGUSER`, `PGPASSWORD`, `PGDATABASE` и общий DSN внутри сети.
Goose получает эти же параметры, `GOOSE_DRIVER` и `GOOSE_MIGRATION_DIR`.

```sh
docker compose run --rm --no-deps --entrypoint /app/goose backend status
docker compose logs --tail=100 backend frontend postgres
```

PostgreSQL имеет healthcheck; backend запускается после готовности БД, frontend — после
готовности backend. Это штатный механизм [depends_on/service_healthy в Compose](https://docs.docker.com/compose/how-tos/startup-order/).
Первый вход требует доступной Flida. DADATA_API_KEY не обязателен для уже сохранённых точек.

## Nginx и административный API

Ваш Nginx должен направлять frontend на `http://127.0.0.1:3000`, backend — на
`http://127.0.0.1:8000`, включая `/storage/` и пути Connect RPC. Для чата нужен streaming:
не буферизуйте ответы backend и задайте достаточный read timeout.
Публичные URL должны быть HTTPS, поскольку сессионная cookie имеет `Secure; HttpOnly; SameSite=None`.
Не нужно передавать сертификаты в контейнеры при завершении TLS на Nginx.
Если Nginx работает не на хосте, настройте маршрутизацию и bind-адреса под его сеть.

AdminService расположен на том же адресе и порту, что остальные RPC, без отдельного 8081.
Он принимает только серверный административный токен; пользовательская cookie его не заменяет.
Пример обращения к локальному backend на самом сервере:

```sh
curl --fail-with-body http://127.0.0.1:8000/admin.v1.AdminService/FindVolunteers \
  -H "Authorization: Bearer $ADMIN_API_KEY" \
  -H 'Content-Type: application/json' \
  --data '{"email_address":"volunteer@example.com"}'
```

В этом примере ADMIN_API_KEY должен быть экспортирован в текущую оболочку.
Полный сценарий создания координатора и подготовки учебных данных — в [backend/README.md](backend/README.md).

## Сохранение данных и обновление

Именованные volumes `postgres-data` и `storage-data` хранят соответственно БД и изображения.
При стандартном имени проекта `hahaton` фактические имена — `hahaton_postgres-data` и
`hahaton_storage-data`. [Volumes живут отдельно от контейнеров](https://docs.docker.com/engine/storage/volumes/):
перезапуск, пересоздание контейнера и обычный `docker compose down` данные сохраняют.
`docker compose down -v` удаляет volumes вместе с данными — для обычной остановки его не используйте.
Имя Compose-проекта тоже сохраняйте: другое имя создаёт другой набор volumes.
Корневой Compose не переносит автоматически БД из старого dev-Compose в `backend/`.

Ключ DATA_ENCRYPTION_KEY сохраняйте вместе с данными: им зашифрованы сообщения и приглашения.
SESSION_TOKEN_V4_SYMMETRIC_KEY сохраняет возможность использовать выданные cookie после перезапуска.
Изменение POSTGRES_PASSWORD в `.env` не меняет пароль уже инициализированной БД — для этого нужен ALTER ROLE.

Обновление с ручной миграцией без работающих приложений:

```sh
docker compose build frontend
docker compose build backend
docker compose stop frontend backend
docker compose run --rm --no-deps --entrypoint /app/goose backend up
docker compose up -d --wait
```

Простая остановка и повторный запуск:

```sh
docker compose down
docker compose up -d --wait
```

Проверки backend, репозиториев и всех RPC описаны в [backend/README.md](backend/README.md).

## CD

После merge в `main` GitHub Actions собирает образы, передаёт их на VPS, сохраняет дамп БД,
применяет миграции Goose и обновляет приложения с проверкой готовности.
Серверный `.env`, PostgreSQL и изображения сохраняются. Настройка — в [deploy/README.md](deploy/README.md).
