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

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run p2p.go <server|client> <port|address:port>")
		return
	}

	mode := os.Args[1]
	arg := os.Args[2]

	

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancellation is called to free resources if main exits

	fmt.Println("Hello, playground")
	// Start TCP server in a goroutine and broadcast the TCP server's address
	var wg sync.WaitGroup

	wg.Add(1)
	go broadcasting(ctx, cancel, &wg)

	

	// Wait for the broadcasting goroutine to finish
	wg.Wait()
	// Start the TCP server and connect to the peer
	if mode == "server" {
		startServer(arg)
	} else if mode == "client" {
		startClient()
	} else {
		fmt.Println("Invalid mode. Use 'server' or 'client'.")
	}
	
	// select{}
	fmt.Println("Main function exiting.")
}
// TCP server
func setReusePort(conn *net.TCPListener) error {
	file, err := conn.File()
	if err != nil {
		return err
	}
	defer file.Close()
	fd := int(file.Fd())
	err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
	if err != nil {
		return err
	}
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Connection closed:", err)
			return
		}
		fmt.Println("Received:", string(buffer[:n]))
		conn.Write([]byte("Message received"))
	}
}

func startServer(port string) {

	addr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	if err := setReusePort(listener); err != nil {
		fmt.Println("Error setting SO_REUSEPORT:", err)
		return
	}

	fmt.Println("Server started on port", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func startClient() {
	
	preData := <-data

	parts := strings.Split(preData, ":")
	if len(parts) < 2 {
		fmt.Println("Invalid data format")
		return
	}
	ipWithSubnet := parts[0]
	port_remote := parts[1]

	// Split the IP address and subnet
	ipParts := strings.Split(ipWithSubnet, "/")
	if len(ipParts) < 2 {
		fmt.Println("Invalid IP format")
		return
	}
	ip := ipParts[0] // This is the IP address without the subnet

	fmt.Println("IP: ", ip)
	fmt.Println("Port: ", port_remote)
	fmt.Println("IP address & port:", ip + ":" + port_remote)
	address := ip + ":" + port_remote
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server:", address)
	for {
		_, err := conn.Write([]byte("Hello from client"))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Println("Received:", string(buffer[:n]))

		time.Sleep(2 * time.Second)
	}
}

// UDP Broadcasting
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
