package Pokemon

type Pokemon struct {
	General           General `json:"General"`
	Profile           Profile	`json:"Profile"`
	DamgeWhenAttacked DamgeWhenAttacked	`json:"DamgeWhenAttacked"`
	Evolutions        []Evolutions `json:"Evolutions"`
	Natural_Moves	 []Natural_Moves `json:"Natural_Moves"`
	Machine_Move	 []Machine_Move `json:"Machine_Move"`
	Tutor_Move		 []Tutor_Move `json:"Tutor_Move"`
}

type General struct {
	Index        int
	Name         string
	Types         []string
	HP           int
	Attack       int
	Defense      int
	Speed        int
	Sp_Atk       int
	Sp_Def       int
	Monster_Species string
	Monster_Description   string
}

type Profile struct {
	Height       float32 // meter
	Weight       float32 // kilogram
	Catch_Rate   float32 // 0.0%
	Gender_Ratio []float32 // 0.0%
	Egg_Groups   []string
	Hatch_Steps  int
	Abilities    []string
	EVs          []string
}

type DamgeWhenAttacked struct {
	Ice      float32 
	Electric float32
	Fire     float32
	Water    float32
	Flying   float32
	Fairy    float32
	Psychic  float32
	Fighting float32
	Ground   float32
	Grass    float32
	Rock	 float32
	Steel   float32
	Poison  float32
	Ghost   float32
	Dark   float32
	Dragon float32
	Bug   float32
	Normal float32
}

type Evolutions struct {
	Level int
	From  string
	To    string
}

type Natural_Moves struct {
	Index int // if the this damge is unknow than the value is -1 as default
	Move              string
	Type              Type
}

type Machine_Move struct {
	Move string
	Type Type
}

type Tutor_Move struct {
	Move string
	Type Type
}

type Type struct {
	Power int // -1 as N/A
	Acc int	// -1 as N/A
	PP int	// -1 as N/A
	Description string
}
// Methods to get Pokemon data from the database
// Create constructor functions to create a new Pokemon
func GetPokemon() Pokemon {
	return Pokemon{}
}

func GetProfile(pokemon Pokemon) Profile {
	return pokemon.Profile
}

func GetDamgeWhenAttacked(pokemon Pokemon) DamgeWhenAttacked {
	return pokemon.DamgeWhenAttacked
}

func GetEvolutions(pokemon Pokemon) []Evolutions {
	return pokemon.Evolutions
}

