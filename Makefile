DB_MIGRATE_URL = postgres://login:pass@localhost:5432/postgres?sslmode=disable
MIGRATE_PATH = ./migration/postgres

run:
	go run ./cmd/app

lint:
	golangci-lint run

up:
	docker compose up --build --force-recreate

down:
	docker compose down