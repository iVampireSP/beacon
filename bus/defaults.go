package bus

import "github.com/iVampireSP/foundation/config"

// NewDefaultConfig returns a bus Config populated from the application config.
func NewDefaultConfig() Config {
	var cfg Config
	if err := config.Unmarshal("bus", &cfg); err != nil {
		panic(err)
	}
	return cfg
}
