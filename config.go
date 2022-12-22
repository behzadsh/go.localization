package lang

type Config struct {
	// Path to the where translations files are stored.
	TranslationPath string

	// The locale used for translation by default.
	DefaultLocale string

	// The fallback locale used when no translation found for default locale.
	FallbackLocale string
}

const (
	defaultPath = "resources/lang"
)

// DefaultConfigs creates a new translation config instance with default config values.
func DefaultConfigs() Config {
	return Config{
		TranslationPath: defaultPath,
		DefaultLocale:   "en",
		FallbackLocale:  "en",
	}
}
