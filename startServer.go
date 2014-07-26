package main

import (
	"bitbucket.org/ridentbyte/fallout-hacking-game/game"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	port := ":2160"
	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	checkError(err)

	listener, err := net.ListenTCP(tcpAddr.Network(), tcpAddr)
	checkError(err)

	log.Println("Starting server on port", port)

	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}
		log.Println("Client connected", connection.LocalAddr())
		go game.StartGame("dict.txt", connection, 4, 7, 10)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}
}
