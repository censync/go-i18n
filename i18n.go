package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

const dictExtension = "json"

// DictionaryEntry "key" => "translation"
type DictionaryEntry map[string]string

// Dictionary "section" => "key" => "translation"
type Dictionary map[string]DictionaryEntry

// DictionaryCollection "locale" => "section" => "key" => "translation"
type DictionaryCollection map[string]Dictionary

type M map[string]interface{}

type Translator struct {
	localeDictionary Dictionary
}

var (
	defLocale             string
	availableLocales      []string
	translatorsCollection map[string]*Translator
)

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
func InitFromDir(defaultLocale, translationsPath string, locales ...string) error {
	defLocale = defaultLocale
	if len(locales) > 0 {
		availableLocales = locales
	} else {
		availableLocales = getFilesFromDir(translationsPath)
	}

	localePath := translationsPath
	translatorsCollection = make(map[string]*Translator)
	for _, locale := range availableLocales {
		file, err := os.Open(localePath + `/` + locale + `.json`)
		if err != nil {
			return err
		}
		tmp := make(Dictionary)
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

	if _, ok := translatorsCollection[defLocale]; !ok {
		return errors.New("no dictionary for default language")
	}

	return nil
}

// Initialize translator with DictionaryCollection structure
//
// For example:
//
//	collection := i18n.DictionaryCollection{
//		"en" : i18n.DictionaryEntry{
//			"key" : "value",
//		},
//		"cz" : i18n.DictionaryEntry{
//			"key" : "hodnota",
//		},
//	}
func Init(defaultLocale string, dictCollection DictionaryCollection, locales ...string) error {
	defLocale = defaultLocale
	if _, ok := dictCollection[defLocale]; !ok {
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

	translatorsCollection = make(map[string]*Translator)
	for _, locale := range availableLocales {
		if dict, ok := dictCollection[locale]; ok {
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
func getFilesFromDir(path string) (locales []string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			suffixLen := len(dictExtension) + 1
			fNameLen := len(file.Name())
			if fNameLen > suffixLen && file.Name()[fNameLen-suffixLen+1:] == dictExtension {
				locales = append(locales, file.Name()[:fNameLen-suffixLen])
			}
		}
	}
	return
}

// New Returns Translator instance, if `locale` translatorsCollection exists.
// If translatorsCollection does not exist, returns translatorsCollection for default locale.
func New(locale string) *Translator {
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

// AvailableLocales Returns loaded locales
func AvailableLocales() []string {
	return availableLocales
}

// DefaultLocale Returns configured default locale
func DefaultLocale() string {
	return defLocale
}

// T Returns translated string
func (tr *Translator) T(section string, key string) string {
	if _, ok := tr.localeDictionary[section]; ok {
		if tr, ok := tr.localeDictionary[section][key]; ok {
			return tr
		} else {
			return section + `.` + key
		}
	} else {
		return section + `.` + key
	}
}

// Tf Returns translated formatted string
func (tr *Translator) Tf(section string, key string, values M) string {
	if tr, ok := tr.localeDictionary[section][key]; ok {
		for key, value := range values {
			switch reflect.TypeOf(value).Kind() {
			case reflect.String:
				tr = strings.Replace(tr, key, value.(string), -1)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				tr = strings.Replace(tr, key, fmt.Sprintf("%d", value), -1)
			case reflect.Float32, reflect.Float64:
				tr = strings.Replace(tr, key, fmt.Sprintf("%f", value), -1)
			}
		}
		return tr
	} else {
		return section + `.` + key
	}
}

// ErrT Returns translated error
func (tr *Translator) ErrT(section string, key string) error {
	return errors.New(tr.T(section, key))
}

// ErrTf Returns translated formatted error
func (tr *Translator) ErrTf(section string, key string, values M) error {
	return errors.New(tr.Tf(section, key, values))
}
