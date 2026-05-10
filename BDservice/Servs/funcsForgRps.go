package servs

import (
	"context"
	"dbService/user_pb"
	"fmt"

	"github.com/jackc/pgx/v5"
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
	VALUES ($1, $2, $3)
	`
	_, err := conn.Exec(ctx, sqlQuery, user_name, last_name, email)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServer) Avtorization(ctx context.Context, req *user_pb.AvtorizationRequest) (*user_pb.AvtorizationResponse, error) {
	UserName, LastName := Avtorization(req.UserName, req.LastName, ctx, s.DB)
	return &user_pb.AvtorizationResponse{IsUserExists: UserName, IsLactNameIsCorrect: LastName}, nil
}

func Avtorization(user_name string, last_name string, ctx context.Context, conn *pgx.Conn) (bool, bool) {
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

func (s *UserServer) GetCart(ctx context.Context, req *user_pb.GetCartRequest) (*user_pb.GetCartResponse, error) {
	sqlQuery := `
	SELECT c.item_name_in_cart, i.item_price 
        FROM cart c
        JOIN items i ON c.item_name_in_cart = i.item_name
        WHERE c.user_id = $1`

	rows, err := s.DB.Query(ctx, sqlQuery, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*user_pb.CartItem
	for rows.Next() {
		var name string
		var price int32
		if err := rows.Scan(&name, &price); err != nil {
			return nil, err
		}
		items = append(items, &user_pb.CartItem{ItemName: name, ItemPrice: price})
	}
	return &user_pb.GetCartResponse{Items: items}, nil
}

func (s *UserServer) DeleteAccount(ctx context.Context, req *user_pb.DeleteAccountRequest) (*user_pb.DeleteAccountResponse, error) {
	err := DeleteAccount(req.UserName, req.LastName, ctx, s.DB)
	if err != nil {
		return &user_pb.DeleteAccountResponse{Succes: false}, err
	}
	return &user_pb.DeleteAccountResponse{Succes: true}, nil
}

func DeleteAccount(user_name string, last_name string, ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	DELETE FROM users
	WHERE user_name = $1 AND last_name = $2
	`
	tag, err := conn.Exec(ctx, sqlQuery, user_name, last_name)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("пользователь не найден")
	}
	return nil
}
