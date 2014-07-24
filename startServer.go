package main

import (
	"bitbucket.org/ridentbyte/fallout-hacking-game/game"
	"fmt"
	"net"
	"os"
)

func main() {
	port := ":2160"
	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	checkError(err)

	listener, err := net.ListenTCP(tcpAddr.Network(), tcpAddr)
	checkError(err)

	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}
		go game.StartGame("dict.txt", connection)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}
}
