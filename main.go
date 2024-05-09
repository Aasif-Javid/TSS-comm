package main

import (
	"fmt"

	client "github.com/Aasif-Javid/TSS-comm/Client"
	server "github.com/Aasif-Javid/TSS-comm/Server"
)

func main() {
	var i int

	fmt.Println("Enter 1 for server and 2 for client")
	fmt.Scanf("%d", &i)
	if i == 1 {
		server.Server()
	} else {
		client.Client()
	}
}
