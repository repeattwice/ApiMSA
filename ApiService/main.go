package main

import (
	servs "Api/Servs"
	"Api/Servs/broker"
)

func main() {
	client, conn := servs.InitGRPCClient()
	defer conn.Close()

	App := &servs.App{UserClient: client}

	k := broker.NewProducer("localhost:9092", "cart_events")

	kaf := servs.CartHandler{Kafka: k}
	servs.Createserver(App, kaf)

}
