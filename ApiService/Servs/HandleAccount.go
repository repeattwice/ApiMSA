package servs

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserAccount struct {
	UserName string `json:"user_name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

func HandleAccountCreation(w http.ResponseWriter, r *http.Request) { // володя, не доделанно
	var user UserAccount
	err := json.NewDecoder(r.Body).Decode(&user)
	WriteErrorBadReq(err, w, r)
}

func HandleAvtorization(w http.ResponseWriter, r *http.Request) { // володя, не доделанно
	var user UserAccount
	err := json.NewDecoder(r.Body).Decode(&user)
	WriteErrorBadReq(err, w, r)

}

func HandleAccoutDelet(w http.ResponseWriter, r *http.Request) { // вадим

}

func WriteErrorBadReq(err error, w http.ResponseWriter, r *http.Request) { //Функция для вывода ошибка json данных, вроде гучи
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		b := []byte(err.Error())
		w.Write(b)
		fmt.Println("Ошибка записи данных. Не верный формат присланных данных:", err)
	}
}
