package servs

import (
	"context"
	"encoding/json"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

type UserCart struct {
	UserId    int    `json:"user_id"`
	ItemName  string `json:"item_name"`
	ItemPrice int    `json:"item_price"`
}

func StartConsumer(conn *pgx.Conn, ctx context.Context) {
	kafkaAddr := os.Getenv("KAFKA_BROKERS")
	if kafkaAddr == "" {
		kafkaAddr = "localhost:9092"
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaAddr},
		Topic:   "cart_events",
		GroupID: "db-service-group",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		k := string(msg.Key)
		if k == "1" {
			// AddToCart
			cart := UserCart{}
			json.Unmarshal(msg.Value, &cart)
			sqlQuery := `
			INSERT INTO items (item_name, item_price)
			VALUES ($1, $2)
			ON CONFLICT (item_name) DO UPDATE
			SET item_price = EXCLUDED.item_price;

			INSERT INTO cart (user_id, item_name_in_cart)
			VALUES ($3, $1);
			`
			conn.Exec(ctx, sqlQuery, cart.ItemName, cart.ItemPrice, cart.UserId)

		} else if k == "2" {
			// DeleteBuyFromCart
			type DeleteCartItem struct {
				UserId   int    `json:"user_id"`
				ItemName string `json:"item_name"`
			}
			item := DeleteCartItem{}
			json.Unmarshal(msg.Value, &item)
			sqlQuery := `
			DELETE FROM cart
			WHERE user_id = $1 AND item_name_in_cart = $2
			`
			conn.Exec(ctx, sqlQuery, item.UserId, item.ItemName)

		} else if k == "3" {
			// ChangePrice
			type ChangePriceItem struct {
				ItemName string `json:"item_name"`
				NewPrice int    `json:"item_price"`
			}
			item := ChangePriceItem{}
			json.Unmarshal(msg.Value, &item)
			sqlQuery := `
			UPDATE items
			SET item_price = $1
			WHERE item_name = $2
			`
			conn.Exec(ctx, sqlQuery, item.NewPrice, item.ItemName)
		}
	}
}
