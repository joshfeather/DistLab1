package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

// Accept new connections constantly from channel
func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, _ := ln.Accept()
		conns <- conn
	}
}

// Read messages constantly from client and store in array when received
func handleClient(client net.Conn, clientid int, msgs chan Message) {
	reader := bufio.NewReader(client)
	for {
		msg, _ := reader.ReadString('\n')
		tidied := Message{clientid, msg}
		msgs <- tidied
	}
}

func main() {
	// Read in the network port we listen on with Default 8030 else from the commandline argument.
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()
	// Listener for that port
	ln, _ := net.Listen("tcp", *portPtr)

	//Create a channels for connections and messages
	conns := make(chan net.Conn)
	msgs := make(chan Message)

	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)
	id := 0

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			clients[id] = conn
			// Handles client new asynchronously
			go handleClient(conn, id, msgs)
			id++
		case msg := <-msgs:
			for clientid, conn := range clients {
				// Send message to other clients
				if clientid != msg.sender {
					_, err := fmt.Fprint(conn, fmt.Sprintf("Client%d : %s", msg.sender, msg.message))
					if err != nil {
						fmt.Print("Server Error sending message: ", msg)
					}
				}
			}
		}
	}
}
