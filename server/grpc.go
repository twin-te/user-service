package server

import (
	"github.com/twin-te/user-service/server/pb"
	"github.com/twin-te/user-service/usecase"
	"google.golang.org/grpc"
)

func NewServer(u usecase.UseCase) *grpc.Server {
	s := grpc.NewServer()
	srv := &userServiceServer{u: u}
	pb.RegisterUserServiceServer(s, srv)
	return s
}
