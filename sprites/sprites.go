// Package sprites provides types and helper methods for retrieving the sprite
// you need.
package sprites

// SpriteURL stores the URL the given sprite is hosted at. Most URLs are
// optional and may or may not be Present.
//
// While the api returns `null` for non-present URLs, the Go JSON decoder
// safely converts such fields to `""`.
type SpriteURL string

func (o SpriteURL) Present() bool { return o != "" }

type (
	// The PokemonDefaults are the minimum set of pokemon sprite resources. They
	// are the ones presented by the api for general use, and are likely to suit a
	// majority of needs. They are the only resources guaranteed to be present.
	PokemonDefaults struct {
		FrontDefault string `json:"front_default"`
		FrontShiny   string `json:"front_shiny"`
		BackDefault  string `json:"back_default"`
		BackShiny    string `json:"back_shiny"`

		FrontFemale      SpriteURL `json:"front_female"`
		FrontShinyFemale SpriteURL `json:"front_shiny_female"`
		BackFemale       SpriteURL `json:"back_female"`
		BackShinyFemale  SpriteURL `json:"back_shiny_female"`
	}

	// DreamWorld stores the sprites sourced from the pokémon dream world (Gen V).
	DreamWorld struct {
		FrontDefault SpriteURL `json:"front_default"`
		FrontFemale  SpriteURL `json:"front_female"`
	}
	// Home stores the sprites sourced from pokémon home.
	Home struct {
		FrontDefault     SpriteURL `json:"front_default"`
		FrontShiny       SpriteURL `json:"front_shiny"`
		FrontFemale      SpriteURL `json:"front_female"`
		FrontShinyFemale SpriteURL `json:"front_shiny_female"`
	}
	// OfficialArtwork stores the sprites sourced from officially released digital
	// artworks.
	OfficialArtwork struct {
		FrontDefault SpriteURL `json:"front_default"`
		FrontShiny   SpriteURL `json:"front_shiny"`
	}
	// Showdown stores the sprites used by pokemon showdown.
	Showdown     PokemonDefaults
	OtherSources struct {
		DreamWorld      DreamWorld      `json:"dream_world"`
		Home            Home            `json:"home"`
		OfficialArtwork OfficialArtwork `json:"official_artwork"`
		Showdown        Showdown        `json:"showdown"`
	}

	IconSprites struct {
		FrontDefault SpriteURL `json:"front_default"`
		FrontFemale  SpriteURL `json:"front_female"`
	}

	GameSpritesKey string

	// GameSprites is a map of GameSpritesKey to the URLs.
	// Acceptable usage is documented on the GameSpritesKey definition itself.
	//
	// Not that some keys may be present in the map but have value nil.
	GameSprites map[GameSpritesKey]any

	// GameVersions is a map of pokeapi.VersionGroup (specifically its
	// [pokeapi.NamedIdentifier.Name]) to GameSprites sprite sets. Exceptions
	// follow:
	//
	// - Gold and Silver have separate entries in the map. Key by game instead.
	//
	// - Sun and Moon are not present in the map. Instead, only Ultra Sun and
	// Ultra Moon's data is returned. Use that instead.
	GameVersions map[string]GameSprites

	// GameSources is a map of pokeapi.Generation (specifically its
	// [pokeapi.NamedIdentifier.Name]) to that generation's GameVersions' sprite
	// lists.
	GameSources map[string]GameVersions

	// Pokemon is a set of URLs pointing to where the sprite images representing
	// the pokeapi.Pokemon are hosted.
	Pokemon struct {
		PokemonDefaults

		Other    OtherSources `json:"other"`
		Versions GameSources  `json:"versions"`
	}

	// PokemonForm is a set of URLs pointing to where the sprite images
	// representing the pokeapi.PokemonForm are hosted.
	PokemonForm struct {
		PokemonDefaults
	}

	// Item is a set of URLs pointing to where the sprite images representing the
	// pokeapi.Item are hosted.
	Item struct {
		Default string `json:"default"`
	}
)

// Cast to string. Typically optional.
const (
	FrontDefault     GameSpritesKey = "front_default"
	FrontTransparent GameSpritesKey = "front_transparent" // Only present in generation-i & generation-ii.
	FrontGray        GameSpritesKey = "front_gray"        // Only present in generation-i.
	FrontShiny       GameSpritesKey = "front_shiny"
	FrontFemale      GameSpritesKey = "front_female"
	FrontShinyFemale GameSpritesKey = "front_shiny_female"

	BackDefault     GameSpritesKey = "back_default"
	BackTransparent GameSpritesKey = "back_transparent" // Only present in generation-ii & generation-ii.
	BackGray        GameSpritesKey = "back_gray"        // Only present in generation-i.
	BackShiny       GameSpritesKey = "back_shiny"
	BackFemale      GameSpritesKey = "back_female"
	BackShinyFemale GameSpritesKey = "back_shiny_female"
)

const (
	// Animated is only present in generation-v. Cast to a PokemonDefaults.
	Animated GameSpritesKey = "animated"
	// Icons is present from generation-vii onwards. Cast to an IconSprites.
	Icons GameSpritesKey = "icons"
)
