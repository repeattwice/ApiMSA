package servs

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

type UserCart struct {
	UserId    int    `json:"user_id"`
	ItemName  string `json:"item_name"`
	ItemPrice int    `json:"item_price"`
}

func StartConsumer(conn *pgx.Conn, ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "cart_events",
		GroupID: "db-service-group",
	})
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		k := string(msg.Key)
		if k == "1" {
			cart := UserCart{}
			json.Unmarshal(msg.Value, &cart)
			sqlQuery := `
			INSERT INTO cart (user_id, item_name_in_cart)
			VALUES ($1, $2);
			`
			conn.Exec(ctx, sqlQuery, cart.UserId, cart.ItemName)
		} else if k == "2" {
			type DeleteCartItem struct {
				UserId   string `json:"user_id"`
				ItemName string `json:"item_name"`
			}
			item := DeleteCartItem{}
			json.Unmarshal(msg.Value, &item)
			sqlQuery := `
			DELETE FROM cart
			WHERE user_id = $1 AND item_name_in_cart = $2
			`
			conn.Exec(ctx, sqlQuery, item.UserId, item.ItemName)
		}
	}
}
