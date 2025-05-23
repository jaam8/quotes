# Quotes

**Quotes** — REST API сервис для управления цитатами.  
Возможности:
- Создание цитаты
- Получение случайной цитаты
- Получение цитат по автору
- Получение всех цитат
- Удаление цитаты

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/jaam8/quotes.git
cd quotes
```

2. Скопируйте файл `.env.example` в `.env`, по желанию измените переменные окружения:
```bash
cp .env.example .env
```

3. Установите зависимости:
```bash
go mod download
```

4. Запустите проект:
```bash
go run cmd/main.go
```

#### По умолчанию сервис будет доступен по адресу: [`http://localhost:8080`](http://localhost:8080)

## Тестирование
Для запуска тестов используйте команду:

```bash
go test ./...
```

## Примеры и эндпоинты

### `POST /quotes`
Создание новой цитаты.
```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
```

#### Пример успешного ответа:
```json
{"id": 1}
```

### `GET /quotes/random`
Получение случайной цитаты.
```bash
curl http://localhost:8080/quotes/random
```

#### Пример успешного ответа:
```json
{
  "id": 1,
  "author": "Confucius",
  "quote": "Life is simple, but we insist on making it complicated."
}
```

### `GET /quotes`

Получение всех цитат.
```bash
curl http://localhost:8080/quotes
```

#### Пример успешного ответа:
```json
[
  {
    "id": 1,
    "author": "Confucius",
    "quote": "Life is simple, but we insist on making it complicated."
  },
  {
    "id": 2,
    "author": "Albert Einstein",
    "quote": "Life is like riding a bicycle. To keep your balance you must keep moving."
  }
]
```

### `GET /quotes?author={author}`
Получение цитат по автору.
```bash
curl http://localhost:8080/quotes?author=Confucius
```
#### Пример успешного ответа:
```json
[
  {
    "id": 1,
    "author": "Confucius",
    "quote": "Life is simple, but we insist on making it complicated."
  }
]
```

### `DELETE /quotes/{id}`
Удаление цитаты по id.
```bash
curl -X DELETE http://localhost:8080/quotes/1
```
