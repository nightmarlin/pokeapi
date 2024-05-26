// Code generated by github.com/nightmarlin/pokeapi/cmd/gettergen@v0; DO NOT EDIT.

package pokeapi

import "context"

const AbilityResource ResourceName[NamedAPIResource[Ability], Ability] = "ability"

func (c *Client) GetAbility(ctx context.Context, ident string) (*Ability, error) {
	return AbilityResource.Get(ctx, c, ident)
}
func (c *Client) ListAbilities(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Ability], Ability], error) {
	return AbilityResource.List(ctx, c, opts)
}

const BerryResource ResourceName[NamedAPIResource[Berry], Berry] = "berry"

func (c *Client) GetBerry(ctx context.Context, ident string) (*Berry, error) {
	return BerryResource.Get(ctx, c, ident)
}
func (c *Client) ListBerries(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Berry], Berry], error) {
	return BerryResource.List(ctx, c, opts)
}

const BerryFirmnessResource ResourceName[NamedAPIResource[BerryFirmness], BerryFirmness] = "berry-firmness"

func (c *Client) GetBerryFirmness(ctx context.Context, ident string) (*BerryFirmness, error) {
	return BerryFirmnessResource.Get(ctx, c, ident)
}
func (c *Client) ListBerryFirmnesses(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[BerryFirmness], BerryFirmness], error) {
	return BerryFirmnessResource.List(ctx, c, opts)
}

const BerryFlavorResource ResourceName[NamedAPIResource[BerryFlavor], BerryFlavor] = "berry-flavor"

func (c *Client) GetBerryFlavor(ctx context.Context, ident string) (*BerryFlavor, error) {
	return BerryFlavorResource.Get(ctx, c, ident)
}
func (c *Client) ListBerryFlavors(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[BerryFlavor], BerryFlavor], error) {
	return BerryFlavorResource.List(ctx, c, opts)
}

const CharacteristicResource ResourceName[APIResource[Characteristic], Characteristic] = "characteristic"

// GetCharacteristic only accepts the ID of the desired Characteristic.
func (c *Client) GetCharacteristic(ctx context.Context, id string) (*Characteristic, error) {
	return CharacteristicResource.Get(ctx, c, id)
}
func (c *Client) ListCharacteristics(ctx context.Context, opts *ListOpts) (*Page[APIResource[Characteristic], Characteristic], error) {
	return CharacteristicResource.List(ctx, c, opts)
}

const ContestEffectResource ResourceName[APIResource[ContestEffect], ContestEffect] = "contest-effect"

// GetContestEffect only accepts the ID of the desired ContestEffect.
func (c *Client) GetContestEffect(ctx context.Context, id string) (*ContestEffect, error) {
	return ContestEffectResource.Get(ctx, c, id)
}
func (c *Client) ListContestEffects(ctx context.Context, opts *ListOpts) (*Page[APIResource[ContestEffect], ContestEffect], error) {
	return ContestEffectResource.List(ctx, c, opts)
}

const ContestTypeResource ResourceName[NamedAPIResource[ContestType], ContestType] = "contest-type"

func (c *Client) GetContestType(ctx context.Context, ident string) (*ContestType, error) {
	return ContestTypeResource.Get(ctx, c, ident)
}
func (c *Client) ListContestTypes(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[ContestType], ContestType], error) {
	return ContestTypeResource.List(ctx, c, opts)
}

const EggGroupResource ResourceName[NamedAPIResource[EggGroup], EggGroup] = "egg-group"

func (c *Client) GetEggGroup(ctx context.Context, ident string) (*EggGroup, error) {
	return EggGroupResource.Get(ctx, c, ident)
}
func (c *Client) ListEggGroups(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[EggGroup], EggGroup], error) {
	return EggGroupResource.List(ctx, c, opts)
}

const EncounterConditionResource ResourceName[NamedAPIResource[EncounterCondition], EncounterCondition] = "encounter-condition"

func (c *Client) GetEncounterCondition(ctx context.Context, ident string) (*EncounterCondition, error) {
	return EncounterConditionResource.Get(ctx, c, ident)
}
func (c *Client) ListEncounterConditions(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[EncounterCondition], EncounterCondition], error) {
	return EncounterConditionResource.List(ctx, c, opts)
}

const EncounterConditionValueResource ResourceName[NamedAPIResource[EncounterConditionValue], EncounterConditionValue] = "encounter-condition-value"

func (c *Client) GetEncounterConditionValue(ctx context.Context, ident string) (*EncounterConditionValue, error) {
	return EncounterConditionValueResource.Get(ctx, c, ident)
}
func (c *Client) ListEncounterConditionValues(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[EncounterConditionValue], EncounterConditionValue], error) {
	return EncounterConditionValueResource.List(ctx, c, opts)
}

const EncounterMethodResource ResourceName[NamedAPIResource[EncounterMethod], EncounterMethod] = "encounter-method"

func (c *Client) GetEncounterMethod(ctx context.Context, ident string) (*EncounterMethod, error) {
	return EncounterMethodResource.Get(ctx, c, ident)
}
func (c *Client) ListEncounterMethods(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[EncounterMethod], EncounterMethod], error) {
	return EncounterMethodResource.List(ctx, c, opts)
}

const EvolutionChainResource ResourceName[APIResource[EvolutionChain], EvolutionChain] = "evolution-chain"

// GetEvolutionChain only accepts the ID of the desired EvolutionChain.
func (c *Client) GetEvolutionChain(ctx context.Context, id string) (*EvolutionChain, error) {
	return EvolutionChainResource.Get(ctx, c, id)
}
func (c *Client) ListEvolutionChains(ctx context.Context, opts *ListOpts) (*Page[APIResource[EvolutionChain], EvolutionChain], error) {
	return EvolutionChainResource.List(ctx, c, opts)
}

const EvolutionTriggerResource ResourceName[NamedAPIResource[EvolutionTrigger], EvolutionTrigger] = "evolution-trigger"

func (c *Client) GetEvolutionTrigger(ctx context.Context, ident string) (*EvolutionTrigger, error) {
	return EvolutionTriggerResource.Get(ctx, c, ident)
}
func (c *Client) ListEvolutionTriggers(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[EvolutionTrigger], EvolutionTrigger], error) {
	return EvolutionTriggerResource.List(ctx, c, opts)
}

const GenderResource ResourceName[NamedAPIResource[Gender], Gender] = "gender"

func (c *Client) GetGender(ctx context.Context, ident string) (*Gender, error) {
	return GenderResource.Get(ctx, c, ident)
}
func (c *Client) ListGenders(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Gender], Gender], error) {
	return GenderResource.List(ctx, c, opts)
}

const GenerationResource ResourceName[NamedAPIResource[Generation], Generation] = "generation"

func (c *Client) GetGeneration(ctx context.Context, ident string) (*Generation, error) {
	return GenerationResource.Get(ctx, c, ident)
}
func (c *Client) ListGenerations(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Generation], Generation], error) {
	return GenerationResource.List(ctx, c, opts)
}

const GrowthRateResource ResourceName[NamedAPIResource[GrowthRate], GrowthRate] = "growth-rate"

func (c *Client) GetGrowthRate(ctx context.Context, ident string) (*GrowthRate, error) {
	return GrowthRateResource.Get(ctx, c, ident)
}
func (c *Client) ListGrowthRates(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[GrowthRate], GrowthRate], error) {
	return GrowthRateResource.List(ctx, c, opts)
}

const ItemResource ResourceName[NamedAPIResource[Item], Item] = "item"

func (c *Client) GetItem(ctx context.Context, ident string) (*Item, error) {
	return ItemResource.Get(ctx, c, ident)
}
func (c *Client) ListItems(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Item], Item], error) {
	return ItemResource.List(ctx, c, opts)
}

const ItemAttributeResource ResourceName[NamedAPIResource[ItemAttribute], ItemAttribute] = "item-attribute"

func (c *Client) GetItemAttribute(ctx context.Context, ident string) (*ItemAttribute, error) {
	return ItemAttributeResource.Get(ctx, c, ident)
}
func (c *Client) ListItemAttributes(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[ItemAttribute], ItemAttribute], error) {
	return ItemAttributeResource.List(ctx, c, opts)
}

const ItemCategoryResource ResourceName[NamedAPIResource[ItemCategory], ItemCategory] = "item-category"

func (c *Client) GetItemCategory(ctx context.Context, ident string) (*ItemCategory, error) {
	return ItemCategoryResource.Get(ctx, c, ident)
}
func (c *Client) ListItemCategories(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[ItemCategory], ItemCategory], error) {
	return ItemCategoryResource.List(ctx, c, opts)
}

const ItemFlingEffectResource ResourceName[NamedAPIResource[ItemFlingEffect], ItemFlingEffect] = "item-fling-effect"

func (c *Client) GetItemFlingEffect(ctx context.Context, ident string) (*ItemFlingEffect, error) {
	return ItemFlingEffectResource.Get(ctx, c, ident)
}
func (c *Client) ListItemFlingEffects(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[ItemFlingEffect], ItemFlingEffect], error) {
	return ItemFlingEffectResource.List(ctx, c, opts)
}

const ItemPocketResource ResourceName[NamedAPIResource[ItemPocket], ItemPocket] = "item-pocket"

func (c *Client) GetItemPocket(ctx context.Context, ident string) (*ItemPocket, error) {
	return ItemPocketResource.Get(ctx, c, ident)
}
func (c *Client) ListItemPockets(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[ItemPocket], ItemPocket], error) {
	return ItemPocketResource.List(ctx, c, opts)
}

const LanguageResource ResourceName[NamedAPIResource[Language], Language] = "language"

func (c *Client) GetLanguage(ctx context.Context, ident string) (*Language, error) {
	return LanguageResource.Get(ctx, c, ident)
}
func (c *Client) ListLanguages(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Language], Language], error) {
	return LanguageResource.List(ctx, c, opts)
}

const LocationResource ResourceName[NamedAPIResource[Location], Location] = "location"

func (c *Client) GetLocation(ctx context.Context, ident string) (*Location, error) {
	return LocationResource.Get(ctx, c, ident)
}
func (c *Client) ListLocations(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Location], Location], error) {
	return LocationResource.List(ctx, c, opts)
}

const LocationAreaResource ResourceName[NamedAPIResource[LocationArea], LocationArea] = "location-area"

func (c *Client) GetLocationArea(ctx context.Context, ident string) (*LocationArea, error) {
	return LocationAreaResource.Get(ctx, c, ident)
}
func (c *Client) ListLocationAreas(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[LocationArea], LocationArea], error) {
	return LocationAreaResource.List(ctx, c, opts)
}

const MachineResource ResourceName[APIResource[Machine], Machine] = "machine"

// GetMachine only accepts the ID of the desired Machine.
func (c *Client) GetMachine(ctx context.Context, id string) (*Machine, error) {
	return MachineResource.Get(ctx, c, id)
}
func (c *Client) ListMachines(ctx context.Context, opts *ListOpts) (*Page[APIResource[Machine], Machine], error) {
	return MachineResource.List(ctx, c, opts)
}

const MoveResource ResourceName[NamedAPIResource[Move], Move] = "move"

func (c *Client) GetMove(ctx context.Context, ident string) (*Move, error) {
	return MoveResource.Get(ctx, c, ident)
}
func (c *Client) ListMoves(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Move], Move], error) {
	return MoveResource.List(ctx, c, opts)
}

const MoveAilmentResource ResourceName[NamedAPIResource[MoveAilment], MoveAilment] = "move-ailment"

func (c *Client) GetMoveAilment(ctx context.Context, ident string) (*MoveAilment, error) {
	return MoveAilmentResource.Get(ctx, c, ident)
}
func (c *Client) ListMoveAilments(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[MoveAilment], MoveAilment], error) {
	return MoveAilmentResource.List(ctx, c, opts)
}

const MoveBattleStyleResource ResourceName[NamedAPIResource[MoveBattleStyle], MoveBattleStyle] = "move-battle-style"

func (c *Client) GetMoveBattleStyle(ctx context.Context, ident string) (*MoveBattleStyle, error) {
	return MoveBattleStyleResource.Get(ctx, c, ident)
}
func (c *Client) ListMoveBattleStyles(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[MoveBattleStyle], MoveBattleStyle], error) {
	return MoveBattleStyleResource.List(ctx, c, opts)
}

const MoveCategoryResource ResourceName[NamedAPIResource[MoveCategory], MoveCategory] = "move-category"

func (c *Client) GetMoveCategory(ctx context.Context, ident string) (*MoveCategory, error) {
	return MoveCategoryResource.Get(ctx, c, ident)
}
func (c *Client) ListMoveCategories(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[MoveCategory], MoveCategory], error) {
	return MoveCategoryResource.List(ctx, c, opts)
}

const MoveDamageClassResource ResourceName[NamedAPIResource[MoveDamageClass], MoveDamageClass] = "move-damage-class"

func (c *Client) GetMoveDamageClass(ctx context.Context, ident string) (*MoveDamageClass, error) {
	return MoveDamageClassResource.Get(ctx, c, ident)
}
func (c *Client) ListMoveDamageClasses(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[MoveDamageClass], MoveDamageClass], error) {
	return MoveDamageClassResource.List(ctx, c, opts)
}

const MoveLearnMethodResource ResourceName[NamedAPIResource[MoveLearnMethod], MoveLearnMethod] = "move-learn-method"

func (c *Client) GetMoveLearnMethod(ctx context.Context, ident string) (*MoveLearnMethod, error) {
	return MoveLearnMethodResource.Get(ctx, c, ident)
}
func (c *Client) ListMoveLearnMethods(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[MoveLearnMethod], MoveLearnMethod], error) {
	return MoveLearnMethodResource.List(ctx, c, opts)
}

const MoveTargetResource ResourceName[NamedAPIResource[MoveTarget], MoveTarget] = "move-target"

func (c *Client) GetMoveTarget(ctx context.Context, ident string) (*MoveTarget, error) {
	return MoveTargetResource.Get(ctx, c, ident)
}
func (c *Client) ListMoveTargets(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[MoveTarget], MoveTarget], error) {
	return MoveTargetResource.List(ctx, c, opts)
}

const NatureResource ResourceName[NamedAPIResource[Nature], Nature] = "nature"

func (c *Client) GetNature(ctx context.Context, ident string) (*Nature, error) {
	return NatureResource.Get(ctx, c, ident)
}
func (c *Client) ListNatures(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Nature], Nature], error) {
	return NatureResource.List(ctx, c, opts)
}

const PalParkAreaResource ResourceName[NamedAPIResource[PalParkArea], PalParkArea] = "pal-park-area"

func (c *Client) GetPalParkArea(ctx context.Context, ident string) (*PalParkArea, error) {
	return PalParkAreaResource.Get(ctx, c, ident)
}
func (c *Client) ListPalParkAreas(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[PalParkArea], PalParkArea], error) {
	return PalParkAreaResource.List(ctx, c, opts)
}

const PokeathlonStatResource ResourceName[NamedAPIResource[PokeathlonStat], PokeathlonStat] = "pokeathlon-stat"

func (c *Client) GetPokeathlonStat(ctx context.Context, ident string) (*PokeathlonStat, error) {
	return PokeathlonStatResource.Get(ctx, c, ident)
}
func (c *Client) ListPokeathlonStats(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[PokeathlonStat], PokeathlonStat], error) {
	return PokeathlonStatResource.List(ctx, c, opts)
}

const PokedexResource ResourceName[NamedAPIResource[Pokedex], Pokedex] = "pokedex"

func (c *Client) GetPokedex(ctx context.Context, ident string) (*Pokedex, error) {
	return PokedexResource.Get(ctx, c, ident)
}
func (c *Client) ListPokedexs(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Pokedex], Pokedex], error) {
	return PokedexResource.List(ctx, c, opts)
}

const PokemonResource ResourceName[NamedAPIResource[Pokemon], Pokemon] = "pokemon"

func (c *Client) GetPokemon(ctx context.Context, ident string) (*Pokemon, error) {
	return PokemonResource.Get(ctx, c, ident)
}
func (c *Client) ListPokemons(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Pokemon], Pokemon], error) {
	return PokemonResource.List(ctx, c, opts)
}

const PokemonColorResource ResourceName[NamedAPIResource[PokemonColor], PokemonColor] = "pokemon-color"

func (c *Client) GetPokemonColor(ctx context.Context, ident string) (*PokemonColor, error) {
	return PokemonColorResource.Get(ctx, c, ident)
}
func (c *Client) ListPokemonColors(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[PokemonColor], PokemonColor], error) {
	return PokemonColorResource.List(ctx, c, opts)
}

const PokemonFormResource ResourceName[NamedAPIResource[PokemonForm], PokemonForm] = "pokemon-form"

func (c *Client) GetPokemonForm(ctx context.Context, ident string) (*PokemonForm, error) {
	return PokemonFormResource.Get(ctx, c, ident)
}
func (c *Client) ListPokemonForms(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[PokemonForm], PokemonForm], error) {
	return PokemonFormResource.List(ctx, c, opts)
}

const PokemonHabitatResource ResourceName[NamedAPIResource[PokemonHabitat], PokemonHabitat] = "pokemon-habitat"

func (c *Client) GetPokemonHabitat(ctx context.Context, ident string) (*PokemonHabitat, error) {
	return PokemonHabitatResource.Get(ctx, c, ident)
}
func (c *Client) ListPokemonHabitats(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[PokemonHabitat], PokemonHabitat], error) {
	return PokemonHabitatResource.List(ctx, c, opts)
}

const PokemonShapeResource ResourceName[NamedAPIResource[PokemonShape], PokemonShape] = "pokemon-shape"

func (c *Client) GetPokemonShape(ctx context.Context, ident string) (*PokemonShape, error) {
	return PokemonShapeResource.Get(ctx, c, ident)
}
func (c *Client) ListPokemonShapes(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[PokemonShape], PokemonShape], error) {
	return PokemonShapeResource.List(ctx, c, opts)
}

const PokemonSpeciesResource ResourceName[NamedAPIResource[PokemonSpecies], PokemonSpecies] = "pokemon-species"

func (c *Client) GetPokemonSpecies(ctx context.Context, ident string) (*PokemonSpecies, error) {
	return PokemonSpeciesResource.Get(ctx, c, ident)
}
func (c *Client) ListPokemonSpecieses(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[PokemonSpecies], PokemonSpecies], error) {
	return PokemonSpeciesResource.List(ctx, c, opts)
}

const RegionResource ResourceName[NamedAPIResource[Region], Region] = "region"

func (c *Client) GetRegion(ctx context.Context, ident string) (*Region, error) {
	return RegionResource.Get(ctx, c, ident)
}
func (c *Client) ListRegions(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Region], Region], error) {
	return RegionResource.List(ctx, c, opts)
}

const StatResource ResourceName[NamedAPIResource[Stat], Stat] = "stat"

func (c *Client) GetStat(ctx context.Context, ident string) (*Stat, error) {
	return StatResource.Get(ctx, c, ident)
}
func (c *Client) ListStats(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Stat], Stat], error) {
	return StatResource.List(ctx, c, opts)
}

const SuperContestEffectResource ResourceName[APIResource[SuperContestEffect], SuperContestEffect] = "super-contest-effect"

// GetSuperContestEffect only accepts the ID of the desired SuperContestEffect.
func (c *Client) GetSuperContestEffect(ctx context.Context, id string) (*SuperContestEffect, error) {
	return SuperContestEffectResource.Get(ctx, c, id)
}
func (c *Client) ListSuperContestEffects(ctx context.Context, opts *ListOpts) (*Page[APIResource[SuperContestEffect], SuperContestEffect], error) {
	return SuperContestEffectResource.List(ctx, c, opts)
}

const TypeResource ResourceName[NamedAPIResource[Type], Type] = "type"

func (c *Client) GetType(ctx context.Context, ident string) (*Type, error) {
	return TypeResource.Get(ctx, c, ident)
}
func (c *Client) ListTypes(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Type], Type], error) {
	return TypeResource.List(ctx, c, opts)
}

const VersionResource ResourceName[NamedAPIResource[Version], Version] = "version"

func (c *Client) GetVersion(ctx context.Context, ident string) (*Version, error) {
	return VersionResource.Get(ctx, c, ident)
}
func (c *Client) ListVersions(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[Version], Version], error) {
	return VersionResource.List(ctx, c, opts)
}

const VersionGroupResource ResourceName[NamedAPIResource[VersionGroup], VersionGroup] = "version-group"

func (c *Client) GetVersionGroup(ctx context.Context, ident string) (*VersionGroup, error) {
	return VersionGroupResource.Get(ctx, c, ident)
}
func (c *Client) ListVersionGroups(ctx context.Context, opts *ListOpts) (*Page[NamedAPIResource[VersionGroup], VersionGroup], error) {
	return VersionGroupResource.List(ctx, c, opts)
}
