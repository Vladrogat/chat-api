# Chat API

REST API для управления чатами и сообщениями на Go с использованием PostgreSQL.

## Описание проекта

Это тестовое задание представляет собой API для работы с чатами и сообщениями. Проект реализован на чистом Go с использованием стандартной библиотеки `net/http`, ORM GORM для работы с PostgreSQL, и системы миграций Goose.

### Архитектурные решения

1. **Модульная структура**: Код разделен на логические слои (handlers, database, domain, respository, service, config)
2. **Repository pattern**: Изоляция бизнес-логики от работы с БД
3. **Валидация на уровне моделей**: Централизованная валидация данных
4. **Middleware**: Переиспользуемое логирование запросов
5. **Каскадное удаление**: Реализовано на уровне БД через FOREIGN KEY с ON DELETE CASCADE

## Технологический стек

- **Go 1.25** - язык программирования
- **net/http** - стандартная библиотека для HTTP сервера
- **GORM** - ORM для работы с PostgreSQL
- **PostgreSQL 15** - база данных
- **Goose** - система миграций
- **Testify** - библиотека для тестирования
- **Docker & Docker Compose** - контейнеризация

## API Endpoints

### 1. Создать чат
```
POST /chats/
```

**Request Body:**
```json
{
  "title": "Название чата"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "title": "Название чата",
  "created_at": "2026-01-20T10:00:00Z"
}
```

### 2. Отправить сообщение
```
POST /chats/{id}/messages/
```

**Request Body:**
```json
{
  "text": "Текст сообщения"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "chat_id": 1,
  "text": "Текст сообщения",
  "created_at": "2026-01-20T10:01:00Z"
}
```

**Errors:**
- `404 Not Found` - если чат не существует

### 3. Получить чат с сообщениями
```
GET /chats/{id}?limit=20
```

**Query Parameters:**
- `limit` (optional) - количество последних сообщений (по умолчанию 20, максимум 100)

**Response (200 OK):**
```json
{
  "id": 1,
  "title": "Название чата",
  "created_at": "2026-01-20T10:00:00Z",
  "messages": [
    {
      "id": 1,
      "chat_id": 1,
      "text": "Первое сообщение",
      "created_at": "2026-01-20T10:01:00Z"
    },
    {
      "id": 2,
      "chat_id": 1,
      "text": "Второе сообщение",
      "created_at": "2026-01-20T10:02:00Z"
    }
  ]
}
```

**Note:** Сообщения отсортированы по `created_at` в хронологическом порядке.

### 4. Удалить чат
```
DELETE /chats/{id}
```

**Response:**
- `204 No Content` - успешное удаление
- `404 Not Found` - чат не найден

**Note:** Все сообщения чата удаляются каскадно.

## Запуск проекта

### Вариант 1: Docker Compose (рекомендуется)

1. Убедитесь, что у вас установлены Docker и Docker Compose

2. Запустите проект:
```bash
docker-compose up --build
```

3. API будет доступен по адресу: `http://localhost:8080`

4. Для остановки:
```bash
docker-compose down
```

5. Для полной очистки (включая данные БД):
```bash
docker-compose down -v
```

### Вариант 2: Локальный запуск

1. Установите PostgreSQL и создайте базу данных:
```sql
CREATE DATABASE chat_db;
```

2. Установите зависимости Go:
```bash
go mod download
```

3. Установите goose для миграций:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

4. Примените миграции:
```bash
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=chat_db sslmode=disable" up
```

5. Установите переменные окружения (опционально):
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=chat_db
export SERVER_PORT=8080
```

6. Запустите приложение:
```bash
go run cmd/api/main.go
```