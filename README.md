Для запуска бэкенда нужно предварительно установить docker и docker-compose.
```
docker-compose build
docker-compose up
```

### Шаблон env файла

```
PGUSER=
PGPASSWORD=
PGHOST=
PGPORT=
PGDATABASE=
PGSSLMODE=

ENVIRONMENT=debug
HTTP_PORT=8080
LOG_FILE_PATH=logs.log
SECRET_KEY=
```
