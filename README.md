# Event_Booker

Event_Booker — это лёгкая и понятная система бронирования мест на мероприятия, с уведомлениями по Email, созданная на Go.

Задание на проект хранится в фале **TASK.md** в корневом репозитории

Проект включает веб-интерфейс, REST API, а также фоновый воркер, который автоматически отменяет просроченные брони.

##  Основные особенности
```
✅  Создание и управление мероприятиями
✅  Бронирование мест с TTL (авто-отмена неподтверждённых броней)
✅  Подтверждение брони (для событий с оплатой)
✅  Разделение ролей: пользователь и администратор
```
## Функциональность

### Пользователь
Пользовательская часть позволяет:

- Просматривать список доступных мероприятий
- Смотреть подробную информацию о событии
- Бронировать место на мероприятии
- Подтверждать бронь (если требуется оплата)
- Видеть, как истекают неподтверждённые брони (TTL)

### Администратор
Администратор может:
- Создавать новые мероприятия
- Создавать новых пользователей
- Просматривать все бронирования
- Управлять параметрами события (места, TTL, необходимость оплаты)

### Системная часть
В фоне работает воркер, который:
- регулярно проверяет истёкшие брони
- автоматически отменяет неподтверждённые записи
- возвращает места в доступные
Все данные хранятся в PostgreSQL
Авторизация работает на JWT
Чёткое разделение ролей и обработка прав доступа

## Структура проекта
Проект организован по принципам чистой архитектуры и разделён на независимые логические модули.

```

event_booker/
│
├── cmd/
│   └── server/            # Точка входа в приложение
│
├── internal/
│   ├── app/               # Инициализация зависимостей
│   ├── handlers/          # HTTP-обработчики (REST API)
│   ├── middleware/        # JWT и проверки ролей
│   ├── model/             # Бизнес-модели (Event, Booking, User)
│   ├── notifier/          # Отправка email-уведомлений
│   ├── service/           # Бизнес-логика: бронирование, TTL-воркер, события
│   └── storage/           # Хранилище (PostgreSQL)
│
├── db/
│   └── dumps/             # Дампы базы данных
│
├── web/
│   └── index.html         # Полный веб-интерфейс (юзерская + админская части)
│
├── TASK.md                # Условие задания
├── README.md              # Документация проекта
└── go.mod / go.sum        # Go-зависимости


```

## Используемые в проекте технологии
### Backend
```
Go 1.25 — основной язык разработки
Gin — быстрый HTTP-фреймворк для создания REST API
JWT (github.com/golang-jwt/jwt) — аутентификация и авторизация
PostgreSQL — основная база данных
pgx — драйвер и адаптер PostgreSQL
bcrypt — безопасное хеширование паролей

### Storage
Чистая реализация через интерфейсы
PostgreSQL в качестве единственного хранилища
SQL-дампы в db/dumps для быстрого восстановления БД

### Services
Внутренняя бизнес-логика (слой service)
TTL-механизм для автоматической отмены броней
Email-уведомления — модуль notifier

### Frontend
Vanilla HTML + JS использует API напрямую
Реализованы:
 - регистрация и вход
 - панель пользователя
 - панель администратора
 - создание событий
 - бронирование
 - подтверждение брони
```

## Устновка и запуск
1. Установка зависимостей
Убедитесь, что у вас установлено
Go 1.25+
Docker и Docker Compose
PostgreSQL-клиент (либо контейнер в докере)
2. Клонирование репозитория
git clone https://github.com/Vladimirmoscow84/event_booker.git
cd event_booker

3. Создайте .env в корне проекта
DATABASE_URI="host=localhost port=5440 user=vladimir password=password dbname=event_booker sslmode=disable"
SERVER_ADDRESS=":6060"
JWT_SECRET="ваш код"
EMAIL_HOST="smtp.yandex.ru" если клиент яндекс
EMAIL_PORT="587"
EMAIL_USER="почта отправителя"
EMAIL_PASS="ваш пароль с системы"
EMAIL_FROM="почта отправителя"
EMAIL_TO="почта получателя,почта получателя"

Укажите реальные SMTP-данные, чтобы тестировать отправку писем.

4. Запуск PostgreSQL через Docker
Проект рассчитан на работу с PostgreSQL.
Поднимите контейнер:
docker run -d \
  --name pg-event_booker \
  -p 5440:5432 \
  -e POSTGRES_USER=vladimir \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=event_booker \
  postgres:latest

Проверка работы контейнера:
docker ps

5. Применение sql. миграций
В проекте используется sql-migrate (локальные .sql файлы в db/dumps)
создание таблиц:
env -u PGLOCALEDIR migrate -database "postgres://localhost:5440/event_booker?sslmode=disable&user=vladimir&password=password" -path /Users/mac/Dev/WBTech/event_booker/db/dumps up

удаление таблиц:
 env -u PGLOCALEDIR migrate -database "postgres://localhost:5440/event_booker?sslmode=disable&user=vladimir&password=password" -path /Users/mac/Dev/WBTech/event_booker/db/dumps down

6. Запуск сервера
go run cmd/server/main.go

Если запуск успешен, в консоли будет:
[postgres] successful connect to DB
[app] storage initialized successfully
[app] email client initialized successfully
[app] service initialized successfully
[app] starting server on :6060

7. Запуск фронта
Фронт расположен в папке:  web/index.html
Доступен по адресу: http://localhost:6060/

8. Создание администратора (через Postman)
По умолчанию в системе нет админа, его нужно создать вручную.
Шаг 1: логин обычным пользователем (его можно создать через /auth/register)
Шаг 2: в админском методе /users можно создать нового ADMIN

Запрос:

POST http://localhost:6060/users
Headers: 
Authorization: Bearer <JWT_токен_админа>
Content-Type: application/json

Body:
{
  "email": "admin@example.com",
  "password": "password",
  "role": "admin"
}
После этого можно создавать события через /events.

 ```

 9. Примеры запросов

Ниже приведены примеры основных API-запросов для работы с системой. Все примеры приведены в формате JSON и подходят для Postman, Insomnia или cURL.

Авторизация
Регистрация пользователя
POST /auth/register
{
  "email": "user@example.com",
  "password": "password"
}

Логин
POST /auth/login
{
  "email": "user@example.com",
  "password": "password"
}

Пример успешного ответа:
{
  "token": "<jwt_token>"
}
Используйте токен в заголовке:
Authorization: Bearer <jwt_token>

Управление пользователями
Создание администратора (только admin)
POST /users
Headers:
Authorization: Bearer <админский токен>
Body:
{
  "email": "admin@example.com",
  "password": "password",
  "role": "admin"
}

Мероприятия (Events)
Создать событие (admin)
POST /events
{
  "title": "Rock Concert",
  "description": "Best concert ever",
  "total_seats": 50,
  "ttl_minutes": 15,
  "require_confirmation": true
}

Получить список событий
GET /events
Ответ:

[
  {
    "id": 1,
    "title": "Rock Concert",
    "available_seats": 42,
    "ttl_minutes": 15
  }
]

Получить одно событие
GET /events/{id}

Бронирования (Bookings)
Создать бронь
POST /bookings
{
  "event_id": 1
}
Ответ:
{
  "booking_id": 12,
  "expires_at": "2025-01-15T12:30:00Z"
}

Подтвердить бронь (если требуется оплата)
POST /bookings/confirm
{
  "booking_id": 12
}

Список всех броней (admin)
GET /bookings

```

## Email-уведомления
Система автоматически отправляет email:
при создании брони
при подтверждении
при авто-отмене (TTL истёк)
Почтовая конфигурация задаётся в .env.



Автор
Разработчик: Vladimirmoscow84
Контакт: ccr1@yandex.ru
GitHub: github.com/Vladimirmoscow84





