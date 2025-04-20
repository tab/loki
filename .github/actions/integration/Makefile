LOKI_HOME ?= ""
LOKI_REPO ?= ${LOKI_HOME}/loki
LOKI_BACKOFFICE_REPO ?= ${LOKI_HOME}/loki-backoffice

LOKI_DB_NAME = loki-test
BACKOFFICE_DB_NAME = loki-backoffice-test
DB_USER = postgres
DB_PASSWORD = postgres
DB_HOST = localhost
DB_PORT = 5432

GOOSE_DRIVER = postgres
LOKI_GOOSE_MIGRATION_DIR = $(LOKI_REPO)/db/migrate
BACKOFFICE_GOOSE_MIGRATION_DIR = $(LOKI_BACKOFFICE_REPO)/db/migrate

NETWORK_NAME = loki-network

ifneq (,$(wildcard $(LOKI_REPO)/.env.test))
	include $(LOKI_REPO)/.env.test
	export $(shell sed 's/=.*//' $(LOKI_REPO)/.env.test)
endif

ifneq (,$(wildcard $(LOKI_BACKOFFICE_REPO)/.env.test))
	include $(LOKI_BACKOFFICE_REPO)/.env.test
	export $(shell sed 's/=.*//' $(LOKI_BACKOFFICE_REPO)/.env.test)
endif

.PHONY: setup
setup: db\:setup certs\:generate docker\:network docker\:start check\:services

.PHONY: db\:setup
db\:setup: db\:create db\:migrate

.PHONY: db\:create
db\:create:
	@echo "Creating databases for integration tests..."
	@echo "PostgreSQL: $(DB_HOST):$(DB_PORT)"
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -c "DROP DATABASE IF EXISTS \"$(LOKI_DB_NAME)\";" postgres
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -c "CREATE DATABASE \"$(LOKI_DB_NAME)\";" postgres
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -c "DROP DATABASE IF EXISTS \"$(BACKOFFICE_DB_NAME)\";" postgres
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -c "CREATE DATABASE \"$(BACKOFFICE_DB_NAME)\";" postgres
	@echo "Databases created successfully"

.PHONY: db\:migrate
db\:migrate:
	@echo "Running migrations..."
	@if [ -d "$(LOKI_GOOSE_MIGRATION_DIR)" ]; then \
		echo "Running loki migrations..."; \
		GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(LOKI_DB_NAME) sslmode=disable" goose -dir $(LOKI_GOOSE_MIGRATION_DIR) up || echo "Note: Some loki migrations might fail if tables already exist from schema"; \
	else \
		echo "Warning: Loki migrations directory not found at $(LOKI_GOOSE_MIGRATION_DIR)"; \
	fi

	@if [ -d "$(BACKOFFICE_GOOSE_MIGRATION_DIR)" ]; then \
		echo "Running loki-backoffice migrations..."; \
		GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(BACKOFFICE_DB_NAME) sslmode=disable" goose -dir $(BACKOFFICE_GOOSE_MIGRATION_DIR) up || echo "Note: Some loki-backoffice migrations might fail if tables already exist from schema"; \
	else \
		echo "Warning: Loki-backoffice migrations directory not found at $(BACKOFFICE_GOOSE_MIGRATION_DIR)"; \
	fi
	@echo "Migrations completed"

.PHONY: certs\:generate
certs\:generate:
	@echo "Generating JWT keys and mTLS certificates..."
	./generate-certs.sh "$(LOKI_REPO)" "$(LOKI_BACKOFFICE_REPO)"
	@echo "Certificate generation completed successfully"

.PHONY: docker\:network
docker\:network:
	@echo "Creating Docker network..."
	docker network inspect $(NETWORK_NAME) >/dev/null 2>&1 || docker network create $(NETWORK_NAME)
	@echo "Docker network ready"

.PHONY: docker\:start
docker\:start:
	@echo "Starting services..."
	cp loki-compose.override.yaml $(LOKI_REPO)/compose.override.yaml
	cp loki-backoffice-compose.override.yaml $(LOKI_BACKOFFICE_REPO)/compose.override.yaml

	cd $(LOKI_REPO) && docker compose up -d
	cd $(LOKI_BACKOFFICE_REPO) && docker compose up -d

	@echo "Services started"

.PHONY: check\:services
check\:services:
	@echo "Waiting for services to be ready..."
	@echo "Displaying initial container logs to help with debugging..."
	@echo "Loki logs:" && docker logs loki
	@echo "Loki-backoffice logs:" && docker logs loki-backoffice

	@echo "Testing connection to services..."
	@for i in $$(seq 1 5); do \
		echo "Attempt $$i/5:"; \
		if curl -s --max-time 5 http://localhost:8080/live 2>&1 | grep -q "alive"; then \
			echo "✅ Loki service is up"; \
			LOKI_UP=1; \
		else \
			echo "❌ Loki service not responding yet"; \
			LOKI_UP=0; \
			docker logs --tail 20 loki; \
		fi; \
		if curl -s --max-time 5 http://localhost:8081/live 2>&1 | grep -q "alive"; then \
			echo "✅ Loki-backoffice service is up"; \
			BACKOFFICE_UP=1; \
		else \
			echo "❌ Loki-backoffice service not responding yet"; \
			BACKOFFICE_UP=0; \
			docker logs --tail 20 loki-backoffice; \
		fi; \
		if [ "$$LOKI_UP" = "1" ] && [ "$$BACKOFFICE_UP" = "1" ]; then \
			break; \
		fi; \
		if [ $$i -eq 5 ]; then \
			echo "⚠️ Timed out waiting for services"; \
			echo "Full Loki logs:"; \
			docker logs loki; \
			echo "Full Loki-backoffice logs:"; \
			docker logs loki-backoffice; \
			exit 1; \
		fi; \
		echo "Waiting for services to start (attempt $$i/5)... retrying in 3 seconds"; \
		sleep 3; \
	done

	@echo "All services are ready!"

.PHONY: run
run:
	@echo "Running integration tests..."
	lua run.lua

.PHONY: cleanup
cleanup:
	@echo "Cleaning up..."
	cd $(LOKI_REPO) && docker compose down || true
	cd $(LOKI_BACKOFFICE_REPO) && docker compose down || true

	rm -f $(LOKI_REPO)/compose.override.yaml
	rm -f $(LOKI_BACKOFFICE_REPO)/compose.override.yaml

	docker network rm $(NETWORK_NAME) || true
	@echo "Cleanup complete"

.PHONY: all
all: setup run cleanup
