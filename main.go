package main

import (
	"chat_server_golang/config"
	"chat_server_golang/network"
	"chat_server_golang/repository"
	"chat_server_golang/service"
	"flag"
	"fmt"
)

//cli 에서 셋팅되는 환경변수 값들.
var pathFlg = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", ":1010", "port set")

//public static void main()
func init() {
}
func main() {
	flag.Parse()
	fmt.Println(*pathFlg, *port)
	
	c := config.NewConfig(*pathFlg)

	if rep, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {
		n := network.NewServer(service.NewService(rep), *port);
		n.StartServer()
		
	}
	
}