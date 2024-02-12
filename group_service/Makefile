postgres:
	docker run --name group_service_db -e POSTGRES_USER=group_service -e POSTGRES_PASSWORD=group_service -p 5435:5432 -d postgres
postgresrm:
	docker stop group_service_db
	docker rm group_service_db

createdb:
	docker exec -it group_service_db createdb --username=group_service --owner=group_service group_service

dropdb:
	docker exec -it group_service_db dropdb --username=group_service simple_bank

migrateup:
	migrate -path internal/adapter/db/migration -database "postgres://group_service:group_service@localhost:5435/group_service?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/adapter/db/migration -database "postgres://group_service:group_service@localhost:5435/group_service?sslmode=disable" -verbose down
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