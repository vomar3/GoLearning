# Order System

Учебный сервис заказов на Go: API принимает заказ, отправляет событие в Kafka, а отдельный worker читает событие и сохраняет заказ в Postgres.

## Конфигурация

Настройки читаются из env. Файл `.env.example` показывает пример значений для запуска через `docker-compose`.

## Запуск

Поднять Postgres, Kafka и Kafka UI:

```bash
docker compose up -d
```

Запустить API:

```bash
go run ./cmd/api
```

Запустить worker во втором терминале:

```bash
go run ./cmd/worker
```

## Что можно делать

Создать заказ:

```bash
curl -X POST http://localhost:8080/order -H "Content-Type: application/json" -d '{"id":"order-1","item":"book","price":1200}'
```

Получить заказ:

```bash
curl http://localhost:8080/orders/order-1
```

Удалить заказ:

```bash
curl -X DELETE http://localhost:8080/orders/order-1
```

Прогнать тесты:

```bash
go test ./...
```
