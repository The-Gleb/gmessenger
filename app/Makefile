postgres:
	docker run --name gmessenger_db -e POSTGRES_USER=gmessenger_gateway -e POSTGRES_PASSWORD=gmessenger_gateway -p 5434:5432 -d postgres
postgresrm:
	docker stop gmessenger_db
	docker rm gmessenger_db

createdb:
	docker exec -it gmessenger_db createdb --username=gmessenger_gateway --owner=gmessenger_gateway gmessenger_gateway

dropdb:
	docker exec -it gmessenger_db dropdb --username=gmessenger_gateway simple_bank

migrateup:
	migrate -path internal/adapter/db/migration -database "postgres://gmessenger_gateway:gmessenger_gateway@localhost:5434/gmessenger_gateway?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/adapter/db/migration -database "postgres://gmessenger_gateway:gmessenger_gateway@localhost:5434/gmessenger_gateway?sslmode=disable" -verbose down
.PHONY: postgres createdb dropdb migrateup migratedown
# postgres:
# 	docker run --name postgres123 -p 5433:5433 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

# createdb:
# 	docker exec -it postgres123 createdb --username=root --owner=root simple_bank

# dropdb:
# 	docker exec -it postgres123 dropdb simple_bank

# migrateup:
# 	migrate -path internal/adapter/db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up

# migratedown:
# 	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down

# .PHONY: postgres createdb dropdb migrateup migratedown