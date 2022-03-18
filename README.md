# Test task Junior BackDev

### Token

Маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса.

Формат GUID : xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx | 8-4-4-4-12

```
GET http://localhost:5000/api/v1/auth/token?guid=5825a481-f13e-4f80-8ee2-c16ad096709b
```

Успешный запрос вернет:

```json
201 Created

{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHdWlkIjoiZjQ2OWExYjYtNWJlNy00MTNkLWE1ZDktZjFkZmZhODNmMWU5IiwiZXhwIjoxNjQ3Mzg1NDU4LCJpYXQiOjE2NDczODQ1NTh9.T7vBeN1DgB2rFvtGv6co_XswL_MmENTLvHJ7AX8m6irWO48eqTb3uDMzjjSTIB4BYqbkb_aHc1eM4hyZr_cV0Q",
  "token_type": "Bearer",
  "expires_at": 1647387956,
  "refresh_token": "MGE5ZjQ0OTItYWI5Ni00NGY0LTgxOGYtNTVhZTEzN2Q2NWNkr_cV0Q",
  "created_at": 1647384356
}
```

Если GUID будет не соответствовать формату, то вернется ошибка:

```json
400 BadRequest
```

Если GUID уже задействован, то вернется ошибка:

```json
401 Unauthorized
```

Если в процессе возникнут другие неприятности, то вернется ошибка:

```json
500 InternalServerError
```

### Refresh

Маршрут выполняет Refresh операцию на пару Access, Refresh токенов.

```
POST http://localhost:5000/api/v1/auth/refresh
content-type: application/json
```

```json
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHdWlkIjoiZjQ2OWExYjYtNWJlNy00MTNkLWE1ZDktZjFkZmZhODNmMWU5IiwiZXhwIjoxNjQ3Mzg1NDU4LCJpYXQiOjE2NDczODQ1NTh9.T7vBeN1DgB2rFvtGv6co_XswL_MmENTLvHJ7AX8m6irWO48eqTb3uDMzjjSTIB4BYqbkb_aHc1eM4hyZr_cV0Q",
  "refresh_token": "MGE5ZjQ0OTItYWI5Ni00NGY0LTgxOGYtNTVhZTEzN2Q2NWNkr_cV0Q"
}
```

Успешный запрос вернет:

```json
200 OK

{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHdWlkIjoiZjQ2OWExYjYtNWJlNy00MTNkLWE1ZDktZjFkZmZhODNmMWU5IiwiZXhwIjoxNjQ3Mzg1NDU4LCJpYXQiOjE2NDczODQ1NTh9.T7vBeN1DgB2rFvtGv6co_XswL_MmENTLvHJ7AX8m6irWO48eqTb3uDMzjjSTIB4BYqbkb_aHc1eM4hyZr_cV0Q",
  "token_type": "Bearer",
  "expires_at": 1647387956,
  "refresh_token": "MGE5ZjQ0OTItYWI5Ni00NGY0LTgxOGYtNTVhZTEzN2Q2NWNkr_cV0Q",
  "created_at": 1647384356
}
```

Если запрос будет не соответствовать формату, то вернется ошибка:

```json
400 BadRequest
```

Если access_token ещё действителен или refresh_token не действителен, то вернется ошибка:

```json
401 Unauthorized
```

Если в процессе возникнут другие неприятности, то вернется ошибка:

```json
500 InternalServerError
```

### Переменные среды установленные по умолчанию для данного проекта

```
LISTEN_ADDRESS=localhost:5000
TEST_MODE=false
SECRET_KEY=X4AfN38_jfWwFGd
DATABASE_URI=mongodb://localhost:27017/
DATABASE_NAME=db
# "ns", "us" (or "µs"), "ms", "s", "m", "h".
ACCESS_TOKEN_TIME_DELTA=15m
REFRESH_TOKEN_TIME_DELTA=1h
```

<details><summary>Screenshots</summary>

/api/v1/auth/token?guid=

![token](img/1.jpg)

/api/v1/auth/refresh

![refresh](img/2.jpg)

</details>

## Тестовое задание на позицию Junior Backend Developer

**Используемые технологии:**

- Go
- JWT
- MongoDB

**Задание:**

Написать часть сервиса аутентификации.

Два REST маршрута:

- Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
- Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

**Требования:**

Access токен тип JWT, алгоритм SHA512, хранить в базе строго запрещено.

Refresh токен тип произвольный, формат передачи base64, хранится в базе исключительно в виде bcrypt хеша, должен быть защищен от изменения на стороне клиента и попыток повторного использования.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним.

**Результат:**

Результат выполнения задания нужно предоставить в виде исходного кода на Github.
