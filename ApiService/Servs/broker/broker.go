package broker

//продюсер
import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(addr string, topic string) *Producer {
	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(addr),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) SendMessage(ctx context.Context, key []byte, value []byte) error {
	return p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}
