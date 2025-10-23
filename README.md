# Shorter — Сервис коротких ссылок

Сервис на Go с поддержкой аналитики кликов.

## Функции
- Сокращение URL
- Редирект 302
- Аналитика: гео, устройство, браузер
- Метрики Prometheus
- Event-driven архитектура (Kafka)

## API
- `POST /api/v1/shorten` — создать ссылку
- `GET /api/v1/stats/{alias}` — статистика
- `GET /{alias}` — редирект
- `GET /metrics` — метрики Prometheus

## Запуск
```bash
docker-compose up --build