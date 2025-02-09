# **Сервис регистрации пользователей и подтверждения email**
## **Обзор**
Данный проект представляет собой API для регистрации пользователей и верификации их email-адресов. Он разработан на языке Go (Golang) с акцентом на простоту, масштабируемость и следование лучшим практикам современной архитектуры веб-приложений.
В проекте используются следующие инструменты и технологии:
- **PostgreSQL** для хранения данных.
- **Gmail API** для отправки писем с подтверждением email.
- **Docker** для контейнеризации и удобного деплоя.
- **JWT** для защиты токенов.
- **Gorilla Mux** как маршрутизатор для обработки HTTP-запросов.

## **Содержание**
1. [Особенности](#особенности)
2. [Архитектура и подходы](#архитектура-и-подходы)
3. [Структура проекта](#структура-проекта)
4. [Переменные окружения](#переменные-окружения)
5. [База данных](#база-данных)
6. [Сервис писем](#сервис-писем)
7. [Как запустить проект](#как-запустить-проект)
    - Режим разработки
    - Режим продакшена
8. [Применённые лучшие практики](#применённые-лучшие-практики)

## **Особенности**
1. **Регистрация пользователей**:
    - Возможность регистрировать пользователя через email и уникальный никнейм.
    - Отправка письма подтверждения регистрации.

2. **Подтверждение email**:
    - Верификация email через токен, переданный в ссылке подтверждения.
    - При успешной верификации пользователь помечается как "подтверждённый" в базе данных.

3. **JWT для защиты токенов**:
    - Генерация безопасных токенов с выставленным временем их действия.
    - Использование токена для подтверждения email.

4. **PostgreSQL**:
    - Хранение данных пользователей в реляционной базе данных.

5. **Docker-окружение**:
    - Полностью контейнеризированная настройка приложения (и приложения, и базы данных).
    - Использование Docker Compose для запуска окружения.

6. **Конфигурация по окружению**:
    - Возможность гибко управлять настройками через разные `.env` файлы для разработки и продакшена.

## **Архитектура и подходы**
Проект разделён на слои в соответствии с принципом **Чистой архитектуры (Clean Architecture)**. Это помогает разделить логику, улучшить читаемость кода и упростить поддержку. Основные слои:
### 1. **Обработчики (Handler)**
Этот слой отвечает за обработку HTTP-запросов, а также за преобразование входящих данных и отправку ответов клиенту.
- Парсинг входящих данных.
- Вызывают сервисы для реализации логики.

### 2. **Сервисы (Service)**
Содержат бизнес-логику приложения.
- Генерация токенов для подтверждения email.
- Сервисы для регистрации пользователей.
- Проверка токенов в процессе подтверждения email.

### 3. **Репозитории (Repository)**
Обеспечивают доступ к базе данных.
- Выполняют SQL-запросы для добавления, поиска или обновления данных.
- Изолируют бизнес-логику от взаимодействия с базой данных.

Эта архитектура позволяет легко изменять один слой, не затрагивая другие.
## **Структура проекта**
Проект имеет следующую файловую структуру:
```jsregexp
├── config/                 # Конфигурация приложения
├── internal/
│   ├── app/                # Инициализация приложения (DB, маршруты и т.д.)
│   ├── dto/                # DTO (структуры запросов и ответов в HTTP)
│   ├── handler/            # HTTP-обработчики
│   ├── model/              # Модели данных для базы и приложения
│   ├── repository/         # Слой взаимодействия с базой данных
│   ├── service/            # Логика приложения (письма, токены, регистрация)
│   └── utils/              # Утилиты и вспомогательные функции
├── cmd/                    # Точка входа в приложение
├── Dockerfile              # Многоэтапная сборка Docker-образа
├── docker-compose.*.yml    # Сборки для разработки и продакшена
├── .env                    # Переменные окружения
├── README.md               # Документация
└── go.mod                  # Модули Go
```

## **Переменные окружения**
Для управления конфигурацией проекта используются переменные окружения из `.env` файлов. Вот основные переменные:

| Имя | Описание | Пример |
| --- | --- | --- |
| `APP_PORT` | Порт сервера | `8080` |
| `APP_ENV` | Режим окружения (development/production) | `development` |
| `APP_API_BASE_PATH` | Префикс для маршрутов API | `/api/v1` |
| `JWT_SECRET_KEY` | Ключ для подписи JWT токенов | `секретный_ключ` |
| `POSTGRES_HOST` | Адрес базы данных | `db` |
| `POSTGRES_PORT` | Порт базы данных | `5432` |
| `POSTGRES_USER` | Логин базы данных | `dev` |
| `POSTGRES_PASSWORD` | Пароль базы данных | `dev` |
| `GOOGLE_CLIENT_ID` | Клиентский ID Google OAuth2 | `ваш_client_id` |
| `GOOGLE_CLIENT_SECRET` | Секретный ключ клиента Google OAuth2 | `ваш_client_secret` |
| `EMAIL_VERIFICATION_URL` | Ссылка для подтверждения email | `http://localhost:8080/api/v1/users/verify` |
## **База данных**
Проект использует PostgreSQL для управления данными пользователей.
### **Таблица `users`**

| Поле | Тип | Описание |
| --- | --- | --- |
| `id` | `uuid` | Уникальный идентификатор |
| `nickname` | `text` | Никнейм |
| `email` | `text` | Email пользователя |
| `verified` | `boolean` | Пометка, подтверждён email или нет |
| `created_at` | `timestamp` | Дата и время создания записи |
Для инициализации базы можно использовать файл `local-init.sql.gz`.
## **Сервис писем**
### **1. Генерация токенов:**
Сервис писем создаёт **JWT токены**, которые используются пользователем для подтверждения email.
- Токен включает в себя email и время жизни (`EMAIL_TOKEN_EXPIRY`, задаётся в минутах).
- Подпись выполняется с использованием `JWT_SECRET_KEY`.

### **2. Отправка писем:**
Для отправки писем используется **Gmail API**.
- Пользователям отправляется письмо с ссылкой вида:
  `http://localhost:8080/api/v1/users/verify?token=<токен>`.
- Если токен действителен, email подтверждается.

## **Как запустить проект**
### **Необходимое ПО**:
Перед запуском убедитесь, что на вашем компьютере установлен **Docker и Docker Compose**.
### **1. Режим разработки**
Шаги для запуска в режиме разработки:
1. Скопируйте файл `.env.development` в `.env`:
    ```shell
      cp .env.development .env
    ```
2. Запустите сборку и контейнеры:
    ```shell
      docker-compose -f docker-compose.development.yml up --build
    ```
### **2. Режим продакшена**
Шаги для запуска в режиме продакшена:
1. Скопируйте файл `.env.production` в `.env`:
    ```shell
      cp .env.production .env
    ```
2. Запустите сборку и контейнеры:
    ```shell
      docker-compose -f docker-compose.production.yml up --build
    ```
Приложение будет работать на `http://localhost:8080`.

## **Применённые лучшие практики**
1. **Модульная структура кода**
   Код разделён на модули (handlers, services, repos), что упрощает поддержку.
2. **Конфигурации по окружениям**
   Использование `.env` файлов для изоляции настроек.
3. **Безопасность**
    - JWT токены со временем жизни.
    - Хранение конфиденциальных данных в переменных окружения.

4. **Масштабируемая структура базы**
   PostgreSQL позволяет легко обрабатывать большие объёмы данных.
5. **Докеризация**
   Использование Docker позволяет запускать проект в любом окружении.
6. **Лёгкая интеграция API (Gmail)**
   Поддерживается работа с Gmail API для отправки email.

## **API Эндпоинты**

| Метод | Адрес | Описание |
| --- | --- | --- |
| POST | `/users/register` | Регистрация нового пользователя |
| GET | `/users/verify?token=...` | Подтверждение email пользователя |