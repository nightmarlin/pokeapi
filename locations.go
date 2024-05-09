package pokeapi

//go:generate go run cmd/gettergen/gettergen.go -- $GOFILE Location LocationArea PalParkArea Region

type Location struct {
	NamedIdentifier

	Region      NamedAPIResource[Region]         `json:"region"`
	Names       []Name                           `json:"names"`
	GameIndices []GenerationGameIndex            `json:"game_indices"`
	Areas       []NamedAPIResource[LocationArea] `json:"areas"`
}

type PokemonEncounter struct {
	Pokemon        NamedAPIResource[Pokemon] `json:"pokemon"`
	VersionDetails []VersionEncounterDetail  `json:"version_details"`
}

type EncounterVersionDetails struct {
	Rate    int                       `json:"rate"`
	Version NamedAPIResource[Version] `json:"version"`
}

type EncounterMethodRate struct {
	EncounterMethod NamedAPIResource[EncounterMethod] `json:"encounter_method"`
	VersionDetails  []EncounterVersionDetails         `json:"version_details"`
}

type LocationArea struct {
	NamedIdentifier

	GameIndex            int                        `json:"game_index"` // The internal id of an API resource within game data.
	EncounterMethodRates []EncounterMethodRate      `json:"encounter_method_rates"`
	Location             NamedAPIResource[Location] `json:"location"`
	Names                []Name                     `json:"names"`
	PokemonEncounters    []PokemonEncounter         `json:"pokemon_encounters"`
}

type PalParkEncounterSpecies struct {
	BaseScore      int                              `json:"base_score"` // The base score given to the player when this Pokémon is caught during a pal park run.
	Rate           int                              `json:"rate"`       // The base rate for encountering this Pokémon in this pal park area.
	PokemonSpecies NamedAPIResource[PokemonSpecies] `json:"pokemon_species"`
}

type PalParkArea struct {
	NamedIdentifier

	Names             []Name                    `json:"names"`
	PokemonEncounters []PalParkEncounterSpecies `json:"pokemon_encounters"`
}

type Region struct {
	NamedIdentifier

	Locations      []NamedAPIResource[Location]     `json:"locations"`
	Names          []Name                           `json:"names"`
	MainGeneration NamedAPIResource[Generation]     `json:"main_generation"` // The generation this region was introduced in.
	Pokedexes      []NamedAPIResource[Pokedex]      `json:"pokedexes"`
	VersionGroups  []NamedAPIResource[VersionGroup] `json:"version_groups"`
}
