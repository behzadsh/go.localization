package lang

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestNormalizeExt(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"toml", ".toml"},
		{".toml", ".toml"},
		{".TOML", ".toml"},
		{"YML", ".yml"},
		{"", ""},
	}
	for _, tc := range tests {
		if got := normalizeExt(tc.in); got != tc.want {
			t.Errorf("normalizeExt(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// fakeUnmarshal parses a trivial "key=value" line format to prove a custom
// decoder can be registered for an arbitrary extension.
func fakeUnmarshal(data []byte, v any) error {
	out, ok := v.(*map[string]any)
	if !ok {
		return nil
	}
	m := make(map[string]any)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		k, val, _ := strings.Cut(line, "=")
		m[strings.TrimSpace(k)] = strings.TrimSpace(val)
	}
	*out = m
	return nil
}

func TestWithDecoderCustomFormat(t *testing.T) {
	fsys := fstest.MapFS{
		"en.custom": &fstest.MapFile{Data: []byte("hello = world\n")},
	}
	tr, err := New(fsys, WithDecoder("custom", fakeUnmarshal))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if got := tr.Translate("hello"); got != "world" {
		t.Errorf("got %q, want %q", got, "world")
	}
}

func TestUnregisteredExtensionSkipped(t *testing.T) {
	fsys := fstest.MapFS{
		"en.custom": &fstest.MapFile{Data: []byte("hello = world\n")},
	}
	tr, err := New(fsys) // no decoder for .custom
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if len(tr.Locales()) != 0 {
		t.Errorf("Locales() = %v, want empty (extension unregistered)", tr.Locales())
	}
}
