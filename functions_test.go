package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetGlobalVars() {
	defaultTranslator = nil
	translationPath = ""
	defaultLocale = ""
	fallbackLocale = ""
}

func TestSetDefaultLocale(t *testing.T) {
	resetGlobalVars()
	SetDefaultLocale("nl")

	assert.Equal(t, "nl", defaultLocale)
}

func TestSetFallbackLocale(t *testing.T) {
	resetGlobalVars()
	SetFallbackLocale("zh")

	assert.Equal(t, "zh", fallbackLocale)
}

func TestSetTranslationFilesPath(t *testing.T) {
	resetGlobalVars()
	SetTranslationFilesPath("somewhere/else")

	assert.Equal(t, "somewhere/else", translationPath)
}

func TestInit(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		resetGlobalVars()
		err := Init()

		assert.Nil(t, err)
		assert.NotNil(t, defaultTranslator)
		assert.Equal(t, DefaultConfigs(), defaultTranslator.config)
	})

	t.Run("configure-by-functions", func(t *testing.T) {
		resetGlobalVars()
		SetTranslationFilesPath("somewhere/else")
		SetDefaultLocale("hi")
		SetFallbackLocale("ar")
		err := Init()

		assert.Nil(t, err)
		assert.NotNil(t, defaultTranslator)
		assert.NotEqual(t, DefaultConfigs(), defaultTranslator.config)
		assert.Equal(t, translationPath, defaultTranslator.config.TranslationPath)
		assert.Equal(t, defaultLocale, defaultTranslator.config.DefaultLocale)
		assert.Equal(t, fallbackLocale, defaultTranslator.config.FallbackLocale)
	})
}

func TestTrans(t *testing.T) {
	t.Run("initiated", func(t *testing.T) {
		resetGlobalVars()
		_ = Init()
		for key, data := range flatDictionary {
			translation := Trans(key, keyParam[key])

			assert.Equal(t, data[defaultTranslator.config.DefaultLocale], translation)
		}
	})

	t.Run("not-initiated", func(t *testing.T) {
		resetGlobalVars()
		for key := range flatDictionary {
			assert.Panics(t, func() {
				Trans(key, keyParam[key])
			})
		}
	})
}

func TestTransBy(t *testing.T) {
	t.Run("initiated", func(t *testing.T) {
		resetGlobalVars()
		_ = Init()
		for _, locale := range existingLanguages {
			for key, data := range flatDictionary {
				t.Run(locale, func(t *testing.T) {
					translation := TransBy(locale, key, keyParam[key])

					assert.Equal(t, data[locale], translation)
				})
			}
		}
	})

	t.Run("not-initiated", func(t *testing.T) {
		resetGlobalVars()
		for _, locale := range existingLanguages {
			for key := range flatDictionary {
				t.Run(locale, func(t *testing.T) {
					assert.Panics(t, func() {
						TransBy(locale, key, keyParam[key])
					})
				})
			}
		}
	})
}
