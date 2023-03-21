package server

import (
	"context"
	"errors"
	"log"

	"github.com/twin-te/user-service/entity"
	"github.com/twin-te/user-service/server/converter"
	"github.com/twin-te/user-service/server/pb"
	"github.com/twin-te/user-service/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		log.Printf("Error userServiceServer.GetOrCreateUser %+v: %v", in, err)
		return nil, status.Error(codes.Internal, "サーバー内で問題が発生しました")
	}
	return &pb.GetOrCreateUserResponse{Id: u.ID}, nil
}

func (s *userServiceServer) AddAuthentication(ctx context.Context, in *pb.AddAuthenticationRequest) (*pb.AddAuthenticationResponse, error) {
	a := &entity.Authentication{Provider: converter.FromPBProvider(in.Provider), SocialID: in.SocialId}
	err := s.u.AddAuthentication(ctx, in.Id, a)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			return nil, status.Error(codes.NotFound, "指定されたユーザーが見つかりませんでした")
		case errors.Is(err, entity.ErrAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "認証情報が既に登録されています")
		default:
			log.Printf("Error userServiceServer.AddAuthentication %+v: %v", in, err)
			return nil, status.Error(codes.Internal, "サーバー内で問題が発生しました")
		}
	}
	return &pb.AddAuthenticationResponse{}, nil
}

func (s *userServiceServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	u, err := s.u.GetUserByID(ctx, in.Id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			return nil, status.Error(codes.NotFound, "指定されたユーザーが見つかりませんでした")
		default:
			log.Printf("Error userServiceServer.GetUser %+v: %v", in, err)
			return nil, status.Error(codes.Internal, "サーバー内で問題が発生しました")
		}
	}
	return &pb.GetUserResponse{Id: u.ID, Authentications: converter.ToPBAuthentications(u.Authentications)}, nil
}

func (s *userServiceServer) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.u.DeleteUserByID(ctx, in.Id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			return nil, status.Error(codes.NotFound, "指定されたユーザーが見つかりませんでした")
		default:
			log.Printf("Error userServiceServer.DeleteUser %+v: %v", in, err)
			return nil, status.Error(codes.Internal, "サーバー内で問題が発生しました")
		}
	}
	return &pb.DeleteUserResponse{}, nil
}
