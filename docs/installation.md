# Installation

## Prerequisites

Before starting the Loki application, you must have the loki-infrastructure running:

```sh
git clone git@github.com/tab/loki-infrastructure.git
cd loki-infrastructure

docker-compose up
```

## Environment Variables

Use `.env` files (e.g., `.env.development`) or provide environment variables for:

- `DATABASE_DSN` for PostgreSQL
- `REDIS_URI` for Redis
- `SMART_ID_API_URL`, `MOBILE_ID_API_URL` and corresponding relying on party credentials
- `TELEMETRY_URI` for OpenTelemetry

Example `.env.development` file:

```
DATABASE_DSN=postgres://postgres:postgres@localhost:5432/loki_development?sslmode=disable
REDIS_URI=redis://localhost:6379/0
SMART_ID_API_URL=https://sid.demo.sk.ee/smart-id-rp/v2/
MOBILE_ID_API_URL=https://tsp.demo.sk.ee/mid-api/
TELEMETRY_URI=http://localhost:4317
```

## Certificate and Key Generation

Before running the services, you need to generate certificates for mTLS and keys for JWT signing.

For more detailed information on certificates, see [Certificates Documentation](certificates.md).

## Database Migrations

Run the following command to apply database migrations:

```sh
GO_ENV=development make db:drop db:create db:migrate
```

## Run application

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

## Related Repositories

The Loki ecosystem consists of the following repositories:

- [Loki Infrastructure](https://github.com/tab/loki-infrastructure) - Infrastructure setup for the Loki ecosystem
- [Loki Backoffice](https://github.com/tab/loki-backoffice) - Backoffice service
- [Loki Proto](https://github.com/tab/loki-proto) - Protocol buffer definitions
- [Loki Frontend](https://github.com/tab/loki-frontend) - Frontend application
- [Smart-ID Client](https://github.com/tab/smartid) - Smart-ID client used for authentication
- [Mobile-ID Client](https://github.com/tab/mobileid) - Mobile-ID client used for authentication
