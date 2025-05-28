
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
