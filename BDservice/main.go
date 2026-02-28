package main

import (
	"context"
	servs "dbService/Servs"
	"dbService/user_pb"
	"net"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	conn, _ := servs.GetBDConnection(ctx)
	defer conn.Close(ctx)

	lis, _ := net.Listen("tcp", ":5051")
	s := grpc.NewServer()

	userServer := &servs.UserServer{
		DB: conn,
	}
	user_pb.RegisterUserServiceServer(s, userServer)
	s.Serve(lis)
}
