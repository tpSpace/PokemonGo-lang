package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

var randomNum int

func main() {
	// gernate random number

	rand.New(rand.NewSource(time.Now().UnixNano()))
	// get random number
	randomNum = rand.Intn(100)

	// Create the socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		fmt.Printf("Error creating socket: %s\n", err)
		os.Exit(1)
	}

	// Enable SO_REUSEPORT to allow multiple applications to bind to the same port
	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		fmt.Printf("Error setting SO_REUSEPORT: %s\n", err)
		syscall.Close(fd)
		os.Exit(1)
	}

	// Enable SO_BROADCAST to allow broadcasting
	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_BROADCAST, 1); err != nil {
		fmt.Printf("Error setting SO_BROADCAST: %s\n", err)
		syscall.Close(fd)
		os.Exit(1)
	}

	// Bind the socket
	addr := &syscall.SockaddrInet4{Port: 3000}
	copy(addr.Addr[:], net.IPv4zero.To4())
	if err := syscall.Bind(fd, addr); err != nil {
		fmt.Printf("Error binding socket: %s\n", err)
		syscall.Close(fd)
		os.Exit(1)
	}

	// Convert the file descriptor to a net.PacketConn to use with Go's net package
	file := os.NewFile(uintptr(fd), "udp")
	conn, err := net.FilePacketConn(file)
	if err != nil {
		fmt.Printf("Error converting file descriptor to net.PacketConn: %s\n", err)
		syscall.Close(fd)
		file.Close()
		os.Exit(1)
	}
	defer conn.Close()

	// Set up the broadcast address
	bcastAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: 3000}

	// Start a goroutine to broadcast messages
			addrs, err := net.InterfaceAddrs()
			if err != nil {
				fmt.Println("Error getting machine IP address:", err)
				os.Exit(1)
			}
	go func() {
		for {
			// get the current ip address
			
			message := "pokemonGo-land:" + addrs[1].String() +":"+strconv.Itoa(randomNum)+ ":end"
			// let's broadcast the tcp conenction to the client
			
			_, err := conn.WriteTo([]byte(message), bcastAddr)
			if err != nil {
				fmt.Printf("Error broadcasting: %s\n", err)
				continue
			}
			fmt.Println("Broadcasted message:", message)
			time.Sleep(1 * time.Second)
		}
	}()

	// Read from the UDP connection
	buffer := make([]byte, 1024)
	fmt.Println("Listening on UDP port 3000")
	for {
		n, addr, err := conn.ReadFrom(buffer)
		// ip, err := net.InterfaceAddrs()
		message := string(buffer[:n])
		if strings.Contains(message[15:len(message)-4],strconv.Itoa(randomNum) ) {
			continue
		}
		if err != nil {
			fmt.Printf("Error reading from UDP: %s\n", err)
			continue
		}
		fmt.Printf("Received '%s' from %s\n", string(buffer[:n]), addr.String())
	}
}
