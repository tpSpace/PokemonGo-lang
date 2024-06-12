package sample

import (
	"fmt"
	"net"
	"sync"
	"time"
)
const (
	tcpPort = "8080"
	udpPort = "8081"
)
type Game struct {
	players []net.Conn
	currentTurn int
	mutex sync.Mutex
}
var result = make(chan string)

func main() {
	game := Game{}

	go listenUdp()
	go broadCasting()
	


}

func broadCasting(connectionInfo string) {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP: net.IPv4bcast,
		Port: 8081,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		_, err = conn.Write([]byte(connectionInfo))
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}
}

func listenUdp() {
	addr, err := net.ResolveUDPAddr("udp", ":8081")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP message: ", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Println("Received from ", remoteAddr, " message: ", message)
	}
}