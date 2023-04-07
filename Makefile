mysql:
	docker run --name db --network bank-network -e MYSQL_ROOT_PASSWORD=123 -e MYSQL_USER=root -p 3306:3306 -d mysql:latest

exec:
	docker exec -it db bash -c "mysql -u root -p"

createdb:
	docker exec -it db bash -c "create database easy_bank default CHARACTER SET UTF8;"

dropdb:
	docker exec -it db drop database easy_bank;

migrateup:
	migrate -path db/migration -database "mysql://root:123@tcp(localhost:3306)/easy_bank" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:123@tcp(localhost:3306)/easy_bank" -verbose down

migrateup1:
	migrate -path db/migration -database "mysql://root:123@tcp(localhost:3306)/easy_bank" -verbose up 1

migratedown1:
	migrate -path db/migration -database "mysql://root:123@tcp(localhost:3306)/easy_bank" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go easybank/db/sqlc Store

.PHONY: mysql exec createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock