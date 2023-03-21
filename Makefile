bash-app:
	docker exec -it twinte-user-service-app /bin/bash

protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    server/pb/UserService.proto

psql:
	psql -h ${PG_HOST} -p ${PG_PORT} -U ${PG_USERNAME} -d ${PG_DATABASE}

setup-db:
	psql -h ${PG_HOST} -p ${PG_PORT} -U ${PG_USERNAME} -d ${PG_DATABASE} -f repository/setup.sql

test:
	go test -count=1 ./...
