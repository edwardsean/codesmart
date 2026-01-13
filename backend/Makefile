build:
	@go build -o bin/auth-service cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/auth-service

#Database Migration

MIGRATION_DIR = ./cmd/migrate/migrations
DB_URL = postgres://dev:SamsungEd123.@postgres:5432/codesmart?sslmode=disable

migrate_create :
	@docker run --rm -v $(PWD)/$(MIGRATION_DIR):/migrations migrate/migrate \
		create -ext sql -dir /migrations $(name)

migrate_up :
	@docker run --rm -v $(PWD)/$(MIGRATION_DIR):/migrations --network code-smart_default migrate/migrate \
		-path=/migrations/ -database "$(DB_URL)" up

migrate_down :
	@docker run --rm -v $(PWD)/$(MIGRATION_DIR):/migrations --network code-smart_default migrate/migrate \
		-path=/migrations/ -database "$(DB_URL)" down 1

migrate_force:
	@docker run --rm -v $(PWD)/$(MIGRATION_DIR):/migrations --network code-smart_default migrate/migrate \
		-path=/migrations -database "$(DB_URL)" force $(version)

access_database:
	@docker exec -it code-smart-postgres-1 psql -U dev -d codesmart
