package servs

import (
	"context"
	"encoding/json"
	"fmt"
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
			fmt.Println("Ошибка чтения сообщения Kafka:", err)
			break
		}
		k := string(msg.Key)
		fmt.Printf("Получено сообщение Kafka: key=%q value=%s\n", k, string(msg.Value))
		if k == "1" {
			// AddToCart
			cart := UserCart{}
			if err := json.Unmarshal(msg.Value, &cart); err != nil {
				fmt.Println("Ошибка разбора AddToCart:", err)
				continue
			}
			itemQuery := `
			INSERT INTO items (item_name, item_price)
			VALUES ($1, $2)
			ON CONFLICT (item_name) DO UPDATE
			SET item_price = EXCLUDED.item_price;
			`
			cartQuery := `
			INSERT INTO cart (user_id, item_name_in_cart)
			VALUES ($1, $2);
			`
			if _, err := conn.Exec(ctx, itemQuery, cart.ItemName, cart.ItemPrice); err != nil {
				fmt.Println("Ошибка записи товара:", err)
				continue
			}
			if _, err := conn.Exec(ctx, cartQuery, cart.UserId, cart.ItemName); err != nil {
				fmt.Println("Ошибка добавления товара в корзину:", err)
				continue
			}
			fmt.Printf("Товар добавлен в корзину: user_id=%d item=%s\n", cart.UserId, cart.ItemName)

		} else if k == "2" {
			// DeleteBuyFromCart
			type DeleteCartItem struct {
				UserId   int    `json:"user_id"`
				ItemName string `json:"item_name"`
			}
			item := DeleteCartItem{}
			if err := json.Unmarshal(msg.Value, &item); err != nil {
				fmt.Println("Ошибка разбора DeleteBuyFromCart:", err)
				continue
			}
			sqlQuery := `
			DELETE FROM cart
			WHERE user_id = $1 AND item_name_in_cart = $2
			`
			if _, err := conn.Exec(ctx, sqlQuery, item.UserId, item.ItemName); err != nil {
				fmt.Println("Ошибка удаления товара из корзины:", err)
			}

		} else if k == "3" {
			// ChangePrice
			type ChangePriceItem struct {
				ItemName string `json:"item_name"`
				NewPrice int    `json:"item_price"`
			}
			item := ChangePriceItem{}
			if err := json.Unmarshal(msg.Value, &item); err != nil {
				fmt.Println("Ошибка разбора ChangePrice:", err)
				continue
			}
			sqlQuery := `
			UPDATE items
			SET item_price = $1
			WHERE item_name = $2
			`
			if _, err := conn.Exec(ctx, sqlQuery, item.NewPrice, item.ItemName); err != nil {
				fmt.Println("Ошибка изменения цены:", err)
			}
		}
	}
}
