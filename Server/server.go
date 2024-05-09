package server

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	ecdsa "github.com/Aasif-Javid/TSS-comm/crypto/ecdsa"
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

type Session struct {
	conn   net.Conn
	logger ecdsa.Logger
	party  *ecdsa.Party
}

var sessions = make(map[string]*Session)
var sessionMutex sync.Mutex

func generateSessionID(conn net.Conn) string {
	remoteAddr := conn.RemoteAddr().String()
	hash := sha256.Sum256([]byte(remoteAddr))
	return hex.EncodeToString(hash[:])
}

func Server() {
	flag.Parse()
	logger := logger("server", "server")

	fmt.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			continue
		}

		sessionID := generateSessionID(conn)
		fmt.Println("Connected to client with session ID:", sessionID)

		sessionMutex.Lock()
		sessions[sessionID] = &Session{
			conn:   conn,
			logger: logger,
			party:  ecdsa.NewParty(1, logger),
		}
		// Handle the connection for a specific session
		go sendMessages(sessions[sessionID].conn, SendChan)
		go receiveMessages(sessions[sessionID].conn, ReceiveChan)
		// sendMsg([]byte("initiated keygen"), false, 1)
		// party := ecdsa.NewParty(1, logger)
		// party.Init([]uint16{1, 2}, 1, sendMsg)
		sessions[sessionID].party.Init([]uint16{1, 2}, 1, sendMsg)
		sessionMutex.Unlock()

		go handleClient(sessionID)
	}
}

func handleClient(sessionID string) {
	session := sessions[sessionID]
	defer session.conn.Close()

	for msg := range ReceiveChan {
		if string(msg) == "initiate keygen" {
			session.logger.Debugf("Initiating key generation for session", sessionID)
			go func() {
				share, err := session.party.KeyGen(context.Background())
				if err != nil {
					fmt.Println("KeyGen error:", err)

				}
				session.logger.Debugf("Key generation complete for session", sessionID, "Share length:", len(share))
			}()

		} else if string(msg) == "initiate sign" {
			session.logger.Debugf("Initiating signature for session", sessionID)
			session.party.LoadLocalPartySaveData()
			go func() {
				sig, err := session.party.Sign(context.Background(), []byte("test"))
				if err != nil {
					fmt.Println("Sign error:", err)
				}
				session.logger.Debugf("Signature generation complete for session", sessionID, "Signature length:", len(sig))
			}()
		} else {
			round, isBroadcast, err := session.party.ClassifyMsg(msg)
			if err != nil {
				fmt.Println("Error in classifying message:", err)
				return
			}
			session.logger.Debugf("Round: %d", round)
			session.party.OnMsg(msg, uint16(2), isBroadcast)
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
