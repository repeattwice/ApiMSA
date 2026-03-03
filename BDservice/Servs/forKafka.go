package servs

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

type UserCort struct {
	UserId    int    `json:"iser_id"`
	ItemName  string `json:"item_name"`
	ItemPrice int    `json:"item_price"`
}

func StartConsumer(conn *pgx.Conn, ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "cort_events",
		GroupID: "db-service-group",
	})
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		k := string(msg.Key)
		if k == "1" {
			cort := UserCort{}
			json.Unmarshal(msg.Value, &cort)
			sqlQuery := `
			INSERT INTO cort (item_name, item_price)
			VALUES ($1, $2, $3);
			`
			conn.Exec(ctx, sqlQuery, cort.UserId, cort.ItemName, cort.ItemPrice)
		}

	}
}
