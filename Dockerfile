# Create executable file.
FROM golang:1.20-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum Makefile ./
COPY server/pb/UserService.proto ./server/pb/
RUN apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install protobuf-compiler \
  && go mod download \
  && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
  && make protoc

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# Deploy
FROM debian:bullseye-slim as deploy

RUN apt-get -y update \
  && apt-get -y install curl tar \
  && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
  && mv ./migrate /usr/local/bin

COPY --from=deploy-builder /app/app .

EXPOSE 50051

CMD ["./app"]