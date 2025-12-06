# Health Map Backend

Бэкенд-сервис, который обслуживает платформу «Health Map» — публичный каталог здоровой / нездоровой инфраструктуры с отзывами, рейтингами и элементами геймификации, создаваемыми пользователями.

## Обзор
- Пользователи регистрируются/входят через JWT, управляют профилем и видят лидерборды.
- Места можно искать, просматривать, а также создавать/обновлять/удалять (только админ) и оставлять отзывы.
- Отзывы включают текст, фотографии, детальные метрики рейтинга и комментарии.
- Геймификация начисляет очки/уровни/ачивки за действия.
- REST API полностью описан в `docs/openapi.yaml`.

## Технологии
- Go 1.23, Gin, pgx.
- PostgreSQL (SQL-миграции в `migrations/`).
- JWT (github.com/golang-jwt/jwt/v5), bcrypt для хеширования паролей.
- OpenAPI 3.0 для contract-first разработки API.

## Структура проекта
```
cmd/server/              # точка входа
docs/openapi.yaml        # контракт API
internal/
  api/                   # DTO, хендлеры, middleware, роутер
  config/                # конфиги (JWT и т. д.)
  domain/                # сущности, ошибки, интерфейсы репозиториев/сервисов
  repository/postgres/   # репозитории на pgx
  services/              # бизнес-логика (users, places, reviews, achievements)
  pkg/auth/              # JWT + вспомогательные функции для паролей
migrations/              # SQL-миграции (формат golang-migrate)
```

## Настройка и запуск
1. Установите зависимости: `go mod download`.
2. Создайте `.env` (см. `.env.example`) с `POSTGRES_DSN`, `JWT_SECRET` и др.
3. Примените миграции (`migrate -path migrations ...` или любым совместимым инструментом).
4. Запустите сервис: `go run cmd/server/main.go` (или `docker-compose up --build`).
5. Базовый URL по умолчанию: `http://localhost:8080`.

## Ключевые эндпоинты API
- Auth: `POST /api/v1/auth/register`, `POST /api/v1/auth/login`.
- Users: профиль (`GET/PUT /api/v1/users/profile`), смена пароля, лидерборд, достижения.
- Places: поиск, объекты рядом, получение по ID, а также админские create/update/delete.
- Reviews: CRUD, список по месту, список по пользователю, комментарии к отзывам.
- Achievements: `GET /api/v1/achievements` (активный каталог) и `/api/v1/users/achievements` (достижения пользователя).

Подробности и коды ошибок см. в `docs/openapi.yaml`.

## Как начисляются очки и достижения
1. Действия начисляют очки через `UserService.AddPoints`.
2. Каждое действие генерирует `AchievementEvent` с типом и метриками.
3. `AchievementService.GrantIfEligible` сопоставляет события с активными достижениями по JSON-условиям и записывает новые значки.
4. Пользователь запрашивает разблокированные достижения через `/api/v1/users/achievements`.

## Безопасность
- Все защищённые эндпоинты требуют JWT (`Authorization: Bearer`).
- Пароли хранятся в виде bcrypt-хешей.
- Админские маршруты должны защищаться middleware с проверкой роли (роль хранится в токене).
- CORS настраивается в `internal/api/routes.go`.

## Docker Compose
- Сборка и запуск всего стека (`api` + `postgres + авто-миграции`):
  ```bash
  docker-compose up --build
  ```
- Остановка и очистка контейнеров:
  ```bash
  docker-compose down -v
  ```

## Лицензия
MIT License.
