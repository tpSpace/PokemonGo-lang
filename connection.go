package main

import (
	"fmt"
	"net"
	"time"
)

func Connection() {
	fmt.Println("Connection")
}

func QuickConnection(userId string) {
	broadcastAddr := "255.255.255.255:9999"

	con, err := net.Dial("udp", broadcastAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer con.Close()
	for {
		message := []byte("findplayer:"+userId)

		_, err = con.Write(message)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		time.Sleep(2 * time.Second)
	}
}

