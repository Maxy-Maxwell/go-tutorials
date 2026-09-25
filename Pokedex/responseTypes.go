package main

// Response structures for map and mapb callouts
type mapResponse struct {
	Count    int
	Next     *string
	Previous *string
	Results  []mapLocation
}

type mapLocation struct {
	Name string
	Url  string
}

// response structures for explore
type exploreResponse struct {
	Id                 int
	Name               string
	Pokemon_encounters []encounters
}

type encounters struct {
	Pokemon pokemonDetails
}

type pokemonDetails struct {
	Name string
	Url  string
}

// response structures for catch
type catchResponse struct {
	Id              int
	Name            string
	Base_experience int
	Height          int
	Weight          int
	Stats           []pokemonStats
	Types           []pokemonTypes
}

type pokemonStats struct {
	Base_stat int
	Effort    int
	Stat      pokemonStatDetails
}

type pokemonStatDetails struct {
	Name string
	Url  string
}

type pokemonTypes struct {
	Slot int
	Type pokemonType
}

type pokemonType struct {
	Name string
	Url  string
}
