package main

import (
	"fmt"
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
        fmt.Println("Enter username: ")
        fmt.Scanln(&username)
        fmt.Println("Enter password: ")
        fmt.Scanln(&password)
        if login(username, password) {
            logined = true
        } else {
            logined = false
            return // exit the program
        }
    case 2:
        register()
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
            break
        } else {
            break
        }
    }
    Crawler(398,649)
}

func login(username string, password string) bool{
    if username == "admin" && password == "admin" {
        fmt.Println("Login successful")
        return true
    } else {
        fmt.Println("Login failed")
        return false
    }
}

func register() {
    fmt.Println("Register")
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