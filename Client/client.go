package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	// mpc "github.com/Aasif-Javid/TSS-comm/binance"
)

var SendChan = make(chan string)
var ReceiveChan = make(chan string)

func Client() {
	var host = "172.20.10.4"
	var port = 8000

	dest := host + ":" + strconv.Itoa(port)
	fmt.Printf("Connecting to %s...\n", dest)

	conn, err := net.Dial("tcp", dest)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	// logger:=mpc.Logger
	// party:=mpc.

	go sendMessages(conn, SendChan)
	go receiveMessages(conn, ReceiveChan)
	msg := "hello ths is client "
	SendChan <- msg

	for {
		select {
		case msg := <-ReceiveChan:
			fmt.Println("Received:", msg)

		}
	}
}

func sendMsg(msg []byte, isBroadcast bool, to uint16) {
	message := string(msg)
	SendChan <- message
	fmt.Println("Message passed to network sending function")

}

func sendMessages(conn net.Conn, sendChan <-chan string) {
	for {
		text, ok := <-sendChan
		fmt.Println("Preparing to send")
		if !ok {
			fmt.Println("Channel closed, stopping message sending.")
			break
		}

		_, err := conn.Write([]byte(text + "\n")) // Add a newline at the end of each message
		if err != nil {
			fmt.Println("Error writing to stream.")
			break
		}
	}
}

func receiveMessages(conn net.Conn, receiveChan chan<- string) {
	scanner := bufio.NewScanner(conn)
	for {
		ok := scanner.Scan()
		if !ok {
			fmt.Println("Reached EOF on server connection.")
			break
		}

		text := scanner.Text()

		receiveChan <- text
	}
}

func handleCommands(text string) bool {
	r, err := regexp.Compile("^%.*%$")
	if err == nil {
		if r.MatchString(text) {
			switch {
			case text == "%quit%":
				fmt.Println("\b\bServer is leaving. Hanging up.")
				os.Exit(0)
			}
			return true
		}
	}
	return false
}
