package lang_test

import (
	"fmt"

	lang "github.com/behzadsh/go.localization/v2"
)

func ExampleTranslator() {
	t, err := lang.NewFromDir("testdata/lang")
	if err != nil {
		panic(err)
	}

	fmt.Println(t.Translate("user.not_found"))
	fmt.Println(t.Translate("user.greeting", map[string]string{"name": "Sam"}))
	fmt.Println(t.TranslateBy("fr", "user.greeting", map[string]string{"name": "Léa"}))
	// "fr" has no validation.required, so it falls back to the default locale.
	fmt.Println(t.TranslateBy("fr", "validation.required", map[string]string{"field": "email"}))
	// Output:
	// User not found
	// Hello Sam
	// Bonjour Léa
	// The email field is required.
}
