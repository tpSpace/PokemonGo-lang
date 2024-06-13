package main

import (
	"context"
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
var data = make(chan string, 1)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // Ensure cancellation is called to free resources if main exits

	fmt.Println("Hello, playground")
	// Start TCP server in a goroutine and broadcast the TCP server's address
	port := 8080
	go StartTCPServer(port)
	broadcasting(ctx, cancel)
}

func StartTCPServer(port int) {
	// fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	// if err != nil {
	// 	fmt.Printf("Error creating TCP socket: %s\n", err)
	// 	os.Exit(1)
	// }

	// if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
	// 	fmt.Printf("Error setting SO_REUSEPORT on TCP socket: %s\n", err)
	// 	syscall.Close(fd)
	// 	os.Exit(1)
	// }

	// sockAddr := &syscall.SockaddrInet4{Port: port}
	// copy(sockAddr.Addr[:], net.ParseIP(""))
	// if err := syscall.Bind(fd, sockAddr); err != nil {
	// 	fmt.Printf("Error binding TCP socket: %s\n", err)
	// 	syscall.Close(fd)
	// 	os.Exit(1)
	// }
	fmt.Print("TCP server listening on port ", port)

}

func broadcasting(ctx context.Context, cancel context.CancelFunc) {
	// gernate random number

	rand.New(rand.NewSource(time.Now().UnixNano()))
	// get random number between 10 - 99

	randomNum = rand.Intn(90) +10 

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
			select {
				case <-ctx.Done(): // Checks if the context has been cancelled
					fmt.Println("Stopping broadcasting due to cancellation.")
					return
				default:
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
		}
	}()

	// Read from the UDP connection
	buffer := make([]byte, 1024)
	fmt.Println("Listening on UDP port 3000")
	for {
		n, addr, err := conn.ReadFrom(buffer)
		// ip, err := net.InterfaceAddrs()
		message := string(buffer[:n])
		// get the random number from the message it should be before ":end" and after the second ":"
		randomNumber, _ := strconv.Atoi(message[strings.Index(message, ":")+16:strings.Index(message, ":end")])
		fmt.Println(randomNumber)
		fmt.Println(randomNum)
		if  randomNumber == randomNum {
			continue
		}
		if err != nil {
			fmt.Printf("Error reading from UDP: %s\n", err)
			continue
		}
		// get the ip from the message and save it to data channel
		data <- message[15:len(message)-4]
		cancel()
		fmt.Printf("Received '%s' from %s\n", string(buffer[:n]), addr.String())
	}
}
