postgres:
	docker run --name db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -d postgres:15-alpine

createdb:
	docker exec -it db createdb --username=root --owner=root easy_bank

dropdb:
	docker exec -it db dropdb easy_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:123@localhost:5432/easy_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123@localhost:5432/easy_bank?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown