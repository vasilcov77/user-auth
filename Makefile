DB_MIGRATE_URL = postgres://login:pass@localhost:5432/postgres?sslmode=disable
MIGRATE_PATH = ./migration/postgres

run:
	go run ./cmd/app

lint:
	golangci-lint run

up:
	docker compose up -d --build --force-recreate

down:
	docker compose down

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-create:
	@read -p "Name:" name; \
	migrate create -ext sql -dir "$(MIGRATE_PATH)" $$name

migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all

.PHONY: test
test:
	go test -v -cover ./...

integration-test:
	go test -count=1 -v -tags=integration ./test/integration