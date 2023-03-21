package main

import (
	"context"

	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/twin-te/user-service/repository"
	"github.com/twin-te/user-service/server"
	"github.com/twin-te/user-service/server/pb"
	"github.com/twin-te/user-service/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func createUser(t *testing.T, c pb.UserServiceClient) (string, *pb.Authentication) {
	t.Helper()

	a := pb.Authentication{Provider: pb.Provider_Google, SocialId: uuid.NewString()}
	res, err := c.GetOrCreateUser(context.Background(), &pb.GetOrCreateUserRequest{Provider: a.Provider, SocialId: a.SocialId})
	if err != nil {
		t.Fatalf("GetOrCreateUser failed: %v", err)
	}
	return res.GetId(), &a
}

func getUserByID(t *testing.T, c pb.UserServiceClient, id string) []*pb.Authentication {
	t.Helper()

	res, err := c.GetUser(context.Background(), &pb.GetUserRequest{Id: id})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	return res.Authentications
}

func compareAuthentications(t *testing.T, sa1, sa2 []*pb.Authentication) {
	t.Helper()

	if len(sa1) != len(sa2) {
		t.Fatalf("The given authentications are not same: %d != %d", len(sa1), len(sa2))
	}

	for i := range sa1 {
		if sa1[i].Provider != sa2[i].Provider || sa1[i].SocialId != sa2[i].SocialId {
			t.Fatalf("The given authentications are not same: %+v %+v", sa1[i], sa2[i])
		}
	}
}

func TestMain(t *testing.T) {
	// Start Server
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	r, close, err := repository.NewRepository()
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	u := usecase.NewUseCase(r)
	s := server.NewServer(u)

	t.Logf("server listening at %v", l.Addr())
	go s.Serve(l)

	// Prepare Client
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	t.Run("GetOrCreateUser", func(t *testing.T) {
		createUser(t, c)
	})

	t.Run("GetUser", func(t *testing.T) {
		_, err = c.GetUser(context.Background(), &pb.GetUserRequest{Id: uuid.NewString()})
		if status.Code(err) != codes.NotFound {
			t.Fatalf("not created user: %s %s", status.Code(err), codes.NotFound)
		}

		id, a := createUser(t, c)

		res, err := c.GetUser(context.Background(), &pb.GetUserRequest{Id: id})
		if err != nil {
			t.Fatalf("GetUser failed: %v", err)
		}

		if res.GetId() != id {
			t.Fatalf("Invalid ID, %s %s", id, res.GetId())
		}
		compareAuthentications(t, []*pb.Authentication{a}, res.Authentications)
	})

	t.Run("AddAuthentication", func(t *testing.T) {
		a1 := &pb.Authentication{Provider: pb.Provider_Apple, SocialId: uuid.NewString()}

		_, err := c.AddAuthentication(context.Background(), &pb.AddAuthenticationRequest{Id: uuid.NewString(), Provider: a1.Provider, SocialId: a1.SocialId})
		if status.Code(err) != codes.NotFound {
			t.Fatalf("not existing user: %s %s", status.Code(err), codes.NotFound)
		}

		id, a0 := createUser(t, c)
		_, err = c.AddAuthentication(context.Background(), &pb.AddAuthenticationRequest{Id: id, Provider: a1.Provider, SocialId: a1.SocialId})
		if err != nil {
			t.Fatalf("AddAuthentication failed: %v", err)
		}

		sa := getUserByID(t, c, id)
		compareAuthentications(t, []*pb.Authentication{a0, a1}, sa)

		_, err = c.AddAuthentication(context.Background(), &pb.AddAuthenticationRequest{Id: id, Provider: a1.Provider, SocialId: a1.SocialId})
		if status.Code(err) != codes.AlreadyExists {
			t.Fatalf("registered authentication: %s %s", status.Code(err), codes.AlreadyExists)
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		id, _ := createUser(t, c)

		_, err := c.DeleteUser(context.Background(), &pb.DeleteUserRequest{Id: id})
		if err != nil {
			t.Fatalf("DeleteUser failed: %v", err)
		}

		_, err = c.GetUser(context.Background(), &pb.GetUserRequest{Id: id})
		if status.Code(err) != codes.NotFound {
			t.Fatalf("deleted user: %s %s", status.Code(err), codes.NotFound)
		}

		_, err = c.DeleteUser(context.Background(), &pb.DeleteUserRequest{Id: id})
		if status.Code(err) != codes.NotFound {
			t.Fatalf("deleted user: %s %s", status.Code(err), codes.NotFound)
		}
	})
}
