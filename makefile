define setup_env
	$(eval ENV_FILE := ./deploy/env/.env.$(1))
	@echo "- setup env $(ENV_FILE)"
	$(eval include ./deploy/env/.env.$(1))
	$(eval export)
endef

setup-local-env:
	$(call setup_env,local)

APP_NAME=delivery

.PHONY: build test
build: test ## Build application
	mkdir -p build
	go build -o build/${APP_NAME} cmd/app/main.go

test: ## Run tests
	go test ./...

generate-server:
	@go tool oapi-codegen -config configs/server.cfg.yaml https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/delivery/contracts/openapi.yml

generate-geo-client:
	@rm -rf internal/generated/clients/geosrv
	@curl -s -o configs/geo.proto https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/geo/contracts/contract.proto
	@protoc --go_out=internal/generated/clients --go-grpc_out=internal/generated/clients configs/geo.proto

generate-basket-queue:
	@rm -rf internal/generated/queues/basketconfirmedpb
	@curl -s -o configs/basket_confirmed.proto https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/basket/contracts/basket_confirmed.proto
	@protoc --go_out=internal/generated --go-grpc_out=internal/generated configs/basket_confirmed.proto

generate-order-queue:
	@rm -rf internal/generated/queues/orderstatuschangedpb
	@curl -s -o configs/order_status_changed.proto https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/delivery/contracts/order_status_changed.proto
	@protoc --go_out=internal/generated --go-grpc_out=internal/generated configs/order_status_changed.proto


# Команды для работы с миграциями

migration-status:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} status -v

migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} down -v

create-migration:
	goose -dir ${MIGRATION_DIR} create $(migration_name) sql


# Команды для работы с миграциями в локальном окружении
local-create-new-migration: setup-local-env create-migration

local-migration-status: setup-local-env migration-status

local-migration-up: setup-local-env migration-up

local-migration-down: setup-local-env migration-down

# Локальный старт окружения (само приложение не стартует)
local-down-app:
	docker-compose --env-file deploy/env/.env.local -f docker-compose.local.yaml down -v

local-start-app:
	docker-compose --env-file deploy/env/.env.local -f docker-compose.local.yaml up -d --build