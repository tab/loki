DB_USER=postgres
DB_DEV_NAME=loki-development
DB_TEST_NAME=loki-test
DB_HOST=localhost
DB_PORT=5432

DB_SERVICE_NAME=database

GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=db/migrate

GO_ENV ?= development
DOCKER_COMPOSE_CMD=docker compose
INFRA_NETWORK=loki-network

ifeq ($(GO_ENV),test)
	ENV_FILE=.env.test
	LOCAL_ENV_FILE=.env.test.local
	DB_NAME=$(DB_TEST_NAME)
	GOOSE_DBSTRING="host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_USER) dbname=$(DB_TEST_NAME) sslmode=disable"
else
	ENV_FILE=.env.development
	LOCAL_ENV_FILE=.env.development.local
	DB_NAME=$(DB_DEV_NAME)
	GOOSE_DBSTRING="host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_USER) dbname=$(DB_DEV_NAME) sslmode=disable"
endif

ifneq (,$(wildcard $(ENV_FILE)))
	include $(ENV_FILE)
	export $(shell sed 's/=.*//' $(ENV_FILE))
endif

ifneq (,$(wildcard $(LOCAL_ENV_FILE)))
	include $(LOCAL_ENV_FILE)
	export $(shell sed 's/=.*//' $(LOCAL_ENV_FILE))
endif

.PHONY: check-infra
check-infra:
	@echo "Checking if infrastructure is running..."
	@docker network inspect $(INFRA_NETWORK) > /dev/null 2>&1 || (echo "Error: $(INFRA_NETWORK) not found. Please make sure loki-infrastructure is running first"; exit 1)
	@echo "Infrastructure is running"

.PHONY: db\:create
db\:create: check-infra
	@echo "Creating $(DB_NAME) database..."
	@docker exec -i $$(docker ps -q -f name=loki-infrastructure-database) createdb -U $(DB_USER) $(DB_NAME) 2>/dev/null || echo "Database $(DB_NAME) already exists."

.PHONY: db\:drop
db\:drop: check-infra
	@echo "Dropping $(DB_NAME) database..."
	@docker exec -i $$(docker ps -q -f name=loki-infrastructure-database) dropdb -U $(DB_USER) --if-exists $(DB_NAME)

.PHONY: db\:migrate
db\:migrate: check-infra
	@echo "Applying migrations to $(DB_NAME) database..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) up

.PHONY: db\:migrate\:status
db\:migrate\:status: check-infra
	@echo "Migration status in $(DB_NAME) database..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) status

.PHONY: db\:rollback
db\:rollback: check-infra
	@echo "Rolling back last migration in $(DB_NAME) database..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) down

.PHONY: db\:schema\:dump
db\:schema\:dump: check-infra
	@echo "Dumping database schema from $(DB_NAME) to db/schema.sql..."
	@docker exec -i $$(docker ps -q -f name=loki-infrastructure-database) pg_dump -U $(DB_USER) -d $(DB_NAME) --schema-only --exclude-table=goose_db_version > db/schema.sql

.PHONY: db\:schema\:load
db\:schema\:load: check-infra
	@echo "Loading schema from db/schema.sql into $(DB_NAME) database..."
	@cat db/schema.sql | docker exec -i $$(docker ps -q -f name=loki-infrastructure-database) psql -U $(DB_USER) -d $(DB_NAME)

.PHONY: run
run: check-infra
	@echo "Starting loki app..."
	$(DOCKER_COMPOSE_CMD) up -d

.PHONY: stop
stop:
	@echo "Stopping loki app..."
	$(DOCKER_COMPOSE_CMD) down

.PHONY: logs
logs:
	@echo "Showing logs..."
	$(DOCKER_COMPOSE_CMD) logs -f

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

.PHONY: lint\:fix
lint\:fix:
	@echo "Running golangci-lint --fix ..."
	golangci-lint run --fix

.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

.PHONY: staticcheck
staticcheck:
	@echo "Running staticcheck..."
	staticcheck ./...

.PHONY: test
test:
	@echo "Running tests..."
	go test -cover ./...

.PHONY: coverage
coverage:
	@echo "Generating test coverage report..."
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

.PHONY: all
all: lint vet staticcheck test
