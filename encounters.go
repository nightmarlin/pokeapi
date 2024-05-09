package pokeapi

//go:generate go run cmd/gettergen/gettergen.go -- $GOFILE EncounterMethod EncounterCondition EncounterConditionValue

type Encounter struct {
	MinLevel        int                                         `json:"min_level"`
	MaxLevel        int                                         `json:"max_level"`
	ConditionValues []NamedAPIResource[EncounterConditionValue] `json:"condition_values"` // A list of EncounterConditionValue s that must be in effect for this encounter to occur.
	Chance          int                                         `json:"chance"`           // % chance for this encounter to occur.
	Method          NamedAPIResource[EncounterMethod]           `json:"method"`
}

type EncounterMethod struct {
	NamedIdentifier

	Order int    `json:"order"` // A good value for sorting.
	Names []Name `json:"names"`
}

type EncounterCondition struct {
	NamedIdentifier

	Names  []Name                                      `json:"names"`
	Values []NamedAPIResource[EncounterConditionValue] `json:"values"`
}

type EncounterConditionValue struct {
	NamedIdentifier

	Condition NamedAPIResource[EncounterCondition] `json:"condition"`
	Names     []Name                               `json:"names"`
}

type VersionEncounterDetail struct {
	Version          NamedAPIResource[Version] `json:"version"`
	MaxChance        int                       `json:"max_chance"` // Total % chance of all encounter potentials.
	EncounterDetails []Encounter               `json:"encounter_details"`
}
