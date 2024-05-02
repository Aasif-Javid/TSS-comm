/*
A very simple TCP server written in Go.

This is a toy project that I used to learn the fundamentals of writing
Go code and doing some really basic network stuff.

Maybe it will be fun for you to read. It's not meant to be
particularly idiomatic, or well-written for that matter.
*/
package server

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"
)

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

var SendChan = make(chan string)
var ReceiveChan = make(chan string)

func Server() {
	flag.Parse()

	fmt.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go sendMessages(conn, SendChan)
		go receiveMessages(conn, ReceiveChan)
		msg := "hello ths is server "
		SendChan <- msg
		for {
			select {
			case msg := <-ReceiveChan:
				fmt.Println("Received:", msg)
			case <-time.After(time.Second * 10):
				fmt.Println("No activity for 10 seconds, exiting.")
				return
			}
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
