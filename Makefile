PROJECT_NAME="fizz-and-buzz"
PKG="gitlab.com/emixam23/$(PROJECT_NAME)"

.PHONY: all

all: build test lint
test: unit-tests integration-tests

lint: ## Lint the files
	@golint -set_exit_status ./...

unit-tests: ## Run unit tests with global code coverage report
	@go test -v ./internal/... -v -coverprofile .testCoverage.txt -p 1

integration-tests: ## Run integration tests
	docker-compose -f docker-compose.test-integration.yml up --abort-on-container-exit

build: ## Build the binary file
	@go build -o ./.build/$(PROJECT_NAME)

vendor: ## Vendor the dependencies to deploy to App Engine
	@go mod vendor

run: ## Run the project using Go SDK
	@go run main.go

start: ## Start the project using Docker
	docker-compose up --build

mock: ## Mock all the interface for integration/unit tests
	mockgen -source=internal/domain/infra/dal.go -destination=tests/mocks/mock_infra_dal.go -package=mocks
	mockgen -source=internal/domain/services/fnbservice/fnbservice.go -destination=tests/mocks/mock_service_fnb.go -package=mocks
	mockgen -source=internal/domain/services/statsservice/statsservice.go -destination=tests/mocks/mock_service_stats.go -package=mocks
	mockgen -source=internal/domain/ui/restapi.go -destination=tests/mocks/mock_ui_restapi.go -package=mocks
