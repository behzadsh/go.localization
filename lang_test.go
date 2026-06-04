package lang

import (
	"testing"
	"testing/fstest"
)

// fixtureFS is the in-memory translation set used by most tests. fstest.MapFS
// exercises the same fs.FS path the library uses for disk and embedded files.
func fixtureFS() fstest.MapFS {
	return fstest.MapFS{
		"en.yaml": &fstest.MapFile{
			Data: []byte(
				"user:\n" +
					"  not_found: User not found\n" +
					"  greeting: \"Hello {name}\"\n" +
					"  profile:\n" +
					"    title: \"{name}'s profile\"\n" +
					"raw: \"{a}{b}\"\n",
			),
		},
		"fr.yaml": &fstest.MapFile{
			Data: []byte(
				"user:\n" +
					"  greeting: \"Bonjour {name}\"\n",
			),
		},
		"README.md": &fstest.MapFile{Data: []byte("ignored")},
	}
}

func newFixtureTranslator(t *testing.T, opts ...Option) *Translator {
	t.Helper()
	tr, err := New(fixtureFS(), opts...)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestTranslateBy(t *testing.T) {
	tr := newFixtureTranslator(t)

	tests := []struct {
		name   string
		locale string
		key    string
		params map[string]string
		want   string
	}{
		{"simple hit", "en", "user.not_found", nil, "User not found"},
		{"deep nested key", "en", "user.profile.title", map[string]string{"name": "Sam"}, "Sam's profile"},
		{"param substitution", "en", "user.greeting", map[string]string{"name": "Sam"}, "Hello Sam"},
		{"other locale", "fr", "user.greeting", map[string]string{"name": "Léa"}, "Bonjour Léa"},
		{"fallback to default", "fr", "user.not_found", nil, "User not found"},
		{"missing key echoes key", "en", "user.unknown", nil, "user.unknown"},
		{"missing locale falls back", "de", "user.not_found", nil, "User not found"},
		{"unprovided placeholder stays literal", "en", "user.greeting", nil, "Hello {name}"},
		{"value with braces not re-expanded", "en", "raw", map[string]string{"a": "{b}", "b": "X"}, "{b}X"},
	}

	for _, tc := range tests {
		t.Run(
			tc.name, func(t *testing.T) {
				got := tr.TranslateBy(tc.locale, tc.key, paramsArg(tc.params)...)
				if got != tc.want {
					t.Errorf("TranslateBy(%q, %q, %v) = %q, want %q", tc.locale, tc.key, tc.params, got, tc.want)
				}
			},
		)
	}
}

func TestTranslateUsesDefaultLocale(t *testing.T) {
	tr := newFixtureTranslator(t, WithDefaultLocale("fr"))
	if got := tr.Translate("user.greeting", map[string]string{"name": "Léa"}); got != "Bonjour Léa" {
		t.Errorf("Translate = %q, want %q", got, "Bonjour Léa")
	}
}

func TestWithFallback(t *testing.T) {
	// Default locale "fr" has no not_found; fallback "en" does.
	tr := newFixtureTranslator(t, WithDefaultLocale("fr"), WithFallback("en"))
	if got := tr.Translate("user.not_found"); got != "User not found" {
		t.Errorf("Translate = %q, want %q", got, "User not found")
	}
}

func TestLocales(t *testing.T) {
	tr := newFixtureTranslator(t)
	got := tr.Locales()
	if len(got) != 2 {
		t.Fatalf("Locales() = %v, want 2 entries", got)
	}
	want := map[string]bool{"en": true, "fr": true}
	for _, l := range got {
		if !want[l] {
			t.Errorf("unexpected locale %q in %v", l, got)
		}
	}
}

func TestReadmeIsSkipped(t *testing.T) {
	tr := newFixtureTranslator(t)
	for _, l := range tr.Locales() {
		if l == "README" {
			t.Errorf("README.md was loaded as a locale")
		}
	}
}

func TestNewErrors(t *testing.T) {
	tests := []struct {
		name string
		fsys fstest.MapFS
	}{
		{
			name: "malformed yaml",
			fsys: fstest.MapFS{"en.yaml": &fstest.MapFile{Data: []byte("user:\n  - : :\n bad")}},
		},
		{
			name: "non-string leaf",
			fsys: fstest.MapFS{"en.yaml": &fstest.MapFile{Data: []byte("count: 42\n")}},
		},
		{
			name: "duplicate locale",
			fsys: fstest.MapFS{
				"en.yaml": &fstest.MapFile{Data: []byte("a: x\n")},
				"en.json": &fstest.MapFile{Data: []byte(`{"a":"y"}`)},
			},
		},
	}
	for _, tc := range tests {
		t.Run(
			tc.name, func(t *testing.T) {
				if _, err := New(tc.fsys); err == nil {
					t.Errorf("New() error = nil, want error")
				}
			},
		)
	}
}

func TestEmptyFS(t *testing.T) {
	tr, err := New(fstest.MapFS{})
	if err != nil {
		t.Fatalf("New(empty) error = %v, want nil", err)
	}
	if got := tr.Translate("any.key"); got != "any.key" {
		t.Errorf("Translate = %q, want echoed key", got)
	}
}

func TestNewFromDir(t *testing.T) {
	tr, err := NewFromDir("testdata/lang")
	if err != nil {
		t.Fatalf("NewFromDir: %v", err)
	}
	if got := tr.TranslateBy("en", "validation.required", map[string]string{"field": "email"}); got != "The email field is required." {
		t.Errorf("got %q", got)
	}
}

func TestJSONLocale(t *testing.T) {
	fsys := fstest.MapFS{
		"en.json": &fstest.MapFile{Data: []byte(`{"user":{"greeting":"Hi {name}"}}`)},
	}
	tr, err := New(fsys)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if got := tr.Translate("user.greeting", map[string]string{"name": "Sam"}); got != "Hi Sam" {
		t.Errorf("got %q", got)
	}
}

func TestSubdirIgnored(t *testing.T) {
	fsys := fstest.MapFS{
		"en.yaml":     &fstest.MapFile{Data: []byte("a: x\n")},
		"sub/fr.yaml": &fstest.MapFile{Data: []byte("a: y\n")}, // nested dir entry must be skipped
	}
	tr, err := New(fsys)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if got := tr.Locales(); len(got) != 1 || got[0] != "en" {
		t.Errorf("Locales() = %v, want [en]", got)
	}
}

func TestNestedNonStringLeaf(t *testing.T) {
	// Error must propagate up through a nested map.
	fsys := fstest.MapFS{
		"en.yaml": &fstest.MapFile{Data: []byte("user:\n  age: 30\n")},
	}
	if _, err := New(fsys); err == nil {
		t.Errorf("New() error = nil, want error for nested non-string leaf")
	}
}

func TestEmptyParamsMap(t *testing.T) {
	tr := newFixtureTranslator(t)
	// An explicit empty map must leave placeholders untouched.
	if got := tr.Translate("user.greeting", map[string]string{}); got != "Hello {name}" {
		t.Errorf("got %q, want %q", got, "Hello {name}")
	}
}

// paramsArg adapts an optional params map to the variadic argument.
func paramsArg(p map[string]string) []map[string]string {
	if p == nil {
		return nil
	}
	return []map[string]string{p}
}
