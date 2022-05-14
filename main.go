package main

import (
	"Cauldron/handler"
	"Cauldron/login"
	"github.com/Tnze/go-mc/net"
	"log"
)

func main() {
	handler.InitPing()
	l, err := net.ListenMC(":25566")
	if err != nil {
		panic(err)
	}

	log.Println("Listening on :25566")
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go acceptConn(conn)
	}
}

func acceptConn(conn net.Conn) {
	defer conn.Close()
	_, intention, err := login.Handshake(conn)
	if err != nil {
		panic(err)
		return
	}

	switch intention {
	default: //unknown error
		log.Printf("Get unknown handshake: %v", intention)
	case 1: // ping
		handler.HandleListPing(conn)
	case 2: // handle join
		handler.HandleLogin(conn)
	}
}
