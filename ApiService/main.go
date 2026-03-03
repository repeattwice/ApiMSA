package main

import (
	servs "Api/Servs"
	"Api/Servs/broker"
)

func main() {
	client, conn := servs.InitGRPCClient()
	defer conn.Close()

	App := &servs.App{UserClient: client}

	k := broker.NewProducer("localhost:9092", "cort_events")

	kaf := servs.KafkaHandler{Kafka: k}
	servs.Createserver(App, kaf)

}
