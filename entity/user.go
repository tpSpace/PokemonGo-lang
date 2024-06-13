package Entity

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
	Inventory []Pokemon `json:"inventory"`
}


