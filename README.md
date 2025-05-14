# Shorty Reborn

Production‑сервис для сокращения ссылок на Go по принципам чистой архитектуры.

## Возможности

* Генерация кастомного или случайного безопасного alias (6 символов)
* Перенаправление с заголовками `Cache-Control` для контроля кэша
* Удаление сокращённых ссылок
* Структурированное логирование через Zap (консоль или JSON)
* Автоматические миграции базы данных при старте
* Настройка через YAML и переменные окружения
* Чистая многослойная архитектура: handlers, service, repository
* Юнит‑тесты с моками (GoMock, sqlmock)

## Архитектура проекта

```
├── cmd/url-shortener         # Точка входа приложения
├── config/local.yaml        # Конфигурация по умолчанию
├── internal
│   ├── config               # Загрузка конфигурации (cleanenv)
│   ├── database             # Подключение к БД, миграции, репозиторий
│   ├── handlers             # HTTP‑хендлеры (Chi)
│   ├── httpserver           # Настройка router, middleware, server
│   ├── logger               # Инициализация Zap logger
│   └── service              # Бизнес‑логика
├── go.mod                   # Модуль Go 1.24
└── go.sum                   # Контроль версий зависимостей
```

**Слои:**

1. **Handlers** — парсинг HTTP‑запросов, вызов сервисного слоя, отправка JSON или redirect.
2. **Service** — валидация, бизнес‑правила, координация работы с репозиторием.
3. **Repository** — CRUD‑операции в Postgres через SQLX, миграции.

## Стек технологий

* Go 1.24
* Chi (маршрутизация, middleware)
* Zap (логирование)
* SQLX + PGX (DB‑доступ)
* Cleanenv (конфигурация)
* GoMock, Testify, Go‑sqlmock (тестирование)

## Быстрый запуск

### Требования

* Go ≥ 1.24
* PostgreSQL ≥ 12

### Установка

```bash
git clone https://github.com/finlleyl/shorty_reborn.git
cd shorty_reborn
```

### Конфигурация

Скопируйте `config/local.yaml` и при необходимости измените:

```yaml
env: "local"
http_server:
  address: "0.0.0.0:8080"
  timeout: 4s
  idle_timeout: 60s
database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  name: "shorty"
  ssl_mode: "disable"
  timeout: 5s
```

Или переопределите через переменные окружения (`CONFIG_PATH`, `DB_HOST`, `DB_USER` и др.).

### Запуск локально

```bash
export CONFIG_PATH=$(pwd)/config/local.yaml
go run cmd/url-shortener/main.go
```

Сервис выполнит миграции и стартует на `http://localhost:8080`.

## Развёртывание

### Docker

1. Постройте образ:

   ```bash
   docker build -t shorty_reborn .
   ```
2. Запустите контейнер:

   ```bash
   docker run -d \
     -p 8080:8080 \
     -e CONFIG_PATH=/app/config/local.yaml \
     -e DB_HOST=<host> -e DB_USER=<user> -e DB_PASSWORD=<pass> \
     --name shorty shorty_reborn
   ```

### Docker Compose

Создайте `docker-compose.yml`:

```yaml
version: '3.8'
services:
  db:
    image: postgres:14
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: shorty
    volumes:
      - pgdata:/var/lib/postgresql/data
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      CONFIG_PATH: /app/config/local.yaml
      DB_HOST: db
      DB_USER: postgres
      DB_PASSWORD: postgres
    depends_on:
      - db
volumes:
  pgdata:
```

Запуск:

```bash
docker-compose up -d
```

### Kubernetes (пример)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: shorty-deployment
spec:
  replicas: 2
  selector:
    matchLabels:
      app: shorty
  template:
    metadata:
      labels:
        app: shorty
    spec:
      containers:
      - name: api
        image: shorty_reborn:latest
        ports:
        - containerPort: 8080
        env:
        - name: CONFIG_PATH
          value: /app/config/local.yaml
        - name: DB_HOST
          value: postgres.default.svc.cluster.local
        - name: DB_USER
          value: postgres
        - name: DB_PASSWORD
          value: postgres
---
apiVersion: v1
kind: Service
metadata:
  name: shorty-service
spec:
  selector:
    app: shorty
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

## Использование API

* **Создать ссылку**

  ```bash
  curl -X POST http://localhost:8080/api/urls \
    -H "Content-Type: application/json" \
    -d '{"url":"https://example.com","alias":"myalias"}'
  ```

  Ответ:

  ```json
  {
    "id":1,
    "alias":"myalias",
    "url":"https://example.com"
  }
  ```

* **Перенаправление**

  ```bash
  curl -v http://localhost:8080/api/urls/myalias
  ```

  Вернёт 302 с `Location: https://example.com` и `Cache-Control: public, max-age=60`.

* **Удаление**

  ```bash
  curl -X DELETE http://localhost:8080/api/urls/myalias
  ```

  Вернёт 204 No Content.

## Тестирование

Запуск всех юнит‑тестов:

```bash
go test ./internal/...
```

Тесты охватывают:

* Репозиторий (sqlmock)
* Сервисный слой (GoMock, Testify)
