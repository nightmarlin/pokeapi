package pokeapi

//go:generate go run cmd/gettergen/gettergen.go -unnamed=EvolutionChain -- $GOFILE EvolutionChain EvolutionTrigger

type EvolutionDetail struct {
	Item                  *NamedAPIResource[Item]            `json:"item"`
	Trigger               NamedAPIResource[EvolutionTrigger] `json:"trigger"`
	Gender                *int                               `json:"gender"` // The id of the gender of the evolving Pokémon species must be in order to evolve into this Pokémon species.
	HeldItem              *NamedAPIResource[Item]            `json:"held_item"`
	KnownMove             *NamedAPIResource[Move]            `json:"known_move"`
	KnownMoveType         *NamedAPIResource[Type]            `json:"known_move_type"`
	Location              *NamedAPIResource[Location]        `json:"location"`
	MinLevel              *int                               `json:"min_level"`
	MinHappiness          *int                               `json:"min_happiness"`
	MinBeauty             *int                               `json:"min_beauty"`
	MinAffection          *int                               `json:"min_affection"`
	NeedsOverworldRain    bool                               `json:"needs_overworld_rain"`
	PartySpecies          *NamedAPIResource[PokemonSpecies]  `json:"party_species"`           // The Pokémon species that must be in the players party in order for the evolving Pokémon species to evolve into this Pokémon species.
	PartyType             *NamedAPIResource[Type]            `json:"party_type"`              // The player must have a Pokémon of this type in their party during the evolution trigger event in order for the evolving Pokémon species to evolve into this Pokémon species
	RelativePhysicalStats *int                               `json:"relative_physical_stats"` // The required relation between the Pokémon's Attack and Defense stats. 1 means Attack > Defense. 0 means Attack = Defense. -1 means Attack < Defense.
	TradeSpecies          *NamedAPIResource[PokemonSpecies]  `json:"trade_species"`
	TurnUpsideDown        bool                               `json:"turn_upside_down"`
}

type ChainLink struct {
	IsBaby           bool                             `json:"is_baby"`
	Species          NamedAPIResource[PokemonSpecies] `json:"species"`
	EvolutionDetails []EvolutionDetail                `json:"evolution_details"`
	EvolvesTo        []ChainLink                      `json:"evolves_to"`
}

type EvolutionChain struct {
	Identifier

	BabyTriggerItem *NamedAPIResource[Item] `json:"baby_trigger_item"` // The Item that a Pokémon would be holding when mating that would trigger the egg hatching a baby Pokémon rather than a basic Pokémon.
	Chain           *ChainLink              `json:"chain"`
}

type EvolutionTrigger struct {
	NamedIdentifier

	Names          []Name                             `json:"names"`
	PokemonSpecies []NamedAPIResource[PokemonSpecies] `json:"pokemon_species"`
}
