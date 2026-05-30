package keystore

import (
	"testing"

	"gopkg.in/yaml.v3"
)

// TestConfigYAMLTags replicates config.Unmarshal's internal round-trip
// (yaml.Marshal of the config subtree, then yaml.Unmarshal into Config) to prove
// the yaml tags map the snake_case config keys onto the struct fields — the only
// thing the move from manual extraction to config.Unmarshal could get wrong.
func TestConfigYAMLTags(t *testing.T) {
	subtree := map[string]any{
		"keys": map[string]any{
			"app": map[string]any{
				"type":        "rsa",
				"private_key": "PRIV",
				"public_key":  "PUB",
			},
		},
	}
	raw, err := yaml.Marshal(subtree)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	k, ok := cfg.Keys["app"]
	if !ok {
		t.Fatalf("key 'app' missing: %+v", cfg)
	}
	if k.Type != "rsa" || k.PrivateKey != "PRIV" || k.PublicKey != "PUB" {
		t.Fatalf("yaml tags did not map snake_case keys: %+v", k)
	}
}
