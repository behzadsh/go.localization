package lang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// supportedFileTypes is list of supported translation file types.
var supportedFileTypes = []string{
	".yaml",
	".yml",
	".json",
}

// dictionary is key value variable that the key represents the key and the value
// represent another key value, in which the key represents the locale and
// the value is the translation string for that locale.
type dictionary map[string]map[string]string

// shelf is a key value variable that the key represents the topic (e.g. validation)
// and the value represent the dictionary.
//
// Topics are a group of keys stored in a file.
type shelf map[string]dictionary

type dictionaryManager struct {
	path  string
	shelf shelf
}

func newDictionaryManager(path string) (*dictionaryManager, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &dictionaryManager{
		path:  root + string(filepath.Separator) + path,
		shelf: make(shelf),
	}, nil
}

func (m *dictionaryManager) find(locale, key string) (string, bool) {
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return key, false
	}

	topic := parts[0]
	topicKey := parts[1]

	if tr, ok := m.tryToFindInShelf(topic, topicKey, locale); ok {
		return tr, true
	}

	if err := m.loadTopicFile(topic); err != nil {
		return key, false
	}

	return m.tryToFindInShelf(topic, topicKey, locale)
}

func (m *dictionaryManager) tryToFindInShelf(topic, topicKey, locale string) (string, bool) {
	dict, ok := m.shelf[topic]
	if !ok {
		return topic + "." + topicKey, false
	}

	keyDict, ok := dict[topicKey]
	if !ok {
		return topic + "." + topicKey, false
	}

	tr, ok := keyDict[locale]
	if !ok {
		return topic + "." + topicKey, false
	}

	return tr, true
}

func (m *dictionaryManager) loadTopicFile(topic string) error {
	for _, fileType := range supportedFileTypes {
		path, _ := filepath.Abs(m.path + string(filepath.Separator) + topic + fileType)

		byteVal, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var result dictionary
		switch fileType {
		case ".yaml", ".yml":
			if err = yaml.Unmarshal(byteVal, &result); err != nil {
				return err
			}
		case ".json":
			if err = yaml.Unmarshal(byteVal, &result); err != nil {
				return err
			}
		}

		m.shelf[topic] = result
		return nil
	}

	return fmt.Errorf("translation file not found")
}
