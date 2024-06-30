package i18n

var (
	multipleDefaultErrorField = "_summary"
)

type I18nMultipleError struct {
	code   int
	locale *string
	errors map[string]*baseError
}

func NewMultipleErr(field, section string, key string, values ...M) *I18nMultipleError {
	if len(values) == 0 {
		return &I18nMultipleError{
			errors: map[string]*baseError{
				field: {
					section: section,
					key:     key,
				},
			},
		}
	} else {
		return &I18nMultipleError{
			errors: map[string]*baseError{
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
			errors: map[string]*baseError{
				multipleDefaultErrorField: {
					section: section,
					key:     key,
				},
			},
		}
	} else {
		return &I18nMultipleError{
			errors: map[string]*baseError{
				multipleDefaultErrorField: {
					section: section,
					key:     key,
					values:  values[0],
				},
			},
		}
	}
}

func (e *I18nMultipleError) MarshalJSON() ([]byte, error) {
	if e == nil {
		return nullJSON, nil
	}

	result := make([]byte, 0)

	if e.locale != nil && *e.locale != "" {
		for k, v := range e.errors {
			trStr := ""
			if len(v.values) == 0 {
				trStr = Get(*e.locale).T(v.section, v.key)
			} else {
				trStr = Get(*e.locale).Tf(v.section, v.key, v.values)
			}
			result = append(result, []byte(`"`+k+`":"`+trStr+`",`)...)
		}
	} else {
		for k, v := range e.errors {
			entry := []byte(`"` + k + `":{"` + v.section + `":"` + v.key + `"},`)
			result = append(result, entry...)
		}
	}

	l := len(result)

	if l == 0 {
		return nullJSON, nil
	}

	result = result[:len(result)-1]
	result = append([]byte{'{'}, result...)
	result = append(result, '}')

	return result, nil
}

func (e *I18nMultipleError) HasErrors() bool {
	return len(e.errors) > 0
}
