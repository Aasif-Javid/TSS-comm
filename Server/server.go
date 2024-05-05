package server

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func Server() {
	var addr = flag.String("addr", "10.17.80.29", "The address to listen to; default is the vpn assigned address.")
	// to set another address: go run server.go --addr=address
	var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

	flag.Parse()

	fmt.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, err := net.Listen("tcp", src)
	if err != nil {
		fmt.Printf("Failed to open port on %s: %s\n", src, err)
		os.Exit(1)
	}
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %s\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		handleMessage(scanner.Text(), conn)
	}

	fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func handleMessage(message string, conn net.Conn) {
	fmt.Println("> " + message)

	if message == "/time" {
		resp := "It is " + time.Now().String() + "\n"
		fmt.Print("< " + resp)
		conn.Write([]byte(resp))
	}
}
