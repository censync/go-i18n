package i18n

import (
	"encoding/json"
	"testing"
)

var (
	testDefaultLocale = "en"
)

func initDict() error {
	collection := DictionaryCollection{
		"en": {
			"section.sub_section": {
				"key": "Translated field",
			},
			"fields.errors": {
				"to_short": "Field too short",
				"to_long":  "Field too long",
			},
		},
		"cz": {
			"section.sub_section": {
				"key": "Přeložené pole",
			},
			"fields.errors": {
				"to_short": "Pole je příliš krátké",
				"to_long":  "Pole je příliš dlouhé",
			},
		},
	}
	err := Init(testDefaultLocale, &collection, "en", "cz")

	if err != nil {
		return err
	}

	return nil
}

func TestI18nError_WithLocale_MarshalJSON(t *testing.T) {
	err := initDict()

	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		setErr  func() *I18nError
		want    string
		wantErr bool
	}{
		{
			name: "empty i18n error",
			setErr: func() *I18nError {
				return &I18nError{}
			},
			want: `null`,
		},
		{
			name: "filled i18n error with locale",
			setErr: func() *I18nError {
				anotherLocale := "cz"
				return &I18nError{
					locale: &anotherLocale,
					BaseError: &BaseError{
						section: "section.sub_section",
						key:     "key",
					},
				}
			},
			want: `"Přeložené pole"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.setErr()

			got, err := json.Marshal(e) // e.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("I18nError.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStr := string(got); gotStr != tt.want {
				t.Errorf("I18nError.MarshalJSON() = %v, want %v", gotStr, tt.want)
			}
		})
	}
}

func TestI18nMultipleError_WithLocale_MarshalJSON(t *testing.T) {
	err := initDict()

	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		setErr  func() *I18nMultipleError
		want    string
		wantErr bool
	}{
		{
			name: "empty i18n error",
			setErr: func() *I18nMultipleError {
				return &I18nMultipleError{}
			},
			want: `null`,
		},
		{
			name: "filled i18n multiple error",
			setErr: func() *I18nMultipleError {
				return &I18nMultipleError{
					locale: &testDefaultLocale,
					errors: map[string]*BaseError{
						"field1": {
							section: "fields.errors",
							key:     "to_short",
						},
					},
				}
			},
			want: `{"field1":"Field too short"}`,
		},
		{
			name: "partial fill i18n error",
			setErr: func() *I18nMultipleError {
				return &I18nMultipleError{
					locale: &testDefaultLocale,
					errors: map[string]*BaseError{
						"field1": {
							section: "fields.errors",
							key:     "to_long",
						},
					},
				}
			},
			want: `{"field1":"Field too long"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.setErr()

			got, err := json.Marshal(e) // e.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("I18nMultipleError.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStr := string(got); gotStr != tt.want {
				t.Errorf("I18nMultipleError.MarshalJSON() = %v, want %v", gotStr, tt.want)
			}
		})
	}
}
