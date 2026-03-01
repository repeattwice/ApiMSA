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
	UserName, LastName := Avtorization(req.UserName, req.LastName, ctx, s.DB)
	return &user_pb.AvtorizationResponse{IsUserExists: UserName, IsLactNameIsCorrect: LastName}, nil
}

func Avtorization(user_name string, last_name string, ctx context.Context, conn *pgx.Conn) (bool, bool) { // сначала IsUserExists
	sqlQuery := `
	SELECT EXISTS(
	SELECT 1
	FROM users
	WHERE user_name = $1
	)
	`
	sqlQuery1 := `
	SELECT EXISTS(
	SELECT 1
	FROM users
	WHERE last_name = $1
	)
	`
	var IsUserExists bool
	var IsLastNameCorrect bool

	err := conn.QueryRow(ctx, sqlQuery, user_name).Scan(&IsUserExists)
	err1 := conn.QueryRow(ctx, sqlQuery1, last_name).Scan(&IsLastNameCorrect)

	if err != nil && err1 != nil {
		return false, false
	} else if err != nil && err1 == nil {
		return false, true
	} else if err == nil && err1 != nil {
		return true, false
	}
	return IsUserExists, IsLastNameCorrect
}
