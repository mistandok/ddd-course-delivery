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