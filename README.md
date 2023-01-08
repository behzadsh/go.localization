# go.localization
![Coverage](https://img.shields.io/badge/Coverage-97.3%25-brightgreen)

[![Go Report Card](https://goreportcard.com/badge/github.com/behzadsh/go.localization)](https://goreportcard.com/report/github.com/behzadsh/go.localization)

The `go.localization` package provides a simple and convenient way to retrieve strings
in various languages, allowing you to easily build multilingual applications.

#### Key features
- Support strings with named variable.
- Support `json` and `yaml` translation files.

## Installation
to install `go.localization` package, run the following command

```
go get -u github.com/behzadsh/go.localization
```

## Defining Translation Strings
The `go.localization` package, loads the translations from `json` or `yaml` files, which by default should be
stored in `{ProjectRoot}/resources/lang` directory. The path to translation files are configurable, we discuss it later.
No matter you choose `json` or `yaml` the structure of the translation file must be like the following example.

The `yaml` file:
```yaml
key_string:
  locale1: "translation in locale1"
  locale2: "translation in locale2"
```

The `json` file:
```json
{
  "key_string": {
    "locale1": "translation in locale1"
    "locale2": "translation in locale2"
  }
}
```

### conventions
* It is recommended to use `snake_case` names for the translation file name and the translation key.
* It is recommended to use 2 letter abbreviation for locales, e.g. [ISO 639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes).

## Retrieving Translation Strings

For retrieving translation strings, you need to pass the translation key and the language to the translation function.
The translation key is consist of the translation file name and the translation key concatenated by a dot.
For example "user.not_found" as a translation key, means that the translation is located in a file named
`user` under the key `not_found`. With this in mind lets continue.

### Using default configuration and helper functions
The easiest way to use the `go.localization` package is that you store your translation files in default path
(as mentioned above in the `{ProjectRoot}/resources/lang`). In this way, you can simply translate the keys like this:

main.go:
```go
package main

import (
	"fmt"
	
	lang "github.com/behzadsh/go.localization"
)

func main() {
	lang.Init()

	// Simple translation
	fmt.Println(lang.Trans("user.not_found"))           // User not found!
	fmt.Println(lang.Trans("errors.invalid_payload"))   // The payload is invalid!
	fmt.Println(lang.TransBy("zh", "user.not_found"))   // 找不到用户!
	
	// Translation with parameter
	fmt.Println(lang.Trans("validation.required", map[string]string{"field": "email"}))         // The field 'email' is required.
	fmt.Println(lang.TransBy("zh", "validation.required", map[string]string{"field": "email"})) // 字段 'email' 是必需的。
}
```
./resources/lang/user.yaml:
```yaml
not_found:
  en: "User not found!"
  zh: "找不到用户!"
```
./resources/lang/errors.yaml:
```yaml
invalid_payload:
  en: "The payload is invalid!"
  zh: "负载无效!"
```
./resources/lang/validations.yaml:
```yaml
required:
  en: "The field ':field:' is required."
  zh: "字段 ':field:' 是必需的。"
```

### Using the translator struct directly

```go
package main

import (
	"fmt"
	"log"

	lang "github.com/behzadsh/go.localization"
)

func main() {
	tr, err := lang.NewTranslator(lang.DefaultConfigs())
	if err != nil {
		log.Fatal(err)
	}

	// Simple translation
	fmt.Println(tr.Translate("user.not_found"))         // User not found!
	fmt.Println(tr.Translate("errors.invalid_payload")) // The payload is invalid!
	fmt.Println(tr.TranslateBy("zh", "user.not_found")) // 找不到用户!

	// Translation with parameter
	fmt.Println(tr.Translate("validation.required", map[string]string{"field": "email"}))         // The field 'email' is required.
	fmt.Println(tr.TranslateBy("zh", "validation.required", map[string]string{"field": "email"})) // 字段 'email' 是必需的。
}
```

## Configuration

You can configure the following options:
* Translation files path: The path to where the translation files are stored.
* The default language: The default language used by functions like `Trans` and `Translate`
* The fallback language: The language used when the no translation found for default or given language.

If you rather use helper function, you can customize these options by these functions before calling `Init()`

```go
package main

import (
	lang "github.com/behzadsh/go.localization"
)

func main() {
	lang.SetDefaultLocale("zh")
	lang.SetFallbackLocale("fr")
	lang.SetTranslationFilesPath("relative/path/to/somewhere/else")
	lang.Init()

	// ...
}
```

Or if you want to use the translator struct directly, you can pass your custom configuration to the constructor.

```go
package main

import (
	"log"

	lang "github.com/behzadsh/go.localization"
)

func main() {
	tr, err := lang.NewTranslator(lang.Config{
		TranslationPath: "relative/path/to/somewhere/else",
		DefaultLocale:   "zh",
		FallbackLocale:  "fr",
	})
	if err != nil {
		log.Fatal(err)
	}

	// ...
}
```
