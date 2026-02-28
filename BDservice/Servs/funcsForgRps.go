package servs

import (
	"context"
	"dbService/user_pb"
	"net"

	"google.golang.org/grpc"
)

type UserServer struct {
	user_pb.UnimplementedUserServiceServer
}

func (s *UserServer) CreateAccount(ctx context.Context, req *user_pb.CreateAccountRequest) (*user_pb.CreateAccountResponse, error) {
	err := CreateAccount(req.UserName, req.LastName, req.Email, ctx)
	if err != nil {
		return &user_pb.CreateAccountResponse{Succes: false}, err
	}
	return &user_pb.CreateAccountResponse{Succes: true}, nil
}

func CreateAccount(user_name string, last_name string, email string, ctx context.Context) error {
	conn, err := GetBDConnection(ctx)
	if err != nil {
		return err
	}
	sqlQuery := `
	INSERT INTO users (user_name, last_name, email)
	SELECT $1, $2, $3
	`
	_, err = conn.Exec(ctx, sqlQuery, user_name, last_name, email)
	if err != nil {
		return err
	}
	return nil
}

func StartgRPCserver() {
	listen, _ := net.Listen("tcp", "5051")
	s := grpc.NewServer()
	user_pb.RegisterUserServiceServer(s, &UserServer{})
	s.Serve(listen)
}
