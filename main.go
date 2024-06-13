package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
	Pokemon "tpSpace/PokemonGo-lang/entity"
	User "tpSpace/PokemonGo-lang/entity"

	"golang.org/x/crypto/bcrypt"
)

var logined bool = false

func main() {
    // print out "Welcome to PokemonGo-land" in terminal ascii art style
    banner()
    fmt.Println("Let's catch some Pokemons!")
    fmt.Println("1. Login")
    fmt.Println("2. Register")
    fmt.Println("3. Exit")
    fmt.Println("Please enter your choice: ")
    var choice int
    fmt.Scanln(&choice)
    switch choice {
    case 1:
        var username, password string
        fmt.Print("Enter username: ")
        fmt.Scanln(&username)
        fmt.Print("Enter password: ")
        fmt.Scanln(&password)
        if login(username, password) {
            fmt.Println("Login successful-----------")
            logined = true
            break
        } else {
            logined = false
            return // exit the program
        }
    case 2:
        var username, password string
        fmt.Print("Enter username: ")
        fmt.Scanln(&username)
        fmt.Print("Enter password: ")
        fmt.Scanln(&password)
        register(username, password)
        // 
    case 3:
        fmt.Println("Goodbye!")
        return // exit the program
    default:
        fmt.Println("Invalid choice")
        // exit the program
        return
    }
    // if login successful, start the game
    for {
        if logined {
            // game()
            // clear screen
            clearScreen()
            fmt.Println("Welcome back to PokemonGo-land!")
            fmt.Println("1. Catch Pokemon")
            fmt.Println("2. Inventory")
            fmt.Println("3. Battle")
            fmt.Println("4. Exit")
            choice := 0
            fmt.Print("Please enter your choice: ")
            fmt.Scanln(&choice)
            switch choice {
            case 1:
                // CatchPokemon()
                break
            case 2:
                // Inventory()
                break
            case 3:
                // Battle()
                break
            case 4:
                fmt.Println("Goodbye!")
                return
            default:
                fmt.Println("Invalid choice")
                fmt.Println("Going back to main menu")
                time.Sleep(2 * time.Second)

            }
            // Connection()
            // break
        } else {
            break
        }
    }
    // Crawler(611,649)
}

func login(username string, password string) bool {
    // Read data from the JSON file
    file, err := os.Open("user.json")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return false
    }
    defer file.Close()

    byteValue, err := io.ReadAll(file)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return false
    }

    var users []User.User
    json.Unmarshal(byteValue, &users)

    // Check if the username and password are correct
    for _, user := range users {
        if user.Username == username {
            // Compare the hashed password with the provided password
            err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
            if err == nil {
                fmt.Println("Login successful")
                return true
            } else {
                fmt.Println("Login failed: Incorrect password")
                return false
            }
        }
    }

    fmt.Println("Login failed: Username not found")
    return false
}


func register(username, password string) {
    // Read data from the JSON file
    file, err := os.Open("user.json")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    byteValue, err := io.ReadAll(file)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    var users []User.User
    json.Unmarshal(byteValue, &users)

    // Check if the username already exists
    for _, user := range users {
        if user.Username == username {
            fmt.Println("Username already exists")
            return
        }
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        fmt.Println("Error hashing password:", err)
        return
    }

    // Add the new user
    newUser := User.User{Username: username, Password: string(hashedPassword), Inventory: []Pokemon.Pokemon{}}
    users = append(users, newUser)

    // Write the updated user list back to the JSON file
    newUserData, err := json.Marshal(users)
    if err != nil {
        fmt.Println("Error marshalling new user data:", err)
        return
    }

    err = os.WriteFile("user.json", newUserData, 0644)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }

    fmt.Println("User registered successfully")
}

func banner() {
    fmt.Println(` __        __   _                            _                                       `)
    fmt.Println(`\ \      / /__| | ___ ___  _ __ ___   ___  | |_ ___                                `)
    fmt.Println(` \ \ /\ / / _ \ |/ __/ _ \| '_ \` + "`" + ` _ \ / _ \ | __/ _ \                               `)
    fmt.Println(`  \ V  V /  __/ | (_| (_) | | | | | |  __/ | || (_) |                              `)
    fmt.Println(` __\_/\_/ \___|_|\___\___/|_| |_| |_|\___|  \__\___/         _                    _ `)
    fmt.Println(`|  _ \ ___ | | _____ _ __ ___   ___  _ __  / ___| ___       | |    __ _ _ __   __| |`)
    fmt.Println(`| |_) / _ \| |/ / _ \ '_ \` + "`" + `_ \ / _ \| ' _ \ | | _/ _ \ _____| |   / _\` + "`" + ` | '_ \ / _\` + "`" + `|`)
    fmt.Println(`|  __/ (_) |   <  __/ | | | | | (_) | | | | |_| | (_) |_____| |__| (_| | | | | (_| |`)
    fmt.Println(`|_|   \___/|_|\_\___|_| |_| |_|\___/|_| |_|\____|\___/      |_____\__,_|_| |_|\__,_|`)
}

func clearScreen() {
    cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()
}