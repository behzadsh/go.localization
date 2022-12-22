package lang

import "strings"

const (
	defaultPath = "resources/lang"
)

type Config struct {
	// Path to the where translations files are stored.
	TranslationPath string

	// The locale used for translation by default.
	DefaultLocale string

	// The fallback locale used when no translation found for default locale.
	FallbackLocale string
}

// DefaultConfigs creates a new translation config instance with default config values.
func DefaultConfigs() Config {
	return Config{
		TranslationPath: defaultPath,
		DefaultLocale:   "en",
		FallbackLocale:  "en",
	}
}

func (c *Config) defaultValues() {
	if c.TranslationPath == "" {
		c.TranslationPath = defaultPath
	}

	if c.DefaultLocale == "" {
		c.DefaultLocale = "en"
	}

	if c.FallbackLocale == "" {
		c.FallbackLocale = "en"
	}
}

// Translator is a tool to translate the translation key based on localization settings.
type Translator struct {
	config            Config
	dictionaryManager *dictionaryManager
}

func NewTranslator(cfg Config) (*Translator, error) {
	cfg.defaultValues()

	dm, err := newDictionaryManager(cfg.TranslationPath)
	if err != nil {
		return nil, err
	}

	return &Translator{
		config:            cfg,
		dictionaryManager: dm,
	}, nil
}

func (t *Translator) Translate(key string, params ...map[string]string) string {
	return t.TranslateBy(t.config.DefaultLocale, key, params...)
}

func (t *Translator) TranslateBy(locale, key string, params ...map[string]string) string {
	var p map[string]string
	if len(params) > 0 {
		p = params[0]
	}

	tr, ok := t.dictionaryManager.find(locale, key)
	if !ok {
		tr, _ = t.dictionaryManager.find(t.config.FallbackLocale, key)
	}

	for k, v := range p {
		tr = strings.Replace(tr, ":"+k+":", v, -1)
	}

	return tr
}
