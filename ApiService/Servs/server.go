package servs

import (
	"Api/user_pb"
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
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

func Createserver(a *App) { // недоделанна
	router := mux.NewRouter()
	port := GetPort()
	router.HandleFunc("/CreateAccount", func(w http.ResponseWriter, r *http.Request) {
		HandleAccountCreation(w, r, a)
	}).Methods("POST")
	router.HandleFunc("/Avtorizacion", func(w http.ResponseWriter, r *http.Request) {
		HandleAvtorization(w, r, a)
	}).Methods("GET")
	router.Path("/DeleteAccount").Methods("DELETE").HandlerFunc(HandleAccoutDelet)

	router.Path("/ShowAllItems").Methods("GET").Queries("").HandlerFunc(HandleShowAllItemsInCort)   //надо придумать и записать query параметры
	router.Path("/CreateBuy").Methods("POST").Queries("").HandlerFunc(HandleAddToCort)              //надо придумать и записать query параметры
	router.Path("/GetDiliverySrarus").Methods("PATCH").Queries("").HandlerFunc(HandleChangePrice)   //надо придумать и записать query параметры
	router.Path("/DeleteBuyFromKorzina").Methods("DELETE").Queries("").HandlerFunc(HandleDeleteBuy) //надо придумать и записать query параметры

	http.ListenAndServe(":"+port, router)

}

type App struct {
	UserClient user_pb.UserServiceClient
}

func InitGRPCClient() (user_pb.UserServiceClient, *grpc.ClientConn) {
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println("Ошибка подключения к бд сервису")
	}
	return user_pb.NewUserServiceClient(conn), conn
}
