package i18n

import "encoding/json"

var (
	multipleDefaultErrorField = "_summary"
)

type I18nMultipleError struct {
	code   int
	locale *string
	errors map[string]*BaseError
}

func NewMultipleEmptyErr() *I18nMultipleError {
	return &I18nMultipleError{
		errors: map[string]*BaseError{},
	}
}

func NewMultipleErr(field, section string, key string, values ...M) *I18nMultipleError {
	if len(values) == 0 {
		return &I18nMultipleError{
			errors: map[string]*BaseError{
				field: {
					section: section,
					key:     key,
				},
			},
		}
	} else {
		return &I18nMultipleError{
			errors: map[string]*BaseError{
				field: {
					section: section,
					key:     key,
					values:  values[0],
				},
			},
		}
	}
}

func NewMultipleDefaultErr(section string, key string, values ...M) *I18nMultipleError {
	if len(values) == 0 {
		return &I18nMultipleError{
			errors: map[string]*BaseError{
				multipleDefaultErrorField: {
					section: section,
					key:     key,
				},
			},
		}
	} else {
		return &I18nMultipleError{
			errors: map[string]*BaseError{
				multipleDefaultErrorField: {
					section: section,
					key:     key,
					values:  values[0],
				},
			},
		}
	}
}

func (e *I18nMultipleError) UnmarshalJSON(b []byte) error {
	var r struct {
		Section string                 `json:"s,omitempty"`
		Key     string                 `json:"k,omitempty"`
		Values  map[string]interface{} `json:"v,omitempty"`
	}
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	return nil
}

func (e *I18nMultipleError) Add(field, section string, key string, values ...M) *I18nMultipleError {
	if e.code == 0 {
		e.code = 400
	}

	if len(values) == 0 {
		e.errors[field] = &BaseError{
			section: section,
			key:     key,
		}
	} else {
		e.errors[field] = &BaseError{
			section: section,
			key:     key,
			values:  values[0],
		}
	}
	return e
}

func (e *I18nMultipleError) AddDefault(section string, key string, values ...M) *I18nMultipleError {
	if e.code == 0 {
		e.code = 500
	}

	if len(values) == 0 {
		e.errors[multipleDefaultErrorField] = &BaseError{
			section: section,
			key:     key,
		}
	} else {
		e.errors[multipleDefaultErrorField] = &BaseError{
			section: section,
			key:     key,
			values:  values[0],
		}
	}
	return e
}

func (e *I18nMultipleError) AddDefaultErr(srcErr error) *I18nMultipleError {
	mErr, ok := srcErr.(*I18nMultipleError)
	if ok {
		*e = *mErr
	} else {
		e.errors[multipleDefaultErrorField] = &BaseError{
			section: "_error",
			key:     srcErr.Error(),
		}
	}
	return e
}

func (e *I18nMultipleError) HasErrors() bool {
	return len(e.errors) > 0
}

func (e *I18nMultipleError) Code() int {
	return e.code
}

func (e *I18nMultipleError) Locale() *string {
	return e.locale
}

func (e *I18nMultipleError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "cannot marshal I18nMultipleError: " + err.Error()
	}
	return string(b)
}

func (e *I18nMultipleError) Errors() map[string]*BaseError {
	return e.errors
}

func (e *I18nMultipleError) MarshalJSON() ([]byte, error) {
	if e == nil || len(e.errors) == 0 {
		return nullJSON, nil
	}

	result := make(map[string]json.RawMessage, len(e.errors))
	for field, baseErr := range e.errors {
		if e.locale != nil && *e.locale != "" {
			trStr := ""
			if len(baseErr.values) == 0 {
				trStr = Get(*e.locale).T(baseErr.section, baseErr.key)
			} else {
				trStr = Get(*e.locale).Tf(baseErr.section, baseErr.key, baseErr.values)
			}
			b, err := json.Marshal(trStr)
			if err != nil {
				return nil, err
			}
			result[field] = b
		} else {
			b, err := baseErr.MarshalJSON()
			if err != nil {
				return nil, err
			}
			result[field] = b
		}
	}

	return json.Marshal(result)
}
