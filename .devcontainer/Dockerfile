FROM mcr.microsoft.com/vscode/devcontainers/go:1.20

WORKDIR /go/src

RUN sudo apt-get -y update \
  && sudo apt-get -y upgrade \
  && sudo apt-get -y install protobuf-compiler \
  && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
  && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest