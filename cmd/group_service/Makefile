postgres:
	docker run --name group_db -e POSTGRES_USER=group_db -e POSTGRES_PASSWORD=group_db -p 5433:5432 -d postgres
postgresrm:
	docker stop group_db
	docker rm group_db

createdb:
	docker exec -it group_db createdb --username=group_db --owner=group_db group_db

dropdb:
	docker exec -it group_db dropdb --username=group_db simple_bank

migrateup:
	migrate -path internal/adapter/db/migration -database "postgres://group_db:group_db@localhost:5433/group_db?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/adapter/db/migration -database "postgres://group_db:group_db@localhost:5433/group_db?sslmode=disable" -verbose down
.PHONY: postgres createdb dropdb migrateup migratedown
# postgres:
# 	docker run --name postgres123 -p 5433:5433 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

# createdb:
# 	docker exec -it postgres123 createdb --username=root --owner=root simple_bank

# dropdb:
# 	docker exec -it postgres123 dropdb simple_bank

# migrateup:
# 	migrate -path gateway/adapter/db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up

# migratedown:
# 	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down

# .PHONY: postgres createdb dropdb migrateup migratedown