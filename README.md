# go.localization v2

![Coverage](https://img.shields.io/badge/Coverage-97.8%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/behzadsh/go.localization)](https://goreportcard.com/report/github.com/behzadsh/go.localization)

A small, idiomatic internationalization (i18n) library for Go. Translations live
in one file per locale; lookups are dot-paths; placeholders are `{name}`.

```bash
go get github.com/behzadsh/go.localization/v2
```

> v2 is a published major version. The import path carries `/v2`, but the package
> name is still `lang`. v1 (`github.com/behzadsh/go.localization`) is unaffected.

## Translation files

One file per locale, named by the locale (`en.yaml`, `fr.json`, ...). Each file is
an arbitrarily nested tree of keys whose leaves are strings.

```yaml
# en.yaml
user:
  not_found: User not found
  greeting: "Hello {name}"
  profile:
    title: "{name}'s profile"
validation:
  required: "The {field} field is required."
```

Look a value up by its dot-path: `user.profile.title`. `{name}` placeholders are
replaced from the params map; an unprovided placeholder is left as-is.

YAML (`.yaml`, `.yml`) and JSON (`.json`) are supported out of the box. Files with
any other extension are ignored, so a `README.md` can sit beside your locales.

## Usage

### From a directory

```go
t, err := lang.NewFromDir("resources/lang")
if err != nil {
log.Fatal(err)
}

t.Translate("user.not_found") // "User not found"
t.Translate("user.greeting", map[string]string{"name": "Sam"}) // "Hello Sam"
t.TranslateBy("fr", "user.greeting", map[string]string{"name": "Léa"})
```

### From embedded files

The loader takes any `io/fs.FS`, so translations can be compiled into the binary:

```go
//go:embed lang/*.yaml
var langFS embed.FS

sub, _ := fs.Sub(langFS, "lang")
t, _ := lang.New(sub, lang.WithDefaultLocale("en"))
```

### Package-level default

For app-wide convenience, register a default and use the package helpers. Unlike
v1, calling a helper before setting the default does **not** panic — it returns
the key unchanged.

```go
t, _ := lang.NewFromDir("resources/lang")
lang.SetDefault(t)

lang.Trans("user.not_found")
lang.TransBy("fr", "user.greeting", map[string]string{"name": "Léa"})
```

## Options

```go
lang.New(fsys,
lang.WithDefaultLocale("en"), // locale used by Translate / Trans (default "en")
lang.WithFallback("en"),      // consulted when a key is missing (default = default locale)
lang.WithDecoder(".toml", toml.Unmarshal), // register another file format
)
```

`Decoder` is just `func(data []byte, v any) error` — the same shape as
`yaml.Unmarshal` and `json.Unmarshal`, so most unmarshalers register in one line.

## Behavior

- **Lookup order:** requested locale → fallback locale → the key itself.
- **Loading is eager:** `New` reads every locale file up front, so a malformed
  file surfaces immediately rather than at the first translation.
- **Leaves must be strings:** quote numbers and booleans (`label: "2.0"`), or
  `New` returns an error naming the offending key and file.
- **Concurrency:** a `Translator` is immutable after `New` and safe for use by
  multiple goroutines.

## Migrating from v1

|              | v1                                   | v2                                                          |
|--------------|--------------------------------------|-------------------------------------------------------------|
| File layout  | per-topic, `key: {en: ..., fr: ...}` | per-locale, nested key tree                                 |
| Key depth    | 2 levels only                        | arbitrary nesting                                           |
| Placeholders | `:name:`                             | `{name}`                                                    |
| Source       | path string                          | any `io/fs.FS` (`os.DirFS`, `embed.FS`, ...) or path string |
| Config       | `Config` struct + setters            | functional options                                          |

Before (v1 `user.yaml`):

```yaml
not_found:
  en: User not found
  fr: Utilisateur introuvable
```

After (v2 `en.yaml` / `fr.yaml`):

```yaml
# en.yaml
user:
  not_found: User not found
# fr.yaml
user:
  not_found: Utilisateur introuvable
```
