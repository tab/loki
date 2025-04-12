# Loki

SSO (Single Sign-On) service that provides user authentication functionality using `Smart-ID` and `Mobile-ID`.
It integrates with SK ID Solutions APIs and manages user roles, permissions, and scopes.

Designed to be easily integrated into microservices architectures and provides logging and telemetry for monitoring.

## Key Features

- Create and update user accounts, with role and scope assignments
- Generate and validate JWT tokens
- Authenticate users via `Smart-ID` and `Mobile-ID` through SK ID Solutions provider APIs
- Comprehensive logging and telemetry support (OpenTelemetry) for easier monitoring and tracing
- Easily integrate into a microservices architecture

## Prerequisites

Before starting this application, you must have the loki-infrastructure running:

```sh
git clone git@github.com/tab/loki-infrastructure.git
cd loki-infrastructure

docker-compose up
```

## Setup and Configuration

### Environment Variables

Use `.env` files (e.g., `.env.development`) or provide environment variables for:

- `DATABASE_DSN` for PostgreSQL
- `REDIS_URI` for Redis
- `SMART_ID_API_URL`, `MOBILE_ID_API_URL` and corresponding relying on party credentials
- `TELEMETRY_URI` for OpenTelemetry

### Generate Certificates and Keys

Before running the services, you need to generate certificates for mTLS and keys for JWT signing:

#### JWT Signing Keys

```sh
mkdir -p certs/jwt

openssl genrsa -out certs/jwt/private.key 4096
openssl rsa -in certs/jwt/private.key -pubout -out certs/jwt/public.key
```

#### mTLS Certificates

```sh
# Generate CA
openssl genrsa -out certs/ca.key 4096
openssl req -new -x509 -key certs/ca.key -sha256 -subj "/CN=Loki CA" -out certs/ca.pem -days 3650

# Generate Server Certificate
openssl genrsa -out certs/server.key 4096
openssl req -new -key certs/server.key -out certs/server.csr -config <(
cat <<-EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
CN = loki-backend

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = backend
IP.1 = 127.0.0.1
IP.2 = 0.0.0.0
EOF
)

openssl x509 -req -in certs/server.csr -CA certs/ca.pem -CAkey certs/ca.key -CAcreateserial -out certs/server.pem -days 825 -sha256 -extfile <(
cat <<-EOF
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = backend
IP.1 = 127.0.0.1
IP.2 = 0.0.0.0
EOF
)

# Generate Client Certificate
openssl genrsa -out certs/client.key 4096
openssl req -new -key certs/client.key -out certs/client.csr -config <(
cat <<-EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
distinguished_name = dn

[dn]
CN = loki-backoffice
EOF
)

openssl x509 -req -in certs/client.csr -CA certs/ca.pem -CAkey certs/ca.key -CAcreateserial -out certs/client.pem -days 825 -sha256
```

For more detailed information on certificates, see [Certificates Documentation](docs/certificates.md).

### Database Migrations

Run the following command to apply database migrations:

```sh
GO_ENV=development make db:drop db:create db:migrate
```

### Run application

```sh
docker-compose build
docker-compose up
```

### Check health status

```sh
curl -X GET http://localhost:8080/live
```

```sh
curl -X GET http://localhost:8080/ready
```

## Documentation

[Documentation](https://tab.github.io/loki)

## API Documentation

Swagger file is available at [api/swagger.yaml](https://github.com/tab/loki/blob/master/api/swagger.yaml)

## Related Repositories

- [Loki Infrastructure](https://github.com/tab/loki-infrastructure) - Infrastructure setup for the Loki ecosystem
- [Loki Backoffice](https://github.com/tab/loki-backoffice) - Backoffice service
- [Loki Proto](https://github.com/tab/loki-proto) - Protocol buffer definitions
- [Loki Frontend](https://github.com/tab/loki-frontend) - Frontend application
- [Smart-ID Client](https://github.com/tab/smartid) - Smart-ID client used for authentication
- [Mobile-ID Client](https://github.com/tab/mobileid) - Mobile-ID client used for authentication

## Architecture

The application follows a layered architecture and clean code principles:

- **cmd/loki**: Application entry point
- **internal/app**: Core application logic, including services, controllers, repositories, and DTOs
- **internal/config**: Configuration loading and setup, server startup, middleware, router initialization, and telemetry configuration
- **pkg**: Reusable utilities such as JWT token handling and logging

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgements

- [SK ID Solutions](https://www.skidsolutions.eu/) for providing the Smart-ID and Mobile-ID APIs
