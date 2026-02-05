help:
	@echo ''
	@echo 'Usage: make [TARGET] [EXTRA_ARGUMENTS]'
	@echo 'Targets:'
	@echo '  make dev           - development (docker-compose up)'
	@echo '  make build        - build container'
	@echo '  make production   - docker production build'
	@echo '  make clean        - remove docker volumes/images'
	@echo '  make run          - run application (go run main.go)'
	@echo '  make test         - run unit tests'
	@echo '  make test-coverage - run tests with coverage report'
	@echo ''

dev:
	if [ ! -f .env ]; then cp .env.example .env; fi;
	docker-compose -f docker-compose-dev.yml down
	docker-compose -f docker-compose-dev.yml up

build:
	docker-compose -f docker-compose-prod.yml build
	docker-compose -f docker-compose-dev.yml down build

production:
	docker-compose -f docker-compose-prod.yml up -d --build

clean:
	docker-compose -f docker-compose-prod.yml down -v
	docker-compose -f docker-compose-dev.yml down -v

run:
	go run main.go

test:
	go test ./tests/unit/... ./internal/... ./pkg/...

test-coverage:
	go test -coverprofile=coverage.out ./tests/unit/... ./internal/... ./pkg/...
	go tool cover -func=coverage.out

# test-coverage-check runs tests with coverage and fails if total coverage is below 70%.
# Requires a Go toolchain that supports merged coverage (see tests/README.md).
test-coverage-check:
	go test -coverprofile=coverage.out ./tests/unit/... ./internal/... ./pkg/...
	@go tool cover -func=coverage.out | grep '^total:' | awk '{gsub(/%/,""); if ($$3 < 70) { print "Coverage " $$3 "% is below 70%"; exit 1 } }'
	go tool cover -func=coverage.out
