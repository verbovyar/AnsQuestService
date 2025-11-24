# QA API Service

Сервис вопросов и ответов

## Стек

- Go, net/http
- GORM
- PostgreSQL
- goose (миграции)
- Docker, docker-compose

## Функциональность

- Создание, получение и удаление вопросов
- Добавление ответов к вопросам
- Получение и удаление ответов
- Каскадное удаление ответов при удалении вопроса

## Запуск

1. Собрать и запустить контейнеры:

```bash
docker-compose up --build

2. Эндпоинты

- API доступно на http://localhost:8080
- GET /questions/ — список всех вопросов.
- POST /questions/ — создать вопрос.
- GET /questions/{id} — получить вопрос и все его ответы.
- DELETE /questions/{id} — удалить вопрос (и каскадно ответы).
- POST /questions/{id}/answers/ — добавить ответ к вопросу.
- GET /answers/{id} — получить ответ.
- DELETE /answers/{id} — удалить ответ.

## Ручное тестирование API через curl

1. Создать вопрос

```bash
curl -X POST http://localhost:8080/questions/ \
  -H "Content-Type: application/json" \
  -d '{"text": "Первый вопрос в системе?"}'

Ожидаемый ответ (201 created)

```bash
{
  "id": 1,
  "text": "Первый вопрос в системе?",
  "created_at": "2025-11-24T10:00:00Z"
}

2. Получить список всех вопросов

```bash
curl http://localhost:8080/questions/

Ответ (200 ОК)

```bash
[
  {
    "id": 1,
    "text": "Первый вопрос в системе?",
    "created_at": "2025-11-24T10:00:00Z"
  }
]


3. Получить конкретный вопрос с ответами

```bash
curl http://localhost:8080/questions/1

Ответ (200 ОК). Если вопрос не найден — вернётся 404 Not Found

```bash
{
  "id": 1,
  "text": "Первый вопрос в системе?",
  "created_at": "2025-11-24T10:00:00Z",
  "answers": [
    {
      "id": 1,
      "question_id": 1,
      "user_id": "8b70201b-7c1c-4f5d-9b48-9eaa9c68c111",
      "text": "Мой первый ответ",
      "created_at": "2025-11-24T10:05:00Z"
    }
  ]
}

4. Добавить ответ к вопросу

```bash
curl -X POST http://localhost:8080/questions/1/answers/ \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "8b70201b-7c1c-4f5d-9b48-9eaa9c68c111",
    "text": "Вот мой ответ на вопрос №1"
  }'

Ответ (201 created). Если вопрос с таким id не существует, сервер вернёт 400 Bad Request с сообщением о том, что вопрос не найден

```bash
{
  "id": 1,
  "question_id": 1,
  "user_id": "8b70201b-7c1c-4f5d-9eaa9c68c111",
  "text": "Вот мой ответ на вопрос №1",
  "created_at": "2025-11-24T10:05:00Z"
}

5. Получить ответ по ID

```bash
curl http://localhost:8080/answers/1

Ответ (200 ОК). Если ответ не найден — 404 Not Found

```bash
{
  "id": 1,
  "question_id": 1,
  "user_id": "8b70201b-7c1c-4f5d-9b48-9eaa9c68c111",
  "text": "Вот мой ответ на вопрос №1",
  "created_at": "2025-11-24T10:05:00Z"
}

6. Удалить ответ

```bash
curl -X DELETE http://localhost:8080/answers/1 -v

Ответ (204 No Content)
