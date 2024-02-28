INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
TEST ?= $(shell go list ./... | grep -v -e vendor -e keys -e tmp)

BINARY:=health-worker
SYSTEM:=
BUILDOPTS:=
GO ?= GO111MODULE=on $(SYSTEM) go

dbcreate: ## Crate user and Create database
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Setup database$(RESET)"
	mysql -uroot -h127.0.0.1 -P3306 -e 'create database if not exists health_worker DEFAULT CHARACTER SET ujis;'

dbdrop: ## Drop database
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Drop database$(RESET)"
	mysql -uroot -h127.0.0.1 -P3306 -e 'drop database health_worker;'

dbmigrate: dbcreate
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Migrate database$(RESET)"
	sql-migrate up health_worker
depends:
	which sql-migrate &> /dev/null || go install github.com/rubenv/sql-migrate/...@latest
	which swag &> /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	go get -u github.com/swaggo/echo-swagger@latest
swag: depends
	swag init --parseInternal --parseDependency --parseDepth 2

run_register:
	$(GO) run main.go --config ./misc/develop.toml register
run_worker:
	$(GO) run main.go --config ./misc/develop.toml worker
run_api: swag
	$(GO) run main.go --config ./misc/develop.toml api

tidy:
	go mod tidy

test: tidy ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)

build_binary:
	$(GO) build $(BUILDOPTS) -ldflags="-s -w" -o $(BINARY)
