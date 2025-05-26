include .env.db
export

MIGRATE = migrate -path db/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable"

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

create-migration:
	migrate create -ext sql -dir db/migrations $(name)
