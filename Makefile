mysql:
	docker run --network=bank-network --name easybank -p 4000:4000 -e DB_SOURCE="root:123@tcp(db:3306)/easy_bank?parseTime=true" easybank:latest

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

dockerdown:
	docker compose down&&docker rmi easybank_api

proto:
	rm -rf pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: mysql exec createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock dockerdown proto