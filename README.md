# URL Shortener

Сервис для сокращения ссылок. Тестовое задание Ozon (стажёр-разработчик).

- 10-символьные короткие ссылки из алфавита `[0-9A-Za-z_]` (63 символа).
- Одному оригинальному URL соответствует ровно одна короткая ссылка.
- Два хранилища на выбор при запуске: `postgres` и `local` (in-memory).

## API

### `POST /shorten` — сократить ссылку

Запрос:
```json
{ "url": "https://example.com/some/long/path" }
```

Ответы:
- `200 OK` — `{ "short_link": "aB3_kZ9xQ0" }`
- `400 Bad Request` — невалидное тело или невалидный URL: `{ "error": "..." }`
- `500 Internal Server Error` — ошибка на стороне сервиса

### `GET /:shortLink` — получить оригинальный URL

Ответы:
- `200 OK` — `{ "original_url": "https://example.com/..." }`
- `404 Not Found` — `{ "error": "short link not found" }`
- `500 Internal Server Error` — ошибка на стороне сервиса

### `GET /health` — проверка живости

Проверяет доступность хранилища (для `postgres` — пинг пула соединений).

- `200 OK` — `{ "status": "ok" }`
- `503 Service Unavailable` — хранилище недоступно

## Запуск

### 1. Через docker-compose (postgres)

```bash
cp .env.example .env
# при необходимости отредактировать .env, STORAGE=postgres
make build
```

`docker-compose` сам поднимет Postgres, дождётся его готовности, автоматически накатит
миграции (сервис `migrator`) и только потом запустит сервис. Отдельных шагов не нужно.

Сервис поднимется на `http://localhost:8080`.

### 2. Через docker-compose (in-memory)

В `.env` выставить `STORAGE=local` — Postgres в этом режиме не используется, но контейнер БД всё равно поднимется (можно его остановить: `docker compose stop db`).

### 3. Локально

```bash
export $(cat .env | xargs)
go run ./cmd
```

## Переменные окружения

| Переменная | Описание |
|-|-|
| `SERVER_PORT` | Порт HTTP-сервера |
| `STORAGE` | `postgres` или `local` |
| `DATABASE_HOST` | Хост Postgres (при `STORAGE=postgres`) |
| `DATABASE_PORT` | Порт Postgres |
| `DATABASE_NAME` | Имя БД |
| `DATABASE_USER` | Пользователь |
| `DATABASE_PASSWORD` | Пароль |

## Makefile

```bash
make up      # docker compose up -d
make build   # docker compose up -d --build
make down    # docker compose down
make clean   # down -v --rmi all
make test    # go test -v ./...
```

## Тесты

```bash
go test ./...
```

## Как генерируется короткая ссылка

Детерминированный FNV-1a хэш от `originalLink`, представленный в 10 символах base-63 (`[0-9A-Za-z_]`).

При коллизии (тот же `short` уже занят другим URL) в сервисе выполняется до 5 ретраев с дополнительной солью (`salt = 1..N`), чтобы получить другой хэш от того же URL. Уникальность гарантируется:
- в Postgres — `UNIQUE`-ограничением на `short_link` (нарушение вызывает `ErrLinkCollision`);
- в local — проверкой в map под мьютексом.

Идемпотентность на уровне БД: `INSERT ... ON CONFLICT (original_link) DO UPDATE ... RETURNING *` возвращает уже сохранённую запись, если такой оригинальный URL уже есть.
