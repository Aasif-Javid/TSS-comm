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
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	ecdsa "github.com/Aasif-Javid/TSS-comm/binance/ecdsa"
	"go.uber.org/zap"
)

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

var SendChan = make(chan []byte, 1000)
var ReceiveChan = make(chan []byte, 1000)

func logger(id string, testName string) ecdsa.Logger {
	logConfig := zap.NewDevelopmentConfig()
	logger, _ := logConfig.Build()
	logger = logger.With(zap.String("t", testName)).With(zap.String("id", id))
	return logger.Sugar()
}

func Server() {
	flag.Parse()

	fmt.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		var parties []uint16
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			os.Exit(1)
		} else if conn != nil {
			fmt.Println("Connected to client")
			parties = append(parties, 1, 2)
			fmt.Println("Parties:", parties)
		}
		logger := logger("server", "server")
		party := ecdsa.NewParty(1, logger)

		party.Init(parties, 1, sendMsg)

		go sendMessages(conn, SendChan)
		go receiveMessages(conn, ReceiveChan)

		go func() {
			share, err := party.KeyGen(context.Background())
			if err != nil {
				return
			}
			fmt.Println("Share:", len(share))
		}()
		for msg := range ReceiveChan {
			fmt.Println("Received:", len(msg))
			round, isBroadcast, err := party.ClassifyMsg(msg)

			if err != nil {
				fmt.Println("Error in classifying message")
				os.Exit(2)
			} else {
				fmt.Println("Round:", round)
			}

			party.OnMsg([]byte(msg), uint16(2), isBroadcast)
		}
	}
}

func sendMsg(msg []byte, isBroadcast bool, to uint16) {
	message := msg
	SendChan <- message
	fmt.Println("Message passed to network sending function")

}
func sendMessages(conn net.Conn, sendChan <-chan []byte) {
	for {
		text, ok := <-sendChan
		if len(text) == 0 {
			fmt.Println("Empty message received, stopping message sending.")
			os.Exit(1)
		}
		fmt.Println("Preparing to send:", len(text))
		if !ok {
			fmt.Println("Channel closed, stopping message sending.")
			break
		}

		// Append the delimiter to the message
		text = append(text, []byte("*****")...)

		_, err := conn.Write(text)
		if err != nil {
			fmt.Println("Error writing to stream.")
			break
		}
	}
}
func receiveMessages(conn net.Conn, receiveChan chan<- []byte) {
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer
	for {
		char, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Reached EOF on server connection.")
			} else {
				fmt.Println("Error reading from connection:", err)
			}
			break
		}

		buffer.WriteByte(char)

		if strings.HasSuffix(buffer.String(), "*****") {
			msg := strings.TrimSuffix(buffer.String(), "*****")
			receiveChan <- []byte(msg)
			buffer.Reset()
		}
	}
}
