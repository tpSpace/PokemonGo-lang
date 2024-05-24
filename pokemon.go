package main

import "fmt"

type Pokemon struct {
	General           General
	Profile           Profile
	DamgeWhenAttacked DamgeWhenAttacked
	Evolutions        []Evolutions
	Moves             []Moves
}

type General struct {
	Index        int
	Name         string
	Type         []string
	HP           int
	Attack       int
	Defense      int
	Speed        int
	Sp_Atk       int
	Sp_Def       int
	Seed_Pokemon string
}

type Profile struct {
	Height       float32 // meter
	Weight       float32 // kilogram
	Catch_Rate   float32 // 0.0%
	Gender_Ratio float32 // 0.0%
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
}

type Evolutions struct {
	Level int
	From  string
	To    string
}

type Moves struct {
	DamgeWhenAttacked int // if the this damge is unknow than the value is -1 as default
	Move              string
	Type              string
}

func pokemon() {
	fmt.Println("I choose you, Pikachu!")
}

// Methods to get Pokemon data from the database
func GetGeneral(pokemon Pokemon) General {
	return pokemon.General
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

func GetMoves(pokemon Pokemon) []Moves {
	return pokemon.Moves
}

func GetPokemon(pokemon Pokemon) Pokemon {
	return pokemon
}

// create a new Pokemon
func NewPokemon(general General, profile Profile, damgeWhenAttacked DamgeWhenAttacked, evolutions []Evolutions, moves []Moves) Pokemon {
	return Pokemon{
		General:           general,
		Profile:           profile,
		DamgeWhenAttacked: damgeWhenAttacked,
		Evolutions:        evolutions,
		Moves:             moves,
	}
}
