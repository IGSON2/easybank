mysql:
	docker run --name db -e MYSQL_ROOT_PASSWORD=123 -p 3306:3306 -d mysql:latest

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

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: mysql exec createdb dropdb migrateup migratedown sqlc test server