package pokeapi

import (
	"context"
	"fmt"

	"github.com/nightmarlin/pokeapi/sprites"
)

type AbilityEffectChange struct {
	EffectEntries []Effect                       `json:"effect_entries"`
	VersionGroup  NamedAPIResource[VersionGroup] `json:"version_group"`
}

type AbilityFlavorText struct {
	FlavorText   string                         `json:"flavor_text"`
	Language     NamedAPIResource[Language]     `json:"language"`
	VersionGroup NamedAPIResource[VersionGroup] `json:"version_group"`
}

type AbilityPokemon struct {
	IsHidden bool                      `json:"is_hidden"`
	Slot     int                       `json:"slot"` // Pokémon have 3 ability 'slots' which hold references to possible abilities they could have. This is the slot of this ability for the referenced Pokémon.
	Pokemon  NamedAPIResource[Pokemon] `json:"pokemon"`
}

type Ability struct {
	//gettergen:plural Abilities
	NamedIdentifier

	IsMainSeries      bool                         `json:"is_main_series"` // Whether this Ability originated in the main series of the video games.
	Generation        NamedAPIResource[Generation] `json:"generation"`
	Names             []Name                       `json:"names"`
	EffectEntries     []VerboseEffect              `json:"effect_entries"`
	EffectChanges     []AbilityEffectChange        `json:"effect_changes"` // The list of previous effects this ability has had across version groups.
	FlavorTextEntries []AbilityFlavorText          `json:"flavor_text_entries"`
	Pokemon           []AbilityPokemon             `json:"pokemon"`
}

// A Characteristic indicates which stat contains a Pokémon's highest IV. A
// Pokémon's Characteristic is determined by the remainder of its highest IV
// divided by 5 (gene_modulo). Check out Bulbapedia for greater detail.
type Characteristic struct {
	Identifier

	GeneModulo     int                    `json:"gene_modulo"`     // The remainder of the highest stat/IV divided by 5.
	PossibleValues []int                  `json:"possible_values"` // The possible values of the highest stat that would result in a Pokémon receiving this Characteristic when divided by 5.
	HighestStat    NamedAPIResource[Stat] `json:"highest_stat"`    // The stat which results in this Characteristic.
	Descriptions   []Description          `json:"descriptions"`
}

type EggGroup struct {
	NamedIdentifier

	Names          []Name                             `json:"names"`
	PokemonSpecies []NamedAPIResource[PokemonSpecies] `json:"pokemon_species"`
}

type PokemonSpeciesGender struct {
	Rate           int                              `json:"rate"` // The chance of this Pokémon being female, in eighths; or -1 for genderless.
	PokemonSpecies NamedAPIResource[PokemonSpecies] `json:"pokemon_species"`
}

type Gender struct {
	NamedIdentifier

	PokemonSpeciesDetails []PokemonSpeciesGender             `json:"pokemon_species_details"` // A list of Pokémon species that can be this Gender and how likely it is that they will be.
	RequiredForEvolution  []NamedAPIResource[PokemonSpecies] `json:"required_for_evolution"`
}

type GrowthRateExperienceLevel struct {
	Level      int `json:"level"`
	Experience int `json:"experience"`
}

type GrowthRate struct {
	NamedIdentifier

	Formula        string                             `json:"formula"` // The LaTeX formula used to calculate the rate at which the Pokémon species gains level.
	Descriptions   []Description                      `json:"descriptions"`
	Levels         []GrowthRateExperienceLevel        `json:"levels"` // A list of levels and the amount of experience needed to attain them based on this growth rate.
	PokemonSpecies []NamedAPIResource[PokemonSpecies] `json:"pokemon_species"`
}

type NatureStatChange struct {
	MaxChange      int                              `json:"max_change"`
	PokeathlonStat NamedAPIResource[PokeathlonStat] `json:"pokeathlon_stat"`
}

type MoveBattleStylePreference struct {
	LowHPPreference  int                               `json:"low_hp_preference"`
	HighHPPreference int                               `json:"high_hp_preference"`
	MoveBattleStyle  NamedAPIResource[MoveBattleStyle] `json:"move_battle_style"`
}

type Nature struct {
	NamedIdentifier

	Names []Name `json:"names"`

	DecreasedStat *NamedAPIResource[Stat]        `json:"decreased_stat"`
	IncreasedStat *NamedAPIResource[Stat]        `json:"increased_stat"`
	HatesFlavour  *NamedAPIResource[BerryFlavor] `json:"hates_flavour"`
	LikesFlavour  *NamedAPIResource[BerryFlavor] `json:"likes_flavour"`

	PokeathlonStatChanges      []NatureStatChange          `json:"pokeathlon_stat_changes"`
	MoveBattleStylePreferences []MoveBattleStylePreference `json:"move_battle_style_preferences"`
}

type NaturePokeathlonStatAffect struct {
	MaxChange int                      `json:"max_change"` // The maximum amount of change to the referenced Pokéathlon stat.
	Nature    NamedAPIResource[Nature] `json:"nature"`
}

type NaturePokeathlonStatAffectSets struct {
	Increase []NaturePokeathlonStatAffect `json:"increase"`
	Decrease []NaturePokeathlonStatAffect `json:"decrease"`
}

type PokeathlonStat struct {
	NamedIdentifier

	Names            []Name                           `json:"names"`
	AffectingNatures []NaturePokeathlonStatAffectSets `json:"affecting_natures"`
}

type PokemonAbility struct {
	IsHidden bool                      `json:"is_hidden"`
	Slot     int                       `json:"slot"` // The slot this ability occupies in this Pokémon species.
	Ability  NamedAPIResource[Ability] `json:"ability"`
}

type PokemonType struct {
	Slot int                    `json:"slot"`
	Type NamedAPIResource[Type] `json:"type"`
}

type PokemonFormType struct {
	Slot int                    `json:"slot"`
	Type NamedAPIResource[Type] `json:"type"`
}

type PokemonTypePast struct {
	Generation NamedAPIResource[Generation] `json:"generation"` // The last generation in which the referenced pokémon had the listed types.
	Types      []PokemonType                `json:"types"`      // The types the referenced pokémon had up to and including the listed generation.
}

type PokemonHeldItemVersion struct {
	Version NamedAPIResource[Version] `json:"version"`
	Rarity  int                       `json:"rarity"`
}

type PokemonHeldItem struct {
	Item           NamedAPIResource[Item]   `json:"item"`
	VersionDetails []PokemonHeldItemVersion `json:"version_details"`
}

type PokemonMoveVersion struct {
	MoveLearnMethod NamedAPIResource[MoveLearnMethod] `json:"move_learn_method"`
	VersionGroup    NamedAPIResource[VersionGroup]    `json:"version_group"`
	LevelLearnedAt  int                               `json:"level_learned_at"` // The minimum level to learn the move.
}

type PokemonMove struct {
	Move                NamedAPIResource[Move] `json:"move"`
	VersionGroupDetails []PokemonMoveVersion   `json:"version_group_details"`
}

type PokemonStat struct {
	Stat     NamedAPIResource[Stat] `json:"stat"`
	Effort   int                    `json:"effort"` // The effort points (EV) the Pokémon has in the stat.
	BaseStat int                    `json:"base_stat"`
}

// PokemonCries are a set of URLs pointing to the sound files for the Pokemon's
// cry.
type PokemonCries struct {
	Latest string  `json:"latest"`
	Legacy *string `json:"legacy"`
}

type Pokemon struct {
	//gettergen:plural Pokemon
	NamedIdentifier

	BaseExperience int                              `json:"base_experience"`
	Height         int                              `json:"height"`     // The height of this Pokémon in decimeters.
	IsDefault      bool                             `json:"is_default"` // Set for exactly one Pokémon used as the default for each species.
	Order          int                              `json:"order"`      // Order for sorting. Almost national order, except families are grouped together.
	Weight         int                              `json:"weight"`     // The weight of this Pokémon in hectograms.
	Abilities      []PokemonAbility                 `json:"abilities"`
	Forms          []NamedAPIResource[PokemonForm]  `json:"forms"`
	GameIndices    []VersionGameIndex               `json:"game_indices"`
	HeldItems      []PokemonHeldItem                `json:"held_items"` // A list of items this Pokémon may be holding when encountered.
	Moves          []PokemonMove                    `json:"moves"`
	PastTypes      []PokemonTypePast                `json:"past_types"`
	Sprites        sprites.Pokemon                  `json:"sprites"`
	Cries          PokemonCries                     `json:"cries"`
	Species        NamedAPIResource[PokemonSpecies] `json:"species"`
	Stats          []PokemonStat                    `json:"stats"`
	Types          []PokemonType                    `json:"types"`

	// A URL to the PokemonLocationArea s this Pokemon can be encountered in.
	// To retrieve, use Client.GetPokemonEncounters or Pokemon.GetEncounters.
	LocationAreaEncounters string `json:"location_area_encounters"`
}

// HeightMillimeters converts Pokemon.Height (in decimeters) to the more common measurement millimeters.
func (p Pokemon) HeightMillimeters() int { return p.Height * 100 }

// WeightGrams converts Pokemon.Weight (in hectograms) to the more common measurement grams.
func (p Pokemon) WeightGrams() int { return p.Weight * 100 }

type PokemonLocationArea struct {
	LocationArea   NamedAPIResource[LocationArea] `json:"location_area"`
	VersionDetails []VersionEncounterDetail       `json:"version_details"`
}

// PokemonLocationAreaEndpoint is attached to a pokemon's url to get the set of
// locations it can be encountered in. It is the only sub-resource in the API,
// and as such does not have the ResourceName type as it cannot be listed
// directly. Use Client.GetPokemonEncounters or Pokemon.GetEncounters.
const PokemonLocationAreaEndpoint string = "encounters"

func (c *Client) GetPokemonEncounters(
	ctx context.Context,
	pokeIdent string,
) ([]PokemonLocationArea, error) {
	return do[[]PokemonLocationArea](
		ctx, c,
		fmt.Sprintf(
			"%s/%s",
			trimSlash(c.getURL(PokemonResource, pokeIdent)), PokemonLocationAreaEndpoint,
		),
		nil,
	)
}

func (p Pokemon) GetEncounters(ctx context.Context, c *Client) ([]PokemonLocationArea, error) {
	return c.GetPokemonEncounters(ctx, p.Ident())
}

type PokemonColor struct {
	NamedIdentifier

	Names          []string                           `json:"names"`
	PokemonSpecies []NamedAPIResource[PokemonSpecies] `json:"pokemon_species"`
}

type PokemonForm struct {
	NamedIdentifier

	Order        int                            `json:"order"`      // The order in which forms should be sorted within all forms. Multiple forms may have equal order, in which case they should fall back on sorting by name.
	FormOrder    int                            `json:"form_order"` // The order in which forms should be sorted within a species' forms.
	IsDefault    bool                           `json:"is_default"` // True for exactly one form used as the default for each Pokémon.
	IsBattleOnly bool                           `json:"is_battle_only"`
	IsMega       bool                           `json:"is_mega"`
	FormName     string                         `json:"form_name"`
	Types        []PokemonFormType              `json:"types"`
	Sprites      sprites.PokemonForm            `json:"sprites"`
	VersionGroup NamedAPIResource[VersionGroup] `json:"version_group"` // The version group this Pokémon form was introduced in.
	Names        []Name                         `json:"names"`         // The form specific full name of this Pokémon form, or empty if the form does not have a specific name.
	FormNames    []Name                         `json:"form_names"`    // The form specific full name of this Pokémon form, or empty if the form does not have a specific name.
}

type PokemonHabitat struct {
	NamedIdentifier

	Names          []Name                             `json:"names"`
	PokemonSpecies []NamedAPIResource[PokemonSpecies] `json:"pokemon_species"`
}

type AwesomeName struct {
	AwesomeName string                     `json:"awesome_name"`
	Language    NamedAPIResource[Language] `json:"language"`
}

type PokemonShape struct {
	NamedIdentifier

	AwesomeNames   []AwesomeName                      `json:"awesome_names"` // The "scientific" name of this Pokémon shape listed in different languages.
	Names          []Name                             `json:"names"`
	PokemonSpecies []NamedAPIResource[PokemonSpecies] `json:"pokemon_species"`
}

type Genus struct {
	Genus    string                     `json:"genus"`
	Language NamedAPIResource[Language] `json:"language"`
}

type PokemonSpeciesDexEntry struct {
	EntryNumber int                         `json:"entry_number"`
	Pokedex     []NamedAPIResource[Pokedex] `json:"pokedex"`
}

type PalParkEncounterArea struct {
	BaseScore int                           `json:"base_score"` // The base score given to the player when the referenced Pokémon is caught during a pal park run.
	Rate      int                           `json:"rate"`       // The base rate for encountering the referenced Pokémon in this pal park area.
	Area      NamedAPIResource[PalParkArea] `json:"area"`
}

type PokemonSpeciesVariety struct {
	IsDefault bool                      `json:"is_default"`
	Pokemon   NamedAPIResource[Pokemon] `json:"pokemon"`
}

type PokemonSpecies struct {
	//gettergen:plural PokemonSpecies
	NamedIdentifier

	Order                int                               `json:"order"`          // The order in which species should be sorted. Based on National Dex order, except families are grouped together and sorted by stage.
	GenderRate           int                               `json:"gender_rate"`    // The chance of this Pokémon being female, in eighths; or -1 for genderless.
	CaptureRate          uint8                             `json:"capture_rate"`   // The base capture rate; up to 255. The higher the number, the easier the catch.
	BaseHappiness        uint8                             `json:"base_happiness"` // The happiness when caught by a normal Pokéball; up to 255. The higher the number, the happier the Pokémon.
	IsBaby               bool                              `json:"is_baby"`
	IsLegendary          bool                              `json:"is_legendary"`
	IsMythical           bool                              `json:"is_mythical"`
	HatchCounter         int                               `json:"hatch_counter"`
	HasGenderDifferences bool                              `json:"has_gender_differences"`
	FormsSwitchable      bool                              `json:"forms_switchable"`
	GrowthRate           NamedAPIResource[GrowthRate]      `json:"growth_rate"`
	PokedexNumbers       []PokemonSpeciesDexEntry          `json:"pokedex_numbers"`
	EggGroups            []NamedAPIResource[EggGroup]      `json:"egg_groups"`
	Color                NamedAPIResource[PokemonColor]    `json:"color"`
	Shape                NamedAPIResource[PokemonShape]    `json:"shape"`
	EvolvesFromSpecies   *NamedAPIResource[PokemonSpecies] `json:"evolves_from_species"`
	EvolutionChain       APIResource[EvolutionChain]       `json:"evolution_chain"`
	Habitat              *NamedAPIResource[PokemonHabitat] `json:"habitat"`
	Generation           NamedAPIResource[Generation]      `json:"generation"` // The generation this Pokémon species was introduced in.
	Names                []Name                            `json:"names"`
	PalParkEncounters    []PalParkEncounterArea            `json:"pal_park_encounters"`
	FlavorTextEntries    []FlavorText                      `json:"flavor_text_entries"`
	FormDescriptions     []Description                     `json:"form_descriptions"`
	Genera               []Genus                           `json:"genera"`
	Varieties            []PokemonSpeciesVariety           `json:"varieties"`
}

type MoveStatAffect struct {
	Change int                    `json:"change"`
	Move   NamedAPIResource[Move] `json:"move"`
}

type MoveStatAffectSets struct {
	Increase []MoveStatAffect `json:"increase"`
	Decrease []MoveStatAffect `json:"decrease"`
}

type NatureStatAffectSets struct {
	Increase []NamedAPIResource[Nature] `json:"increase"`
	Decrease []NamedAPIResource[Nature] `json:"decrease"`
}

type Stat struct {
	NamedIdentifier

	GameIndex        int                                `json:"game_index"` // ID the games use for this stat.
	IsBattleOnly     bool                               `json:"is_battle_only"`
	AffectingMoves   MoveStatAffectSets                 `json:"affecting_moves"`
	AffectingNatures NatureStatAffectSets               `json:"affecting_natures"`
	Characteristics  []APIResource[Characteristic]      `json:"characteristics"` // A list of characteristics that are set on a Pokémon when its highest base stat is this stat.
	MoveDamageClass  *NamedAPIResource[MoveDamageClass] `json:"move_damage_class"`
	Names            []Name                             `json:"names"`
}

type TypePokemon struct {
	Slot    int                       `json:"slot"`
	Pokemon NamedAPIResource[Pokemon] `json:"pokemon"`
}

type TypeRelations struct {
	NoDamageTo       []NamedAPIResource[Type] `json:"no_damage_to"`
	HalfDamageTo     []NamedAPIResource[Type] `json:"half_damage_to"`
	DoubleDamageTo   []NamedAPIResource[Type] `json:"double_damage_to"`
	NoDamageFrom     []NamedAPIResource[Type] `json:"no_damage_from"`
	HalfDamageFrom   []NamedAPIResource[Type] `json:"half_damage_from"`
	DoubleDamageFrom []NamedAPIResource[Type] `json:"double_damage_from"`
}

type TypeRelationsPast struct {
	Generation      NamedAPIResource[Generation] `json:"generation"` // The last generation in which the referenced type had the listed damage relations.
	DamageRelations TypeRelations                `json:"damage_relations"`
}

type Type struct {
	NamedIdentifier

	DamageRelations     TypeRelations                      `json:"damage_relations"`
	PastDamageRelations []TypeRelationsPast                `json:"past_damage_relations"`
	GameIndices         []GenerationGameIndex              `json:"game_indices"`
	Generation          NamedAPIResource[Generation]       `json:"generation"` // The generation this type was introduced in.
	MoveDamageClass     *NamedAPIResource[MoveDamageClass] `json:"move_damage_class"`
	Names               []Name                             `json:"names"`
	Pokemon             []TypePokemon                      `json:"pokemon"`
	Moves               []NamedAPIResource[Move]           `json:"moves"`
}
