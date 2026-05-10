# 🛒 ApiMSA — Marketplace

> Микросервисная архитектура на Go. HTTP → gRPC → PostgreSQL + Kafka.

```
┌─────────────┐   HTTP/JSON   ┌──────────────┐   gRPC    ┌──────────────┐
│   Клиент    │ ────────────► │  ApiService  │ ────────► │  BDservice   │
└─────────────┘               └──────┬───────┘           └──────┬───────┘
                                     │                           │
                                     │ Kafka (async)             │ PostgreSQL
                                     ▼                           ▼
                              [ cart_events ]            [ users / items / cart ]
```

---

## 📁 Структура проекта

```
ApiMSA/
├── ApiService/              # HTTP сервер (порт 8080)
│   ├── main.go
│   ├── Servs/
│   │   ├── server.go        # роутер, запуск
│   │   ├── HandleAccount.go # создание / авторизация / удаление аккаунта
│   │   ├── HandleCort.go    # корзина
│   │   └── broker/
│   │       └── broker.go    # Kafka producer
│   └── user_pb/             # сгенерированный gRPC код
│
├── BDservice/               # gRPC сервер + Kafka consumer (порт 5051)
│   ├── main.go
│   ├── Servs/
│   │   ├── TablesAndConnection.go  # подключение к БД, создание таблиц
│   │   ├── funcsForgRps.go         # gRPC handlers
│   │   └── forKafka.go             # Kafka consumer
│   └── user_pb/             # сгенерированный gRPC код
│
└── Kafka/
    └── docker-compose.yaml  # Kafka + PostgreSQL
```

---

## 🚀 Запуск

### 1. Поднять Kafka и PostgreSQL

```bash
cd Kafka
docker-compose up -d
```

---

### 2. Запустить BDservice

**Windows (PowerShell):**
```powershell
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"

cd BDservice
go run main.go
```

**Linux / Mac:**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres

cd BDservice
go run main.go
```

Вывод при успехе:
```
BDservice запущен на порту :5051
```

---

### 3. Запустить ApiService

```bash
cd ApiService
go run main.go --port 8080
```

Вывод при успехе:
```
ApiService запущен на порту :8080
```

---

## 🔌 Эндпоинты

### 👤 Аккаунты

| Метод    | Путь             | Описание              |
|----------|------------------|-----------------------|
| `POST`   | `/CreateAccount` | Создать аккаунт       |
| `GET`    | `/Avtorizacion`  | Авторизация           |
| `DELETE` | `/DeleteAccount` | Удалить аккаунт       |

**Тело запроса (JSON):**
```json
{
  "user_name": "Ivan",
  "last_name": "Ivanov",
  "email": "ivan@example.com"
}
```

**Ответ:**
```json
{
  "response": "Аккаунт успешно добавлен в бд",
  "succes": true
}
```

---

### 🛒 Корзина

| Метод    | Путь                    | Query             | Описание                    |
|----------|-------------------------|-------------------|-----------------------------|
| `GET`    | `/ShowAllItems`         | `?user_id=1`      | Получить корзину пользователя |
| `POST`   | `/CreateBuy`            | —                 | Добавить товар в корзину    |
| `PATCH`  | `/ChangePrice`          | —                 | Изменить цену товара        |
| `DELETE` | `/DeleteBuyFromKorzina` | —                 | Удалить товар из корзины    |

**POST /CreateBuy:**
```json
{
  "user_id": 1,
  "item_name": "Телефон",
  "item_price": 50000
}
```

**PATCH /ChangePrice:**
```json
{
  "item_name": "Телефон",
  "item_price": 45000
}
```

**DELETE /DeleteBuyFromKorzina:**
```json
{
  "user_id": 1,
  "item_name": "Телефон"
}
```

---

## Схема базы данных

```sql
-- Пользователи
CREATE TABLE users (
    user_id   SERIAL PRIMARY KEY,
    user_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    email     VARCHAR
);

-- Товары
CREATE TABLE items (
    item_name  VARCHAR NOT NULL PRIMARY KEY,
    item_price INTEGER NOT NULL
);

-- Корзина
CREATE TABLE cart (
    user_id          INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    item_name_in_cart VARCHAR REFERENCES items(item_name)
);
```

> Таблицы создаются **автоматически** при запуске BDservice.

---

## Технологии

| Что          | Чем                          |
|--------------|------------------------------|
| HTTP сервер  | `gorilla/mux`                |
| Межсервисное | `gRPC` + `protobuf`          |
| Асинхронность| `Apache Kafka` (kafka-go)    |
| База данных  | `PostgreSQL` (pgx/v5)        |
| Инфраструктура | `Docker` + `docker-compose` |

---

## Переменные окружения (BDservice)

| Переменная    | Описание         | Дефолт      |
|---------------|------------------|-------------|
| `DB_HOST`     | Хост PostgreSQL  | `localhost` |
| `DB_PORT`     | Порт PostgreSQL  | `5432`      |
| `DB_USER`     | Пользователь БД  | `postgres`  |
| `DB_PASSWORD` | Пароль БД        | `postgres`  |

---

## Kafka Topics

| Topic         | Ключ | Операция             |
|---------------|------|----------------------|
| `cart_events` | `1`  | Добавить в корзину   |
| `cart_events` | `2`  | Удалить из корзины   |
| `cart_events` | `3`  | Изменить цену товара |
