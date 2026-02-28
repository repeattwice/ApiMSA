package servs

import (
	"Api/user_pb"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserAccount struct {
	UserName string `json:"user_name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

func HandleAccountCreation(w http.ResponseWriter, r *http.Request, a *App) { // володя, готово
	var user UserAccount
	err := json.NewDecoder(r.Body).Decode(&user)
	WriteErrorBadReq(err, w, r)
	client := a.UserClient

	grpcReq := &user_pb.CreateAccountRequest{
		UserName: user.UserName,
		LastName: user.LastName,
		Email:    user.Email,
	}
	response, err := client.CreateAccount(r.Context(), grpcReq)

	if response.Succes == true {
		w.WriteHeader(http.StatusOK)
		resp := "Аккаунт успешно добавлен в бд"
		b := []byte(resp)
		w.Write(b)
	} else {
		w.WriteHeader(http.StatusConflict)
		resp := "Аккаунт не зарегестрирован"
		b := []byte(resp)
		w.Write(b)
	}

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
