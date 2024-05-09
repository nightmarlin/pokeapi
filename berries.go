package pokeapi

import (
	"time"
)

//go:generate go run cmd/gettergen/gettergen.go -- $GOFILE Berry BerryFirmness BerryFlavor

type BerryFlavorMap struct {
	Potency int                           `json:"potency"`
	Flavor  NamedAPIResource[BerryFlavor] `json:"flavor"`
}

// A Berry !
type Berry struct {
	NamedIdentifier

	Item    NamedAPIResource[Item] `json:"item"`
	Flavors []BerryFlavorMap       `json:"flavors"`

	GrowthTime  int `json:"growth_time"` // In hours.
	Size        int `json:"size"`        // In millimeters.
	MaxHarvest  int `json:"max_harvest"`
	SoilDryness int `json:"soil_dryness"`

	NaturalGiftPower int                    `json:"natural_gift_power"`
	NaturalGiftType  NamedAPIResource[Type] `json:"natural_gift_type"`

	Smoothness int                             `json:"smoothness"`
	Firmness   NamedAPIResource[BerryFirmness] `json:"firmness"`
}

// GrowthTimeDuration converts the Berry.GrowthTime (how long, in hours, it
// takes for a berry tree to grow 1 stage) to its corresponding time.Duration.
func (b Berry) GrowthTimeDuration() time.Duration { return time.Duration(b.GrowthTime) * time.Hour }

type BerryFirmness struct {
	NamedIdentifier

	Berries []NamedAPIResource[Berry] `json:"berries"`
	Names   []Name                    `json:"names"`
}

type FlavorBerryMap struct {
	Potency int                     `json:"potency"`
	Berry   NamedAPIResource[Berry] `json:"berry"`
}

type BerryFlavor struct {
	NamedIdentifier

	Berries     []FlavorBerryMap              `json:"berries"`
	ContestType NamedAPIResource[ContestType] `json:"contest_type"`
	Names       []Name                        `json:"names"`
}
