
###Перед запуском

Настроить файлы конфигурации и переменные окружения:
`configs/config.yaml` 
`.env` 
или оставить настройки по умолчанию
Настройка в переменных окружения преоритетнее чем config.yaml

###Запуск

`docker compose up -d`
`make run`
или 
```bash
make up
```

логи с приложения печатаются в stdout
логи об обработанном запросе выводятся в stdout и в файл с меткой запуска в папке `log_folder: "./logs"` в файле конфигурации

Файлы миграции
`migrations/`

2 реализации memory и postgresql
при выборе `DB_TYPE=postgres` автоматически подтянутся миграции с помощью goose
для миграций нужно задать переменные окружения, в случае изменения параметров по умолчанию,
можно изменить в Makefile
```Makefile
GOOSE_MIGRATION_DIR?=./migrations
GOOSE_DRIVER?=postgres
GOOSE_DBSTRING?="postgres://calendar:calendar@localhost:5432/calendar"
```

в качестве роутера выбрал gorilla mux
для постгреса:  
sqlx + pgx

## Примеры запросов:

### POST /events  
URL: http://localhost:8080/events  
Body:
```json
    {
      "user_id":      "123",
      "title":        "example",
      "description":  "example",
      "starts_at":    "2025-06-20T20:00:00Z",
      "ends_at":      "2025-06-20T21:00:00Z",
      "notify_offset": 90000000000
    }
```


### PATCH /events/{event_id}  
URL: http://localhost:8080/events/{event_id}  
Body:
```json
    {
      "title":         "example_upd",
      "description":   "example_upd",
      "starts_at":     "2025-06-20T20:30:00Z",
      "ends_at":       "2025-06-20T21:30:00Z",
      "notify_offset": 90000000000
    }
```
### DELETE /events/{event_id}  
URL: http://localhost:8080/events/{event_id}

### GET /events/day/{user_id}?date=YYYY-MM-DD  
URL: http://localhost:8080/events/day/{user_id}?date=2025-06-20

### GET /events/week/{user_id}?date=YYYY-MM-DD  
URL: http://localhost:8080/events/week/{user_id}?date=2025-06-16

### GET /events/month/{user_id}?date=YYYY-MM-DD  
URL: http://localhost:8080/events/month/{user_id}?date=2025-06-01