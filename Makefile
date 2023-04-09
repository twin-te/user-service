bash-app:
	docker exec -it twinte-user-service-app /bin/bash

protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    server/pb/UserService.proto

test:
	go test -count=1 ./...

psql:
	psql -h ${PG_HOST} -p ${PG_PORT} -U ${PG_USERNAME} -d ${PG_DATABASE}

migrate:
	migrate -database "postgres://${PG_USERNAME}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DATABASE}?sslmode=${PG_SSLMODE}" -path db/migrations up

migrate-up-one:
	migrate -database "postgres://${PG_USERNAME}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DATABASE}?sslmode=${PG_SSLMODE}" -path db/migrations up 1

migrate-down-one:
	migrate -database "postgres://${PG_USERNAME}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DATABASE}?sslmode=${PG_SSLMODE}" -path db/migrations down 1
