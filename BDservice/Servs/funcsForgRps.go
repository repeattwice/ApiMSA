package servs

import (
	"context"
	"dbService/user_pb"
	"net"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

type UserServer struct {
	user_pb.UnimplementedUserServiceServer
	DB *pgx.Conn
}

func (s *UserServer) CreateAccount(ctx context.Context, req *user_pb.CreateAccountRequest) (*user_pb.CreateAccountResponse, error) {
	err := CreateAccount(req.UserName, req.LastName, req.Email, ctx, s.DB)
	if err != nil {
		return &user_pb.CreateAccountResponse{Succes: false}, err
	}
	return &user_pb.CreateAccountResponse{Succes: true}, nil
}

func CreateAccount(user_name string, last_name string, email string, ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	INSERT INTO users (user_name, last_name, email)
	SELECT $1, $2, $3
	`
	_, err := conn.Exec(ctx, sqlQuery, user_name, last_name, email)
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

func (s *UserServer) Avtorization(ctx context.Context, req *user_pb.AvtorizationRequest) (*user_pb.AvtorizationResponse, error) {
	Avtorization(req.UserName, req.LastName, ctx)
	return &user_pb.AvtorizationResponse{IsUserExists: true, IsLactNameIsCorrect: true}, nil
}

func Avtorization(user_name string, last_name string, ctx context.Context) {

}
