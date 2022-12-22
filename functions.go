package lang

var defaultTranslator *Translator

var translationPath, defaultLocale, fallbackLocale string

// SetTranslationFilesPath set translation files path.
func SetTranslationFilesPath(p string) {
	translationPath = p
}

// SetDefaultLocale set translation default locale.
func SetDefaultLocale(locale string) {
	defaultLocale = locale
}

// SetFallbackLocale set translation fallback locale.
func SetFallbackLocale(locales string) {
	fallbackLocale = locales
}

// Init initiates the global translator instance.
func Init() error {
	cfg := DefaultConfigs()

	if translationPath != "" {
		cfg.TranslationPath = translationPath
	}

	if defaultLocale != "" {
		cfg.DefaultLocale = defaultLocale
	}

	if fallbackLocale != "" {
		cfg.FallbackLocale = fallbackLocale
	}

	var err error
	defaultTranslator, err = NewTranslator(cfg)

	return err
}

// Trans is an alias for Translator.Translate() function.
func Trans(key string, params ...map[string]string) string {
	return defaultTranslator.Translate(key, params...)
}

// TransBy is an alias for Translator.TranslateBy() function.
func TransBy(lang, key string, params ...map[string]string) string {
	return defaultTranslator.TranslateBy(lang, key, params...)
}
