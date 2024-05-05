package main

import (
	"fmt"

	client "github.com/Aasif-Javid/TSS-comm/Client"
	server "github.com/Aasif-Javid/TSS-comm/Server"
)

func main() {
	var i int

	fmt.Println(len("61475942569042929552575879927626297992497946025064426225017111835525262582257"))
	fmt.Println(len("32818041060790185734426536913487979853126176099195031170027992934316596810628"))
	fmt.Println("Enter 1 for server and 2 for client")
	fmt.Scanf("%d", &i)
	if i == 1 {
		server.Server()
	} else {
		client.Client()
	}
}
