package eflag

var (
	// default config
	defaultConfig = Config{
		TagName: "flag",
		ItemSep: "@", // item1@item2@item3
		MapSep:  "=", // key1=value1@key2=value2
	}
)

// The configuration
type Config struct {
	// struct tag name
	TagName string
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
