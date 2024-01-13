sqlc:
	go get github.com/sqlc-dev/sqlc/cmd/sqlc
devdb:
	docker run --rm -d --name sms_db -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres postgres:16-alpine3.19
dropdb:
	docker stop sms_db
	docker rm sms_db
migrateinit:
	migrate create -ext sql -dir db/migrations -seq init_schema
migrateup:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/sms_db?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/sms_db?sslmode=disable" -verbose down