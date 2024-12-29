# Usage

## API Documentation

Swagger file is available at [api/swagger.yaml](https://github.com/tab/loki/blob/master/api/swagger.yaml)

## Endpoints

* `POST /api/auth/smart_id`

body:
```json
{
  "country": "EE",
  "personal_code": "50001029996"
}
```

example:
```
curl -X POST http://localhost:8080/api/auth/smart_id \
  -H "Content-Type: application/json" \
  -H "X-Trace-ID: f4c28fec-07fd-415f-900c-37be7fb705fa" \
  -d '{ "country": "EE", "personal_code": "50001029996" }'
```

response:
```json
{
  "id": "a658556f-f2ec-42f5-86dc-2665f011d5f7",
  "code": "8317"
}
```

* `GET /api/sessions/{id}`

example:
```
curl -X GET http://localhost:8080/api/sessions/a658556f-f2ec-42f5-86dc-2665f011d5f7 \
  -H "Content-Type: application/json" \
  -H "X-Trace-ID: f4c28fec-07fd-415f-900c-37be7fb705fa"
```

response:
```json
{
  "id": "a658556f-f2ec-42f5-86dc-2665f011d5f7",
  "status": "SUCCESS"
}
```

* `POST /api/sessions/{id}`

example:
```
curl -X POST http://localhost:8080/api/sessions/a658556f-f2ec-42f5-86dc-2665f011d5f7 \
  -H "Content-Type: application/json" \
  -H "X-Trace-ID: f4c28fec-07fd-415f-900c-37be7fb705fa"
```

response:
```json
{
  "id": "f4c28fec-07fd-415f-900c-37be7fb705fe",
  "identity_number": "PNOEE-50001029996",
  "personal_code": "50001029996",
  "first_name": "TESTNUMBER",
  "last_name": "ADULT",
  "access_token": "ey-Access-Token...",
  "refresh_token": "ey-Refresh-Token..."
}
```

* `GET /api/me`

example:
```
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -H "X-Trace-ID: 0cbc1fe0-c29c-44d5-84a1-4ec5ddb9e08f"
```

* `POST /api/tokens/refresh`

body:
```json
{
  "refresh_token": "ey-Refresh-Token..."
}
```

response:
```json
{
  "access_token": "ey-New-Access-Token...",
  "refresh_token": "ey-New-Refresh-Token..."
}
```

example:
```
curl -X POST http://localhost:8080/api/tokens/refresh \
  -H "Content-Type: application/json" \
  -H "X-Trace-ID: 754cfd21-69b2-436a-af5f-737932cfd874"
  -d '{ "refresh_token": "<REFRESH_TOKEN>" }'
```

JWT access token example:
```
{
  "exp": 1734454731,
  "jti": "PNOEE-50001029996",
  "roles": [
    "user"
  ],
  "permissions": [
    "write:self",
    "read:self"
  ],
  "scope": [
    "self-service"
  ]
}
```

JWT refresh token example:
```
{
  "exp": 1734454731,
  "jti": "PNOEE-50001029996",
}
```
