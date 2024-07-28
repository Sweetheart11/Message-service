# Микросервис обработки сообщений

Приложение получает сообщение по запросу к REST API, сохраняет его в бд, затем отправляет в Kafka, где симулируется обработка запроса. После чего сообщение будет отмечано как обработанное.

Есть возможность запросить статистику по обработанным сообщениям.

## Запуск приложения

```console
docker-compose up
```

## Тестирование работы приложения

Тестировать работу можно по средствам, например, curl или postman.

Для просмотра статистики
```console
curl http://localhost:8080/message
```

Для добавления нового сообщения(Для postman еще нужно настроить Basic Auth)
```console
curl -u user:pass \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{"message":"test_message"}' \
  http://localhost:8080/message
```

Пример .env файла находмтся в .env_example:

```console
ENV=local
HTTP_SERVER_ADDRESS=localhost:8080
HTTP_SERVER_USER=user
HTTP_SERVER_PASSWORD=pass

DB_USER=postgres
DB_PASSWORD=example
DB_HOST=localhost
DB_PORT=5432
DB_SSLMODE=disable
DB_NAME=postgres
```