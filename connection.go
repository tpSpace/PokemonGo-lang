package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

var randomNum int
var data = make(chan string, 1)
const (
	tcpPort = "8080"
	udpPort = "3000"
)
func Connection() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancellation is called to free resources if main exits

	fmt.Println("Hello, playground")
	// Start TCP server in a goroutine and broadcast the TCP server's address
	port := 8080
	var wg sync.WaitGroup

	wg.Add(1)
	go broadcasting(ctx, cancel, &wg)

	

	// Wait for the broadcasting goroutine to finish
	wg.Wait()
	StartTCPServer(port)
	fmt.Println("Main function exiting.")
}

func StartTCPServer(port int) {
	fmt.Println("TCP server listening on", port)
	// Create a TCP listener
	// get ip address from data channel
	// get the current ip address of the machine
	// netip, eer := net.InterfaceAddrs()
	preData := <-data

	parts := strings.Split(preData, ":")
	if len(parts) < 2 {
		fmt.Println("Invalid data format")
		return
	}
	ip := parts[0]
	port_remote := parts[1]

	fmt.Println("IP: ", ip)
	fmt.Println("Port: ", port_remote)
	fmt.Println("IP address & port:", preData)
	
}

func broadcasting(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	// Generate random number
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum = rand.Intn(90) + 10

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
				message := "pokemonGo-land:" + addrs[1].String() + ":"+ tcpPort +":"+ strconv.Itoa(randomNum) + ":end"
				// Let's broadcast the TCP connection to the client
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
		select {
		case <-ctx.Done():
			fmt.Println("Stopping reading due to cancellation.")
			return
		default:
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				fmt.Printf("Error reading from UDP: %s\n", err)
				continue
			}
			message := string(buffer[:n])
			
			// Check if the message contains the random number
			randomNumber, _ := strconv.Atoi(message[strings.Index(message, ":")+21 : strings.Index(message, ":end")])
			fmt.Println(randomNum)
			fmt.Println(randomNumber)
			if randomNumber == randomNum {
				continue
			}
			// Get the IP from the message and save it to data channel
			data <- message[15 : len(message)-4]
			fmt.Printf("Received '%s' from %s\n", string(buffer[:n]), addr.String())
			// send a message to other clients to stop broadcasting and listening
			for i:= 0; i < 2; i++ {
				message := "pokemonGo-land:" + addrs[1].String() + ":"+ tcpPort +":"+ strconv.Itoa(randomNum) + ":end"
				// Let's broadcast the TCP connection to the client
				_, err := conn.WriteTo([]byte(message), bcastAddr)
				if err != nil {
					fmt.Printf("Error broadcasting: %s\n", err)
					continue
				}
				fmt.Println("Broadcasted message:", message)
				time.Sleep(1 * time.Second)
			}
			cancel()
		}
	}
}
