// Package lang is a small, idiomatic internationalization (i18n) library.
//
// Translations live in one file per locale (en.yaml, fr.json, ...), each holding an arbitrarily nested tree of keys.
// Files are read from any io/fs.FS, so the same code works against the disk (os.DirFS), embedded files (embed.FS) or an
// in-memory fixture (fstest.MapFS).
//
//	//go:embed lang/*.yaml
//	var langFS embed.FS
//
//	sub, _ := fs.Sub(langFS, "lang")
//	t, _ := lang.New(sub)
//	t.Translate("user.greeting", map[string]string{"name": "Sam"})
//
// Keys are looked up by dot path (user.greeting maps to the greeting field under user). Placeholders are written {name}
// and substituted from the params map.
package lang

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

// Translator holds the loaded translations. It is immutable after New and safe for concurrent use by multiple
// goroutines.
type Translator struct {
	store          map[string]map[string]string
	defaultLocale  string
	fallbackLocale string
}

// config holds construction-time settings mutated by Options.
type config struct {
	defaultLocale  string
	fallbackLocale string
	decoders       map[string]Decoder
}

// Option configures a Translator at construction time.
type Option func(*config)

// WithDefaultLocale sets the locale used by Translate. Defaults to "en".
func WithDefaultLocale(locale string) Option {
	return func(c *config) { c.defaultLocale = locale }
}

// WithFallback sets the locale consulted when a key is missing in the requested locale. Defaults to the default locale.
func WithFallback(locale string) Option {
	return func(c *config) { c.fallbackLocale = locale }
}

// WithDecoder registers (or overrides) the decoder used for files with the given extension,
// e.g. WithDecoder(".toml", toml.Unmarshal). The leading dot and case are normalized, so "toml" and ".TOML" are
// equivalent.
func WithDecoder(ext string, dec Decoder) Option {
	return func(c *config) { c.decoders[normalizeExt(ext)] = dec }
}

// New loads every recognized locale file from fsys and returns a Translator. Loading is eager: a malformed file
// surfaces here rather than at lookup time. Files whose extension has no registered decoder are skipped.
func New(fsys fs.FS, opts ...Option) (*Translator, error) {
	cfg := &config{
		defaultLocale: "en",
		decoders:      defaultDecoders(),
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.fallbackLocale == "" {
		cfg.fallbackLocale = cfg.defaultLocale
	}

	store, err := loadStore(fsys, cfg.decoders)
	if err != nil {
		return nil, err
	}

	return &Translator{
		store:          store,
		defaultLocale:  cfg.defaultLocale,
		fallbackLocale: cfg.fallbackLocale,
	}, nil
}

// NewFromDir is a convenience wrapper around New that reads from the directory at dir on the local filesystem.
func NewFromDir(dir string, opts ...Option) (*Translator, error) {
	return New(os.DirFS(dir), opts...)
}

// loadStore reads and flattens every decodable locale file in the root of fsys.
func loadStore(fsys fs.FS, decoders map[string]Decoder) (map[string]map[string]string, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("lang: read translations dir: %w", err)
	}

	store := make(map[string]map[string]string)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		ext := path.Ext(name)
		dec, ok := decoders[ext]
		if !ok {
			continue
		}

		locale := strings.TrimSuffix(name, path.Ext(name))
		if _, exists := store[locale]; exists {
			return nil, fmt.Errorf("lang: duplicate locale %q (multiple files named %q.*)", locale, locale)
		}

		var data []byte
		data, err = fs.ReadFile(fsys, name)
		if err != nil {
			return nil, fmt.Errorf("lang: read %q: %w", name, err)
		}

		var tree map[string]any
		if err = dec(data, &tree); err != nil {
			return nil, fmt.Errorf("lang: decode %q: %w", name, err)
		}

		flat := make(map[string]string)
		if err = flatten("", tree, flat); err != nil {
			return nil, fmt.Errorf("lang: %q: %w", name, err)
		}
		store[locale] = flat
	}

	return store, nil
}

// flatten walks a nested translation tree, writing each string leaf into dst keyed by its dot path. A non-string,
// non-map leaf is an error.
func flatten(prefix string, tree map[string]any, dst map[string]string) error {
	for k, v := range tree {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case string:
			dst[key] = val
		case map[string]any:
			if err := flatten(key, val, dst); err != nil {
				return err
			}
		default:
			return fmt.Errorf("key %q: want string or nested map, got %T", key, v)
		}
	}

	return nil
}

// Translate looks up key in the default locale, falling back as configured, and substitutes any params. A missing key
// returns the key itself.
func (t *Translator) Translate(key string, params ...map[string]string) string {
	return t.TranslateBy(t.defaultLocale, key, params...)
}

// TranslateBy looks up key in the given locale, then the fallback locale, and substitutes any params. A missing key
// returns the key itself.
func (t *Translator) TranslateBy(locale, key string, params ...map[string]string) string {
	tr := t.lookup(locale, key)

	if len(params) > 0 {
		tr = substitute(tr, params[0])
	}

	return tr
}

// Locales returns the loaded locales in no particular order.
func (t *Translator) Locales() []string {
	locales := make([]string, 0, len(t.store))
	for l := range t.store {
		locales = append(locales, l)
	}
	return locales
}

// lookup resolves key across the requested then fallback locale, returning the key itself when neither has it.
func (t *Translator) lookup(locale, key string) string {
	if m, ok := t.store[locale]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}

	if locale != t.fallbackLocale {
		if m, ok := t.store[t.fallbackLocale]; ok {
			if v, ok := m[key]; ok {
				return v
			}
		}
	}

	return key
}

// substitute replaces {name} placeholders with the matching param in a single pass, so a substituted value containing
// braces is not re-expanded.
func substitute(s string, params map[string]string) string {
	if len(params) == 0 {
		return s
	}

	pairs := make([]string, 0, len(params)*2)
	for k, v := range params {
		pairs = append(pairs, "{"+k+"}", v)
	}

	return strings.NewReplacer(pairs...).Replace(s)
}
