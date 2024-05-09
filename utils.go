package pokeapi

import (
	"strings"
)

//go:generate go run cmd/gettergen/gettergen.go -- "$GOFILE" "Language"

type Language struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Official bool   `json:"official"`
	IS369    string `json:"is_369"`
	IS3166   string `json:"is_3166"`
	Names    []Name `json:"names"`
}

// A Description is a contextual description of the resource in the Language referenced.
type Description struct {
	Description string                     `json:"description"`
	Language    NamedAPIResource[Language] `json:"language"`
}

// An Effect is a localized text effect of the resource in the Language referenced.
type Effect struct {
	Description string                     `json:"description"`
	Language    NamedAPIResource[Language] `json:"language"`
}

type FlavorText struct {
	FlavorText string                     `json:"flavor_text"` // The localized flavor text for an API resource in the Language referenced.
	Language   NamedAPIResource[Language] `json:"language"`
	Version    NamedAPIResource[Version]  `json:"version"` // The game version this flavor text is extracted from.
}

var flavorTextReplacer = strings.NewReplacer(
	"\f", "\n",
	"\u00ad\n", "",
	"\u00ad", "",
	" -\n", " - ",
	"-\n", "-",
	"\n", " ",
)

// NormalizedFlavorText implements the recommendation at
// https://github.com/veekun/pokedex/issues/218 to correctly render
// FlavorText.FlavorText as the expected string.
func (ft FlavorText) NormalizedFlavorText() string {
	return flavorTextReplacer.Replace(ft.FlavorText)
}

type GenerationGameIndex struct {
	GameIndex  int                          `json:"game_index"` // An internal ID of a resource within game data.
	Generation NamedAPIResource[Generation] `json:"generation"`
}

type MachineVersionDetail struct {
	Machine      APIResource[Machine]           `json:"machine"`
	VersionGroup NamedAPIResource[VersionGroup] `json:"version_group"`
}

// A Name is a localized representation of the resource's name in the Language referenced.
type Name struct {
	Name     string                     `json:"name"`
	Language NamedAPIResource[Language] `json:"language"`
}

type VerboseEffect struct {
	Effect      string                     `json:"effect"`
	ShortEffect string                     `json:"short_effect"`
	Language    NamedAPIResource[Language] `json:"language"`
}

type VersionGameIndex struct {
	GameIndex int                       `json:"game_index"`
	Version   NamedAPIResource[Version] `json:"version"`
}

type VersionGroupFlavorText struct {
	Text         string                         `json:"text"` // The localized name of the API resource in the referenced Language.
	Language     NamedAPIResource[Language]     `json:"language"`
	VersionGroup NamedAPIResource[VersionGroup] `json:"version_group"` // The version group which uses this flavor text.
}
