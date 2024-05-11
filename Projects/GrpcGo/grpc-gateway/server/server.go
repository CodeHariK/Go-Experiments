package server

import (
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	gw "grpcgateway/proto"
)

type server struct {
	mu    *sync.RWMutex
	users []*gw.User
	gw.UnimplementedGreeterServer
	gw.UnimplementedUserServiceServer
}

func New() *server {
	return &server{
		mu: &sync.RWMutex{},
	}
}

func (s *server) SayHello(ctx context.Context, in *gw.HelloRequest) (*gw.HelloReply, error) {
	return &gw.HelloReply{Message: in.Name + " world"}, nil
}

// AddUser adds a user to the in-memory store.
func (s *server) AddUser(ctx context.Context, _ *gw.AddUserRequest) (*gw.AddUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &gw.User{
		Id: uuid.Must(uuid.NewUUID()).String(),
	}
	s.users = append(s.users, user)

	return &gw.AddUserResponse{
		User: user,
	}, nil
}

func (s *server) GetUser(req *gw.GetUserRequest, srv gw.UserService_GetUserServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx := slices.IndexFunc(s.users, func(c *gw.User) bool { return c.Id == req.Id })

	if idx != -1 {
		err := srv.Send(&gw.GetUserResponse{
			User: s.users[idx],
		})
		if err != nil {
			return err
		}
		return nil
	}
	return status.Errorf(codes.NotFound, "user not found")
}

// ListUsers lists all users in the store.
func (s *server) ListUsers(_ *gw.ListUsersRequest, srv gw.UserService_ListUsersServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		err := srv.Send(&gw.ListUsersResponse{User: user})
		if err != nil {
			return err
		}
	}

	return nil
}
