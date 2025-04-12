# Usage

## API Documentation

Swagger file is available at [api/swagger.yaml](https://github.com/tab/loki/blob/master/api/swagger.yaml)

## Endpoints

### Smart-ID

#### Create smart-id session

* `POST /api/auth/smart_id`

body:
```json
{
  "country": "EE",
  "personal_code": "50001029996"
}
```

example:
```sh
curl -X POST http://localhost:8080/api/auth/smart_id \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: 4de2f35d-7e30-466e-923b-aab80a424b34" \
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

#### Fetch smart-id session status

* `GET /api/sessions/{id}`

example:
```sh
curl -X GET http://localhost:8080/api/sessions/a658556f-f2ec-42f5-86dc-2665f011d5f7 \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: 7877796b-54b9-4737-a44f-0b0bb4f5eb88" \
  -H "X-Trace-ID: f4c28fec-07fd-415f-900c-37be7fb705fa"
```

response:
```json
{
  "id": "a658556f-f2ec-42f5-86dc-2665f011d5f7",
  "status": "SUCCESS"
}
```

#### Complete smart-id session

* `POST /api/sessions/{id}`

example:
```sh
curl -X POST http://localhost:8080/api/sessions/a658556f-f2ec-42f5-86dc-2665f011d5f7 \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: 2aeb8bca-8af0-498f-8136-c179d3a6f1bd" \
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

### Mobile-ID

#### Create mobile-id session

* `POST /api/auth/mobile_id`

body:
```json
{
  "locale": "ENG",
  "phone_number": "+37268000769",
  "personal_code": "60001017869"
}
```

response:
```json
{
  "id": "a658556f-f2ec-42f5-86dc-2665f011d5f7",
  "code": "8317"
}
```

#### Fetch mobile-id session status

* `GET /api/sessions/{id}`

response:
```json
{
  "id": "a658556f-f2ec-42f5-86dc-2665f011d5f7",
  "status": "SUCCESS"
}
```

#### Complete mobile-id session

* `POST /api/sessions/{id}`

response:
```json
{
  "id": "f4c28fec-07fd-415f-900c-37be7fb705fe",
  "identity_number": "PNOEE-60001017869",
  "personal_code": "60001017869",
  "first_name": "EID2016",
  "last_name": "TESTNUMBER",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### User

#### Fetch user information

* `GET /api/me`

example:
```sh
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: 4844c624-4c3f-4cdf-96dd-01bc53915e02" \
  -H "X-Trace-ID: 0cbc1fe0-c29c-44d5-84a1-4ec5ddb9e08f"
```
response:
```json
{
  "id": "f4c28fec-07fd-415f-900c-37be7fb705fe",
  "identity_number": "PNOEE-50001029996",
  "personal_code": "50001029996",
  "first_name": "TESTNUMBER",
  "last_name": "ADULT"
}
```

### Tokens

### Refresh access token using refresh token

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
```sh
curl -X POST http://localhost:8080/api/tokens/refresh \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: 83bc86e9-1a29-46a8-b358-6db39ab7c2f5" \
  -H "X-Trace-ID: 754cfd21-69b2-436a-af5f-737932cfd874"
  -d '{ "refresh_token": "<REFRESH_TOKEN>" }'
```

#### JWT access token examples

##### Admin

```json
{
  "exp": 1734879499,
  "jti": "PNOEE-50001029996",
  "roles": [
    "admin",
    "user"
  ],
  "permissions": [
    "read:users",
    "write:users",
    "write:self",
    "read:self"
  ],
  "scope": [
    "self-service",
    "sso-service"
  ]
}
```

##### Manager

```json
{
  "exp": 1734879550,
  "jti": "PNOBE-00010299944",
  "roles": [
    "manager",
    "user"
  ],
  "permissions": [
    "read:users",
    "write:self",
    "read:self"
  ],
  "scope": [
    "self-service",
    "sso-service"
  ]
}
```

##### User

```json
{
  "exp": 1734879566,
  "jti": "PNOEE-60001017869",
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

#### JWT refresh token example

```json
{
  "exp": 1734454731,
  "jti": "PNOEE-50001029996"
}
```
