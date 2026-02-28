package main

import servs "Api/Servs"

func main() {
	client, conn := servs.InitGRPCClient()
	defer conn.Close()

	App := &servs.App{UserClient: client}

	servs.Createserver(App)

}
