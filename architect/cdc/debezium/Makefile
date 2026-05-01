# =============================================================================
# Debezium CDC — Makefile
# =============================================================================
# Usage:  make <target>
#         make help        → show all targets
# =============================================================================

COMPOSE         := docker compose
CONNECT_URL     := http://localhost:8083
CONNECTOR_FILE  := connectors/register-postgres.json
CONNECTOR_NAME  := inventory-connector
KAFKA_TOPIC     := cdc.public.customers
KAFKA_BROKER    := localhost:29092       # host access
KAFKA_INTERNAL  := kafka:9092           # docker-internal
PG_HOST         := localhost
PG_PORT         := 5433
PG_USER         := postgres
PG_DB           := inventory

# Colours
CYAN   := \033[0;36m
GREEN  := \033[0;32m
YELLOW := \033[1;33m
RED    := \033[0;31m
NC     := \033[0m

.DEFAULT_GOAL := help
.PHONY: help up down restart logs ps \
        wait-connect \
        connector-register connector-status connector-delete connector-pause connector-resume \
        topics consume consume-raw \
        migrate seed sql-shell sql-show \
        pg-shell pg-insert \
        verify clean

# ── Help ──────────────────────────────────────────────────────────────────────
help:
	@echo ""
	@echo "$(CYAN)Debezium CDC — available targets$(NC)"
	@echo ""
	@echo "  $(GREEN)Stack$(NC)"
	@echo "    make up               Start all containers (detached)"
	@echo "    make down             Stop and remove containers"
	@echo "    make down-volumes     Stop, remove containers AND volumes"
	@echo "    make restart          Restart all containers"
	@echo "    make ps               Show container status"
	@echo "    make logs             Tail logs for all services"
	@echo "    make logs-connect     Tail Kafka Connect logs"
	@echo "    make logs-kafka       Tail Kafka logs"
	@echo "    make logs-postgres    Tail Postgres logs"
	@echo ""
	@echo "  $(GREEN)Connector$(NC)"
	@echo "    make connector-register   Register the Postgres connector"
	@echo "    make connector-status     Show connector & task status"
	@echo "    make connector-delete     Delete the connector"
	@echo "    make connector-pause      Pause the connector"
	@echo "    make connector-resume     Resume a paused connector"
	@echo "    make connector-list       List all registered connectors"
	@echo ""
	@echo "  $(GREEN)Kafka$(NC)"
	@echo "    make topics           List all Kafka topics"
	@echo "    make consume          Consume CDC events (pretty JSON)"
	@echo "    make consume-raw      Consume CDC events (raw)"
	@echo ""
	@echo "  $(GREEN)PostgreSQL$(NC)"
	@echo "    make pg-shell         Open psql shell"
	@echo "    make pg-insert        Insert a test row into customers"
	@echo ""
	@echo "  $(GREEN)SQL Migrations$(NC)"
	@echo "    make migrate          Run all V*.sql files in sql/migrations/ (ordered)"
	@echo "    make seed             Insert seed data (V2)"
	@echo "    make sql-show         SELECT * FROM customers"
	@echo "    make sql-shell        Open psql shell (alias)"
	@echo ""
	@echo "  $(GREEN)Utilities$(NC)"
	@echo "    make wait-connect     Wait until Kafka Connect REST API is ready"
	@echo "    make verify           Run end-to-end verification"
	@echo "    make setup            Full setup: up + wait + register connector"
	@echo "    make clean            Stop stack and remove all volumes"
	@echo ""

# ── Stack ─────────────────────────────────────────────────────────────────────
up:
	@echo "$(CYAN)Starting Debezium stack...$(NC)"
	$(COMPOSE) up -d
	@echo "$(GREEN)Stack is up. Run 'make ps' to check status.$(NC)"

down:
	@echo "$(YELLOW)Stopping Debezium stack...$(NC)"
	$(COMPOSE) down

down-volumes:
	@echo "$(RED)Stopping stack and removing volumes...$(NC)"
	$(COMPOSE) down -v

restart:
	@echo "$(CYAN)Restarting Debezium stack...$(NC)"
	$(COMPOSE) restart

ps:
	$(COMPOSE) ps

logs:
	$(COMPOSE) logs -f

logs-connect:
	$(COMPOSE) logs -f connect

logs-kafka:
	$(COMPOSE) logs -f kafka

logs-postgres:
	$(COMPOSE) logs -f postgres

# ── Wait for Kafka Connect ────────────────────────────────────────────────────
wait-connect:
	@echo "$(CYAN)Waiting for Kafka Connect REST API...$(NC)"
	@retries=30; \
	until curl -sf $(CONNECT_URL)/connectors > /dev/null 2>&1; do \
		retries=$$((retries - 1)); \
		if [ $$retries -le 0 ]; then \
			echo "$(RED)Kafka Connect did not become ready. Check: make logs-connect$(NC)"; \
			exit 1; \
		fi; \
		printf "."; \
		sleep 5; \
	done; \
	echo ""; \
	echo "$(GREEN)Kafka Connect is ready.$(NC)"

# ── Connector ─────────────────────────────────────────────────────────────────
connector-register: wait-connect
	@echo "$(CYAN)Registering connector: $(CONNECTOR_NAME)$(NC)"
	@EXISTING=$$(curl -sf $(CONNECT_URL)/connectors/$(CONNECTOR_NAME) 2>/dev/null || true); \
	if [ -n "$$EXISTING" ]; then \
		echo "$(YELLOW)Connector already exists — updating config...$(NC)"; \
		curl -sf -X PUT \
			-H "Content-Type: application/json" \
			-d "$$(jq '.config' $(CONNECTOR_FILE))" \
			$(CONNECT_URL)/connectors/$(CONNECTOR_NAME)/config | jq .; \
	else \
		curl -sf -X POST \
			-H "Content-Type: application/json" \
			-d @$(CONNECTOR_FILE) \
			$(CONNECT_URL)/connectors | jq .; \
	fi
	@echo "$(GREEN)Done. Run 'make connector-status' to verify.$(NC)"

connector-status:
	@echo "$(CYAN)Connector status:$(NC)"
	@curl -sf $(CONNECT_URL)/connectors/$(CONNECTOR_NAME)/status | jq .

connector-delete:
	@echo "$(RED)Deleting connector: $(CONNECTOR_NAME)$(NC)"
	@curl -sf -X DELETE $(CONNECT_URL)/connectors/$(CONNECTOR_NAME)
	@echo "$(GREEN)Deleted.$(NC)"

connector-pause:
	@echo "$(YELLOW)Pausing connector: $(CONNECTOR_NAME)$(NC)"
	@curl -sf -X PUT $(CONNECT_URL)/connectors/$(CONNECTOR_NAME)/pause
	@echo "$(GREEN)Paused.$(NC)"

connector-resume:
	@echo "$(CYAN)Resuming connector: $(CONNECTOR_NAME)$(NC)"
	@curl -sf -X PUT $(CONNECT_URL)/connectors/$(CONNECTOR_NAME)/resume
	@echo "$(GREEN)Resumed.$(NC)"

connector-list:
	@echo "$(CYAN)Registered connectors:$(NC)"
	@curl -sf $(CONNECT_URL)/connectors | jq .

# ── Kafka ─────────────────────────────────────────────────────────────────────
topics:
	@echo "$(CYAN)Kafka topics:$(NC)"
	@docker exec debezium-kafka /kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 --list

consume:
	@echo "$(CYAN)Consuming from topic: $(KAFKA_TOPIC) (Ctrl+C to stop)$(NC)"
	@docker exec debezium-kafka /kafka/bin/kafka-console-consumer.sh \
		--bootstrap-server localhost:9092 \
		--topic $(KAFKA_TOPIC) \
		--from-beginning \
		--property print.key=true \
		--property key.separator=" → " \
	| while IFS= read -r line; do echo "$$line" | python3 -m json.tool 2>/dev/null || echo "$$line"; done

consume-raw:
	@echo "$(CYAN)Consuming raw from topic: $(KAFKA_TOPIC) (Ctrl+C to stop)$(NC)"
	@docker exec debezium-kafka /kafka/bin/kafka-console-consumer.sh \
		--bootstrap-server localhost:9092 \
		--topic $(KAFKA_TOPIC) \
		--from-beginning

# ── PostgreSQL ────────────────────────────────────────────────────────────────
pg-shell:
	@echo "$(CYAN)Opening psql shell (inventory db)...$(NC)"
	@docker exec -it debezium-postgres psql -U $(PG_USER) -d $(PG_DB)

pg-insert:
	@echo "$(CYAN)Inserting test row into customers...$(NC)"
	$(eval RND := $(shell echo $$RANDOM))
	@docker exec debezium-postgres psql -U $(PG_USER) -d $(PG_DB) -c \
		"INSERT INTO public.customers (name, email) \
		 VALUES ('CDC Test $(RND)', 'cdc$(RND)@test.com') \
		 ON CONFLICT (email) DO NOTHING;"
	@echo "$(GREEN)Inserted test customer. Watch 'make consume' for the CDC event.$(NC)"

# ── SQL Migrations ────────────────────────────────────────────────────────────
migrate:
	@echo "$(CYAN)Running SQL migrations...$(NC)"
	@for f in $$(ls sql/migrations/V*.sql | sort); do \
		echo "  → $$f"; \
		docker exec -i debezium-postgres psql -U $(PG_USER) -d $(PG_DB) < $$f; \
	 done
	@echo "$(GREEN)Migrations complete.$(NC)"

seed:
	@echo "$(CYAN)Running seed data...$(NC)"
	@docker exec -i debezium-postgres psql -U $(PG_USER) -d $(PG_DB) \
		< sql/migrations/V2__seed_customers.sql
	@echo "$(GREEN)Seed complete.$(NC)"

sql-show:
	@echo "$(CYAN)Customers table:$(NC)"
	@docker exec debezium-postgres psql -U $(PG_USER) -d $(PG_DB) \
		-c "SELECT * FROM public.customers ORDER BY id;"

sql-shell: pg-shell

# ── Utilities ─────────────────────────────────────────────────────────────────
verify:
	@bash scripts/verify.sh

setup: up wait-connect connector-register
	@echo ""
	@echo "$(GREEN)=== Setup complete ===$(NC)"
	@echo "  PostgreSQL    → localhost:$(PG_PORT)"
	@echo "  Kafka (host)  → $(KAFKA_BROKER)"
	@echo "  Kafka Connect → $(CONNECT_URL)"
	@echo "  Kafka UI      → http://localhost:8090"
	@echo ""
	@$(MAKE) connector-status

clean: down-volumes
	@echo "$(GREEN)All containers, networks, and volumes removed.$(NC)"

# ── Reset (fixes InconsistentClusterIdException) ──────────────────────────────
# Use when Kafka fails with:
#   "InconsistentClusterIdException: The Cluster ID doesn't match stored clusterId"
# Cause: Zookeeper & Kafka volumes fall out of sync after a partial restart.
# Fix:   Wipe ALL volumes so both start fresh with the same cluster ID.
reset:
	@echo "$(RED)Resetting stack — wiping all volumes and restarting...$(NC)"
	$(COMPOSE) down -v
	$(COMPOSE) up -d
	@echo "$(GREEN)Stack reset complete. Run 'make migrate' to recreate DB schema.$(NC)"

