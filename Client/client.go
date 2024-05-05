package client

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func Client() {
	var host = flag.String("host", "10.17.80.29", "The hostname or IP to connect to; defaults to VPN assigned address.")
	// to set another host: go run client.go --host=host's address
	var port = flag.Int("port", 8000, "The port to connect to; defaults to 8000.")

	flag.Parse()

	dest := *host + ":" + strconv.Itoa(*port)
	fmt.Printf("Connecting to %s...\n", dest)

	conn, err := net.Dial("tcp", dest)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	go readConnection(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error writing to stream:", err)
			break
		}
	}
}

func readConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Printf("\b\b** %s\n> ", scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("Error reading from server:", scanner.Err())
	}
}
