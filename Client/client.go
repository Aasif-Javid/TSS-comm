package client

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	ecdsa "github.com/Aasif-Javid/TSS-comm/binance/ecdsa"
	"go.uber.org/zap"
)

var SendChan = make(chan []byte, 1000)
var ReceiveChan = make(chan []byte, 1000)

func logger(id string, testName string) ecdsa.Logger {
	logConfig := zap.NewDevelopmentConfig()
	logger, _ := logConfig.Build()
	logger = logger.With(zap.String("t", testName)).With(zap.String("id", id))
	return logger.Sugar()
}

func Client() {
	var host = "172.20.10.4"
	var port = 8000

	dest := host + ":" + strconv.Itoa(port)
	fmt.Printf("Connecting to %s...\n", dest)

	conn, err := net.Dial("tcp", dest)
	var parties []uint16
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	} else if conn != nil {
		fmt.Println("Connected to server")
		parties = append(parties, 1, 2)
		fmt.Println("Parties:", parties)
	}
	logger := logger("client", "client")
	party := ecdsa.NewParty(2, logger)

	party.Init(parties, 1, sendMsg)

	go sendMessages(conn, SendChan)
	go receiveMessages(conn, ReceiveChan)

	go func() {
		share, err := party.KeyGen(context.Background())
		if err != nil {
			return
		}
		fmt.Println("Share:", share)
	}()
	for msg := range ReceiveChan {
		fmt.Println("Received:", len(msg))
		round, isBroadcast, err := party.ClassifyMsg(msg)

		if err != nil {
			fmt.Println("Error in classifying message")
			os.Exit(1)
		} else {
			fmt.Println("Round:", round)
		}

		party.OnMsg([]byte(msg), uint16(1), isBroadcast)
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
