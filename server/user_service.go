package server

import (
	"context"
	"log"

	"github.com/twin-te/user-service/entity"
	"github.com/twin-te/user-service/server/converter"
	"github.com/twin-te/user-service/server/pb"
	"github.com/twin-te/user-service/usecase"
)

// userServiceServer is used to implement pb.UserServiceServer.
type userServiceServer struct {
	u usecase.UseCase
	pb.UnimplementedUserServiceServer
}

func (s *userServiceServer) GetOrCreateUser(ctx context.Context, in *pb.GetOrCreateUserRequest) (*pb.GetOrCreateUserResponse, error) {
	a := &entity.Authentication{Provider: converter.FromPBProvider(in.Provider), SocialID: in.SocialId}
	u, err := s.u.GetOrCreateUser(ctx, a)
	if err != nil {
		return nil, converter.ToGRPCError(err, func(err error) { log.Printf("Error userServiceServer.GetOrCreateUser %+v -> %v", in, err) })
	}
	return &pb.GetOrCreateUserResponse{Id: u.ID}, nil
}

func (s *userServiceServer) AddAuthentication(ctx context.Context, in *pb.AddAuthenticationRequest) (*pb.AddAuthenticationResponse, error) {
	a := &entity.Authentication{Provider: converter.FromPBProvider(in.Provider), SocialID: in.SocialId}
	err := s.u.AddAuthentication(ctx, in.Id, a)
	if err != nil {
		return nil, converter.ToGRPCError(err, func(err error) { log.Printf("Error userServiceServer.AddAuthentication %+v -> %v", in, err) })
	}
	return &pb.AddAuthenticationResponse{}, nil
}

func (s *userServiceServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	u, err := s.u.GetUserByID(ctx, in.Id)
	if err != nil {
		return nil, converter.ToGRPCError(err, func(err error) { log.Printf("Error userServiceServer.GetUser %+v -> %v", in, err) })
	}
	return &pb.GetUserResponse{Id: u.ID, Authentications: converter.ToPBAuthentications(u.Authentications)}, nil
}

func (s *userServiceServer) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.u.DeleteUserByID(ctx, in.Id)
	if err != nil {
		return nil, converter.ToGRPCError(err, func(err error) { log.Printf("Error userServiceServer.DeleteUser %+v -> %v", in, err) })
	}
	return &pb.DeleteUserResponse{}, nil
}

func (s *userServiceServer) DeleteAuthentication(ctx context.Context, in *pb.DeleteAuthenticationRequest) (*pb.DeleteAuthenticationResponse, error) {
	err := s.u.DeleteAuthentication(ctx, in.Id, converter.FromPBProvider(in.Provider))
	if err != nil {
		return nil, converter.ToGRPCError(err, func(err error) { log.Printf("Error userServiceServer.DeleteAuthentication %+v -> %v", in, err) })
	}
	return &pb.DeleteAuthenticationResponse{}, nil
}
