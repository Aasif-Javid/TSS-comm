package main

import (
	"fmt"

	client "github.com/Aasif-Javid/comm-TSS/Client"
	server "github.com/Aasif-Javid/comm-TSS/Server"
	// "./server"
)

func main() {
	var i int
	// ip := "127.0.0.1"

	fmt.Println("HELLO")
	fmt.Scanf("%d", &i)
	if i == 1 {
		server.Server()
	} else {
		client.Client()
	}
}
