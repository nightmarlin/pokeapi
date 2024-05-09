package pokeapi

//go:generate go run cmd/gettergen/gettergen.go -unnamed=ContestEffect,SuperContestEffect -- $GOFILE ContestType ContestEffect SuperContestEffect

type ContestName struct {
	Name     string                     `json:"name"`
	Color    string                     `json:"color"` // The color associated with this contest's name.
	Language NamedAPIResource[Language] `json:"language"`
}

type ContestType struct {
	NamedIdentifier

	BerryFlavour NamedAPIResource[BerryFlavor] `json:"berry_flavour"` // The BerryFlavor that correlates with this contest type.
	Names        []ContestName                 `json:"names"`
}

type ContestEffect struct {
	Identifier

	Appeal            int          `json:"appeal"`         // The base number of hearts the user of this move gets.
	Jam               int          `json:"jam"`            // The base number of hearts the user's opponent loses.
	EffectEntries     []Effect     `json:"effect_entries"` // The result of this contest effect listed in different languages.
	FlavorTextEntries []FlavorText `json:"flavor_text_entries"`
}

type SuperContestEffect struct {
	Identifier

	Appeal            int                      `json:"appeal"`
	FlavorTextEntries []FlavorText             `json:"flavor_text_entries"`
	Moves             []NamedAPIResource[Move] `json:"moves"`
}
