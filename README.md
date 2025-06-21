# Practical-assignment
HTTP API для управления I/O-bound задачами с in-memory хранилищем. Позволяет создавать задачи, получать их статус, удалять задачи и просматривать список задач. Каждая задача обрабатывается асинхронно (3-5 минут) и может находиться в одном из статусов: `pending`, `completed` или `failed`.

## Функциональность

- **Создание задачи**: `POST /tasks` — создаёт новую задачу и возвращает её ID.
- **Получение статуса задачи**: `GET /tasks/{id}` — возвращает информацию о задаче, включая статус, дату создания и продолжительность обработки.
- **Список задач**: `GET /tasks` — возвращает список всех задач.
- **Удаление задачи**: `DELETE /tasks/{id}` — удаляет задачу по ID.

Формат ответа для задач:
```json
{
  "id": "5bffba08-2439-46af-b7cd-cc8a229a11b7",
  "status": "pending | completed | failed",
  "created_at": "2025-06-21T01:38:32Z",
  "duration": "3m0s",
  "result": "Task completed",
  "completed_at": "2025-06-21T01:41:32Z"
}
```

## Требования

- Go 1.22 или выше
- Зависимости:
    - `github.com/gorilla/mux v1.8.1`
    - `github.com/google/uuid v1.6.0`

## Установка

1. Склонируйте репозиторий:
   ```bash
   git clone <repository-url>
   cd practical-assignment
   ```
2. Установите зависимости:
   ```bash
   go mod tidy
   ```

## Запуск

Запустите сервер:
```bash
go run main.go
```

Сервер будет доступен на `http://localhost:8080`.

## Примеры использования API

1. **Создать задачу**:
   ```bash
   curl -X POST http://localhost:8080/tasks
   ```
   Ответ:
   ```json
   {"id": "5bffba08-2439-46af-b7cd-cc8a229a11b7"}
   ```

2. **Получить статус задачи**:
   ```bash
   curl http://localhost:8080/tasks/5bffba08-2439-46af-b7cd-cc8a229a11b7
   ```
   Ответ (в процессе):
   ```json
   {
     "id": "5bffba08-2439-46af-b7cd-cc8a229a11b7",
     "status": "pending",
     "created_at": "2025-06-21T01:38:32Z"
   }
   ```
   Ответ (завершена):
   ```json
   {
     "id": "5bffba08-2439-46af-b7cd-cc8a229a11b7",
     "status": "completed",
     "created_at": "2025-06-21T01:38:32Z",
     "duration": "3m0s",
     "result": "Task completed successfully",
     "completed_at": "2025-06-21T01:41:32Z"
   }
   ```

3. **Получить список задач**:
   ```bash
   curl http://localhost:8080/tasks
   ```
   Ответ:
   ```json
   [
     {
       "id": "5bffba08-2439-46af-b7cd-cc8a229a11b7",
       "status": "completed",
       "created_at": "2025-06-21T01:38:32Z",
       "duration": "3m0s",
       "result": "Task completed successfully",
       "completed_at": "2025-06-21T01:41:32Z"
     }
   ]
   ```

4. **Удалить задачу**:
   ```bash
   curl -X DELETE http://localhost:8080/tasks/5bffba08-2439-46af-b7cd-cc8a229a11b7
   ```
   Ответ: `204 No Content`

5. **Ошибки**:
    - Невалидный ID:
      ```bash
      curl http://localhost:8080/tasks/invalid
      ```
      Ответ: `400 Bad Request`, `{"error": "invalid task ID"}`
    - Задача не найдена:
      ```bash
      curl http://localhost:8080/tasks/nonexistent
      ```
      Ответ: `404 Not Found`, `{"error": "task not found"}`

## Архитектура

Проект построен с использованием чистой архитектуры для упрощения поддержки и расширения:

- **`handler`**: HTTP-обработчики для REST API (`POST /tasks`, `GET /tasks`, etc.).
- **`service`**: Бизнес-логика, включая асинхронную обработку задач через горутины.
- **`storage`**: In-memory хранилище для задач.
- **`model`**: Определение структуры задачи (`Task`).

Ключевые компоненты:
- Интерфейс `TaskProcessor` для модульности и тестируемости.
- Асинхронная обработка задач с использованием `context.Background()` для независимости от HTTP-запросов.
- Потокобезопасное хранилище (предполагается использование мьютексов в `storage/memory.go`).

## Инженерные практики

- **Логирование**: Используется `log/slog` для структурированных логов с контекстом (`task_id`, `error`, `stage`).
- **Обработка ошибок**: HTTP-статусы (`400`, `404`, `500`) и JSON-ответы для ошибок.
- **Graceful Shutdown**: Сервер корректно завершает работу по SIGINT с таймаутом 5 секунд.
- **Модульность**: Интерфейсы и DI упрощают замену компонентов (например, хранилища).
- **Минимальные зависимости**: Только `gorilla/mux` и `uuid`.