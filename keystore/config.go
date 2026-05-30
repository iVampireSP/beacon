package keystore

// Config holds the key definitions loaded from configuration.
type Config struct {
	Keys map[string]KeyConfig `yaml:"keys"`
}

// KeyConfig defines a single named key.
type KeyConfig struct {
	Type       string `yaml:"type"`        // "rsa" or "ecdsa"
	PrivateKey string `yaml:"private_key"` // PEM or base64-encoded
	PublicKey  string `yaml:"public_key"`  // PEM or base64-encoded
}
