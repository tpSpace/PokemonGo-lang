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
type Pokemon struct {
	Name               string
	Level              int
	BaseExp            int
	HP                 int
	Attack             int
	Defense            int
	SpecialAttack      int
	SpecialDefense     int
	Speed              int
	ElementalMultiplier float64
}

type Player struct {
	Name    string
	Pokemons []Pokemon
}
var randomNum int
var data = make(chan string, 1)
// var state = 0
const (
	tcpPort = "8080"
	udpPort = "3000"
)
func main() {
	// Define two sample players with 3 Pokémon each
	player1 := Player{
		Name: "Ash",
		Pokemons: []Pokemon{
			{Name: "Pikachu", Level: 5, BaseExp: 112, HP: 35, Attack: 55, Defense: 40, SpecialAttack: 50, SpecialDefense: 50, Speed: 90, ElementalMultiplier: 1.2},
			{Name: "Bulbasaur", Level: 5, BaseExp: 64, HP: 45, Attack: 49, Defense: 49, SpecialAttack: 65, SpecialDefense: 65, Speed: 45, ElementalMultiplier: 1.0},
			{Name: "Charmander", Level: 5, BaseExp: 62, HP: 39, Attack: 52, Defense: 43, SpecialAttack: 60, SpecialDefense: 50, Speed: 65, ElementalMultiplier: 1.1},
		},
	}

	player2 := Player{
		Name: "Gary",
		Pokemons: []Pokemon{
			{Name: "Squirtle", Level: 5, BaseExp: 63, HP: 44, Attack: 48, Defense: 65, SpecialAttack: 50, SpecialDefense: 64, Speed: 43, ElementalMultiplier: 1.0},
			{Name: "Pidgey", Level: 5, BaseExp: 50, HP: 40, Attack: 45, Defense: 40, SpecialAttack: 35, SpecialDefense: 35, Speed: 56, ElementalMultiplier: 1.0},
			{Name: "Rattata", Level: 5, BaseExp: 51, HP: 30, Attack: 56, Defense: 35, SpecialAttack: 25, SpecialDefense: 35, Speed: 72, ElementalMultiplier: 1.0},
		},
	}
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run p2p.go <server|client> <port|address:port>")
		return
	}

	mode := os.Args[1]
	// arg := os.Args[2]

	

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancellation is called to free resources if main exits

	fmt.Println("Hello, playground")
	// Start TCP server in a goroutine and broadcast the TCP server's address
	var wg sync.WaitGroup

	wg.Add(1)
	go broadcasting(ctx, cancel, &wg)

	

	// Wait for the broadcasting goroutine to finish
	wg.Wait()
	// time.Sleep(2 * time.Second)
	// Start the TCP server and connect to the peer
	if mode == "server" {
		startServer(player1)
	} else if mode == "client" {
		startClient(player2)
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

func startClient(player Player) {
	
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
		fmt.Println("Error connecting to s--erver:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server:", address)
	state_client := 0
	var activePoke Pokemon
	for {
		
		if state_client == 0 {
			// send the fastest pokemon to the server
			fastest := fastestPokemon(player.Pokemons)
			activePoke = fastest
			_, err = conn.Write([]byte(strconv.Itoa(fastest.Speed)))
			if err != nil {
				fmt.Println("Error writing to server:", err)
				return
			}
			fmt.Println("Sent:", fastest.Speed)
			// Read the server's response
			buffer := make([]byte, 1024)
			_, err = conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading from server:", err)
				return
			}
			fmt.Println("Received:", string(buffer))
			state_client, _ = strconv.Atoi(string(buffer[0]))
			fmt.Println("State:", state_client)
			time.Sleep(1 * time.Second)
		} else if state_client == 1 {
			fmt.Println("State 1")
			// Read the server's response
			buffer := make([]byte, 1024)
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading from server:", err)
				return
			}
			fmt.Println("Received:", string(buffer))
			// get the attack power of the pokemon
			attack := string(buffer)
			attackParts := strings.Split(attack, ":")
			if attackParts[0] == "normalAttak" {
				fmt.Println("Normal Attack")
				fmt.Println("Attack Power:", attackParts[1])
				// minus the attack powerint from the pokemon's HP
				attackPower, err := strconv.ParseFloat(attackParts[1], 64)
				if err != nil {
					fmt.Println("Error converting attack power to float:", err)
					return // or handle the error appropriately
				}
				activePoke.HP = activePoke.HP- int(attackPower)
				fmt.Println("HP:", activePoke.HP)
			} else {
				fmt.Println("Special Attack")
				fmt.Println("Attack Power:", attackParts[1])
			}
			// send the state to the server
			for _, pokemon := range player.Pokemons {
				fmt.Println(pokemon.Name)
			}
			var choice int
			fmt.Println("Choose the pokemon to fight")
			fmt.Scan(&choice)
			// send the attack power of the pokemon to the client

			atk, isNormalAttack := randomMove(player.Pokemons[choice])
			if isNormalAttack {
				_, err := conn.Write([]byte("normalAttak:"+strconv.Itoa(int(atk))))
				if err != nil {
					fmt.Println("Error writing to connection:", err)
					return
				}
			} else {
				attack := int(atk)
				_, err := conn.Write([]byte("specialAttack:"+strconv.Itoa(attack)))
				if err != nil {
					fmt.Println("Error writing to connection:", err)
					return
				}
			} 
		} else if state_client == 2 {
			fmt.Println("State 2")
			for _, pokemon := range player.Pokemons {
				fmt.Println(pokemon.Name)
			}
			var choice int
			fmt.Println("Choose the pokemon to fight")
			fmt.Scan(&choice)
			// send the attack power of the pokemon to the client

			atk, isNormalAttack := randomMove(player.Pokemons[choice])
			if isNormalAttack {
				_, err := conn.Write([]byte("normalAttak:"+strconv.Itoa(int(atk))))
				if err != nil {
					fmt.Println("Error writing to connection:", err)
					return
				}
			} else {
				_, err := conn.Write([]byte("specialAttack:"+strconv.FormatFloat(atk, 'f', 6, 64)))
				if err != nil {
					fmt.Println("Error writing to connection:", err)
					return
				}
			} 

		}
		
	}
}

// find the fastest pokemon
func fastestPokemon(pokemons []Pokemon) Pokemon {
	fastest := pokemons[0]
	for _, pokemon := range pokemons {
		if pokemon.Speed > fastest.Speed {
			fastest = pokemon
		}
	}
	return fastest
}

func handleConnection(conn net.Conn, player Player) {
	defer conn.Close()
	state := 0
	time.Sleep(1 * time.Second)
	for {
		if state == 0 {
			// Read the incoming connection
			buffer := make([]byte, 1024)
			fmt.Println("Waiting for message")
			n, err :=conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading from connection:", err)
				return
			}
			fmt.Println("Received:", string(buffer[:n]))
			mess, _ := strconv.Atoi(string(buffer[:n]))
			// check if the fastest pokemon is faster than the received message
			fastest := fastestPokemon(player.Pokemons)
			max := fastest.Speed
			if max > mess {
				state = 1
			} else {
				state = 2
			}
			fmt.Println("Sending state:", state)
			//send the state to the client
			_, err = conn.Write([]byte(strconv.Itoa(state)))
			if err != nil {
				fmt.Println("Error writing to connection:", err)
				return
			}
			time.Sleep(1 * time.Second)
		} else if state == 1 {
			fmt.Println("State 1======")
			// player choose the pokemon to fight
			// send the pokemon to the client
			for _, pokemon := range player.Pokemons {
				fmt.Println(pokemon.Name)
			}
			var choice int
			fmt.Println("Choose the pokemon to fight")
			fmt.Scan(&choice)
			// send the attack power of the pokemon to the client

			atk, isNormalAttack := randomMove(player.Pokemons[choice])
			if isNormalAttack {
				_, err := conn.Write([]byte("normalAttak:"+strconv.Itoa(int(atk))))
				fmt.Println("Normal Attacked")
				if err != nil {
					fmt.Println("Error writing to connection:", err)
					return
				}
			} else {
				attack := int(atk)
				fmt.Println("Special Attacked")
				_, err := conn.Write([]byte("specialAttack:"+strconv.Itoa(attack)))
				if err != nil {
					fmt.Println("Error writing to connection:", err)
					return
				}
			} 
			// Read the server's response
			buffer := make([]byte, 1024)
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading from server:", err)
				return
			}
			fmt.Println("Received:", string(buffer))
			fmt.Println("Her")
			// state, _ = strconv.Atoi(string(buffer))
			fmt.Println("Her2")
			// state =1
			time.Sleep(1 * time.Second)
		} else if state == 2 {
			fmt.Println("State 2")

		}
	}
	
}

func randomMove(poke Pokemon) (float64, bool) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	if rand.Intn(2) == 0 {
		normalAttack := poke.Attack
		return float64(normalAttack), false
	} else {
		specialAttack := float64(poke.SpecialAttack)*poke.ElementalMultiplier
		return specialAttack, true
	}
}

func startServer(player Player) {
	port := "8080"
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
		go handleConnection(conn, player)
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
