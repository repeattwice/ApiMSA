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

func GetPort() string {
	PortPtr := flag.String("port", "", "Введите номер порта")
	flag.Parse()
	port := *PortPtr
	port = strings.TrimSpace(port)
	if port == "" {
		fmt.Println("Ошибка подключения, порт пуст, введите порт снова, без флага:")
		scaner := bufio.NewScanner(os.Stdin)
		scaner.Scan()
		port = scaner.Text()
		return port
	}
	return port
}

func Createserver(a *App, k CartHandler) {
	router := mux.NewRouter()
	port := GetPort()

	router.HandleFunc("/CreateAccount", func(w http.ResponseWriter, r *http.Request) {
		HandleAccountCreation(w, r, a)
	}).Methods("POST")

	router.HandleFunc("/Avtorizacion", func(w http.ResponseWriter, r *http.Request) {
		HandleAvtorization(w, r, a)
	}).Methods("GET")

	router.Path("/DeleteAccount").Methods("DELETE").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleAccoutDelet(w, r, a)
	})

	router.Path("/ShowAllItems").Methods("GET").HandlerFunc(k.HandleShowAllItemsInCart)
	router.Path("/CreateBuy").Methods("POST").HandlerFunc(k.HandleAddToCart)
	router.Path("/ChangePrice").Methods("PATCH").HandlerFunc(k.HandleChangePrice)
	router.Path("/DeleteBuyFromKorzina").Methods("DELETE").HandlerFunc(k.HandleDeleteBuy)

	fmt.Println("ApiService запущен на порту :" + port)
	http.ListenAndServe(":"+port, router)
}

type App struct {
	UserClient user_pb.UserServiceClient
}

func InitGRPCClient() (user_pb.UserServiceClient, *grpc.ClientConn) {
	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = "localhost:5051"
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Ошибка подключения к бд сервису:", err)
	}
	return user_pb.NewUserServiceClient(conn), conn
}
