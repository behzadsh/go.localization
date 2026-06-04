package lang

// defaultTranslator backs the package-level Trans/TransBy helpers.
var defaultTranslator *Translator

// SetDefault sets the Translator used by the package-level Trans and TransBy.
func SetDefault(t *Translator) {
	defaultTranslator = t
}

// Trans translates key with the default Translator's default locale. If no default Translator has been set, it returns
// the key unchanged (no panic).
func Trans(key string, params ...map[string]string) string {
	if defaultTranslator == nil {
		return key
	}

	return defaultTranslator.Translate(key, params...)
}

// TransBy translates key in locale with the default Translator. If no default Translator has been set, it returns the
// key unchanged (no panic).
func TransBy(locale, key string, params ...map[string]string) string {
	if defaultTranslator == nil {
		return key
	}

	return defaultTranslator.TranslateBy(locale, key, params...)
}
