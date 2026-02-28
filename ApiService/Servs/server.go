package servs

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func GetPort() string { // валидировать порт после ввода и улчшить валидацию (сделать проверку не только на пустоту)
	PortPtr := flag.String("port", "", "Введите номер порта")
	flag.Parse()
	port := *PortPtr
	port = strings.TrimSpace(port)
	if port == "" {
		fmt.Println("Ошибка подключения, порт пуст, введите прот снова, без флага:")
		scaner := bufio.NewScanner(os.Stdin)
		scaner.Scan()
		port = scaner.Text()
		return port
	}
	return port
}

func Createserver() { // недоделанна
	router := mux.NewRouter()
	port := GetPort()
	router.Path("/CreateAccount").Methods("POST").HandlerFunc(HandleAccountCreation)
	router.Path("/Avtorizacion").Methods("GET").HandlerFunc(HandleAvtorization)
	router.Path("/DeleteAccount").Methods("DELETE").HandlerFunc(HandleAccoutDelet)

	router.Path("/ShowAllItems").Methods("GET").Queries("").HandlerFunc(HandleShowAllItemsInCort)   //надо придумать и записать query параметры
	router.Path("/CreateBuy").Methods("POST").Queries("").HandlerFunc(HandleBuyCreation)            //надо придумать и записать query параметры
	router.Path("/GetDiliverySrarus").Methods("PATCH").Queries("").HandlerFunc(HandleChangePrice)   //надо придумать и записать query параметры
	router.Path("/DeleteBuyFromKorzina").Methods("DELETE").Queries("").HandlerFunc(HandleDeleteBuy) //надо придумать и записать query параметры

	http.ListenAndServe(":"+port, router)

}
