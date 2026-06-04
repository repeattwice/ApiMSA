package main

import (
	servs "Api/Servs"
	"Api/Servs/broker"
	"os"
)

func main() {
	client, conn := servs.InitGRPCClient()
	defer conn.Close()

	App := &servs.App{UserClient: client}

	kafkaAddr := os.Getenv("KAFKA_BROKERS")
	if kafkaAddr == "" {
		kafkaAddr = "localhost:9092"
	}
	k := broker.NewProducer(kafkaAddr, "cart_events")

	kaf := servs.CartHandler{
		Kafka:      k,
		GrpcClient: client,
	}
	servs.Createserver(App, kaf)
}
