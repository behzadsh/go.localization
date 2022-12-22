package lang

import "strings"

// Translator is a tool to translate the translation key based on localization settings.
type Translator struct {
	config            Config
	dictionaryManager *dictionaryManager
}

func NewTranslator(cfg Config) (*Translator, error) {
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

func (t *Translator) TranslateBy(locale string, key string, params ...map[string]string) string {
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
