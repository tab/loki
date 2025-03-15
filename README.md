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
```

```sh
docker-compose up
```

## Setup and Configuration

**Environment Variables**:

Use `.env` files (e.g., `.env.development`) or provide environment variables for:

- `SECRET_KEY` for JWT signing
- `DATABASE_DSN` for PostgreSQL
- `REDIS_URI` for Redis
- `SMART_ID_API_URL`, `MOBILE_ID_API_URL` and corresponding relying on party credentials
- `TELEMETRY_URI` for OpenTelemetry

**Database Migrations**:

Run the following command to apply database migrations:

```sh
GO_ENV=development make db:drop db:create db:migrate
```

**Run the Services**:

```sh
docker-compose build
docker-compose up
```

**Check health status**:

```sh
curl -X GET http://localhost:8080/health
```

## Documentation

[Documentation](https://tab.github.io/loki)

## API Documentation

Swagger file is available at [api/swagger.yaml](https://github.com/tab/loki/blob/master/api/swagger.yaml)

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
