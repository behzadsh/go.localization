package lang

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var existingLanguages = []string{
	"en", "zh", "hi", "es", "fr", "ar", "fa",
}

func TestNewDictionaryManager(t *testing.T) {
	t.Run("correct-path", func(t *testing.T) {
		dm, err := newDictionaryManager(defaultPath)

		assert.Nil(t, err)
		assert.NotNil(t, dm)
		assert.NotEmpty(t, dm.path)

		// since we know the defaultPath exists dm.path should exist.
		i, err := os.Stat(dm.path)
		assert.Nil(t, err)
		assert.True(t, i.IsDir())
	})

	t.Run("incorrect-path", func(t *testing.T) {
		dm, err := newDictionaryManager("path/to/somewhere/else")

		assert.Nil(t, err)
		assert.NotNil(t, dm)
		assert.NotEmpty(t, dm.path)

		// since we know the defaultPath exists dm.path should exist.
		i, err := os.Stat(dm.path)
		assert.NotNil(t, err)
		assert.True(t, os.IsNotExist(err))
		assert.Nil(t, i)
	})
}

func TestDictionaryManager_find(t *testing.T) {
	dm, _ := newDictionaryManager(defaultPath)

	t.Run("success", func(t *testing.T) {
		for _, lang := range existingLanguages {
			t.Run(".yml>"+lang, func(t *testing.T) {
				key := "user.not_found"
				str, ok := dm.find(lang, key)

				assert.True(t, ok)
				assert.NotEqual(t, key, str)
			})

			t.Run(".yaml>"+lang, func(t *testing.T) {
				key := "validation.password_length_is_insufficient"
				str, ok := dm.find(lang, key)

				assert.True(t, ok)
				assert.NotEqual(t, key, str)
			})

			t.Run(".json>"+lang, func(t *testing.T) {
				key := "messages.insufficient_balance"
				str, ok := dm.find(lang, key)

				assert.True(t, ok)
				assert.NotEqual(t, key, str)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("incorrect-key", func(t *testing.T) {
			key := "incorrect_key"
			str, ok := dm.find(existingLanguages[0], key)

			assert.False(t, ok)
			assert.Equal(t, key, str)
		})

		t.Run("unknown-topic", func(t *testing.T) {
			key := "unknown_topic.some_key"
			str, ok := dm.find(existingLanguages[0], key)

			assert.False(t, ok)
			assert.Equal(t, key, str)
		})

		t.Run("unknown-topic-key", func(t *testing.T) {
			key := "user.unknown_key"
			str, ok := dm.find(existingLanguages[0], key)

			assert.False(t, ok)
			assert.Equal(t, key, str)
		})

		t.Run("unknown-locale", func(t *testing.T) {
			key := "user.not_found"
			str, ok := dm.find("nl", key)

			assert.False(t, ok)
			assert.Equal(t, key, str)
		})

		t.Run("malformed-yaml", func(t *testing.T) {
			key := "malformed_yaml.some_key"
			str, ok := dm.find(existingLanguages[0], key)

			assert.False(t, ok)
			assert.Equal(t, key, str)
		})

		t.Run("malformed-json", func(t *testing.T) {
			key := "malformed_json.some_key"
			str, ok := dm.find(existingLanguages[0], key)

			assert.False(t, ok)
			assert.Equal(t, key, str)
		})
	})
}
