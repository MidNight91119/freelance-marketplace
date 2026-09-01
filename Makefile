DB_URL=postgresql://root:secret@localhost:5432/freelance_marketplace?sslmode=disable

postgres:
	docker run --name freelance-marketplace-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine

createdb:
	docker exec -it freelance-marketplace-postgres createdb --username=root --owner=root freelance_marketplace

dropdb:
	docker exec -it freelance-marketplace-postgres dropdb freelance_marketplace

migrateup:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir internal/db/migration -seq $(name)

db_schema:
	dbml2sql doc/db.dbml --postgres -o doc/schema.sql

sqlc:
	sqlc generate

# always run AFTER `make sqlc` — a changed Querier makes the mock stale
mock:
	mockgen -package mockdb -destination internal/db/mock/store.go \
		github.com/MidNight91119/freelance-marketplace/internal/db/sqlc Store

test:
	go test -v -cover -short ./...

server:
	go run ./cmd/api


.PHONY: postgres createdb dropdb migrateup migratedown db_schema sqlc mock test server new_migration