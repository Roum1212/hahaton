# Автоматическое развёртывание

Workflow `Deploy` запускается после изменения `main`, включая merge Pull Request,
и вручную через GitHub Actions → Deploy → Run workflow. Прямой push в `main` тоже запускает деплой.
Сборки идут на GitHub, на VPS передаются готовые образы `linux/amd64` с тегом полного SHA коммита.
Registry и доступ VPS к GitHub не требуются. Одновременно выполняется один деплой.

## Настройки GitHub

Repository → Settings → Secrets and variables → Actions:

| Вид | Имя | Значение |
| --- | --- | --- |
| Variable | `DEPLOY_HOST` | Адрес VPS |
| Variable | `DEPLOY_USER` | По умолчанию `root` |
| Variable | `DEPLOY_PORT` | По умолчанию `22` |
| Variable | `SITE_NAME` | Название сайта для сборки, как `SITE_NAME` в серверном `.env` |
| Variable | `SITE_DESCRIPTION` | Описание сайта для сборки, как `SITE_DESCRIPTION` в `.env` |
| Secret | `DEPLOY_SSH_KEY` | Приватный ключ отдельного доступа для деплоя |
| Secret | `DEPLOY_KNOWN_HOSTS` | Проверенная запись SSH host key VPS |

Приватный ключ приложения, пароль PostgreSQL и остальные значения серверного `.env`
в GitHub не передаются. Не используйте личный SSH-ключ разработчика для Actions.

## Подготовка VPS

Ожидается существующая установка Docker Compose, Python 3, `flock`, запущенный PostgreSQL,
конфигурация `/root/hahaton/.env` и исходный `/root/hahaton/docker-compose.yaml`.
Первый деплой может восстановить предыдущие приложения из этой установки.

Установите `deploy/receive.sh` как `/usr/local/sbin/hahaton-deploy`, владелец `root`, права `755`.
Для отдельного публичного ключа Actions добавьте в `/root/.ssh/authorized_keys` строку:

```text
restrict,command="/usr/local/sbin/hahaton-deploy" ssh-ed25519 <public-key> hahaton-github-actions
```

Этот ключ принимает только SHA коммита и архив релиза; интерактивная оболочка, SCP,
туннели и forwarding не используются. При изменении `receive.sh` установленный файл
обновляется отдельно. Workflow передаёт прикладной `deploy.sh` с каждым релизом.

## Что происходит при деплое

1. GitHub собирает backend и frontend, загружая открытый справочник городов для образа backend.
2. Образы, Compose и скрипт развёртывания передаются через SSH.
3. VPS сохраняет релиз в `/root/hahaton-deploy/releases/<SHA>` и загружает образы в Docker.
4. Проверяет Compose с существующим `/root/hahaton/.env` и готовность PostgreSQL.
5. Останавливает frontend и backend, сохраняет `database.dump` через `pg_dump`.
6. Выполняет `goose up` из нового образа backend.
7. Запускает приложения с новыми образами и ждёт успешных healthcheck.
8. Обновляет ссылку `/root/hahaton-deploy/current` на успешный релиз.

Nginx, сертификаты, `.env`, volume `hahaton_postgres-data` и `hahaton_storage-data`
сохраняются. Небольшой перерыв в доступности приходится на миграции и запуск приложений.
При ошибке сборки работающий сайт не затрагивается. При ошибке обновления после остановки
скрипт пробует поднять прежние образы. Успешно применённые миграции автоматически не откатываются:
старые приложения должны оставаться совместимыми с новой схемой.

Дампы содержат данные приложения, доступны только root и не уходят в GitHub.
Релизы, дампы и образы сохраняются на VPS; следите за свободным местом и удаляйте старые вручную,
оставляя текущий и предыдущий релиз. Build cache на VPS для CD не нужен.

Если миграция не прошла, исправьте её и повторите workflow. При ручном запуске выбирайте `main`:
деплой из других веток запрещён условием workflow. Автоматического отката БД и вызова `down -v` нет.
