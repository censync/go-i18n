package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

const dictExtension = "json"

// DictionaryEntry "key" => "translation"
type DictionaryEntry map[string]string

// Dictionary "section" => "key" => "translation"
type Dictionary map[string]*DictionaryEntry

// DictionaryCollection "locale" => "section" => "key" => "translation"
type DictionaryCollection map[string]*Dictionary

type M map[string]interface{}

type Translator struct {
	localeDictionary *Dictionary
}

type TranslatorCollection map[string]*Translator

var (
	mu                    sync.RWMutex
	defLocale             string
	availableLocales      []string
	translatorsCollection TranslatorCollection
)

// InitFromDir
// Initialize dictionaries from JSON files
// Where file name is translation name:
// "en_US" =>"en_US.json", "cz" => "cz.json"
//
//		Sample translatorsCollection structure:
//
//		{
//			"en_US": {
//				"errors": {
//					"unknown": "Unknown error"
//				},
//				"errors.connections": {
//					"connections_limit": "Connections limit is {count}"
//				},
//				"form.signup": {
//					"welcome": "Welcome to registration"
//	             "disabled": "Registration is temporarily unavailable"
//				},
//				"form.login": {
//					"title": "Hello, {name}"
//				}
//			},
//			"cs_CZ": {
//				"errors": {
//					"unknown": "Neznámá chyba"
//				},
//				"errors.connections": {
//					"connections_limit": "Limit připojení je {count}"
//				},
//				"form.signup": {
//					"welcome": "Vítejte v registraci"
//	             "disabled": "Registrace je dočasně nedostupná"
//				},
//				"form.login": {
//					"title": "Ahoj, {name}"
//				}
//			}
//		}
func InitFromDir(defaultLocale, translationsPath string, locales ...string) (err error) {
	mu.Lock()
	defer mu.Unlock()

	defLocale = defaultLocale
	if len(locales) > 0 {
		availableLocales = locales
	} else {
		availableLocales, err = getFilesFromDir(translationsPath)
		if err != nil {
			return err
		}
	}

	translatorsCollection = make(map[string]*Translator)
	for _, locale := range availableLocales {
		file, err := os.Open(filepath.Join(translationsPath, locale+".json"))
		if err != nil {
			return err
		}
		tmp := &Dictionary{}
		err = json.NewDecoder(file).Decode(&tmp)
		_ = file.Close()
		if err != nil {
			return err
		}
		tr := &Translator{
			localeDictionary: tmp,
		}
		translatorsCollection[locale] = tr
	}

	if _, ok := translatorsCollection[defaultLocale]; !ok {
		return errors.New("no dictionary for default language")
	}

	return nil
}

// Init
// Initialize translator with DictionaryCollection structure
//
// For example:
//
//	collection := i18n.DictionaryCollection{
//		"en": {
//			"section": {
//				"key": "value",
//			},
//		},
//		"cz": {
//			"section": {
//				"key": "hodnota",
//			},
//		},
//	}
func Init(defaultLocale string, dictCollection *DictionaryCollection, locales ...string) error {
	mu.Lock()
	defer mu.Unlock()

	defLocale = defaultLocale
	if _, ok := (*dictCollection)[defaultLocale]; !ok {
		return errors.New("no dictionary for default language")
	}

	if len(locales) > 0 {
		availableLocales = locales
	} else {
		availableLocales = dictCollection.getLocales()
	}

	if len(availableLocales) == 0 {
		return errors.New("available locales not set")
	}

	translatorsCollection = make(TranslatorCollection)
	for _, locale := range availableLocales {
		if dict, ok := (*dictCollection)[locale]; ok {
			tr := &Translator{
				localeDictionary: dict,
			}
			translatorsCollection[locale] = tr
		}
	}
	return nil
}

// getLocales Returns available locales for dictionaries collection
func (c *DictionaryCollection) getLocales() (locales []string) {
	if c != nil {
		for locale := range *c {
			locales = append(locales, locale)
		}
	}
	return
}

// getFilesFromDir Returns available locales for dictionary
func getFilesFromDir(path string) (locales []string, err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	ext := "." + dictExtension
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.HasSuffix(name, ext) && len(name) > len(ext) {
				locales = append(locales, strings.TrimSuffix(name, ext))
			}
		}
	}
	return
}

// Get Returns Translator instance, if `locale` translatorsCollection exists.
// If translatorsCollection does not exist, returns translatorsCollection for default locale.
func Get(locale string) *Translator {
	mu.RLock()
	defer mu.RUnlock()

	if translatorsCollection == nil {
		panic("translator not initialized")
	}
	if _, ok := translatorsCollection[locale]; ok {
		return translatorsCollection[locale]
	} else {
		if _, ok := translatorsCollection[defLocale]; ok {
			return translatorsCollection[defLocale]
		} else {
			return &Translator{}
		}
	}
}

// New
// Deprecated: must be used Get
func New(locale string) *Translator {
	return Get(locale)
}

// AvailableLocales Returns loaded locales
func AvailableLocales() []string {
	mu.RLock()
	defer mu.RUnlock()

	return availableLocales
}

// DefaultLocale Returns configured default locale
func DefaultLocale() string {
	mu.RLock()
	defer mu.RUnlock()

	return defLocale
}

// T Returns translated string
func (tr *Translator) T(section string, key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if tr.localeDictionary == nil {
		return section + "." + key
	}
	if entries, ok := (*tr.localeDictionary)[section]; ok {
		if entry, ok := (*entries)[key]; ok {
			return entry
		}
	}
	return section + "." + key
}

// Tf Returns translated formatted string
func (tr *Translator) Tf(section string, key string, values M) string {
	mu.RLock()
	defer mu.RUnlock()

	if tr.localeDictionary == nil {
		return section + "." + key
	}
	entries, ok := (*tr.localeDictionary)[section]
	if !ok {
		return section + "." + key
	}
	result, ok := (*entries)[key]
	if !ok {
		return section + "." + key
	}
	for placeholder, value := range values {
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			result = strings.Replace(result, placeholder, value.(string), -1)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result = strings.Replace(result, placeholder, fmt.Sprintf("%d", value), -1)
		case reflect.Float32, reflect.Float64:
			result = strings.Replace(result, placeholder, fmt.Sprintf("%f", value), -1)
		default:
			result = strings.Replace(result, placeholder, fmt.Sprintf("%v", value), -1)
		}
	}
	return result
}

// ErrT Returns translated error
func (tr *Translator) ErrT(section string, key string) error {
	return errors.New(tr.T(section, key))
}

// ErrTf Returns translated formatted error
func (tr *Translator) ErrTf(section string, key string, values M) error {
	return errors.New(tr.Tf(section, key, values))
}
