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

type DeleteCartRequest struct {
	UserID   int    `json:"user_id"`
	ItemName string `json:"item_name"`
}

type ChangePriceRequest struct {
	ItemName string `json:"item_name"`
	NewPrice int    `json:"item_price"`
}

func (h *CartHandler) HandleShowAllItemsInCart(w http.ResponseWriter, r *http.Request) {
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

func (k *CartHandler) HandleAddToCart(w http.ResponseWriter, r *http.Request) {
	var cart UserCart
	err := json.NewDecoder(r.Body).Decode(&cart)
	WriteErrorBadReq(err, w, r)
	if err != nil {
		return
	}
	if cart.UserID <= 0 || cart.ItemName == "" || cart.ItemPrice <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректные данные товара"))
		return
	}
	payload, _ := json.Marshal(cart)
	key := []byte("1")
	err = k.Kafka.SendMessage(r.Context(), key, payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка брокера"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (k *CartHandler) HandleChangePrice(w http.ResponseWriter, r *http.Request) {
	var req ChangePriceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	WriteErrorBadReq(err, w, r)
	if err != nil {
		return
	}
	if req.ItemName == "" || req.NewPrice <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректные данные цены"))
		return
	}
	payload, _ := json.Marshal(req)
	key := []byte("3")
	err = k.Kafka.SendMessage(r.Context(), key, payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка брокера"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (k *CartHandler) HandleDeleteBuy(w http.ResponseWriter, r *http.Request) {
	var req DeleteCartRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	WriteErrorBadReq(err, w, r)
	if err != nil {
		return
	}
	payload, _ := json.Marshal(req)
	key := []byte("2")
	err = k.Kafka.SendMessage(r.Context(), key, payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка брокера"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
