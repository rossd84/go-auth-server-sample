include ./environments/.env.db
export

# Go commands
run:
	go run ./cmd/server

test:
	go test ./...

# Migrate
MIGRATE = migrate -path internal/db/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable"

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

create-migration:
	migrate create -ext sql -dir internal/db/migrations $(name)
