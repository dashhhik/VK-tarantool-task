## Инструкция по запуску решения

```commandline
docker-compose up --build
```

### Общие ответы API
Все эндпоинты возвращают ответы в формате JSON.

- `200 OK`: Запрос выполнен успешно.
- `400 Bad Request`: Ошибка в структуре или формате запроса.
- `500 Internal Server Error`: Внутренняя ошибка сервера.

## Endpoints

### Авторизация

#### `POST /login`
Авторизует пользователя по переданным данным.

**Request:**
```json
{
  "username": "example_username",
  "password": "example_password"
}
```

Response (200):
```json
{
  "token": "jwt_token"
}
```

Response (400):
```json
{
   "error": "invalid request body"
}

```

Response (500):
```json
{
   "error": "invalid request body"
}

```

Response (401):

```json
{
"error": "authentication failed"
}
```

Response (500):

```json
{
   "error": "internal server error"
}

```

Описание:

Параметры передаются в теле запроса и должны содержать username и password.
При успешной авторизации возвращается JWT-токен. Имеется пользователь admin с паролем presale.

### Чтение данных
#### POST /read
Читает данные по указанным ключам.

**Request:**
```json
{
   "keys": ["key1", "key2", "key3"]
}

```

Response (200):
```json
{
   "data": {
      "key1": "value1",
      "key2": "value2",
      "key3": "value3"
   }
}
```

Response (400):
```json
{
   "error": "Invalid request format"
}
```

Response (500):
```json
{
   "error": "Internal server error"
}

```



Описание:

Параметры передаются в теле запроса и должны содержать массив ключей keys.
Возвращает данные по указанным ключам в случае успешного выполнения.
Если в запросе отсутствуют ключи, возвращается ошибка с кодом 400.
При возникновении внутренней ошибки сервера возвращается код 500.

### Запись данных
#### POST /write
Читает данные по указанным ключам.

**Request:**
```json
{
   "data": {
      "key1": "value1",
      "key2": "value2",
      "key3": "value3"
   }
}


```

Response (200):
```json
{
   "status": "success"
}

```

Response (400):
```json
{
   "error": "Invalid request structure"
}
```

Response (500):
```json
{
   "error": "Internal server error"
}

```



Описание:

Параметры передаются в теле запроса и должны содержать объект data, в котором ключи соответствуют значениям, которые необходимо записать.
При успешной записи данных возвращается статус success.
Если структура запроса некорректна, возвращается ошибка с кодом 400.
В случае внутренней ошибки сервера возвращается код 500.