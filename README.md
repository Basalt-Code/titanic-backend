Для запуска бэкенда нужно предварительно установить docker и docker-compose.
```
docker-compose build
docker-compose up
```


Для создания миграции использовать в терминале контейнера команду:
```
migrate create -ext sql -dir migrations <имя_Файла> 
```
Имя файла может быть, например, create_user, если в рамках этой миграции добавляется модель user

Для применения миграции использовать команду:
```
migrate -path migrations -database "postgres://postgres:password@db:5432/titanic_db?sslmode=disable" up
```
