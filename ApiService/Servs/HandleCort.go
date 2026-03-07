package servs

import (
	"Api/Servs/broker"
	"Api/user_pb"
	"encoding/json"
	"net/http"
	"strconv"
)

type CartHandler struct {
	Kafka      *broker.Producer
	GrpcClient user_pb.UserServiceClient
}

type UserCart struct {
	UserID    int    `json:"user_id"`
	ItemName  string `json:"item_name"`
	ItemPrice int    `json:"item_price"`
}

func (h *CartHandler) HandleShowAllItemsInCart(w http.ResponseWriter, r *http.Request) { //санек
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректный user_id"))
		return
	}
	req := &user_pb.GetCartRequest{
		UserId: int32(userID),
	}
	resp, err := h.GrpcClient.GetCart(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка при получении данных из БД-сервиса"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Items)
}

func (k *CartHandler) HandleAddToCart(w http.ResponseWriter, r *http.Request) { //Володя
	var cart UserCart
	err := json.NewDecoder(r.Body).Decode(&cart)
	WriteErrorBadReq(err, w, r)
	payload, _ := json.Marshal(cart)
	key := []byte("1")
	err = k.Kafka.SendMessage(r.Context(), key, payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b := []byte("Ошибка брокера")
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusAccepted)

}

func HandleChangePrice(w http.ResponseWriter, r *http.Request) { // саня

}

func HandleDeleteBuy(w http.ResponseWriter, r *http.Request) { //Вадим

}
