package eflag

var (
	// default config
	defaultConfig = Config{
		TagName:      "flag",
		TagNameShort: "flag_short",
		ItemSep:      "@", // item1@item2@item3
		MapSep:       "=", // key1=value1@key2=value2
	}
)

// The configuration
type Config struct {
	// struct tag name
	TagName string
	// short tag name
	TagNameShort string
	// array element separator
	ItemSep string
	// map element separator
	MapSep string
}

// EFlagOption
type EFlagOption func(*Config)

// Specify struct tag name.
func WithTagName(tag string) EFlagOption {
	return func(c *Config) {
		c.TagName = tag
	}
}

// Specify struct short tag name.
func WithTagNameShort(tag string) EFlagOption {
	return func(c *Config) {
		c.TagNameShort = tag
	}
}

// Specify array separator
func WithItemSep(sep string) EFlagOption {
	return func(c *Config) {
		c.ItemSep = sep
	}
}

// Specify map separator
func WithMapSep(sep string) EFlagOption {
	return func(c *Config) {
		c.MapSep = sep
	}
}
