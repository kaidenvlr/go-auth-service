# Auth Service

Сервис авторизации на Go с использованием PostgreSQL для хранения данных и Redis для кэширования.

## Технологии

- Go 1.21
- PostgreSQL 14
- Redis 7
- Docker & Docker Compose
- JWT для авторизации
- Gin Web Framework

## Структура проекта

```bash
.
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── README.md
├── internal
│ ├── config
│ │ └── config.go
│ ├── handler
│ │ └── auth_handler.go
│ ├── model
│ │ └── user.go
│ ├── repository
│ │ └── user_repository.go
│ └── service
│ └── auth_service.go
└── migrations
└── 001_create_users_table.sql
```

## Установка и запуск

### Предварительные требования

- Docker
- Docker Compose
- Go 1.21 (для локальной разработки)

### Запуск сервиса

1. Клонировать репозиторий:

```bash
git clone <repository-url>
cd auth-service
```

2. Создать файл `.env` с необходимыми переменными окружения:
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auth_db
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=your-secret-key
```

3. Запустить сервис:
```bash
docker-compose up --build
```

Сервис будет доступен по адресу: `http://localhost:8080`

## API Endpoints

### Регистрация

```http
POST /register
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

### Авторизация

```http
POST /login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

Успешный ответ:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

## Разработка

### Локальный запуск

1. Установить зависимости:
```bash
go mod download
```

2. Запустить базы данных:
```bash
docker-compose up postgres redis
```

3. Запустить сервис:
```bash
go run main.go
```

### Тестирование

Запуск тестов:
```bash
go test ./...
```

## Мониторинг и логирование

- Логи сервиса доступны через `docker logs`
- Все ошибки логируются с подробным описанием
- Метрики и статусы сервисов можно отслеживать через Docker Compose

## Безопасность

- Пароли хешируются с использованием bcrypt
- Используется JWT для авторизации
- Все чувствительные данные хранятся в переменных окружения
- Реализована защита от основных атак (SQL-инъекции, XSS)

## Кэширование

- Redis используется для кэширования данных пользователей
- Время жизни кэша - 30 минут
- Автоматическая инвалидация при обновлении данных
