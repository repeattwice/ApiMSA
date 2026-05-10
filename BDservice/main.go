package main

import (
	"context"
	servs "dbService/Servs"
	"dbService/user_pb"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	conn, err := servs.GetBDConnection(ctx)
	if err != nil {
		fmt.Println("Ошибка подключения к БД:", err)
		return
	}
	defer conn.Close(ctx)

	servs.CreateTables(ctx, conn)

	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		fmt.Println("Ошибка запуска listener:", err)
		return
	}
	s := grpc.NewServer()

	userServer := &servs.UserServer{
		DB: conn,
	}
	user_pb.RegisterUserServiceServer(s, userServer)

	// Kafka consumer runs concurrently
	go servs.StartConsumer(userServer.DB, ctx)

	fmt.Println("BDservice запущен на порту :5051")
	s.Serve(lis)
}
