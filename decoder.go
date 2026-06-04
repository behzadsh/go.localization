package lang

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
)

// Decoder unmarshals the raw bytes of a translation file into v, which is always
// a *map[string]any. Its signature matches yaml.Unmarshal and json.Unmarshal, so
// those (and most third-party unmarshalers, e.g. toml.Unmarshal) can be registered
// directly without an adapter.
type Decoder func(data []byte, v any) error

// defaultDecoders returns the decoder registry every Translator starts with.
// YAML and JSON are supported out of the box; register more via WithDecoder.
func defaultDecoders() map[string]Decoder {
	return map[string]Decoder{
		".yaml": yaml.Unmarshal,
		".yml":  yaml.Unmarshal,
		".json": json.Unmarshal,
	}
}

// normalizeExt lowercases ext and ensures it has a leading dot, so callers may
// pass either "toml" or ".TOML".
func normalizeExt(ext string) string {
	ext = strings.ToLower(ext)
	if ext == "" {
		return ext
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return ext
}
