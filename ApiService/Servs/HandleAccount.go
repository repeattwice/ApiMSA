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
type ResponseJsonAccount struct {
	Response string `json:"response"`
	Succes   bool   `json:"succes"`
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
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseJsonAccount{
			Response: "Ошибка BDservice",
			Succes:   false,
		})
		return
	}
	var resp ResponseJsonAccount

	if response.Succes == true {
		w.WriteHeader(http.StatusOK)
		resp.Response = "Аккаунт успешно добавлен в бд"
		resp.Succes = true
		json.NewEncoder(w).Encode(resp)
	} else {
		w.WriteHeader(http.StatusConflict)
		resp.Response = "Аккаунт не зарегестрирован"
		resp.Succes = false
		json.NewEncoder(w).Encode(resp)
	}

}

func HandleAvtorization(w http.ResponseWriter, r *http.Request, a *App) { // володя, готово
	user := UserAccount{
		UserName: r.URL.Query().Get("user_name"),
		LastName: r.URL.Query().Get("last_name"),
		Email:    r.URL.Query().Get("email"),
	}
	client := a.UserClient

	gRPCreq := &user_pb.AvtorizationRequest{
		UserName: user.UserName,
		LastName: user.LastName,
		Email:    user.Email}

	SMT, err := client.Avtorization(r.Context(), gRPCreq)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b := []byte(err.Error())
		w.Write(b)
		return
	}

	var resp ResponseJsonAccount

	if SMT.IsLactNameIsCorrect && SMT.IsUserExists {
		resp.Response = "Авторизация прошла успешно"
		resp.Succes = true
		json.NewEncoder(w).Encode(resp)

	} else if SMT.IsUserExists == false {
		resp.Response = "Такого пользователя не существует"
		resp.Succes = false
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)

	} else if !SMT.IsLactNameIsCorrect {
		resp.Response = "Фамилия введена не правильно"
		resp.Succes = false
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func HandleAccoutDelet(w http.ResponseWriter, r *http.Request, a *App) { // вадим, готово
	var user UserAccount
	err := json.NewDecoder(r.Body).Decode(&user)
	WriteErrorBadReq(err, w, r)
	client := a.UserClient

	grpcReq := &user_pb.DeleteAccountRequest{
		UserName: user.UserName,
		LastName: user.LastName,
	}
	response, err := client.DeleteAccount(r.Context(), grpcReq)

	var resp ResponseJsonAccount

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Response = "Ошибка при удалении аккаунта"
		resp.Succes = false
		json.NewEncoder(w).Encode(resp)
		return
	}

	if response.Succes {
		w.WriteHeader(http.StatusOK)
		resp.Response = "Аккаунт успешно удалён"
		resp.Succes = true
		json.NewEncoder(w).Encode(resp)
	} else {
		w.WriteHeader(http.StatusNotFound)
		resp.Response = "Аккаунт не найден"
		resp.Succes = false
		json.NewEncoder(w).Encode(resp)
	}
}

func WriteErrorBadReq(err error, w http.ResponseWriter, r *http.Request) { //Функция для вывода ошибка json данных, вроде гучи
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		b := []byte(err.Error())
		w.Write(b)
		fmt.Println("Ошибка записи данных. Не верный формат присланных данных:", err)
	}
}
