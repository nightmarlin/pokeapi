package pokeapi

//go:generate go run cmd/gettergen/gettergen.go -- $GOFILE Item ItemAttribute ItemCategory ItemFlingEffect ItemPocket

type ItemSprites struct {
	Default string `json:"default"`
}

type ItemHolderPokemonVersionDetail struct {
	Rarity  int                       `json:"rarity"`
	Version NamedAPIResource[Version] `json:"version"`
}

type ItemHolderPokemon struct {
	Pokemon        NamedAPIResource[Pokemon]        `json:"pokemon"`
	VersionDetails []ItemHolderPokemonVersionDetail `json:"version_details"`
}

type Item struct {
	NamedIdentifier

	FlingPower  int                               `json:"fling_power"`
	FlingEffect NamedAPIResource[ItemFlingEffect] `json:"fling_effect"`

	Cost              int                               `json:"cost"`
	Attributes        []NamedAPIResource[ItemAttribute] `json:"attributes"`
	Category          NamedAPIResource[ItemCategory]    `json:"category"`
	EffectEntries     []VerboseEffect                   `json:"effect_entries"`
	FlavorTextEntries []VersionGroupFlavorText          `json:"flavor_text_entries"`
	GameIndices       []GenerationGameIndex             `json:"game_indices"`
	Names             []Name                            `json:"names"`
	Sprites           ItemSprites                       `json:"sprites"`
	HeldByPokemon     []ItemHolderPokemon               `json:"held_by_pokemon"`
	BabyTriggerFor    *APIResource[EvolutionChain]      `json:"baby_trigger_for"`
	Machines          []MachineVersionDetail            `json:"machines"`
}

type ItemAttribute struct {
	NamedIdentifier

	Names        []Name        `json:"names"`
	Descriptions []Description `json:"descriptions"`
}

type ItemCategory struct {
	NamedIdentifier

	Items  []NamedAPIResource[Item]     `json:"items"`
	Names  []Name                       `json:"names"`
	Pocket NamedAPIResource[ItemPocket] `json:"pocket"`
}

type ItemFlingEffect struct {
	NamedIdentifier

	EffectEntries []Effect                 `json:"effect_entries"`
	Items         []NamedAPIResource[Item] `json:"items"`
}

type ItemPocket struct {
	NamedIdentifier

	Categories []NamedAPIResource[ItemCategory] `json:"categories"`
	Names      []Name                           `json:"names"`
}
