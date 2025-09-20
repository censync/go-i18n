package i18n

import "encoding/json"

var (
	nullJSON = []byte("null")
)

type BaseError struct {
	section string                 `json:"s,omitempty"`
	key     string                 `json:"k,omitempty"`
	values  map[string]interface{} `json:"v,omitempty"`
}

func (e *BaseError) Section() string {
	return e.section
}

func (e *BaseError) Key() string {
	return e.key
}

func (e *BaseError) Values() map[string]interface{} {
	return e.values
}

// I18nError Base error type
type I18nError struct {
	*BaseError
	code   int
	locale *string
}

func (e *BaseError) MarshalJSON() ([]byte, error) {
	if e == nil {
		return nullJSON, nil
	}

	return []byte("{\"" + e.section + "\":\"" + e.key + "\"}"), nil
}

func (e *BaseError) UnmarshalJSON(b []byte) error {
	var r struct {
		Section string                 `json:"s,omitempty"`
		Key     string                 `json:"k,omitempty"`
		Values  map[string]interface{} `json:"v,omitempty"`
	}
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	*e = BaseError{
		section: r.Section,
		key:     r.Key,
		values:  r.Values,
	}

	return nil
}

func (e *I18nError) MarshalJSON() ([]byte, error) {
	if e == nil {
		return nullJSON, nil
	}

	if e.locale != nil && *e.locale != "" {
		trStr := ""
		if len(e.values) == 0 {
			trStr = Get(*e.locale).T(e.section, e.key)
		} else {
			trStr = Get(*e.locale).Tf(e.section, e.key, e.values)
		}
		return []byte(`"` + trStr + `"`), nil
	}

	return e.BaseError.MarshalJSON()
}

// Error Returns concatenated string "section.key"
func (e *I18nError) Error() string {
	if e.key != "" {
		return e.section + "." + e.key
	} else {
		return e.section
	}
}

// NewErr Creates *I18nError object
func NewErr(section string, key string, values ...M) *I18nError {
	if len(values) == 0 {
		return &I18nError{
			BaseError: &BaseError{
				section: section,
				key:     key,
			},
		}
	} else {
		return &I18nError{
			BaseError: &BaseError{
				section: section,
				key:     key,
				values:  values[0],
			},
		}
	}

}

// NewErrWithCode Creates *I18nError object with status code
func NewErrWithCode(code int, section string, key string, values ...M) *I18nError {
	if len(values) == 0 {
		return &I18nError{
			code: code,
			BaseError: &BaseError{
				section: section,
				key:     key,
			},
		}
	} else {
		return &I18nError{
			code: code,
			BaseError: &BaseError{
				section: section,
				key:     key,
				values:  values[0],
			},
		}
	}

}

// Errors setters

// SetCode Set status code, e.g. `err.SetCode(http.StatusBadRequest)`
func (e *I18nError) SetCode(code int) {
	e.code = code
}

// SetLocale Defined priority locale for an error message
func (e *I18nError) SetLocale(locale string) {
	e.locale = &locale
}

// SetSection Set translatorsCollection section
func (e *I18nError) SetSection(section string) {
	e.section = section
}

// SetKey Set translatorsCollection key
func (e *I18nError) SetKey(key string) {
	e.key = key
}

// SetValues Set values for formatted output
func (e *I18nError) SetValues(values M) {
	e.values = values
}

// Errors builders

// WithCode Returns error with status code
func (e *I18nError) WithCode(code int) *I18nError {
	e.code = code
	return e
}

// WithLocale Returns error with locale
func (e *I18nError) WithLocale(locale string) *I18nError {
	e.locale = &locale
	return e
}

// WithSection Returns error with translatorsCollection section
func (e *I18nError) WithSection(section string) *I18nError {
	e.section = section
	return e
}

// WithKey Returns error with translatorsCollection key
func (e *I18nError) WithKey(key string) *I18nError {
	e.key = key
	return e
}

// WithValues Returns error with values for formatted output
func (e *I18nError) WithValues(values M) *I18nError {
	e.values = values
	return e
}

// Errors getters

// Code Returns status code
func (e *I18nError) Code() int {
	return e.code
}

// Section Returns translatorsCollection section
func (e *I18nError) Section() string {
	return e.section
}

// Key Returns translatorsCollection key
func (e *I18nError) Key() string {
	return e.key
}

// Values Returns values for formatted output
func (e *I18nError) Values() string {
	return e.key
}

// Errors translator functions

// T Returns translated string from I18nError
func (e *I18nError) T(tr *Translator) string {
	return tr.T(e.section, e.key)
}

// Tf Returns translated formatted string from I18nError
func (e *I18nError) Tf(tr *Translator) string {
	return tr.Tf(e.section, e.key, e.values)
}

// ErrT Returns translated error from I18nError
func (e *I18nError) ErrT(tr *Translator) error {
	return tr.ErrT(e.section, e.key)
}

// ErrTf Returns translated formatted error from I18nError
func (e *I18nError) ErrTf(tr *Translator) error {
	return tr.ErrTf(e.section, e.key, e.values)
}
