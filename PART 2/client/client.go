package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, _ := reader.ReadString('\n')
		fmt.Print(msg)
	}
}

func write(conn net.Conn) {
	stdin := bufio.NewReader(os.Stdin)
	for {
		msg, _ := stdin.ReadString('\n')
		_, err := fmt.Fprint(conn, msg)
		if err != nil {
			fmt.Print("Client Error sending message: ", msg)
		}
	}
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()
	conn, _ := net.Dial("tcp", *addrPtr)
	go read(conn)
	go write(conn)
	for {
	}
}
