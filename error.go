package i18n

// Base error type for low levels
type I18nError struct {
	code    int
	section string
	key     string
	values  map[string]interface{}
}

// Returns concatenated string "section.key"
func (e *I18nError) Error() string {
	if e.key != "" {
		return e.section + "." + e.key
	} else {
		return e.section
	}
}

// Creating *I18nError object
func NewErr(section string, key string, values ...M) *I18nError {
	if len(values) == 0 {
		return &I18nError{
			section: section,
			key:     key,
		}
	} else {
		return &I18nError{
			section: section,
			key:     key,
			values:  values[0],
		}
	}

}

// Creating *I18nError object with status code
func NewErrWithCode(code int, section string, key string, values ...M) *I18nError {
	if len(values) == 0 {
		return &I18nError{
			code:    code,
			section: section,
			key:     key,
		}
	} else {
		return &I18nError{
			code:    code,
			section: section,
			key:     key,
			values:  values[0],
		}
	}

}

// Set status code
// For example:  err.SetCode(http.StatusBadRequest)
func (e *I18nError) SetCode(code int) {
	e.code = code
}

// Set translatorsCollection section
func (e *I18nError) SetSection(section string) {
	e.section = section
}

// Set translatorsCollection key
func (e *I18nError) SetKey(key string) {
	e.key = key
}

// Set values for formatted output
func (e *I18nError) SetValues(values M) {
	e.values = values
}

// Return status code
func (e *I18nError) Code() int {
	return e.code
}

// Return translatorsCollection section
func (e *I18nError) Section() string {
	return e.section
}

// Return translatorsCollection key
func (e *I18nError) Key() string {
	return e.key
}

// Return values for formatted output
func (e *I18nError) Values() string {
	return e.key
}

// Return translated string from I18nError
func (e *I18nError) T(tr *Translator) string {
	return tr.T(e.section, e.key)
}

// Return translated formatted string from I18nError
func (e *I18nError) Tf(tr *Translator) string {
	return tr.Tf(e.section, e.key, e.values)
}

// Return translated error from I18nError
func (e *I18nError) ErrT(tr *Translator) error {
	return tr.ErrT(e.section, e.key)
}

// Return translated formatted error from I18nError
func (e *I18nError) ErrTf(tr *Translator) error {
	return tr.ErrTf(e.section, e.key, e.values)
}
