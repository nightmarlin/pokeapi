package pokeapi

type Machine struct {
	Identifier

	Item         NamedAPIResource[Item]         `json:"item"`
	Move         NamedAPIResource[Move]         `json:"move"`
	VersionGroup NamedAPIResource[VersionGroup] `json:"version_group"`
}

type ContestComboDetail struct {
	UseBefore []NamedAPIResource[Move] `json:"use_before"`
	UseAfter  []NamedAPIResource[Move] `json:"use_after"`
}

type ContestComboSets struct {
	Normal *ContestComboDetail `json:"normal"`
	Super  *ContestComboDetail `json:"super"`
}

type MoveFlavorText struct {
	FlavorText   string                         `json:"flavor_text"`
	Language     NamedAPIResource[Language]     `json:"language"`
	VersionGroup NamedAPIResource[VersionGroup] `json:"version_group"`
}

type MoveMetaData struct {
	Ailment       NamedAPIResource[MoveAilment]  `json:"ailment"`
	Category      NamedAPIResource[MoveCategory] `json:"category"`
	MinHits       *int                           `json:"min_hits"`  // The minimum number of times this move hits. Null if it always only hits once.
	MaxHits       *int                           `json:"max_hits"`  // The maximum number of times this move hits. Null if it always only hits once.
	MinTurns      *int                           `json:"min_turns"` // The minimum number of turns this move continues to take effect. Null if it always only lasts one turn.
	MaxTurns      *int                           `json:"max_turns"` // The maximum number of turns this move continues to take effect. Null if it always only lasts one turn.
	Drain         int                            `json:"drain"`     // HP drain (if positive) or Recoil damage (if negative), in percent of damage done.
	Healing       int                            `json:"healing"`   // The amount of hp gained by the attacking Pokémon, in percent of it's maximum HP.
	CritRate      int                            `json:"crit_rate"`
	AilmentChance int                            `json:"ailment_chance"`
	FlinchChance  int                            `json:"flinch_chance"`
	StatChance    int                            `json:"stat_chance"`
}

type MoveStatChange struct {
	Change int                    `json:"change"`
	Stat   NamedAPIResource[Stat] `json:"stat"`
}

type PastMoveStatValues struct {
	Accuracy      *int                           `json:"accuracy"`
	EffectChance  *int                           `json:"effect_chance"`
	Power         *int                           `json:"power"`
	PP            *int                           `json:"pp"`
	EffectEntries []VerboseEffect                `json:"effect_entries"`
	Type          *NamedAPIResource[Type]        `json:"type"`
	VersionGroup  NamedAPIResource[VersionGroup] `json:"version_group"`
}

type Move struct {
	NamedIdentifier

	Accuracy         *int                              `json:"accuracy"`
	EffectChance     *int                              `json:"effect_chance"`
	PP               int                               `json:"pp"`
	Priority         int                               `json:"priority"` // -8 <= Priority <= 8
	Power            *int                              `json:"power"`    // May be 0 for moves with variable power.
	DamageClass      NamedAPIResource[MoveDamageClass] `json:"damage_class"`
	Generation       NamedAPIResource[Generation]      `json:"generation"`
	LearnedByPokemon []NamedAPIResource[Pokemon]       `json:"learned_by_pokemon"`
	Machines         []MachineVersionDetail            `json:"machines"`
	Meta             *MoveMetaData                     `json:"meta"`
	PastValues       []PastMoveStatValues              `json:"past_values"`
	StatChanges      []MoveStatChange                  `json:"stat_changes"`
	Target           NamedAPIResource[MoveTarget]      `json:"target"`
	Type             NamedAPIResource[Type]            `json:"type"`

	ContestCombos      *ContestComboSets                `json:"contest_combos"`
	ContestType        *NamedAPIResource[ContestType]   `json:"contest_type"`
	ContestEffect      *APIResource[ContestEffect]      `json:"contest_effect"`
	SuperContestEffect *APIResource[SuperContestEffect] `json:"super_contest_effect"`

	EffectEntries     []VerboseEffect       `json:"effect_entries"`
	EffectChanges     []AbilityEffectChange `json:"effect_changes"`
	FlavorTextEntries []MoveFlavorText      `json:"flavor_text_entries"`
}

type MoveAilment struct {
	NamedIdentifier

	Moves []NamedAPIResource[Move] `json:"moves"`
	Names []Name                   `json:"names"`
}

type MoveBattleStyle struct {
	NamedIdentifier

	Names []Name `json:"names"`
}

type MoveCategory struct {
	NamedIdentifier

	Moves        []NamedAPIResource[Move] `json:"moves"`
	Descriptions []Description            `json:"descriptions"`
}

type MoveDamageClass struct {
	NamedIdentifier

	Descriptions []Description            `json:"descriptions"`
	Moves        []NamedAPIResource[Move] `json:"moves"`
	Names        []Name                   `json:"names"`
}

type MoveLearnMethod struct {
	NamedIdentifier

	Descriptions  []Description                    `json:"descriptions"`
	Names         []Name                           `json:"names"`
	VersionGroups []NamedAPIResource[VersionGroup] `json:"version_groups"`
}

type MoveTarget struct {
	NamedIdentifier

	Descriptions []Description            `json:"descriptions"`
	Moves        []NamedAPIResource[Move] `json:"moves"`
	Names        []Name                   `json:"names"`
}
