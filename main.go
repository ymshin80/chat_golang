package main

import "chat_server_golang/network"

//public static void main()
func init() {
}
func main() {
	n := network.NewServer();
	n.StartServer()
}