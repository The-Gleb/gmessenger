postgres:
	docker run --name gmessenger_db -e POSTGRES_USER=gateway_db -e POSTGRES_PASSWORD=gateway_db -p 5434:5432 -d postgres:alpine
postgresrm:
	docker stop gmessenger_db
	docker rm gmessenger_db

createdb:
	docker exec -it gmessenger_db createdb --username=gateway_db --owner=gateway_db gateway_db

dropdb:
	docker exec -it gmessenger_db dropdb --username=gateway_db simple_bank

migrateup:
	migrate -path internal/adapter/db/migration -database "postgres://gateway_db:gateway_db@localhost:5434/gateway_db?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/adapter/db/migration -database "postgres://gateway_db:gateway_db@localhost:5434/gateway_db?sslmode=disable" -verbose down
.PHONY: postgres createdb dropdb migrateup migratedown
