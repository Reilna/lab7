# Лабораторная работа 8

## Задачи

1. Построение базового REST API:

   - Реализуйте сервер, поддерживающий маршруты:
   - GET /users — получение списка пользователей.
   - GET /users/{id} — получение информации о конкретном пользователе.
   - POST /users — добавление нового пользователя.
   - PUT /users/{id} — обновление информации о пользователе.
   - DELETE /users/{id} — удаление пользователя.

2. Подключение базы данных:

   - Добавьте базу данных (например, PostgreSQL или MongoDB) для хранения информации о пользователях.
   - Модифицируйте сервер для взаимодействия с базой данных.

3. Обработка ошибок и валидация данных:

   - Реализуйте централизованную обработку ошибок.
   - Добавьте валидацию данных при создании и обновлении пользователей.

4. Пагинация и фильтрация:

   - Добавьте поддержку пагинации и фильтрации по параметрам запроса (например, поиск пользователей по имени или возрасту).

5. Тестирование API:

   - Реализуйте unit-тесты для каждого маршрута.
   - Проверьте корректность работы при различных вводных данных.

6. Документация API:
   - Создайте документацию для разработанного API с описанием маршрутов, методов, ожидаемых параметров и примеров запросов.

## Инструкции по запуску и тестированию кода

### Запуск программы на Go:

1. REST api

   - Убедиться, что установлен mux: `cd task-1`, `go get github.com/gorilla/mux`;
   - Запустить сервер: `go run task-1/rest-api.go`;
   - Для тестирования GET-запросов можно открыть в браузере: `http://localhost:8080/`.
   - Тестирование запросов через curl:
     - GET - `curl http://localhost:8080/users`, `http://localhost:8080/users/1`;
     - POST - `curl -X POST -H "Content-Type: application/json" -d '{"name": "Типур", "age": 20}' http://localhost:8080/users`;
     - PUT - `curl -X PUT -H "Content-Type: application/json" -d '{"name": "Новый Типур", "age": 21}' http://localhost:8080/users/1`;
     - DELETE - `curl -X DELETE http://localhost:8080/users/1`.

2. БД с использованием PostgreSQL
   - Убедиться, что установлен постгрес: `brew install postgresql@14`; Запуск постгрес: `brew services start postgresql`.
   - Подготовка БД:
     - `brew services start postgresql`
     - `psql postgres`
     - `CREATE DATABASE lb8_users;`
     - `\c lb8_users;`
     - `CREATE TABLE users ( id SERIAL PRIMARY KEY, name VARCHAR(100), age INT);`
   - Запустить сервер: `cd task-2` `go run bd-postgresql.go`
   - Повторить запросы с предыдушего задания.
3. Работа с ошибками и валидация данных
   - Запустить сервер: `cd task-2` `go run validation.go`
   - Ввести неверный запрос на имя: `curl -X POST -H "Content-Type: application/json" -d '{"name": "", "age": 24}' http://localhost:8080/users`
   - Ввести неверный запрос на возраст: `curl -X POST -H "Content-Type: application/json" -d '{"name": "Типур", "age": -10}' http://localhost:8080/users`
4. Пагинация и фильтрация
   - Запустить сервер: `cd task-3` `go run validation.go`
   - Запрос, сочетаюзий пагинацию и фильтрацию сразу: ввести в браузере `http://localhost:8080/users?limit=5&offset=0&name=Артематрр`
5. Тестирование пустой БД
   - Подготовка БД:
     - `psql postgres`
     - `CREATE DATABASE lb8_users_test;`
     - `\c lb8_users_test;`
     - `CREATE TABLE users ( id SERIAL PRIMARY KEY, name VARCHAR(100), age INT);`
   - Запустить сервер: `cd task-6` `go run bd.go`
   - Запустить тесты: `go test -v`
6. См. этот Readme файл
