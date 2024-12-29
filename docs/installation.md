# Installation

1. **Environment Variables**: Use `.env` files (e.g., `.env.development`) or provide environment variables for:
- `SECRET_KEY` for JWT signing
- `DATABASE_DSN` for PostgreSQL
- `REDIS_URI` for Redis
- `SMART_ID_API_URL`, `MOBILE_ID_API_URL` and corresponding relying on party credentials


2. **Database Migrations**: Run the following command to apply database migrations:
   ```sh
   GO_ENV=development make db:drop db:create db:migrate
   ```


3. **Run the Services**:
   ```sh
   docker-compose build
   docker-compose up
   ```


4. **Check health status**:
   ```sh
   curl -X GET http://localhost:8080/health
   ```

