# Chat API

REST API для управления чатами и сообщениями на Go с использованием PostgreSQL.

## Описание проекта

Это тестовое задание представляет собой API для работы с чатами и сообщениями. Проект реализован на чистом Go с использованием стандартной библиотеки `net/http`, ORM GORM для работы с PostgreSQL, и системы миграций Goose.

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