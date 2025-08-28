# Order Service
Сервис для обработки заказов: принимает сообщения о заказах из Kafka, сохраняет в БД, кеширует и отдаёт через HTTP API.
# Запуск
```bash
$ docker-compose up
```
Сервис стартует на порту `8080`
# HTTP API
## Получение заказа по ID
```bash
GET /orders/:id
```

# Отправка сообщений в Kafka
```bash
$ make kafka-produce FILE=/ПУТЬ К ФАЙЛУ/
```
## Примеры сообщений
- [`model.json`](model.json) — корректный JSON заказа
- [`invalid.json`](invalid.json) — некорректный JSON (синтаксическая ошибка)
- [`invalid_fields.json`](invalid_fields.json) — JSON с неверными значениями полей (например, телефон, email, валюта)

# Валидация сообщений
Все JSON-сообщения проходят валидацию через `go-playground/validator`. Основные правила:

## Delivery:

- `name`, `phone`, `zip`, `city`, `address`, `region`, `email` - обязательные.

- `phone` должен быть в формате E.164 (+123456789).

- `email` должен быть корректным адресом.

## Payment:

- `transaction`, `currency`, `provider`, `amount`, `payment_dt`, `bank`, `delivery_cost`, `goods_total` - обязательные.

- `currency` - строка из 3 заглавных букв (например, USD).

- `amount`, `delivery_cost`, `goods_total` - положительные числа.

## Order:

- `order_uid`, `track_number`, `entry`, `locale`, `customer_id`, `delivery_service`, `shardkey`, `sm_id`, `date_created`, `oof_shard` - обязательные.

## Item:

- Все поля обязательны, кроме `id` и `order_id`.

- `sale` должен быть от 0 до 99.

- `price` и `total_price` - положительные числа.

