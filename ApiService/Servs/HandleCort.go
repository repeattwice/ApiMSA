package servs

import (
	"Api/Servs/broker"
	"encoding/json"
	"net/http"
)

type KafkaHandler struct {
	Kafka *broker.Producer
}

type UserCort struct {
	UserID    int
	ItemName  string
	ItemPtice int
}

func HandleShowAllItemsInCort(w http.ResponseWriter, r *http.Request) { //санек

}

func (k *KafkaHandler) HandleAddToCort(w http.ResponseWriter, r *http.Request) { //Володя
	var cort UserCort
	err := json.NewDecoder(r.Body).Decode(&cort)
	WriteErrorBadReq(err, w, r)
	payload, _ := json.Marshal(cort)
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
